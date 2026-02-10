package main

import (
	"crypto/sha256"
	"encoding/hex"
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
		encodedHashString, hashBytesWritten, fileHashErr := hashFile(path)

		if fileHashErr != nil {
			fmt.Println(fileHashErr)
		}

		resultForHashing := "\nHashed file " + path + " (" + strconv.Itoa(int(hashBytesWritten)) + " bytes)"

		fmt.Println(resultForHashing)

		hashSuccessResult := "Encoded hash for " + path + " is: " + encodedHashString

		fmt.Println(hashSuccessResult)
	}

	return nil
}

func hashFile(path string) (string, int64, error) {
	file, fileOpenErr := os.OpenFile(path, os.O_RDONLY, os.ModePerm)

	if fileOpenErr != nil {
		return "", int64(0), fileOpenErr
	}

	defer file.Close()

	hashSet := sha256.New()

	bytesWritten, copyErr := io.Copy(hashSet, file)

	if copyErr != nil {
		return "", bytesWritten, copyErr
	}

	// Finalize the hash
	hashInBytes := hashSet.Sum(nil)

	// Encode the hashed bytes
	encodedFileHash := hex.EncodeToString(hashInBytes)

	return encodedFileHash, bytesWritten, nil
}
