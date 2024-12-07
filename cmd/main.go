package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	url     string
	workers int
	delay   int
)

func main() {
	flag.StringVar(&url, "url", "", "TODO")
	flag.IntVar(&workers, "w", 1, "")
	flag.IntVar(&delay, "d", 1, "")
	flag.Parse()

	wordCh := make(chan string, 0)
	s := make(chan string, 1)
	f := make(chan string, 0)

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		file, err := os.Open("/usr/local/share/dirb/wordlists/common.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			wordCh <- scanner.Text()
		}

		wg.Done()
	}()

	for w := 1; w <= workers; w++ {
		go checkWord(wordCh, s, f)
	}

	for {
		select {
		case exists := <-s:
			fmt.Println(exists)
		}
	}
}

func checkWord(wordCh <-chan string, s, f chan<- string) {
	for w := range wordCh {
		fmt.Println(w)
		t := time.Millisecond
		time.Sleep(t * 10)

		r, err := http.Get(fmt.Sprintf("%s/%s", url, w))
		if err != nil || r.StatusCode == http.StatusNotFound {
			continue
		}

		s <- w
	}
}
