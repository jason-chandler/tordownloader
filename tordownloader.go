package main

import (
	"io"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	proxy, _ := url.Parse("socks5://127.0.0.1:9150")
	userAgent := "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/115.0"
	fmt.Println("target url: ")
	var targetUrl, fileName string
	_, err := fmt.Scanln(&targetUrl)
	check(err)
	fmt.Println("file name: ")
	_, err2 := fmt.Scanln(&fileName)
	check(err2)

	// targetUrl := "http://check.torproject.org"
	// targetUrl := "https://media.tenor.com/rxi2UfL2LFcAAAPo/happy-birthday.mp4"

	i := 0
	for i < 100 {
		h := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}, CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
		req, err := http.NewRequest("GET", targetUrl, nil)
		check(err)
		req.Header.Set("User-Agent", userAgent)

		resp, err := h.Do(req)
		check(err)
		defer resp.Body.Close()

		if resp.StatusCode == 301 || resp.StatusCode == 302 {
			targetUrl = resp.Header.Get("Location")
			i += 1
		} else {
			file, err := os.Create(fileName)
			check(err)
			size, err := io.Copy(file, resp.Body)
			time.Sleep(7)
			fmt.Printf("Downloaded a file with size %d", size)
			break
		}
	}
	fmt.Println("\nPress any key to continue. And Gopher bless.")
	fmt.Scanln()
}
