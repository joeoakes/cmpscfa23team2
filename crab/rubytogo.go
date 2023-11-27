package main

//
//func GetRandomUserAgent2() string {
//	userAgents := []string{
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
//		"Mozilla/5.0 (iPad; CPU OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
//		"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.58 Mobile Safari/537.36",
//		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36",
//		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:97.0) Gecko/20100101 Firefox/97.0",
//		"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko",
//		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
//		"Opera/9.80 (Windows NT 6.0) Presto/2.12.388 Version/12.14",
//		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/74.0",
//		"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/536.6",
//		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
//		"Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
//		"Mozilla/5.0 (Android 11; Mobile; LG-M255; rv:90.0) Gecko/90.0 Firefox/90.0",
//		"Mozilla/5.0 (iPad; CPU OS 14_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
//		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:91.0) Gecko/20100101 Firefox/91.0",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
//		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/90.0.4430.212 Mobile/15E148 Safari/604.1",
//		"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; Touch; rv:11.0) like Gecko",
//		"Mozilla/5.0 (X11; CrOS x86_64 13729.56.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.95 Safari/537.36",
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; rv:11.0) like Gecko",
//		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
//		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Safari/605.1.15",
//		"Mozilla/5.0 (X11; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0",
//		"Mozilla/5.0 (Android 10; Tablet; rv:68.0) Gecko/68.0 Firefox/68.0",
//		"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.17",
//		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
//	}
//	rand.Seed(time.Now().UnixNano())
//	return userAgents[rand.Intn(len(userAgents))]
//}
//
//type Scraper struct {
//	link      string
//	userAgent string
//}
//
//func NewScraper(link string) *Scraper {
//	return &Scraper{
//		link:      link,
//		userAgent: GetRandomUserAgent(),
//	}
//}
//
//func (s *Scraper) ScrapeJobDetails(detailURL string) (map[string]string, error) {
//	detailResponse, err := http.Get(detailURL)
//	if err != nil {
//		return nil, err
//	}
//	defer detailResponse.Body.Close()
//
//	// Proceed if the GET request succeeded
//	if detailResponse.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("error fetching detail page: %d %s", detailResponse.StatusCode, detailURL)
//	}
//
//	// Load the detail page HTML
//	detailDoc, err := goquery.NewDocumentFromReader(detailResponse.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	// Extract job details based on the provided HTML structure
//	jobDetails := make(map[string]string)
//	jobDetails["title"] = detailDoc.Find(".jobview-containerstyles__JobViewContainer-sc-16af7k7-3 h2").Text()
//	jobDetails["company"] = detailDoc.Find(".jobview-containerstyles__JobViewContainer-sc-16af7k7-3 h2").Next().Text()
//	jobDetails["location"] = detailDoc.Find(".headerstyle__JobViewHeaderLocation-sc-1ijq9nh-4").Text()
//	jobDetails["description"] = detailDoc.Find(".descriptionstyles__DescriptionContainer-sc-13ve12b-0").Text()
//	// ... include more fields as needed
//
//	return jobDetails, nil
//}
//
//type Loop struct {
//	Scraper // Embedded type, gives Loop access to Scraper's methods
//	total   int
//	page    int
//}
//
//func NewLoop(total, page int, link string) *Loop {
//	return &Loop{
//		Scraper: Scraper{link: link, userAgent: GetRandomUserAgent()},
//		total:   total,
//		page:    page,
//	}
//}
//
//func (l *Loop) Start() error {
//	allJobs := make([]map[string]string, 0)
//	maxPages := 10 // Assuming you want to iterate over a constant number of pages
//
//	for l.page <= maxPages {
//		delay := time.Duration(rand.Intn(5)+1) * time.Second
//		fmt.Printf("Sleeping for %v before fetching page %d\n", delay, l.page)
//		time.Sleep(delay)
//
//		// Randomize UserAgent for each request
//		l.Scraper.userAgent = GetRandomUserAgent()
//
//		// Construct the page-specific URL
//		pageLink := fmt.Sprintf("%s&pn=%d", l.link, l.page)
//		l.Scraper.link = pageLink
//
//		// Assume ScrapePage is a method of Scraper that takes the page URL as a string and returns a slice of job maps and an error
//		jobDetails, err := l.Scraper.ScrapeJobDetails(pageLink)
//		if err != nil {
//			log.Printf("Error scraping page: %v\n", err)
//			return err
//		}
//
//		fmt.Printf("Page %d - Job: %+v\n", l.page, jobDetails)
//		allJobs = append(allJobs, jobDetails)
//
//		// Optional: Save jobs to JSON after fetching each page
//		err = WriteJobsToJSON(allJobs, "jobs.json")
//		if err != nil {
//			log.Printf("Error writing to JSON: %v\n", err)
//			return err
//		}
//
//		l.page++
//	}
//
//	// Final JSON write is not needed here if you're saving after each page
//
//	return nil
//}
//
//func WriteJobsToJSON(jobs []map[string]string, filename string) error {
//	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	jsonData, err := json.Marshal(jobs)
//	if err != nil {
//		return err
//	}
//
//	// Write to file
//	_, err = file.Write(jsonData)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func main() {
//	// For Monster.com, you need to identify the search URL's structure. An example might be:
//	// "https://www.monster.com/jobs/search?q=Software+Engineer&stpage=1&endpage=3" where stpage and endpage define your loop's start and end.
//	// You need to adjust the URL according to how Monster.com's pagination works. If unsure, observe the website's behavior in a web browser.
//
//	// Make sure to update the URL with the correct search parameters and pagination structure.
//	loop := NewLoop(3, 1, "https://www.monster.com/jobs/search?q=Software+Engineer")
//	if err := loop.Start(); err != nil {
//		log.Fatal(err)
//	}
//}
