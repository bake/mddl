# MangaDex Downloader

[![Go Report Card](https://goreportcard.com/badge/github.com/bakerolls/mddl)](https://goreportcard.com/report/github.com/bakerolls/mddl)

Download mangas (or chapters) from MangaDex. Provide `-manga` or `-chapter` along with the corresponding manga- or chapter-ID.

```
$ go install github.com/bakerolls/mddl
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
