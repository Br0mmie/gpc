package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func make_folders(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("[ERROR] Failed to create %s: %v", path, err)
	}
}

func create_folders(path string) {
	for _, name := range rootFolders {
		full := filepath.Join(path, name)
		make_folders(full)
	}
	for _, name := range internalFolders {
		full := filepath.Join(path, "internal", name)
		make_folders(full)
	}

	{
		temp_path := filepath.Base(filepath.Clean(path))
		fmt.Println("Cleaned path:", temp_path)
		full := filepath.Join(path, "cmd", temp_path)
		make_folders(full)
	}

	for _, name := range apiFolders {
		full := filepath.Join(path, "api", name)
		make_folders(full)
	}
}

func make_go_project(path string) {
	projectName := filepath.Base(filepath.Clean(path))
	cmd := exec.Command("go", "mod", "init", projectName)

	cmd.Dir = path

	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("[ERROR] Failed to create go project")
		return
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
				make_go_project(wd)
			case "n":
				// Handle no case
				fmt.Print("Enter project name: ")
				if scanner.Scan() {
					folderName := strings.TrimSpace(scanner.Text())
					fmt.Println(folderName)
					path := filepath.Join(wd, folderName)
					fmt.Println("Path: ", path)
					create_folders(path)
					make_go_project(path)
				}
				return
			default:
				// default case for yes
				create_folders(wd)
				make_go_project(wd)
			}
		}
	}
}
