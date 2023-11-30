package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/temoto/robotstxt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetRandomUserAgent2() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
		"Mozilla/5.0 (iPad; CPU OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
		"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.58 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:97.0) Gecko/20100101 Firefox/97.0",
		"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		"Opera/9.80 (Windows NT 6.0) Presto/2.12.388 Version/12.14",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/74.0",
		"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/536.6",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Android 11; Mobile; LG-M255; rv:90.0) Gecko/90.0 Firefox/90.0",
		"Mozilla/5.0 (iPad; CPU OS 14_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/90.0.4430.212 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; Touch; rv:11.0) like Gecko",
		"Mozilla/5.0 (X11; CrOS x86_64 13729.56.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.95 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0",
		"Mozilla/5.0 (Android 10; Tablet; rv:68.0) Gecko/68.0 Firefox/68.0",
		"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.17",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
	}
	rand.Seed(int64(uint64(time.Now().UnixNano())))
	index := rand.Intn(len(userAgents))
	return userAgents[index]
}
func isURLAllowedByRobotsTXT2(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return false
	}

	if parsedURL.Host == "" {
		log.Println("Invalid URL, no host found:", urlStr)
		return false
	}

	robotsURL := "http://" + parsedURL.Host + "/robots.txt"

	resp, err := http.Get(robotsURL)
	if err != nil {
		log.Println("Error fetching robots.txt:", err)
		return true
	}

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		log.Println("Error parsing robots.txt:", err)
		return true
	}

	return data.TestAgent(urlStr, "GoEngine")
}

// Define the IndeedJobData struct to store job data
type IndeedJobData struct {
	Domain string `json:"domain"`
	URL    string `json:"url"`
	Data   []struct {
		Title      string `json:"title"`
		Company    string `json:"company"`
		Location   string `json:"location"`
		Salary     string `json:"salary,omitempty"`
		Summary    string `json:"summary"`
		DatePosted string `json:"date_posted"`
		Urgency    string `json:"urgency,omitempty"`
		Metadata   struct {
			Source    string `json:"source"`
			Timestamp string `json:"timestamp"`
		} `json:"metadata"`
	} `json:"data"`
}

// Helper function to clean string
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// Scraping function for Indeed
func scrapeIndeedJobs(baseURL string) {
	var indeedData IndeedJobData
	indeedData.Domain = "indeed"
	indeedData.URL = baseURL

	c := colly.NewCollector(
		colly.AllowedDomains("indeed.com", "www.indeed.com"),
	)
	c.OnHTML("div.mosaic-zone", func(e *colly.HTMLElement) {
		e.ForEach("div.jobsearch-SerpJobCard", func(_ int, el *colly.HTMLElement) {
			job := struct {
				Title      string `json:"title"`
				Company    string `json:"company"`
				Location   string `json:"location"`
				Salary     string `json:"salary,omitempty"`
				Summary    string `json:"summary"`
				DatePosted string `json:"date_posted"`
				Urgency    string `json:"urgency,omitempty"`
				Metadata   struct {
					Source    string `json:"source"`
					Timestamp string `json:"timestamp"`
				} `json:"metadata"`
			}{
				Title:      cleanString(el.ChildText("h2.jobTitle")),
				Company:    cleanString(el.ChildText("div.company_location")),
				Location:   cleanString(el.ChildText("div.company_location")),
				Salary:     cleanString(el.ChildText("div.heading6")),
				Summary:    cleanString(el.ChildText("div.job-snippet")),
				DatePosted: cleanString(el.ChildText("span.date")),
				Urgency:    cleanString(el.ChildText("div.css-tvvxwd")),
				Metadata: struct {
					Source    string `json:"source"`
					Timestamp string `json:"timestamp"`
				}{
					Source:    baseURL,
					Timestamp: time.Now().Format(time.RFC3339),
				},
			}
			indeedData.Data = append(indeedData.Data, job)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(baseURL)

	// Save the scraped data as JSON
	file, err := os.Create("indeed_jobs.json")
	if err != nil {
		log.Fatalf("Failed to create JSON file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(indeedData); err != nil {
		log.Fatalf("Failed to encode JSON data: %v", err)
	}

	log.Println("Indeed jobs data written to indeed_jobs.json")
}

type Job struct {
	Title    string
	Location string
	Company  string
	Salary   string
	Synopsis string
}

func parse(url string) ([]Job, error) {
	if !isURLAllowedByRobotsTXT(url) {
		return nil, fmt.Errorf("scraping not allowed by robots.txt: %s", url)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", GetRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		job := Job{
			Title:    strings.TrimSpace(s.Find("h2.jobTitle").Text()),
			Location: strings.TrimSpace(s.Find(".location").Text()),
			Company:  strings.TrimSpace(s.Find(".company").Text()),
			Salary:   strings.TrimSpace(s.Find(".no-wrap").Text()),
			Synopsis: strings.TrimSpace(s.Find(".summary").Text()),
		}

		if job.Title == "" {
			job.Title = "None"
		}
		if job.Location == "" {
			job.Location = "None"
		}
		if job.Company == "" {
			job.Company = "None"
		}
		if job.Salary == "" {
			job.Salary = "None"
		}

		jobs = append(jobs, job)
	})

	return jobs, nil
}

type JobListing struct {
	Title          string `json:"title"`
	Location       string `json:"location"`
	Company        string `json:"company"`
	ApplyURL       string `json:"apply_url"`
	Qualifications string `json:"qualifications"`
	Description    string `json:"description"`
}

type GoogleJobData struct {
	Domain   string       `json:"domain"`
	URL      string       `json:"url"`
	Data     []JobListing `json:"data"`
	Metadata struct {
		Source    string `json:"source"`
		Timestamp string `json:"timestamp"`
	} `json:"metadata"`
}

func ScrapeJobDescription(url string) string {
	c := colly.NewCollector(
		colly.UserAgent(GetRandomUserAgent()),
	)

	var description string
	c.OnHTML("div.aG5W3", func(e *colly.HTMLElement) {
		description = e.Text
	})

	c.Visit(url)
	return description
}

func ScrapeGoogleJobs(urls []string) GoogleJobData {
	var jobData GoogleJobData
	jobData.Domain = "Google"
	jobData.Metadata.Timestamp = time.Now().Format(time.RFC3339)

	c := colly.NewCollector(
		colly.UserAgent(GetRandomUserAgent()),
	)

	detailedCollector := c.Clone() // Collector for detailed job page

	for _, url := range urls {
		if !isURLAllowedByRobotsTXT(url) {
			log.Printf("Scraping blocked by robots.txt: %s\n", url)
			continue
		}
		jobData.URL = url
		jobData.Metadata.Source = url

		for page := 0; page < 2; page++ { // Limit to 2 pages
			pageURL := fmt.Sprintf("%s&page=%d", url, page)

			c.OnHTML("div.sMn82b", func(e *colly.HTMLElement) {
				title := strings.TrimSpace(e.ChildText("h3.QJPWVe"))
				applyURL := "https://www.google.com/about/careers/applications/" + e.ChildAttr("a", "href")

				qualificationsText := e.ChildText("div.Xsxa1e")
				lines := strings.Split(qualificationsText, "\n")
				if len(lines) > 1 && strings.Contains(lines[0], "Minimum qualifications") {
					lines = lines[1:] // Skip the first line
				}
				finalQualifications := strings.Join(lines, " ")

				locationText := strings.TrimSpace(e.ChildText("div.EAcu5e.Gx4ovb"))
				locationText = strings.Replace(locationText, "Google | ", "", 1) // Remove "Google |" prefix

				locationParts := strings.Split(locationText, "; ")
				filteredLocationParts := []string{}
				for _, part := range locationParts {
					if strings.Contains(part, "+") && len(filteredLocationParts) > 0 {
						filteredLocationParts[len(filteredLocationParts)-1] = part
					} else {
						filteredLocationParts = append(filteredLocationParts, part)
					}
				}
				locationText = strings.Join(filteredLocationParts, "; ")

				// Visiting each job's apply URL for detailed data
				detailedCollector.OnHTML("div.DkhPwc", func(e *colly.HTMLElement) {
					description := strings.TrimSpace(e.ChildText("div.aG5W3"))
					description = strings.TrimPrefix(description, "About the job") // Remove "About the job"

					jobData.Data = append(jobData.Data, JobListing{
						Title:          title,
						Location:       locationText,
						Company:        "Google",
						ApplyURL:       applyURL,
						Qualifications: finalQualifications,
						Description:    description,
					})
				})

				detailedCollector.Visit(applyURL)
			})

			c.Visit(pageURL)
		}
	}

	return jobData
}

func main2() {
	urls := []string{
		"https://www.google.com/about/careers/applications/jobs/results/?location=USA",
		// Add more URLs as needed
	}
	jobData := ScrapeGoogleJobs(urls)

	file, err := json.MarshalIndent(jobData, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	_ = os.WriteFile("googlejobs.json", file, 0644)
	log.Println("Google job data written to google_jobs.json")
}

//func main() {
//	url := "https://www.indeed.com/jobs?q=software+developer&l="
//	jobs, err := parse(url)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Marshal the jobs into JSON
//	jsonData, err := json.MarshalIndent(jobs, "", "  ")
//	if err != nil {
//		log.Fatalf("Error marshalling data to JSON: %s", err)
//	}
//
//	// Write the JSON data to a file
//	file, err := os.Create("jobs.json")
//	if err != nil {
//		log.Fatalf("Error creating JSON file: %s", err)
//	}
//	defer file.Close()
//
//	_, err = file.Write(jsonData)
//	if err != nil {
//		log.Fatalf("Error writing data to JSON file: %s", err)
//	}
//
//	fmt.Println("Job data written to jobs.json")
//}
