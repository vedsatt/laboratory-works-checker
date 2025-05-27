package checker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	labCongif   *labConfig
	lab         *models.LabRequest
	tempDirPath string
	result      *models.CheckerResponse
}

type outputs struct {
	expected string
	got      string
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

	checker.labCongif = labCfg

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

func (c *Checker) createAndCompile() (string, error) {
	codeFilePath := fmt.Sprintf("./%v/code.%v", c.tempDirPath, c.labCongif.Language)
	codeFile, err := os.Create(codeFilePath)
	if err != nil {
		return "", err
	}

	code, err := strconv.Unquote(`"` + c.lab.Code + `"`)
	if err != nil {
		return "", err
	}

	_, err = codeFile.WriteString(code)
	if err != nil {
		return "", nil
	}

	codeFile.Close()

	var cmd *exec.Cmd

	switch c.labCongif.Language {
	case "c":
		cmd = exec.Command("gcc", codeFilePath, "-o", fmt.Sprintf("./%v/code", c.tempDirPath))
	case "cpp", "cxx", "cc":
		cmd = exec.Command("g++", codeFilePath, "-o", fmt.Sprintf("./%v/code", c.tempDirPath))
	default:
		err := fmt.Errorf("unsupported file extension: %s", c.labCongif.Language)
		log.Println(err)
		return "", err
	}

	// запуск компиляции
	if err := cmd.Run(); err != nil {
		log.Println(fmt.Errorf("compilation failed: %w", err))
		return err.Error(), nil
	}
	os.Remove(codeFilePath)

	return "", nil
}

func (c *Checker) Check() (string, error) {
	dirPath := fmt.Sprintf("tmp-%v", c.lab.ID)
	err := os.Mkdir(dirPath, 0770)
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(dirPath)

	c.tempDirPath = dirPath

	msg, err := c.createAndCompile()
	if err != nil {
		return "", err
	}
	if msg != "" {
		return msg, nil
	}

	switch c.labCongif.LabType {
	case "splited":
		err = c.splitedTests()
	case "monolit":
		msg, err = c.monolitTests()
	}

	if err != nil {
		return "", err
	}
	return msg, nil
}
