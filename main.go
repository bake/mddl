package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/bakerolls/httpcache"
	"github.com/bakerolls/httpcache/diskcache"
	"github.com/bakerolls/mangadex"
)

var (
	md  *mangadex.Client
	out = flag.String("out", ".", "Download directory")
	mid = flag.String("manga", "", "Manga ID")
	cid = flag.String("chapter", "", "Chapter ID")
)

func main() {
	flag.Parse()

	cache := httpcache.New(diskcache.New("cache", diskcache.NoExpiration))
	md = mangadex.New(mangadex.WithHTTPClient(cache.Client()))

	var err error
	if *mid != "" {
		err = downloadManga(*out, *mid)
	}
	if *cid != "" {
		err = downloadChapter(*out, *cid)
	}
	if err != nil {
		log.Println(err)
	}
}

func downloadManga(out, id string) error {
	m, cs, err := md.Manga(id)
	if err != nil {
		return err
	}
	mOut := path.Join(out, m.Title)
	for _, c := range cs {
		if err := downloadChapter(mOut, c.ID.String()); err != nil {
			return err
		}
	}
	return nil
}

func downloadChapter(out, id string) error {
	c, err := md.Chapter(id)
	if err != nil {
		return err
	}
	title := fmt.Sprintf("Vol. %s Ch. %s", c.Volume, c.Chapter)
	if c.Title != "" {
		title = fmt.Sprintf("%s - %s", title, c.Title)
	}
	cOut := path.Join(out, title)
	if err := os.MkdirAll(cOut, 0744); err != nil {
		return err
	}
	fmt.Println(title)
	for i, url := range c.Images() {
		if err := download(i, url, cOut); err != nil {
			return err
		}
	}
	return nil
}

func download(i int, url, out string) error {
	out = path.Join(out, path.Base(url))
	if _, err := os.Stat(out); err == nil {
		return nil
	}

	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}
	return nil
}
