package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if err := do(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

type godeps struct {
	ImportPath string
	GoVersion  string
	Packages   []string `json:",omitempty"`
	Deps       []dependency
}

type dependency struct {
	ImportPath string
	Comment    string `json:",omitempty"`
	Rev        string
}

func do() (retErr error) {
	file, err := os.Open("Godeps/Godeps.json")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	var godeps godeps
	if err := json.NewDecoder(file).Decode(&godeps); err != nil {
		return err
	}
	successImportPaths := make(map[string]bool)
	for _, dep := range godeps.Deps {
		if dep.ImportPath != "" {
			importPath := dep.ImportPath
			for importPath != "" {
				if _, ok := successImportPaths[importPath]; ok {
					break
				}
				if err := runCommand("godep", "update", importPath+"/..."); err == nil {
					fmt.Println(importPath)
					successImportPaths[importPath] = true
					break
				}
				importPath, _ = filepath.Split(importPath)
				if importPath != "" && strings.HasSuffix(importPath, "/") {
					importPath = importPath[:len(importPath)-1]
				}
			}
		}
	}
	return nil
}

func runCommand(args ...string) error {
	if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
		return fmt.Errorf("%v: %v", args, err)
	}
	return nil
}
