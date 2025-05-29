package checker

import (
	"slices"
	"strings"
)

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
	progLines := strings.FieldsFunc(programOutput, func(r rune) bool { return r == '\n' })
	correctOutLines := strings.Split(correctAns, "\n")

	i, j := 0, 0
	for i < len(progLines) {
		// удаляет лишние пробелы и прочие знаки кодировки
		currLine := strings.TrimSpace(progLines[i])
		correctLine := strings.TrimSpace(correctOutLines[j])
		// Если в строке есть префикс #> - это шум и он ни на что не влияет
		// Также пропускаются любые пустые строки
		if isNoise(currLine) {
			i++
			continue
		}

		// Сравниваем вывод программы с выводом эталонного решения
		if currLine != correctLine && !isReport(currLine, correctLine) {
			return false
		}
		i++
		j++
		if j >= len(correctOutLines) {
			return false
		}
	}

	return true
}

// Форматирует некорректный вывод ученика таким образом, чтобы убрать шум и пустые строки
// Это нужно, чтобы ученики не подумали, что решение не прошло из-за шума
func (c *Checker) redact(code string) string {
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
