package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/bakerolls/httpcache"
	"github.com/bakerolls/httpcache/diskcache"
	"github.com/bakerolls/mangadex"
	"github.com/bakerolls/retry"
	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb"
)

var md *mangadex.Client

func main() {
	out := flag.String("out", ".", "Download directory")
	mid := flag.String("manga", "", "Manga ID")
	cid := flag.String("chapter", "", "Chapter ID")
	retries := flag.Int("retries", 5, "Retries in case a request fails")
	cacheDir := flag.String("cache", "cache", "Diretory to store cached API responses")
	backoff := flag.Duration("backoff", 100*time.Millisecond, "Backoff time between retries")
	workers := flag.Int("workers", 10, "Concurrent download workers")
	flag.Parse()

	cache := httpcache.New(
		diskcache.New(*cacheDir, diskcache.NoExpiration),
		httpcache.WithVerifier(httpcache.StatusInTwoHundreds),
		httpcache.WithTransport(retry.New(*retries, *backoff, nil)),
	)
	md = mangadex.New(mangadex.WithHTTPClient(cache.Client()))

	var reqs []*grab.Request
	var err error
	if *mid != "" {
		reqs, err = downloadManga(*out, *mid)
	}
	if *cid != "" {
		reqs, err = downloadChapter(*out, *cid)
	}
	if err != nil {
		log.Fatal(err)
	}

	bar := pb.StartNew(len(reqs))
	gc := grab.NewClient()
	resch := gc.DoBatch(*workers, reqs...)
	for res := range resch {
		if err := res.Err(); err != nil {
			log.Println(err)
			continue
		}
		bar.Increment()
	}
}

func downloadManga(out, id string) ([]*grab.Request, error) {
	m, cs, err := md.Manga(id)
	if err != nil {
		return nil, err
	}
	out = path.Join(out, m.Title)

	var reqs []*grab.Request
	for _, c := range cs {
		chReqs, err := downloadChapter(out, c.ID.String())
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, chReqs...)
	}
	return reqs, nil
}

func downloadChapter(out, id string) ([]*grab.Request, error) {
	c, err := md.Chapter(id)
	if err != nil {
		return nil, err
	}
	title := fmt.Sprintf("Vol. %s Ch. %s", c.Volume, c.Chapter)
	if c.Title != "" {
		title = fmt.Sprintf("%s - %s", title, c.Title)
	}

	var reqs []*grab.Request
	for _, image := range c.Images() {
		dst := path.Join(out, title, path.Base(image))
		req, err := grab.NewRequest(dst, image)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}
