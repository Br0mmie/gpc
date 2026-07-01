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

// projectLayout defines the folder structure for a new Go project.
type projectLayout struct {
	root     []string
	internal []string
	api      []string
}

// defaultLayout is the standard project structure used when no custom layout is provided.
var defaultLayout = projectLayout{
	root:     []string{"api", "cmd", "internal", "logs"},
	internal: []string{"auth", "http", "metrics", "models"},
	api:      []string{"reads", "writes"},
}

// validName restricts project names to alphanumeric characters, hypens and underscores.
var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// validateProjectName ensure the given name is safe to use as a directory and module name
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("Project name cannot be empty")
	}
	if !validName.MatchString(name) {
		return fmt.Errorf("Project name %q contains invalid characters, only letters, numbers, hyphens and underscores are allowed.", name)
	}
	return nil
}

// makeFolders creates a directory and all necessary parts.
func makeFolders(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("[ERROR] Failed to create %s: %v", path, err)
	}
}

// createFolders builds the full project directory structure at the given path.
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

// checkProjectExists returns an error if a go.mod file already exists at the given path,
// indicating the directory has already been initialised as a Go project.
func checkProjectExists(path string) error {
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return fmt.Errorf("A go project already exists at %s", path)
	}
	return nil
}

// makeGoProject initialises a Go module at the given path using the directory name as the
// module name
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

func initializeAir(path string) {
	cmd := exec.Command("go", "install", "github.com/air-verse/air@latest")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("[ERROR] Failed to install air: %v\n%s", err, out)
	}

	gopath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		log.Fatalf("[ERROR] Failed to get GOPATH: %v", err)
	}
	gopathBin := filepath.Join(strings.TrimSpace(string(gopath)), "bin")

	cmd = exec.Command("air", "init")
	cmd.Dir = path
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH")+":"+gopathBin)

	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize air: %v\n%s", err, out)
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
				createFolders(wd, defaultLayout)
				makeGoProject(wd)
				fmt.Println("Do you wanna install and initialize air? [y/N]")
				if scanner.Scan() {
					air := strings.TrimSpace(scanner.Text())
					switch air {
					case "y", "Y":
						initializeAir(wd)
					case "n", "N", "":
						return
					}
				}
			case "n":
				fmt.Print("Enter project name: ")
				if scanner.Scan() {
					folderName := strings.TrimSpace(scanner.Text())

					if err := validateProjectName(folderName); err != nil {
						log.Fatalf("[ERROR] Invalid project name: %v", err)
					}

					path := filepath.Join(wd, folderName)
					createFolders(path, defaultLayout)
					makeGoProject(path)
					fmt.Println("Do you wanna install and initialize air? [y/N]")
					if scanner.Scan() {
						air := strings.TrimSpace(scanner.Text())
						switch air {
						case "y", "Y":
							initializeAir(path)
						case "n", "N", "":
							return
						}
					}
				}
				return
			default:
				createFolders(wd, defaultLayout)
				makeGoProject(wd)
			}
		}
	} else {
		log.Fatalf("[ERROR] This program does not accept arguments, run it without any.")
	}
}
