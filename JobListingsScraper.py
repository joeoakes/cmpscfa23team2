
import os
from urllib.robotparser import RobotFileParser
import et as et
from selenium.common.exceptions import NoSuchElementException, TimeoutException
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import json
import random
import datetime
from bs4 import BeautifulSoup
from lxml import etree as et
from selenium import webdriver
from selenium.common import NoSuchElementException
from selenium.webdriver.common.by import By
from selenium.webdriver.firefox.options import Options
import time

def get_random_user_agent():
    user_agents = [
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
        "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"
    ]
    return random.choice(user_agents)

# Function to check robots.txt
def is_url_allowed_by_robots_txt(url):
    rp = RobotFileParser()
    rp.set_url(url + "/robots.txt")
    rp.read()
    return rp.can_fetch("*", url)

class JobListing2:
    def __init__(self, title, location, company, apply_url, qualifications, description):
        self.title = title
        self.location = location
        self.company = company
        self.apply_url = apply_url
        self.qualifications = qualifications
        self.description = description

class GoogleJobData:
    def __init__(self, domain, url, data, source, timestamp):
        self.domain = domain
        self.url = url
        self.data = data
        self.metadata = {"source": source, "timestamp": timestamp}
class SimplyHiredJobData:
    def __init__(self, domain, url, data, metadata):
        self.domain = domain
        self.url = url
        self.data = data
        self.metadata = metadata
class JobListing23:
    def __init__(self, title, company, location, salary, job_type, job_link, description):
        self.title = title
        self.company = company
        self.location = location
        self.salary = salary
        self.job_type = job_type
        self.job_link = job_link
        self.description = description
# def scrape_google_jobs(urls):
#     job_data = GoogleJobData("Google", "", [], "", datetime.datetime.now().isoformat())
#     headers = {"User-Agent": get_random_user_agent()}
#     job_listings = []
#     for url in urls:
#         if not is_url_allowed_by_robots_txt(url):
#             print(f"Scraping blocked by robots.txt: {url}")
#             continue
#
#         job_data.url = url
#         job_data.metadata["source"] = url
#
#         # Scraping logic for the main job listing page
#         headers = {'User-Agent': get_random_user_agent()}
#         response = requests.get(url, headers=headers)
#         soup = BeautifulSoup(response.text, 'html.parser')
#
#         for div in soup.find_all("div", class_="sMn82b"):
#             title = div.find("h3", class_="QJPWVe").text.strip()
#             location = div.find("div", class_="EAcu5e").text.strip().split("; ")[0]  # Adjust according to location parsing logic
#             apply_url = "https://www.google.com/about/careers/applications/" + div.find("a")['href']
#             qualifications = div.find("div", class_="Xsxa1e").text.strip()
#
#             # Fetching detailed job description
#             detail_response = requests.get(apply_url, headers=headers)
#             detail_soup = BeautifulSoup(detail_response.text, 'html.parser')
#             description = detail_soup.find("div", class_="aG5W3").text.strip()
#
#             job_listings.append(JobListing(title, location, "Google", apply_url, qualifications, description))
#
#     job_data.data = job_listings
#     return job_data

# Function to scrape SimplyHired
class JobListing:
    def __init__(self, url, title, company, location, salary, description, searched_job, searched_location):
        self.url = url
        self.title = title
        self.company = company
        self.location = location
        self.salary = salary
        self.description = description
        self.searched_job = searched_job
        self.searched_location = searched_location

    def to_dict(self):
        return vars(self)

def __get_dom(driver):
    """Gather the data of the current page from web driver and return it."""
    page_content = driver.page_source
    product_soup = BeautifulSoup(page_content, 'html.parser')
    dom = et.HTML(str(product_soup))
    return dom

def __scrape_indeed(driver, query, location):
    try:
        print("Starting to scrape jobs...")
        job_listings = []
        indeed_pagination_url = f"https://www.indeed.com/jobs?q={query}&l={location}"
        job_urls = []

        driver.get(indeed_pagination_url)
        print(f"Fetching URL: {indeed_pagination_url}")

        while True:
            job_page = driver.find_element(By.ID, "mosaic-jobResults")
            jobs = job_page.find_elements(By.CLASS_NAME, "job_seen_beacon")

            if not jobs:
                print("No job cards found on the page.")
                break
            else:
                print(f"Found {len(jobs)} job cards.")

            for job_card in jobs:

                job_url_element = job_card.find_element(By.XPATH, './/a[@data-jk]')
                job_url = job_url_element.get_attribute('href')
                job_urls.append(job_url)

            try:
                next_button = driver.find_element(By.XPATH, '//a[@aria-label="Next"]')
                next_button.click()
                time.sleep(5)
            except NoSuchElementException:
                break

        # Now process each job URL
        for job_url in job_urls:
            driver.get(job_url)
            wait = WebDriverWait(driver, 10)  # Wait for the page to load

            try:
                job_title_element = driver.find_element(By.XPATH, '//h1[contains(@class, "jobsearch-JobInfoHeader-title")]')
                job_title = job_title_element.text
            except NoSuchElementException:
                job_title = "Title Not Found"

            # Add code here to extract other details such as company, location, etc.

            # Example placeholders, replace with actual code to extract company and location
            try:
                company_element = driver.find_element(By.CLASS_NAME, "css-1saizt3.e1wnkr790")
                company = company_element.text
            except NoSuchElementException:
                company = "Company Not Found"

            try:
                location = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, "[class*='css-9yl11a']"))).text
            except TimeoutException:
                location = "Location Not Found"

            try:
                salary_element = driver.find_element(By.CLASS_NAME, "css-2iqe2o.eu4oa1w0")
                salary = salary_element.text
            except NoSuchElementException:
                salary = "Not Available"

            try:
                # Locate the container that holds the job description
                description_container = wait.until(EC.presence_of_element_located((By.ID, "jobDescriptionText")))

                # Find all paragraph and list item elements within the container
                description_elements = description_container.find_elements(By.XPATH, './/p|.//ul/li')

                # Concatenate the text from each element to form the full description
                description_texts = []
                for element in description_elements:
                    # For paragraph elements, simply add their text
                    if element.tag_name == 'p':
                        description_texts.append(element.text)
                    # For list items, you might want to prepend with a bullet or some form of marker
                    elif element.tag_name == 'li':
                        description_texts.append('• ' + element.text)

                # Join all the description texts into one string
                description = ' '.join(description_texts)
            except NoSuchElementException:
                description = "Not Available"

            job_listing = JobListing(
                url=job_url,
                title=job_title,
                company=company,
                location=location,
                salary= salary,
                searched_job=query,
                searched_location=location,
                description = description
            )
            job_listings.append(job_listing.to_dict())

        # Save the job listings to a file
        file_path = 'output/indeed_jobs.json'
        with open(file_path, 'w', encoding='utf-8') as f:
            json.dump(job_listings, f, ensure_ascii=False, indent=4)

        print("Scraping completed.")
    finally:
        driver.quit()

def scrape(job_search_keyword='', location_search_keyword='', glassdoor_start_url='', scrape_option=0) -> None:
    """Scrapes job listings from job board sites.

    Initializes a web driver, creates a directory for output files if it does not exist, and scrapes the desired sites
    based on the provided scrape option. The web driver is closed when the scraping is complete.

    Args:
        job_search_keyword: A job title to be used to search Indeed
        location_search_keyword: A location to search for jobs on Indeed
        glassdoor_start_url: The URL of the first page of a Glassdoor search, used to establish initial page before
            traversing pages
        scrape_option: An integer used to determine which job board sites to scrape

            - default: Scrape both Indeed and Glassdoor
            - 1: Only scrape Indeed
            - 2: Only scrape Glassdoor
    """
    # Initialize webdriver
    options = Options()
    options.add_argument("-headless")
    driver = webdriver.Edge(options=options)


    # Make output folder if one does not exist
    if not os.path.exists('output'):
        os.mkdir('output')

    match scrape_option:
        case 1:
            print('Scraping Indeed...')
            __scrape_indeed(driver, job_search_keyword, location_search_keyword)


    # Close web browser
    driver.quit()
    print('Scraping Complete.')
def main():
    print("Script started.")
    query = 'software+developer'
    location = 'New+York'
    scrape('Healthcare', 'Philadelphia', scrape_option=1)

    # Correct usage of datetime to create a timestamp
    timestamp = datetime.datetime.now().strftime('%Y-%m-%d_%H%M%S')
    filename = f'indeed_jobs_{timestamp}.json'

    # if jobs:
    #     with open(filename, 'w') as file:
    #         json.dump(jobs, file, indent=4)
    #     print(f'Scraped data saved to {filename}')
    # else:
    #     print("No jobs were scraped.")

if __name__ == "__main__":
    main()
