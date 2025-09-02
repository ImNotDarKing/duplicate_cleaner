package application

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileInfo struct {
	ModTime time.Time
}

type HashStore struct {
	mu    sync.Mutex
	Files map[string]map[string]FileInfo
}

func walkAll(path string) ([]string, error) {
	var filenames []string

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			filenames = append(filenames, path)
		}

		return nil
	})
	return filenames, err
}

func (hs *HashStore) HashProcessor(path string, workerCount int, printAllFiles bool) {
	startTime := time.Now()

	files, err := walkAll(path)
	if err != nil {
		fmt.Printf("Ошибка при обработке каталога: %v\n", err)
		return
	}

	jobs := make(chan string)
	var wg sync.WaitGroup
	// fmt.Println(strings.Repeat("-", 100))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for name := range jobs {
				file, err := os.Open(name)
				if err != nil {
					fmt.Printf("[worker %d] Не удалось открыть файл: %s, ошибка: %v\n", workerID, name, err)
					continue
				}

				hash := sha256.New()
				if _, err := io.Copy(hash, file); err != nil {
					fmt.Printf("[worker %d] Ошибка при чтении файла: %s, ошибка: %v\n", workerID, name, err)
					file.Close()
					continue
				}
				file.Close()

				sum := fmt.Sprintf("%x", hash.Sum(nil))

				fileInfo, err := os.Stat(name)
				if err != nil {
					fmt.Printf("[worder %d] Ошибка при получении информации о файле: %s, ошибка: %v\n", workerID, name, err)
					continue
				}
				modTime := fileInfo.ModTime()

				hs.mu.Lock()
				if hs.Files[sum] == nil {
					hs.Files[sum] = make(map[string]FileInfo)
				}
				hs.Files[sum][name] = FileInfo{ModTime: modTime}
				hs.mu.Unlock()

				if printAllFiles {
					fmt.Printf("Файл: %s\nSHA256: %s\n", name, sum)
					fmt.Println(strings.Repeat("-", 100))
				}
			}
		}(i + 1)
	}

	var fileSum int
	for _, name := range files {
		jobs <- name
		fileSum += 1
	}
	close(jobs)

	wg.Wait()
	elapsed := time.Since(startTime)
	fmt.Println("Операция анализа файлов завершена за", elapsed.Seconds())
	fmt.Println("Всего найдено файлов:", fileSum)
}
