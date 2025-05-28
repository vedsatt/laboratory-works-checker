package checker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (c *Checker) createSolution(solution *os.File) error {
	for i := 1; i <= c.labCongif.TasksCount; i++ {
		path := fmt.Sprintf("./labs/lab%v/task%v/var%v/solution.py",
			c.lab.LabNum, i, c.lab.Tasks[fmt.Sprintf("task%v", i)])
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		solutionPart, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		_, err = solution.WriteString(string(solutionPart))
		if err != nil {
			return err
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

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error with running ref solution: %v\nStderr: %s", err, stderr.String())
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
		return "", err
	}
	defer inputFile.Close()

	outputPath := fmt.Sprintf("./%v/output.txt", c.tempDirPath)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
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
		refOut, err := c.runRefSolution(absRefPath, testCase)
		if err != nil {
			return "", err
		}
		fmt.Printf("ref:\n%v", refOut)

		if err := os.WriteFile(inputPath, []byte(testCase), 0644); err != nil {
			return "", err
		}

		// Сбрасываем позицию в начало и очищаем файл
		if _, err := inputFile.Seek(0, 0); err != nil {
			return "", err
		}
		if err := outputFile.Truncate(0); err != nil { // Очищаем файл
			return "", err
		}
		if _, err := outputFile.Seek(0, 0); err != nil {
			return "", err
		}

		stdErr, err := c.runCode(absCodePath, inputFile, outputFile)
		if err != nil {
			return "", err
		}
		if stdErr != "" {
			return stdErr, nil
		}

		if err := outputFile.Sync(); err != nil {
			return "", err
		}
		if _, err := outputFile.Seek(0, 0); err != nil {
			return "", err
		}
		codeOut, _ := io.ReadAll(outputFile)

		if isCorrect := c.validate(string(codeOut), refOut); !isCorrect {
			message := fmt.Sprintf("Неверный ответ.\nТест-кейс:\n%v\nОжидалось:\n%v\nПолучено:\n%v", testCase, refOut, string(codeOut))
			return message, nil
		}
	}

	return "OK", err
}

func (c *Checker) monolitTests() (string, error) {
	refSolutionPath := fmt.Sprintf("./%v/ref-solution.py", c.tempDirPath)
	refSolution, err := os.Create(refSolutionPath)
	if err != nil {
		return "", err
	}

	err = c.createSolution(refSolution)
	if err != nil {
		return "", err
	}
	refSolution.Close()

	msg, err := c.runTests()
	if err != nil {
		return "", err
	}

	return msg, nil
}
