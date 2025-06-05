package checker

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Создает готовое эталонное решение из кусков кода для каждого задания на питоне
// Каждый кусок устроен так, что он может быть совместим с любым другим куском
// Благодаря чему обеспечивается гибкость в создании решения
func (c *Checker) createSolution() error {
	// Создаем пустой файл для итогового решения
	refSolutionPath := fmt.Sprintf("./%v/ref-solution.py", c.tempDirPath)
	refSolution, err := os.Create(refSolutionPath)
	if err != nil {
		e := fmt.Errorf("error with creating solution file: %v", err)
		log.Println(e)
		return e
	}
	defer refSolution.Close()

	// Проходимся по нужным вариациям и заполняем итоговый файл
	for i := 1; i <= c.labConfig.TasksCount; i++ {
		// Форматируем путь и открываем нужный кусок кода
		path := fmt.Sprintf("./labs/lab%v/task%v/var%v/solution.py",
			c.lab.LabNum, i, c.lab.Tasks[fmt.Sprintf("task%v", i)])
		file, err := os.Open(path)
		if err != nil {
			e := fmt.Errorf("error with opening part-solution file: %v", err)
			log.Println(e)
			return e
		}
		solutionPart, err := io.ReadAll(file)
		if err != nil {
			e := fmt.Errorf("error with reading code from part-solution file: %v", err)
			log.Println(e)
			return e
		}

		// Записываем полученный кусок в итоговый файл
		_, err = refSolution.WriteString(string(solutionPart))
		if err != nil {
			e := fmt.Errorf("error with creating solution file: %v", err)
			log.Println(e)
			return e
		}
	}

	return nil
}

// Запускает эталонное решение с конкретными входными данными (тест-кейсом)
func (c *Checker) runRefSolution(absPath, input string) (string, error) {
	// Создаем команду и перенаправляем ввод и вывод
	var cmd *exec.Cmd
	switch c.OS {
	case "darwin":
		cmd = exec.Command("python3", absPath)
	default:
		cmd = exec.Command("python", absPath)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Принудительно изменяем кодировку для корректного отображения русского языка
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	// Запускаем решение и возвращаем вывод
	err := cmd.Run()
	if err != nil {
		e := fmt.Errorf("error with running ref solution: %v\nStderr: %s", err, stderr.String())
		log.Println(e)
		return "", e
	}
	return stdout.String(), nil
}

// Запускаем тесты для проверки кода ученика
func (c *Checker) runTests() (string, error) {
	// Создаем input и output файлы для дальнейшего заполнения вводом и выводом соответственно
	inputPath := fmt.Sprintf("./%v/input.txt", c.tempDirPath)
	inputFile, err := os.Create(inputPath)
	if err != nil {
		e := fmt.Errorf("error with creating input file: %v", err)
		log.Println(e)
		return "", e
	}
	defer inputFile.Close()

	outputPath := fmt.Sprintf("./%v/output.txt", c.tempDirPath)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		e := fmt.Errorf("error with creating output file: %v", err)
		log.Println(e)
		return "", e
	}
	defer outputFile.Close()

	// Получаем абсолютный путь для кода в зависимости от ОС
	var execPath string
	switch c.OS {
	case "windows":
		execPath = fmt.Sprintf("./%v/code.exe", c.tempDirPath)
	case "linux":
		execPath = fmt.Sprintf("./%v/code", c.tempDirPath)
	case "darwin":
		execPath = fmt.Sprintf("./%v/code", c.tempDirPath)
	default:
		err := fmt.Errorf("please use the appropriate OS: Windows/Linux")
		log.Printf("the program does not support the current OS: %v", c.OS)
		return "", err
	}

	absCodePath, _ := filepath.Abs(execPath)
	refPath := fmt.Sprintf("./%v/ref-solution.py", c.tempDirPath)
	absRefPath, _ := filepath.Abs(refPath)

	// Проходимся по тест-кейсам и запускаем программы
	for i := range c.labConfig.TestCases {
		testCase := c.labConfig.TestCases[i]

		// При запуске программ используется многопоточность, позволяя одновременно
		// прогнать тест через эталонное решение и через код ученика.
		// Даже на небольших работах это экономит ~100мс

		// Здесь создается отдельный канал и горутина (аналогично корутине, но намного легковеснее, т. к. корутины
		// используют уже созданные ранее треды) для эталонного решения
		refSolCh := make(chan struct {
			out string
			err error
		})

		go func() {
			// Запускаем проверку
			refOut, err := c.runRefSolution(absRefPath, testCase)
			refSolCh <- struct {
				out string
				err error
			}{refOut, err}
		}()

		// Здесь то же самое, но для кода ученика
		codeCh := make(chan struct {
			err error
			out string
		})

		go func() {
			// Записываем тест-кейсы по очереди в input файл
			if err := os.WriteFile(inputPath, []byte(testCase), 0644); err != nil {
				e := fmt.Errorf("error with wtriting test-case to the input file: %v", err)
				log.Println(e)
				codeCh <- struct {
					err error
					out string
				}{e, ""}
				return
			}

			// Сбрасываем позицию input в начало и очищаем файл
			_, err = inputFile.Seek(0, 0)
			if err != nil {
				e := fmt.Errorf("error with setting start pos in input file: %v", err)
				log.Println(e)
				codeCh <- struct {
					err error
					out string
				}{e, ""}
				return
			}

			err = outputFile.Truncate(0) // очищаем output файл (чтобы следующий вывод программы не нарушил проверку)
			if err != nil {
				e := fmt.Errorf("error with cleaning output file: %v", err)
				log.Println(e)
				codeCh <- struct {
					err error
					out string
				}{e, ""}
				return
			}

			_, err = outputFile.Seek(0, 0)
			if err != nil {
				e := fmt.Errorf("error with setting start pos in output file: %v", err)
				log.Println(e)
				codeCh <- struct {
					err error
					out string
				}{e, ""}
				return
			}

			// Запускаем код ученика
			err := c.runCode(absCodePath, inputFile, outputFile)
			if err != nil {
				log.Println(fmt.Errorf("error with student code: %v", err))
				codeCh <- struct {
					err error
					out string
				}{err, ""}
				return
			}

			// Сохраняем данные из буфера в файл и смещаем позицию чтения на 0,
			// чтобы получить вывод программы ученика
			err = outputFile.Sync()
			if err != nil {
				e := fmt.Errorf("error with sync output file: %v", err)
				log.Println(e)
				codeCh <- struct {
					err error
					out string
				}{e, ""}
				return
			}

			_, err = outputFile.Seek(0, 0)
			if err != nil {
				e := fmt.Errorf("error with setting start pos in output file: %v", err)
				log.Println(e)
				codeCh <- struct {
					err error
					out string
				}{e, ""}
				return
			}

			codeOut, err := io.ReadAll(outputFile)
			if err != nil {
				e := fmt.Errorf("error with reading from output file: %v", err)
				log.Println(e)
				codeCh <- struct {
					err error
					out string
				}{e, ""}
				return
			}

			codeCh <- struct {
				err error
				out string
			}{nil, string(codeOut)}
		}()

		// Ожидаем выполнения обоих программ
		code := <-codeCh

		refSol := <-refSolCh
		if refSol.err != nil {
			return "", err
		}

		// Проверяем выполнение кода ученика на ошибки
		if code.err != nil {
			message := fmt.Sprintf("Неверный ответ.\nТест-кейс:\n%v\nОжидалось:\n%v\nПолучено:\n%v", testCase, refSol.out, code.err)
			return message, nil
		}

		// Сравниваем вывод программ
		if isCorrect := c.validate(code.out, refSol.out); !isCorrect {
			out := c.redact(code.out)
			message := fmt.Sprintf("Неверный ответ.\nТест-кейс:\n%v\nОжидалось:\n%v\nПолучено:\n%v", testCase, refSol.out, out)
			return message, nil
		}
	}

	return "OK", nil
}

// Модуль проверки для монолитных лабораторных работ
func (c *Checker) monolitTests(compilerCh chan struct {
	msg string
	err error
}) (string, error) {
	// Создаем канал для получения данных о создании эталонного решения
	errSolutionCh := make(chan error)
	go func() {
		err := c.createSolution()
		errSolutionCh <- err
	}()

	// Проверяем данные из каналов на ошибки и другие выводы
	out := <-compilerCh
	if out.err != nil {
		return "", out.err
	}
	if out.msg != "" {
		return out.msg, nil
	}

	err := <-errSolutionCh
	if err != nil {
		return "", err
	}
	// Запускаем тесты
	msg, err := c.runTests()

	if err != nil {
		return "", err
	}

	return msg, nil
}
