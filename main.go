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
	"strings"
	"sync"
)

// Determines number of concurrent go-routines
var CONCURRENT_WORKERS int = 10

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

	var walkDirWaitGroup sync.WaitGroup

	// Initialize the channel
	jobs := make(chan string)

	// Loop deploys the workers
	for i := 0; i < CONCURRENT_WORKERS; i++ {
		// This is the worker function that waits for the job
		go func() {
			for job := range jobs {
				encodedHashString, hashBytesWritten, fileHashErr := hashFile(job)

				if fileHashErr != nil {
					fmt.Println(fileHashErr)
				}

				fmt.Println("FILE:", job)
				fmt.Println("SIZE:", hashBytesWritten, "bytes")
				fmt.Println("HASH:", encodedHashString)
				fmt.Println()

				// Report done
				walkDirWaitGroup.Done()
			}
		}()
	}

	// To get access to the vars in main, we include this in the closure of main()
	walkDirCallback := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			walkDirWaitGroup.Add(1)
			// Send the path to the jobs channel
			jobs <- path
		}

		return nil
	}

	walkDirErr := filepath.WalkDir(*directoryPath, walkDirCallback)

	if walkDirErr != nil {
		fmt.Println(walkDirErr)
		os.Exit(1)
	}

	// Close the channel to let the workers know that no more jobs are coming
	close(jobs)

	// Wait for all the jobs to report Done
	walkDirWaitGroup.Wait()
}

func hashFile(path string) (string, int64, error) {
	file, fileOpenErr := os.Open(path)

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
