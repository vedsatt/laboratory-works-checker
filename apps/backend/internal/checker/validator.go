package checker

import (
	"strings"
)

func isReport(program, correct string) bool {
	return (strings.HasPrefix(program, "Error:") && strings.HasPrefix(correct, "Error:")) ||
		(strings.HasPrefix(program, "Notice:") && strings.HasPrefix(correct, "Notice:"))
}

func (c *Checker) validate(programOutput, correctAns string) bool {
	progLines := strings.FieldsFunc(programOutput, func(r rune) bool { return r == '\n' })
	correctOutLines := strings.Split(correctAns, "\n")

	i, j := 0, 0
	for i < len(progLines) {
		currLine := strings.TrimSpace(progLines[i])
		if strings.HasPrefix(currLine, "#>") {
			i++
			continue
		}

		if currLine != correctOutLines[j] && !isReport(currLine, correctOutLines[j]) {
			return false
		}
		j++
	}

	return true
}
