package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/peter-edge/go-tools/common"
)

func main() {
	common.Main(do)
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

func do() error {
	data, err := common.ReadAll("Godeps/Godeps.json")
	if err != nil {
		return err
	}
	var godeps godeps
	if err := json.Unmarshal(data, &godeps); err != nil {
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
				if err := common.Cmd(nil, nil, "godep", "update", importPath+"/..."); err == nil {
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
