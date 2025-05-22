package checker

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/vedsatt/laboratory-works-checker/apps/backend/internal/models"
)

type labConfig struct {
	Language   string   `json:"language"`
	LabType    string   `json:"type"`
	TasksCount int      `json:"tasks-count"`
	TestCases  []string `json:"test-cases"`
}

type Checker struct {
	tempDirPath string
	labType     string
	lab         *models.LabRequest
	Language    string
	tasksCount  int
	TestCases   []string
	result      *models.CheckerResponse
}

func New(lab *models.LabRequest) (*Checker, error) {
	checker := &Checker{
		lab:    lab,
		result: &models.CheckerResponse{},
	}

	labCfg, err := checker.config()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if labCfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	checker.labType = labCfg.LabType
	checker.tasksCount = labCfg.TasksCount
	checker.Language = labCfg.Language

	return checker, nil
}

func (c *Checker) config() (*labConfig, error) {
	path := fmt.Sprintf("./labs/lab%v/config.json", c.lab.LabNum)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config *labConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Checker) Check() error {
	dirPath := fmt.Sprintf("tmp-%v", c.lab.ID)
	err := os.Mkdir(dirPath, 0770)
	if err != nil {
		return err
	}

	//defer os.RemoveAll(dirPath)
	c.tempDirPath = dirPath

	codeFilePath := fmt.Sprintf("./%v/code.%v", c.tempDirPath, c.Language)
	codeFile, err := os.Create(codeFilePath)
	if err != nil {
		return err
	}

	code, err := strconv.Unquote(`"` + c.lab.Code + `"`)
	if err != nil {
		return err
	}

	_, err = codeFile.WriteString(code)
	if err != nil {
		return nil
	}

	codeFile.Close()

	switch c.labType {
	case "splited":
		err = c.splitedCheck(codeFilePath)
	case "monolit":
		err = c.monolitChecker(codeFilePath)
	}

	if err != nil {
		return err
	}
	return nil
}
