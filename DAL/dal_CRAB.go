package DAL

/*
import (
	"colly"
	_ "errors"
	"fmt"
	_ "github.com/gocolly/colly/..."
	"log"
)

func main() {
	// Create a new collector
	c := colly.NewCollector()
	// Set up a callback for when a visited HTML element is found
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		fmt.Println("First column of a table row:", e.Text)
	})

	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Start the web crawl on a specific website
	err := c.Visit("http://google.com")
	if err != nil {
		log.Fatal(err)
	}
}
*/
