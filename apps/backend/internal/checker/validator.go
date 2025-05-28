package checker

import (
	"strings"
)

func isReport(program, correct string) bool {
	program = strings.TrimSpace(program)
	correct = strings.TrimSpace(correct)
	return (strings.HasPrefix(program, "Error:") && strings.HasPrefix(correct, "Error:")) ||
		(strings.HasPrefix(program, "Notice:") && strings.HasPrefix(correct, "Notice:"))
}

func (c *Checker) validate(programOutput, correctAns string) bool {
	progLines := strings.FieldsFunc(programOutput, func(r rune) bool { return r == '\n' })
	correctOutLines := strings.Split(correctAns, "\n")

	i, j := 0, 0
	for i < len(progLines) {
		currLine := strings.TrimSpace(progLines[i])
		correctLine := strings.TrimSpace(correctOutLines[j])
		if strings.HasPrefix(currLine, "#>") {
			i++
			continue
		}

		if currLine == "" || currLine == "\t" || currLine == " " {
			i++
			continue
		}

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
