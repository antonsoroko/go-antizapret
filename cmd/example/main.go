package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/antonsoroko/go-antizapret"
)

func main() {
	antizapretProxy := antizapret.NewAntizapretProxy()

	// Test cases
	fmt.Println("Testing proxy detection:")
	torrentTrackers := []string{
		"rutracker.org",
		"rutracker.net",
		"rutor.info",
		"rutor.is",
		"nnmclub.to",
		"kinozal.tv",
		"rustorka.com",
		"megapeer.vip",
		"bluebird-hd.org",
		"pirat.one",
		"www.lostfilm.tv",
		"www.lostfilm.download",
		"newstudio.tv",
		"le-production.tv",
		"anilibria.tv",
		"anilibria.top",
		"anilibria.wtf",
		"tr.anidub.com",
		"animaunt.fun",
		"tracker.0day.community",
		"baibako.tv",
		"toloka.to",
		"thepiratebay.org",
		"thepiratebay10.xyz",
		"1337x.to",
		"glodls.to",
		"yts.mx",
		"eztvx.to",
		"bt4gprx.com",
		"bitsearch.to",
		"torrenting.com",
		"www.limetorrents.lol",
		"uindex.org",
	}
	for _, domain := range torrentTrackers {
		fmt.Printf("%s -> %s\n", domain, antizapretProxy.Detect(domain))
	}
	fmt.Println(strings.Repeat("-", 20))

	testURLs := []string{"https://rutracker.org", "https://rutracker.net/forum/viewtopic.php?t=5324346", "https://rutor.info/torrent/472", "https://rutor.is", "https://kinozal.tv", "https://www.lostfilm.tv:443", "https://ifconfig.co:443/ip"}
	for _, urlStr := range testURLs {
		fmt.Printf("\nTesting opening URL: %s\n", urlStr)
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			log.Fatal(err)
			continue
		}
		proxyURL := antizapretProxy.Detect(parsedURL.Hostname())
		fmt.Printf("Detected Antizapret proxy for %s: %s\n", parsedURL.Hostname(), proxyURL)

		client := &http.Client{}
		if proxyURL != "" {
			proxyParsedURL, err := url.Parse(proxyURL)
			if err != nil {
				fmt.Printf("Error parsing proxy URL: %s\n", err)
				continue
			} else {
				client = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyParsedURL),
					},
				}
			}
		}

		response, err := client.Get(urlStr)
		if err != nil {
			fmt.Printf("Can't make a request: %s\n", err)
			continue
		}
		defer response.Body.Close()

		fmt.Printf("Code: %d\n", response.StatusCode)

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading body: %s\n", err)
			continue
		}
		bodyStr := string(body)
		if len(bodyStr) < 100 {
			fmt.Printf("Body: %s\n", bodyStr)
		} else {
			fmt.Printf("Body length: %d\n", len(bodyStr))
		}
	}
}
