package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
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

	_, readDirErr := os.ReadDir(*directoryPath)

	if readDirErr != nil {
		fmt.Println(readDirErr)
		os.Exit(1)
	}

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
		fileHashErr := hashFile(path)

		if fileHashErr != nil {
			fmt.Println(fileHashErr)
		}
	}

	return nil
}

func hashFile(path string) error {
	file, fileOpenErr := os.OpenFile(path, os.O_RDONLY, os.ModePerm)

	if fileOpenErr != nil {
		return fileOpenErr
	}

	defer file.Close()

	hashSet := sha256.New()

	bytesWritten, copyErr := io.Copy(hashSet, file)

	if copyErr != nil {
		return copyErr
	}

	resultForHashing := "Hashed file " + path + "(" + strconv.Itoa(int(bytesWritten)) + " bytes)\n"

	fmt.Println(resultForHashing)

	// TODO: We need to do something with the hashed file here!

	return nil
}
