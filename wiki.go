// This package holds the logic for reading a dir, finding out the main file,
// etc. It's a mere abstraction on top of the wiki's filesystem. When fed with
// a path, it'll return a read file for markdown consumption

package main

import (
	"errors"
	"os"
	"path"
)

func Lookup(wikiPath string, docPath string) ([]byte, error) {
	if docPath == "" {
		docPath = "README.md"
	}

	routePath := path.Join(wikiPath, docPath)
	fileInfo, err := os.Stat(routePath)

	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(routePath)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			switch file.Name() {
			case
				"index.md",
				"INDEX.md",
				"readme.md",
				"README.md":
				source, err := os.ReadFile(path.Join(routePath, file.Name()))

				if err != nil {
					return nil, err
				}

				return source, nil
			}
		}
	} else {
		source, err := os.ReadFile(routePath)

		if err != nil {
			return nil, err
		}

		return source, nil
	}

	return nil, errors.New("Something wen't really wrong")
}
