package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Структура для конкретного тест-кейса
type testCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// Структура для всех тест-кейсов
type testCases struct {
	TestCases []testCase `json:"test-cases"`
}

// Получаем тест-кейсы из конфига
func (c *Checker) getTestCases() (*testCases, error) {
	var tc *testCases

	// Находим и открываем конкретный файл с тест-кейсами
	taskVar := c.lab.Tasks[fmt.Sprintf("task%v", c.lab.CurrTask)]
	path := fmt.Sprintf("./labs/lab%v/task%v/var%v/tests.json", c.lab.LabNum, c.lab.CurrTask, taskVar)
	config, err := os.Open(path)
	if err != nil {
		e := fmt.Errorf("error with getting test-cases: %v", err)
		log.Println(e)
		return nil, e
	}

	// Декодируем данные из файла в структуру
	if err := json.NewDecoder(config).Decode(&tc); err != nil {
		e := fmt.Errorf("error with decoding test-cases: %v", err)
		log.Println(e)
		return nil, e
	}

	return tc, nil
}

// Запускаем тесты для проверки кода ученика
func (c *Checker) runTestCases(tc *testCases) (string, error) {
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

	// Получаем абсолютный путь к файлу с кодом
	var execPath string
	switch c.OS {
	case "Windows":
		execPath = fmt.Sprintf("./%v/code.exe", c.tempDirPath)
	case "Linux":
		execPath = fmt.Sprintf("./%v/code", c.tempDirPath)
	case "darwin":
		execPath = fmt.Sprintf("./%v/code", c.tempDirPath)
	default:
		log.Printf("This system is not suitable, please use a compatible OS.")
	}
	absCodePath, _ := filepath.Abs(execPath)

	// Проходимся по тест-кейсам и запускаем программы
	for i := range tc.TestCases {
		testCase := tc.TestCases[i]

		// Записываем тест-кейсы по очереди в input файл
		if err := os.WriteFile(inputPath, []byte(testCase.Input), 0644); err != nil {
			e := fmt.Errorf("error with wtriting test-case to the input file: %v", err)
			log.Println(e)
			return "", err
		}

		// Сбрасываем позицию input в начало и очищаем файл
		inputFile.Seek(0, 0)
		outputFile.Truncate(0) // очищаем output файл (чтобы следующий вывод программы не нарушил проверку)
		outputFile.Seek(0, 0)

		// Запускаем код ученика
		err := c.runCode(absCodePath, inputFile, outputFile)
		if err != nil {
			return "", err
		}

		// Сохраняем данные из буфера в файл и смещаем позицию чтения на 0,
		// чтобы получить вывод программы ученика
		outputFile.Sync()
		outputFile.Seek(0, 0)

		codeOut, _ := io.ReadAll(outputFile)

		// Сравниваем вывод программ
		if isCorrect := c.validate(string(codeOut), testCase.Output); !isCorrect {
			out := c.redact(string(codeOut))
			message := fmt.Sprintf("Неверный ответ.\nТест-кейс:\n%v\nОжидалось:\n%v\nПолучено:\n%v", testCase.Input, testCase.Output, out)
			return message, nil
		}
	}

	return "OK", nil
}

func (c *Checker) splitedTests(compilerCh chan struct {
	msg string
	err error
}) (string, error) {
	// Проверяем, успешно ли завершилась компиляция
	out := <-compilerCh
	if out.err != nil {
		return "", out.err
	}
	if out.msg != "" {
		return out.msg, nil
	}

	tc, err := c.getTestCases()
	if err != nil {
		return "", nil
	}

	msg, err := c.runTestCases(tc)
	if err != nil {
		return "", err
	}

	return msg, nil
}
