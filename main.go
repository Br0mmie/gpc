package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type projectLayout struct {
	root     []string
	internal []string
	api      []string
}

var defaultLayout = projectLayout{
	root:     []string{"api", "cmd", "internal"},
	internal: []string{"auth", "http", "metrics", "models"},
	api:      []string{"reads", "writes"},
}

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("Project name cannot be empty")
	}
	if !validName.MatchString(name) {
		return fmt.Errorf("Project name %q contains invalid characters, only letters, numbers, hyphens and underscores are allowed.", name)
	}
	return nil
}

func makeFolders(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("[ERROR] Failed to create %s: %v", path, err)
	}
}

func createFolders(path string, layout projectLayout) {
	for _, name := range layout.root {
		full := filepath.Join(path, name)
		makeFolders(full)
	}
	for _, name := range layout.internal {
		full := filepath.Join(path, "internal", name)
		makeFolders(full)
	}

	makeFolders(filepath.Join(path, "cmd", filepath.Base(filepath.Clean(path))))

	for _, name := range layout.api {
		full := filepath.Join(path, "api", name)
		makeFolders(full)
	}
}

func checkProjectExists(path string) error {
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return fmt.Errorf("A go project already exists at %s", path)
	}
	return nil
}

func makeGoProject(path string) {
	projectName := filepath.Base(filepath.Clean(path))

	if err := checkProjectExists(path); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	cmd := exec.Command("go", "mod", "init", projectName)

	cmd.Dir = path

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("[ERROR] Failed to create go project: %v\n%s", err, out)
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
			case "y", "Y", "":
				// default case for yes
				createFolders(wd, defaultLayout)
				makeGoProject(wd)
			case "n":
				// Handle no case
				fmt.Print("Enter project name: ")
				if scanner.Scan() {
					folderName := strings.TrimSpace(scanner.Text())

					if err := validateProjectName(folderName); err != nil {
						log.Fatalf("[ERROR] Invalid project name: %v", err)
					}

					path := filepath.Join(wd, folderName)
					createFolders(path, defaultLayout)
					makeGoProject(path)
				}
				return
			default:
				// default case for yes
				createFolders(wd, defaultLayout)
				makeGoProject(wd)
			}
		}
	} else {
		log.Fatalf("[ERROR] This program does not accept arguments, run it without any.")
	}
}
