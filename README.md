# CRLF replacer

Program for replacing CRLF with LF

## CLI

| flag             | default | usage                                                                                            |
| ---------------- | ------- | ------------------------------------------------------------------------------------------------ |
| `-path`          | `.`     |                                                                                                  |
| `-replace`       | `false` | should program replace CRLF with LF. If it is false, program just prints name of files with CRLF |
| `-ex-files`      |         | list of excluded files separated by comma                                                        |
| `-ex-extensions` |         | list of excluded extensions separated by comma                                                   |
| `-ex-folders`    |         | list of excluded folders separated by comma                                                      |

## Examples

**Note**: all files have CRLF ending.

**Folder structure**:

```
  |-- .gitignore
  |-- .git
    |-- ...
  |-- node_modules
    |-- empty_file.js
    |-- ...
    |-- 5k files in 2-3 lines
    |-- ...
    |-- last.js
  |-- web
    |-- web.go
    |-- api.go
    |-- parser
      |-- testdata
        |-- test.txt
      |-- parser.go
  |-- main.go
  |-- README.md
```

* `./crlf-replacer -ex-files .gitignore,README.md,test.txt -ex-folders node_modules,.git -replace`

  Command will **replace** CRLF with LF in files:
  * web.go
  * api.go
  * parser
  * parser.go
  * main.go

* `./crlf-replacer -ex-files .gitignore,README.md,test.txt -ex-folders node_modules,.git`

  Command will **print** next files:
  * web.go
  * api.go
  * parser
  * parser.go
  * main.go

* `./crlf -ex-extensions .go -ex-folders node_modules,.git`

  Command will **print** next files:
  * text.txt
  * README.md
