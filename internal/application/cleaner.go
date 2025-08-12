package application

import(
	"fmt"
	"strings"
	"sync"
	"os"
)


func (hs *HashStore) Cleaner(printAllDuplicate bool, removeDuplicates bool, workerCount int) {
	fmt.Println(strings.Repeat("-", 100))	

	duplicatesFound := false
	deletedSomething := false
	for hash, files := range hs.Files {
		if len(files) > 1 {
			duplicatesFound = true
			fmt.Println("Дубликаты для хэша:", hash)
			fmt.Println("Найдено файлов:", len(files))

			if printAllDuplicate {
				for _, file := range files {
					fmt.Println(" -", file)
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
							if err != nil {
								fmt.Printf("Ошибка в горутине %d при удалении %s: %v\n", workerID, file, err)
							}
						}
					} (i+1)
				}
				for _, file := range files[1:] {
					deleteJobs <- file
					deletedSomething = true
				}
				close(deleteJobs)
				wg.Wait()
			}
			fmt.Println(strings.Repeat("-", 100))			
		} 
	}
	if removeDuplicates && deletedSomething {
		fmt.Println("Все дубликаты были удалены.")
	}
	if !duplicatesFound {
		fmt.Println("Дубликаты не найдены.")
	}
}

