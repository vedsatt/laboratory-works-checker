package checker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (c *Checker) createSolution() error {
	refSolutionPath := fmt.Sprintf("./%v/ref-solution.py", c.tempDirPath)
	refSolution, err := os.Create(refSolutionPath)
	if err != nil {
		e := fmt.Errorf("error with creating solution file: %v", err)
		log.Println(e)
		return e
	}
	defer refSolution.Close()

	for i := 1; i <= c.labCongif.TasksCount; i++ {
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

		_, err = refSolution.WriteString(string(solutionPart))
		if err != nil {
			e := fmt.Errorf("error with creating solution file: %v", err)
			log.Println(e)
			return e
		}
	}

	return nil
}

func (c *Checker) runRefSolution(absPath, input string) (string, error) {
	cmd := exec.Command("python", absPath)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	err := cmd.Run()
	if err != nil {
		e := fmt.Errorf("error with running ref solution: %v\nStderr: %s", err, stderr.String())
		log.Println(e)
		return "", e
	}
	return stdout.String(), nil
}

func (c *Checker) runCode(absPath string, inputFile, outputFile *os.File) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, absPath)
	var stderr bytes.Buffer
	cmd.Stdin = inputFile
	cmd.Stdout = outputFile
	cmd.Stderr = &stderr

	err := cmd.Run()
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "Превышено максимальное время выполнения.", nil
		}
		return stderr.String(), err
	}

	return "", nil
}

func (c *Checker) runTests() (string, error) {
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

	// Для питона весь ввод и вывод - строка из конфига лаб. работы,
	// т.к. это удобнее и ускоряет работу программы
	execPath := fmt.Sprintf("./%v/code.exe", c.tempDirPath)
	absCodePath, _ := filepath.Abs(execPath)
	refPath := fmt.Sprintf("./%v/ref-solution.py", c.tempDirPath)
	absRefPath, _ := filepath.Abs(refPath)

	for i := range c.labCongif.TestCases {
		testCase := c.labCongif.TestCases[i]

		refSolCh := make(chan struct {
			out string
			err error
		})

		go func() {
			refOut, err := c.runRefSolution(absRefPath, testCase) // горутина 1
			refSolCh <- struct {
				out string
				err error
			}{refOut, err}
		}()

		codeCh := make(chan struct {
			stdErr string
			err    error
			out    string
		})

		go func() {
			// Записываем тест-кейсы по очереди в input файл
			if err := os.WriteFile(inputPath, []byte(testCase), 0644); err != nil {
				e := fmt.Errorf("error with wtriting test-case to the input file: %v", err)
				log.Println(e)
				codeCh <- struct {
					stdErr string
					err    error
					out    string
				}{"", e, ""}
			}

			// Сбрасываем позицию input в начало и очищаем файл
			inputFile.Seek(0, 0)
			outputFile.Truncate(0) // очищаем output файл
			outputFile.Seek(0, 0)

			stdErr, err := c.runCode(absCodePath, inputFile, outputFile) // горутина 2
			if err != nil {
				codeCh <- struct {
					stdErr string
					err    error
					out    string
				}{"", err, ""}
			}

			if stdErr != "" {
				codeCh <- struct {
					stdErr string
					err    error
					out    string
				}{stdErr, nil, ""}
			}

			// Сохраняем данные из буфера в файл и смещаем позицию чтения на 0,
			// чтобы получить вывод программы ученика
			outputFile.Sync()
			outputFile.Seek(0, 0)

			codeOut, _ := io.ReadAll(outputFile)
			codeCh <- struct {
				stdErr string
				err    error
				out    string
			}{"", nil, string(codeOut)}
		}()

		refSol := <-refSolCh
		if refSol.err != nil {
			return "", err
		}

		code := <-codeCh
		if code.err != nil {
			return "", err
		}
		if code.stdErr != "" {
			return code.stdErr, nil
		}

		// Ждем выполнения эталонного решения и кода ученика
		if isCorrect := c.validate(code.out, refSol.out); !isCorrect {
			message := fmt.Sprintf("Неверный ответ.\nТест-кейс:\n%v\nОжидалось:\n%v\nПолучено:\n%v", testCase, refSol.out, code.out)
			return message, nil
		}
	}

	return "OK", nil
}

func (c *Checker) monolitTests(compilerCh chan struct {
	msg string
	err error
}) (string, error) {
	errSolutionCh := make(chan error)
	go func() {
		/////////////
		err := c.createSolution()
		errSolutionCh <- err
	}()

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
	msg, err := c.runTests()

	if err != nil {
		return "", err
	}

	return msg, nil
}
