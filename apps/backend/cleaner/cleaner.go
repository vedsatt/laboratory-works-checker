package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
)

// Cleaner нужен, чтобы избежать ошибок в случае некорректного завершения программы.
// Если программа завершилась с ошибкой и остались временные директории, которые могут
// помешать работе программы - cleaner удалит все эти папки.
func Clean() error {
	dirPath := "../../../"
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Проверяем, что имя начинается с "tmp-"
		if info.IsDir() && filepath.Base(path) != filepath.Base(dirPath) {
			if matched, _ := filepath.Match("tmp-*", info.Name()); matched {
				fmt.Printf("deleting directory: %s\n", path)

				// Удаляем папку и всё её содержимое
				err := os.RemoveAll(path)
				if err != nil {
					return fmt.Errorf("error with removing %s: %v", path, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
