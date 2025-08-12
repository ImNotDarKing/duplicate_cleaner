package main

import (
	"github.com/ImNotDarKing/duplicate_cleaner/internal/application"
)

func main() {
	GOROUTINES := 5
	PRINT_ALL_FILES := false
	PRINT_ALL_DUPLICATES := true
	REMOVE_DUPLICATES := false
	PATH := "C:/Users/DarKing/Desktop/projects_Go/Программирование на Go  24"

	hs := application.HashStore{Files: make(map[string][]string)}
	hs.HashProcessor(PATH, GOROUTINES, PRINT_ALL_FILES)
	hs.Cleaner(PRINT_ALL_DUPLICATES, REMOVE_DUPLICATES, GOROUTINES)
}