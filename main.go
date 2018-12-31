package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	excludeFolders := []string{
		"node_modules",
		".git",
	}

	excludeExtensions := []string{
		".png",
		".jpg",
		".eot",
		".ttf",
		".woff",
		".woff2",
	}

	allPaths := make(chan string)
	crlfFilesPaths := RunPool(5, allPaths)
	done := make(chan struct{})

	// Printer function
	go func() {
		fmt.Print("Start\n\n")

		c := 0
		fmt.Print("Files with CRLF:")
		for path := range crlfFilesPaths {
			if c == 0 {
				fmt.Print("\n")
			}
			fmt.Printf("* %s\n", path)
			c++
		}

		if c == 0 {
			fmt.Println(" Empty")
		}

		fmt.Println("\nDone")

		close(done)
	}()

	filepath.Walk(".", filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, s := range excludeFolders {
			if strings.Contains(path, s) {
				return nil
			}
		}

		ext := filepath.Ext(info.Name())
		for _, s := range excludeExtensions {
			if ext == s {
				return nil
			}
		}

		allPaths <- path

		return nil
	}))

	close(allPaths)

	<-done
}

// RunPool runs passed number of workers
func RunPool(number int, paths <-chan string) <-chan string {
	results := make(chan string, 10)

	go func() {
		wg := new(sync.WaitGroup)

		for i := 0; i < number; i++ {
			wg.Add(1)
			go runWorker(paths, results, wg)
		}
		wg.Wait()

		close(results)
	}()

	return results
}

func runWorker(paths <-chan string, results chan<- string, wg *sync.WaitGroup) {
	for path := range paths {
		if hasCRLF(path) {
			results <- path
		}
	}

	wg.Done()
}

func hasCRLF(path string) bool {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		fmt.Printf("[ERR] can't open file %s: %s\n", path, err)
		return false
	}

	scanner := bufio.NewScanner(file)
	// Have to use custom split function to have \r and \n in result of scanner.Bytes()
	scanner.Split(splitFunction)

	for scanner.Scan() {
		bytes := scanner.Bytes()
		// CR == 13 (0x0D), LF == 10 (0x0A)
		if len(bytes) > 2 && bytes[len(bytes)-2] == 0x0D && bytes[len(bytes)-1] == 0x0A {
			return true
		}
	}

	return false
}

// Secondary functions

// splitFunction saves \r and \n
func splitFunction(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// Return with \r and \n
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
