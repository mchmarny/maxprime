package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	count = flag.Int("count", 5, "client count")
	prime = flag.Int("prime", 999999, "calculate largest prime less than")
	url   = flag.String("url", "https://maxprime.demo.knative.tech", "endpoint to get")
)

type client struct {
	requestCount int
	lastResponse string
	err          error
}

func (c *client) start(stopCh <-chan struct{}) {
	tickerCh := time.NewTicker(time.Millisecond * 100).C
	for {
		select {
		case <-tickerCh:
			urlWithParams := fmt.Sprintf("%v/prime/%v", *url, *prime)
			resp, err := http.Get(urlWithParams)
			if err != nil {
				c.err = err
				c.lastResponse = err.Error()
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.err = err
				c.lastResponse = err.Error()
				resp.Body.Close()
				continue
			}
			resp.Body.Close()
			c.err = nil
			c.lastResponse = strings.TrimSpace(string(body))
			c.requestCount++
		case <-stopCh:
			return
		}
	}
}

func main() {
	flag.Parse()
	if *count < 1 {
		panic("count must be at least 1")
	}
	stopCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		close(stopCh)
	}()
	clients := make([]*client, *count)
	for i := 0; i < *count; i++ {
		c := &client{}
		go c.start((<-chan struct{})(stopCh))
		clients[i] = c
	}
	tickerCh := time.NewTicker(time.Second).C
	for {
		select {
		case <-tickerCh:
			fmt.Printf("ID\tCOUNT\tLAST RESPONSE\n")
			for i, client := range clients {
				fmt.Printf("%v\t%v\t%v\n", i, client.requestCount, client.lastResponse)
			}
			fmt.Printf("\n\n")
		case <-stopCh:
			os.Exit(0)
		}
	}
}
