@ECHO off

rd /s /q build
mkdir build

:: Widnows
set GOARCH=amd64
SET GOOS=windows
go build -o "build\crlf-replacer.exe"

:: Linux
set GOARCH=arm64
set GOOS=linux
go build -o "build\crlf-replacer"
