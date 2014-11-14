package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ActiveState/tail"
	"net/http"
)

func HitUrl(url string) (byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error getting %s: %s", url, err))
	}

	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error initialize scanner for %s: %s", url, err))
	}

	return r.ReadByte()
}

func main() {
	logFilePath := flag.String("file", "varnish-urls.log", "location of log of urls to follow")
	flag.Parse()

	fmt.Println(fmt.Sprintf("Beginning to follow file: %s", *logFilePath))

	offset := tail.SeekInfo{Offset: 0, Whence: 2}
	t, err := tail.TailFile(*logFilePath, tail.Config{Follow: true, Location: &offset})
	if err != nil {
		fmt.Println(fmt.Sprintf("Error encountered opening file: %s", err))
	}

	// defer t.Stop()

	for url := range t.Lines {
		fmt.Println(fmt.Sprintf("Hitting url: %s", url.Text))
		go HitUrl(url.Text)
	}
}
