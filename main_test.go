package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestHasCRLF(t *testing.T) {
	tests := []struct {
		ctlfFile string
		lfFile   string
	}{
		{"testdata/crlf.txt", "testdata/lf.txt"},
		{"testdata/crlf.md", "testdata/lf.md"},
		{"testdata/crlf.json", "testdata/lf.json"},
	}

	for i, tt := range tests {
		if !hasCRLF(tt.ctlfFile) {
			t.Errorf("Test #%d: can't find CRLF in file %s\n", i, tt.ctlfFile)
		}
	}
}

func TestReplaceCRLF(t *testing.T) {
	checkFileErr := func(err error, format string, args ...interface{}) {
		if err != nil {
			msg := fmt.Sprintf(format, args...)
			t.Fatalf("%s: %s", msg, err)
		}
	}

	areFilesEqual := func(path1, path2 string) bool {
		f1, err := ioutil.ReadFile(path1)
		if err != nil {
			checkFileErr(err, "can't open file %s", path1)
		}
		f2, err := ioutil.ReadFile(path2)
		if err != nil {
			checkFileErr(err, "can't open file %s", path1)
		}

		return bytes.Equal(f1, f2)
	}

	const (
		folder     = "testdata/"
		tempFolder = "temp/"
	)

	tests := []struct {
		ctlfFile string
		lfFile   string
	}{
		{"crlf.txt", "lf.txt"},
		{"crlf.md", "lf.md"},
		{"crlf.json", "lf.json"},
	}

	// Copy test files into temp directory
	os.Mkdir(tempFolder, 0777)
	for _, tt := range tests {
		f1, err := os.Open(folder + tt.ctlfFile)
		checkFileErr(err, "can't open file %s", folder+tt.ctlfFile)
		f2, err := os.Create(tempFolder + tt.ctlfFile)
		checkFileErr(err, "can't create file %s", tempFolder+tt.ctlfFile)

		io.Copy(f2, f1)

		f1.Close()
		f2.Close()
	}

	for i, tt := range tests {
		err := replaceCRLF(folder + tt.ctlfFile)
		if err != nil {
			t.Fatalf("Test #%d: can't replace CRLF: %s", i, err)
		}

		if !areFilesEqual(folder+tt.ctlfFile, folder+tt.lfFile) {
			t.Errorf("Test #%d: files are not equal", i)
		}
	}

	// Recover files from temp folder
	for _, tt := range tests {
		f1, err := os.OpenFile(folder+tt.ctlfFile, os.O_TRUNC|os.O_RDWR, 0600)
		checkFileErr(err, "can't open file %s", folder+tt.ctlfFile)
		f2, err := os.Open(tempFolder + tt.ctlfFile)
		checkFileErr(err, "can't open file %s", tempFolder+tt.ctlfFile)

		io.Copy(f1, f2)

		f1.Close()
		f2.Close()
	}

	err := os.RemoveAll(tempFolder)
	if err != nil {
		t.Errorf("can't remove temp folder: %s", err)
	}
}
