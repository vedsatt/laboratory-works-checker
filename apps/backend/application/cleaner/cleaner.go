// Package cleaner предоставляет утилиту для очистки временных директорий в случае завершения программы с ошибкой.
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
	dirPath := "./"

	// Читаем содержимое корневой директории
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	// Проходим по всем элементам в корневой директории
	for _, file := range files {
		// Проверяем, что это директория и имя начинается с "tmp-"
		if file.IsDir() {
			if matched, _ := filepath.Match("tmp-*", file.Name()); matched {
				path := filepath.Join(dirPath, file.Name())
				fmt.Printf("deleting directory: %s\n", path)

				// Удаляем папку и всё её содержимое
				err := os.RemoveAll(path)
				if err != nil {
					return fmt.Errorf("error with removing %s: %v", path, err)
				}
			}
		}
	}

	return nil
}
