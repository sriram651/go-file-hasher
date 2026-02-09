package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	directoryPath := flag.String("dir", "", "The path to the directory that needs to be parsed.")

	flag.Parse()

	if *directoryPath == "" {
		fmt.Println("Flag --dir is required!")
		os.Exit(2)
	}

	fmt.Println("Path to directory as given:", *directoryPath)

	walkDirErr := filepath.WalkDir(*directoryPath, walkDirCallback)

	if walkDirErr != nil {
		fmt.Println(walkDirErr)
		os.Exit(1)
	}
}

func walkDirCallback(path string, d fs.DirEntry, err error) error {
	if d.IsDir() && (d.Name() == "node_modules" || d.Name() == "android" || d.Name() == "dist" || d.Name() == "ios" || d.Name() == "public" || strings.HasPrefix(d.Name(), ".")) {
		return filepath.SkipDir
	}

	if !d.IsDir() {
		fmt.Println(d)
	}

	return nil
}
