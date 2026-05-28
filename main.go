package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var rootFolders = []string{
	"api",
	"cmd",
	"internal",
}

var internalFolders = []string{
	"auth",
	"http",
	"metrics",
	"models",
}

var apiFolders = []string{
	"reads",
	"writes",
}

func create_folders(path string) {
	fmt.Println(path)
	for _, name := range rootFolders {
		full := filepath.Join(path, name)
		if err := os.MkdirAll(full, os.ModePerm); err != nil {
			log.Fatalf("[ERROR] Failed to create %s: %v", full, err)
		}
	}
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("No path specified. Would you like to create your project in the current folder? [Y/n]")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatalf("[ERROR] Failed to get current working directory: %v", err)
			}

			choice := strings.TrimSpace(scanner.Text())
			switch choice {
			case "y", "":
				// default case for yes
				create_folders(wd)
			case "n":
				// Handle no case
				fmt.Print("Enter project name: ")
				if scanner.Scan() {
					folderName := strings.TrimSpace(scanner.Text())
					fmt.Println(folderName)
					path := filepath.Join(wd, folderName)
					fmt.Print("Path: ", path)
					create_folders(path)
				}
				return
			default:
				// default case for yes
				create_folders(wd)
			}
		}
	}
}
