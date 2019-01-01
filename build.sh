#!/bin/sh

rm -r build
mkdir build

# Widnows
export GOARCH=amd64
export GOOS=windows
go build -o "build/crlf-replacer.exe"

# Linux
export GOARCH=arm64
export GOOS=linux
go build -o "build/crlf-replacer"
