package checker

import (
	"fmt"
	"io"
	"os"
)

func (c *Checker) monolitChecker(codeFilePath string) error {
	refSolutionPath := fmt.Sprintf("./%v/ref-solution.py", c.tempDirPath)
	refSolution, err := os.Create(refSolutionPath)
	if err != nil {
		return err
	}

	err = c.createSolution(refSolution)
	if err != nil {
		return err
	}
	refSolution.Close()

	//c.runCode(codeFilePath)
	return nil
}

// func (c *Checker) runCode(path string) {
// 	cmd := exec.Command()
// }

func (c *Checker) createSolution(solution *os.File) error {
	for i := 1; i <= c.tasksCount; i++ {
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
