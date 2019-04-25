# MangaDex Downloader

[![Go Report Card](https://goreportcard.com/badge/github.com/bake/mddl)](https://goreportcard.com/report/github.com/bake/mddl)

Download mangas (or chapters) from MangaDex. Provide `-manga` or `-chapter` along with the corresponding manga- or chapter-ID.

## Installation

```
$ go get github.com/bake/mddl
```

## Usage

```
$ mddl -chapter 517244
 1 / 1 [===========================================================] 100.00% 0s
```

```
$ mddl -help
Usage of mddl:
  -backoff duration
        Backoff time between retries (default 100ms)
  -cache string
        Diretory to store cached API responses (default "cache")
  -chapter string
        Chapter ID
  -manga string
        Manga ID
  -out string
        Download directory (default ".")
  -retries int
        Retries in case a request fails (default 5)
```
