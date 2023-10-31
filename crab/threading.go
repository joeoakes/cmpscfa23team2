/*

First:
go get github.com/PuerkitoBio/goquery
go get github.com/temoto/robotstxt

Explanation:
Domain-Agnostic Element Selection:
We use the variable targetElement to specify which elements we want to scrape.
In this case, .quote .text targets quote texts on the website.
You can change this selector to any valid CSS selector to target other elements.

Robots.txt Handling:
We first fetch the robots.txt file of our target site to ensure we are allowed to access our desired URLs.
The fetchRobotsTxt function gets the robots.txt and the main function checks if our bot (default is * for any bot) is allowed to access the URLs.

Rate Limiting:
We introduce a delay before every request using time.Sleep(delay) to respect the website's request rate.
This is especially important for real-world scraping to avoid overwhelming a server and getting blocked.

Concurrent Scraping:
use Goroutines to fetch multiple pages concurrently.
We manage concurrency using a WaitGroup (wg) and use channels (ch) to communicate results between Goroutines.

Flexible Scraping:
The scrape function accepts a targetElement parameter which is the CSS selector of the elements you want to scrape.
This makes the function domain agnostic, as you can use it to target different elements on different websites.

This implementation is respectful, domain agnostic, and allows for flexible scraping by just specifying the CSS selector of the target elements.
Always ensure you have the legal right to scrape a website, and respect its robots.txt and terms of use.
*/

package main

//import (
//	"flag"
//	"fmt"
//	"log"
//	"net/http"
//	"sync"
//	"time"
//
//	"github.com/PuerkitoBio/goquery"
//	"github.com/temoto/robotstxt"
//)
//
//type ScrapeResult struct {
//	URL   string
//	Data  string
//	Error error
//}
//
//func fetchRobotsTxt(url string) (*robotstxt.Group, error) {
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	data, err := robotstxt.FromResponse(resp)
//	if err != nil {
//		return nil, err
//	}
//
//	return data.FindGroup("*"), nil
//}
//
//func scrape(url string, targetElement string, ch chan<- ScrapeResult, delay time.Duration) {
//	time.Sleep(delay)
//
//	resp, err := http.Get(url)
//	if err != nil {
//		ch <- ScrapeResult{URL: url, Error: err}
//		return
//	}
//	defer resp.Body.Close()
//
//	doc, err := goquery.NewDocumentFromReader(resp.Body)
//	if err != nil {
//		ch <- ScrapeResult{URL: url, Error: err}
//		return
//	}
//
//	doc.Find(targetElement).Each(func(index int, element *goquery.Selection) {
//		ch <- ScrapeResult{URL: url, Data: element.Text()}
//	})
//}
//
//// For testing reasons if this main does not compile, try the other one
//func main() {
//	var bypassRobotsCheck bool
//	flag.BoolVar(&bypassRobotsCheck, "bypass-robots", false, "Bypass the robots.txt check")
//	flag.Parse()
//
//	baseURL := "http://quotes.toscrape.com"
//	robotsURL := baseURL + "/robots.txt"
//	urls := []string{
//		baseURL + "/page/1/",
//	}
//
//	if !bypassRobotsCheck {
//		robotsGroup, err := fetchRobotsTxt(robotsURL)
//		if err != nil {
//			log.Fatalf("Failed fetching robots.txt: %v", err)
//		}
//
//		for _, url := range urls {
//			if !robotsGroup.Test(url) {
//				log.Fatalf("Access denied by robots.txt for URL: %s", url)
//			}
//		}
//	}
//
//	targetElement := ".quote .text"
//
//	ch := make(chan ScrapeResult)
//	var wg sync.WaitGroup
//
//	delay := 1 * time.Second
//	for _, url := range urls {
//		wg.Add(1)
//		go func(u string) {
//			defer wg.Done()
//			scrape(u, targetElement, ch, delay)
//		}(url)
//	}
//
//	go func() {
//		wg.Wait() // Wait for all scraping Goroutines to finish
//		close(ch) // Then close the channel
//	}()
//
//	for result := range ch {
//		if result.Error != nil {
//			log.Printf("Error scraping %s: %v", result.URL, result.Error)
//		} else {
//			fmt.Printf("Scraped %s: %s\n", result.URL, result.Data)
//		}
//	}
//}

//
//func main() {
//	baseURL := "http://quotes.toscrape.com"
//	// robotsURL := baseURL + "/robots.txt" // Commented out for now
//	urls := []string{
//		baseURL + "/page/1/",
//	}
//
//	// Comment out robots.txt check for debugging
//	/*
//	   robotsGroup, err := fetchRobotsTxt(robotsURL)
//	   if err != nil {
//	       log.Fatalf("Failed fetching robots.txt: %v", err)
//	   }
//
//	   for _, url := range urls {
//	       if !robotsGroup.Test(url) {
//	           log.Fatalf("Access denied by robots.txt for URL: %s", url)
//	       }
//	   }
//	*/
//
//	targetElement := ".quote .text"
//
//	ch := make(chan ScrapeResult)
//	var wg sync.WaitGroup
//
//	delay := 1 * time.Second
//	for _, url := range urls {
//		wg.Add(1)
//		go func(u string) {
//			defer wg.Done()
//			scrape(u, targetElement, ch, delay)
//		}(url)
//	}
//
//	go func() {
//		wg.Wait() // Wait for all scraping Goroutines to finish
//		close(ch) // Then close the channel
//	}()
//
//	for result := range ch {
//		if result.Error != nil {
//			log.Printf("Error scraping %s: %v", result.URL, result.Error)
//		} else {
//			fmt.Printf("Scraped %s: %s\n", result.URL, result.Data)
//		}
//	}
//}
