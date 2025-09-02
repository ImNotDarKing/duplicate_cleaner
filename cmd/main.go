package main

import (
	"github.com/ImNotDarKing/duplicate_cleaner/internal/application"
)

func main() {
	GOROUTINES := 5
	PRINT_ALL_FILES := false
	PRINT_ALL_DUPLICATES := true
	SORT_ASCENDING := true   // f-поздний, t-ранний
	REMOVE_DUPLICATES := false
	PATH := "C:/Users/DarKing/Desktop/projects_Go/duplicate_cleaner/test"

	hs := application.HashStore{Files: make(map[string]map[string]application.FileInfo)}
	hs.HashProcessor(PATH, GOROUTINES, PRINT_ALL_FILES)
	hs.Cleaner(PRINT_ALL_DUPLICATES, SORT_ASCENDING, REMOVE_DUPLICATES, GOROUTINES)
}