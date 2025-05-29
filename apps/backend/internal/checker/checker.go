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

// Создает объект структуры Checker, через которую будет проводиться дальнейшая проверка.
// P.S. Для тех, кто не особо знаком с Go - у структур есть свои методы, это аналогично классам в других языках
func New(lab *models.LabRequest) (*Checker, error) {
	checker := &Checker{
		lab:    lab,
		result: &models.CheckerResponse{},
	}

	// Получаем данные их конфига лабы
	labCfg, err := checker.config()
	if err != nil {
		e := fmt.Errorf("failed to load config: %w", err)
		log.Println(e)
		return nil, e
	}

	if labCfg == nil {
		e := fmt.Errorf("config is nil")
		log.Println(e)
		return nil, e
	}

	// Сохраняем конфиг в Checker
	checker.labCongif = labCfg

	return checker, nil
}

// Берет данные из config.json в конкретной лабораторной работе.
// У всех лабораторных обязательно должен быть конфиг с описанием.
// Иначе программа не поймет, что ей делать с конкретной лабой
func (c *Checker) config() (*labConfig, error) {
	// Формируем путь и получаем доступ к конфигу
	path := fmt.Sprintf("./labs/lab%v/config.json", c.lab.LabNum)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Декодируем данные из файла в структуру config
	var config *labConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		e := fmt.Errorf("error with decoding lab config: %v", err)
		log.Println(e)
		return nil, e
	}

	return config, nil
}

// Функция создает файл code.с/cpp... из полученного кода
// и сразу компилирует его, удаляя после этого файл с кодом
func (c *Checker) createAndCompile() (string, error) {
	//Создаем пустой файл
	codeFilePath := fmt.Sprintf("./%v/code.%v", c.tempDirPath, c.labCongif.Language)
	codeFile, err := os.Create(codeFilePath)
	if err != nil {
		e := fmt.Errorf("error with creating code file: %v", err)
		log.Println(e)
		return "", e
	}

	// Парсим строку из json в нормальный код
	code, err := strconv.Unquote(`"` + c.lab.Code + `"`)
	if err != nil {
		e := fmt.Errorf("error with unquoting code: %v", err)
		log.Println(e)
		return "", e
	}

	// Записываем код в файл
	_, err = codeFile.WriteString(code)
	if err != nil {
		e := fmt.Errorf("error writing code to file: %v", err)
		log.Println(e)
		return "", e
	}

	codeFile.Close()

	// Создаем и настраиваем команду компиляции
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

	// Запускаем компиляцию и удаляем временный файл с кодом
	if err := cmd.Run(); err != nil {
		log.Println(fmt.Errorf("compilation failed: %w", err))
		return err.Error(), nil
	}
	os.Remove(codeFilePath)

	return "", nil
}

// Модуль проверки, который, собственно, и занимается проверкой кода учеников
func (c *Checker) Check() (string, error) {
	// Создаем временную папку, в которой будут все файлы, связанные с проверкой
	// После проверки папка удаляется (если программа завершится с ошибкой и папка не удалится, то
	// при следующем запуске Cleaner почистит все, что осталось)
	dirPath := fmt.Sprintf("tmp-%v", c.lab.ID)
	err := os.Mkdir(dirPath, 0770)
	if err != nil {
		e := fmt.Errorf("error with creating tmp directory: %v", err)
		log.Println(e)
		return "", e
	}

	defer os.RemoveAll(dirPath)

	c.tempDirPath = dirPath

	// Создаем канал, чтобы параллельно скомпилировать файл, и в случае монолитного типа лаб собрать эталонное решение
	// По результатам проверки на существующих работах производительность не сильно улучшает, но если работы будут больше,
	// то разница же будет чувствоваться (на маленьких при многопоточном создании ~ -20мс)
	compilerCh := make(chan struct {
		msg string
		err error
	})

	// Компилируем файл
	go func() {
		msg, err := c.createAndCompile()
		compilerCh <- struct {
			msg string
			err error
		}{msg, err}
	}()

	var testsMsg string
	// У лаб есть 2 типа - все части являются подпрограммами одного большого кода (monolit), и все части являются отдельной программой (splited).
	// В зависимости от этого используются разные проверки, т. к. в случае монолитной лабы это упрощает написания тестов, а в случае разбитой лабы
	// ускоряет работу программы, позволяя не создавать для каждого задания эталонное решение, а просто составить тест-кейсы
	switch c.labCongif.LabType {
	case "splited":
		testsMsg, err = c.splitedTests(compilerCh)
	case "monolit":
		testsMsg, err = c.monolitTests(compilerCh)
	}

	if err != nil {
		return "", err
	}

	// Логируем результат проверки, чтобы было проще отслеживать работу программы
	if testsMsg == "OK" {
		log.Printf("request with id %v: all tests passed successfully", c.lab.ID)
	}
	return testsMsg, nil
}
