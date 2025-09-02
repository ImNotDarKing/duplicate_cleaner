package application

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type FileEntry struct {
	Path string
	ModTime time.Time
}


func (hs *HashStore) Cleaner(printAllDuplicate bool, sortAscending bool, removeDuplicates bool, workerCount int) {

	if sortAscending{
		fmt.Println("Включена сортировка по самому раннему файлу")
	} else {
		fmt.Println("Включена сортировка по самому позднему файлу")
	}
	fmt.Println(strings.Repeat("-", 100))	

	duplicatesFound := false
	
	for hash, files := range hs.Files {
		if len(files) > 1 {
			duplicatesFound = true
			fmt.Println("Дубликаты для хэша:", hash)
			fmt.Println("Найдено файлов:", len(files))

			if printAllDuplicate {
				for path, info := range files {
					fmt.Printf(" - %s (Последняя дата изменения: %s)\n", path, info.ModTime.Format(time.RFC3339))
				}
			}
			if removeDuplicates {
				var wg sync.WaitGroup
				deleteJobs := make(chan string)
				for i := 0; i < workerCount; i ++ {
					wg.Add(1)
					go func(workerID int) {
						defer wg.Done()
						for file := range deleteJobs {
							err := os.Remove(file)
							fmt.Println("Удалён:", file)
							if err != nil {
								fmt.Printf("Ошибка в горутине %d при удалении %s: %v\n", workerID, file, err)
							}
						}
					} (i+1)
				}

				var fileEntries []FileEntry
				for path, info  := range files {
					fileEntries = append(fileEntries, FileEntry {
						Path: path,
						ModTime: info.ModTime,
					})
				}

				sort.Slice(fileEntries, func(i, j int) bool {
					if sortAscending {
						return fileEntries[i].ModTime.Before(fileEntries[j].ModTime)
					} else {
						return fileEntries[i].ModTime.After(fileEntries[j].ModTime)
					}
				})

				keep := fileEntries[0].Path
				for _, entry := range fileEntries {
					if entry.Path == keep {
						continue
					}
					deleteJobs <- entry.Path
				}
				close(deleteJobs)
				wg.Wait()
			}
			fmt.Println(strings.Repeat("-", 100))			
		} 
	}

	if !duplicatesFound {
		fmt.Println("Дубликаты не найдены.")
	}
}

