package checker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"
)

// Обрабатывает ошибки процесса (таймаут, stderr и т.д.)
func handleProcessError(err error, ctx context.Context, stderr *bytes.Buffer) error {
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("Превышено максимальное время выполнения")
	}
	if err != nil {
		if stderr.String() != "" {
			return fmt.Errorf("%s", stderr.String())
		}
		return fmt.Errorf("Ошибка выполнения программы: %v", err)
	}
	return nil
}

// Запускает код ученика с конкретным тест-кейсом
func (c *Checker) runCode(absPath string, inputFile, outputFile *os.File) error {
	// Контекст с таймаутом для контроля времени выполнения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем команду
	cmd := exec.CommandContext(ctx, absPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Перенаправляем stdout в файл
	cmd.Stdout = outputFile

	// Получаем Pipe для stdin, чтобы контролировать запись
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Не удалось создать stdin pipe: %v", err)
	}

	// Запускаем процесс
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Не удалось запустить программу: %v", err)
	}

	// Канал для сигнала об окончании записи входных данных
	writeDone := make(chan error, 1)
	// Канал для сигнала о завершении процесса
	processDone := make(chan error, 1)

	// Горутина для записи входных данных
	go func() {
		defer close(writeDone)
		_, err := io.Copy(stdinPipe, inputFile)
		if err != nil {
			writeDone <- fmt.Errorf("Ошибка записи в stdin: %v", err)
			return
		}
		// Закрываем stdin после записи
		if err := stdinPipe.Close(); err != nil {
			writeDone <- fmt.Errorf("Ошибка закрытия stdin: %v", err)
			return
		}
		writeDone <- nil
	}()

	// Горутина для ожидания завершения процесса
	go func() {
		processDone <- cmd.Wait()
	}()

	// Ожидаем либо завершения записи, либо завершения процесса
	select {
	case writeErr := <-writeDone:
		if writeErr != nil {
			return fmt.Errorf("Программа завершилась некорректно: %v", writeErr)
		}
		// Дожидаемся завершения процесса в случае успешной записи
		select {
		case procErr := <-processDone:
			return handleProcessError(procErr, ctx, &stderr)
		case <-ctx.Done():
			return fmt.Errorf("Превышено максимальное время выполнения")
		}
	case procErr := <-processDone:
		// Программа завершилась до окончания записи ввода
		return fmt.Errorf("Программа завершилась некорректно (раннее завершение): %v", procErr)
	case <-ctx.Done():
		return fmt.Errorf("Превышено максимальное время выполнения")
	}
}

// Проверяет, является ли вывод уведомлением об ошибке или граничном случае
// Сделано для того, чтобы не загонять учеников с лишком узкие рамки, чтобы они могли
// Написать свое описание ошибки или уведомления. Главное, чтобы был нужный префикс
func isReport(program, correct string) bool {
	program = strings.TrimSpace(program)
	correct = strings.TrimSpace(correct)
	return (strings.HasPrefix(program, "Error:") && strings.HasPrefix(correct, "Error:")) ||
		(strings.HasPrefix(program, "Notice:") && strings.HasPrefix(correct, "Notice:"))
}

// Проверяет, является ли строка шумом
func isNoise(str string) bool {
	return strings.HasPrefix(str, "#>") || str == "" || str == "\t" || str == " "
}

// Сравнивает вывод эталонного решения и вывод программы ученика
func (c *Checker) validate(programOutput, correctAns string) bool {
	// Разбивает вывод на массив строк
	progLines := strings.Split(programOutput, "\n")
	correctOutLines := strings.Split(c.redact(correctAns), "\n")
	if len(correctOutLines) != len(progLines) {
		return false
	}

	i, j := 0, 0
	for i < len(progLines) {
		// удаляет лишние пробелы и прочие знаки кодировки
		currLine := strings.TrimSpace(progLines[i])
		correctLine := strings.TrimSpace(correctOutLines[j])
		// Если в строке есть префикс #> - это шум и он ни на что не влияет
		// Также пропускаются любые пустые строки

		// Сравниваем вывод программы с выводом эталонного решения
		if currLine != correctLine && !isReport(currLine, correctLine) {
			return false
		}

		// Проверяем, что счетчик не вышел за границы
		if j >= len(correctOutLines) {
			return false
		}
		i++
		j++
	}

	// Проверим, что в выводе эталонного решения не осталось строк
	for j < len(correctOutLines) {
		if strings.TrimSpace(correctOutLines[j]) != "" {
			return false
		}
		j++
	}

	return true
}

// Форматирует некорректный вывод ученика таким образом, чтобы убрать шум и пустые строки
// Это нужно, чтобы ученики не подумали, что решение не прошло из-за шума
func (c *Checker) redact(code string) string {
	code = strings.ReplaceAll(code, "\\n", "\n")
	code = strings.ReplaceAll(code, "\r\n", "\n")
	codeLines := strings.Split(code, "\n")

	i := 0
	for i < len(codeLines) {
		if isNoise(codeLines[i]) {
			codeLines = slices.Delete(codeLines, i, i+1)
		} else {
			i++
		}
	}

	return strings.Join(codeLines, "\n")
}
