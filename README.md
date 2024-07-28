# gocmp

[![build](https://github.com/m0t9/gocmp/actions/workflows/go.yml/badge.svg)](https://github.com/m0t9/gocmp/actions/workflows/go.yml)
[![coverage](https://raw.githubusercontent.com/m0t9/gocmp/badges/.badges/master/coverage.svg)](https://github.com/m0t9/gocmp/actions/workflows/.testcoverage.yml)

A simple implementation of Huffman compression tool on Go-language.



## Build

`go build -o gocmp cmd/gocmp/main.go`

## Usage

### Compression

```sh
./gocmp src-path compressed-path
(=^ ◡ ^=) successfully compressed to file 'compressed-filename'
( ^..^)ﾉ  compression rate is 1.42
(^･o･^)ﾉ  gocmp running time is 339.456083ms
```

### Decompression

```sh
./gocmp -d compressed-path decompressed-path
(=^ ◡ ^=) successfully decompressed to file 'decompressed-filename'
(^･o･^)ﾉ  gocmp running time is 603.004375ms
```
