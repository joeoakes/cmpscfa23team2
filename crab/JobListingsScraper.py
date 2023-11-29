import os
import json
import random
import datetime
from urllib.robotparser import RobotFileParser

from selenium import webdriver
from selenium.webdriver.chrome.options import Options as ChromeOptions
from selenium.webdriver.firefox.options import Options as FirefoxOptions
from selenium.webdriver.edge.options import Options as EdgeOptions
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.edge.service import Service as EdgeService
from selenium.webdriver.firefox.service import Service as FirefoxService
from webdriver_manager.chrome import ChromeDriverManager
from webdriver_manager.firefox import GeckoDriverManager
from webdriver_manager.microsoft import EdgeChromiumDriverManager
from bs4 import BeautifulSoup
from lxml import etree as et

from selenium.common.exceptions import NoSuchElementException, TimeoutException
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
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
    def __init__(self, url, title, company, location, salary, description):
        self.url = url
        self.title = title
        self.company = company
        self.location = location
        self.salary = salary
        self.description = description

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
        page_counter = 0
        while True:
            wait = WebDriverWait(driver, 10)
            jobs = wait.until(EC.presence_of_all_elements_located((By.CLASS_NAME, "job_seen_beacon")))

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
                next_button = driver.find_element(By.XPATH, '//a[@data-testid="pagination-page-next"]')
                next_button.click()
                page_counter += 1
                if page_counter >= 3:
                    break
            except Exception:
                break

        for job_url in job_urls:
            driver.get(job_url)
            wait = WebDriverWait(driver, 3)

            job_title = __extract_element_text(wait, By.XPATH,
                                               '//h1[contains(@class, "jobsearch-JobInfoHeader-title")]',
                                               "Title Not Found")
            company = __extract_element_text(wait, By.CLASS_NAME, "css-1saizt3.e1wnkr790", "Company Not Found")
            location = __extract_element_text(wait, By.CLASS_NAME, "css-9yl11a.eu4oa1w0", "Location Not Found")
            salary = __extract_element_text(wait, By.CLASS_NAME, "css-2iqe2o.eu4oa1w0", "Salary Not Available")
            description = __extract_job_description(wait)

            job_listing = JobListing(
                url=job_url,
                title=job_title,
                company=company,
                location=location,
                salary=salary,
                description=description
            )
            job_listings.append(job_listing.to_dict())

        scraped_data = {
            "domain": query,
            "url": indeed_pagination_url,
            "data": job_listings,
            "metadata": {
                "source": indeed_pagination_url,
                "timestamp": datetime.datetime.now().isoformat()
            }
        }

        if not os.path.exists('output'):
            os.makedirs('output')
        file_path = os.path.join('output', 'indeed_jobs.json')
        with open(file_path, 'w', encoding='utf-8') as f:
            json.dump(scraped_data, f, ensure_ascii=False, indent=4)

        print("Scraping completed.")
    finally:
        driver.quit()


def __extract_element_text(wait, by, locator, default_text):
    try:
        element = wait.until(EC.presence_of_element_located((by, locator)))
        return element.text
    except Exception:
        return default_text


def __extract_job_description(wait):
    try:
        description_container = wait.until(EC.presence_of_element_located((By.ID, "jobDescriptionText")))
        description_elements = description_container.find_elements(By.XPATH, './/p|.//ul/li')
        description_texts = ['• ' + element.text if element.tag_name == 'li' else element.text
                             for element in description_elements]
        return ' '.join(description_texts)
    except Exception:
        return "Not Available"


def detect_browser_driver():
    """
    Detects the installed browser on the system and returns the corresponding Selenium WebDriver in headless mode.
    """
    user_agent = get_random_user_agent()


    try:
        firefox_options = FirefoxOptions()
        firefox_options.add_argument("--headless")
        firefox_options.set_preference("general.useragent.override", user_agent)
        return webdriver.Firefox(service=FirefoxService(GeckoDriverManager().install()), options=firefox_options)
    except:
        pass
    try:
        edge_options = EdgeOptions()
        edge_options.add_argument("--headless")
        edge_options.add_argument(f'user-agent={user_agent}')
        edge_options.add_argument('window-size=1920x1080')
        return webdriver.Edge(service=EdgeService(EdgeChromiumDriverManager().install()), options=edge_options)
    except:
        pass
    try:
        chrome_options = ChromeOptions()
        chrome_options.add_argument("--headless")
        chrome_options.add_argument(f'user-agent={user_agent}')
        chrome_options.add_argument('window-size=1920x1080')
        return webdriver.Chrome(service=ChromeService(ChromeDriverManager().install()), options=chrome_options)
    except:
        raise Exception("No supported browser found on the system.")


def scrape(job_search_keyword='', location_search_keyword='', scrape_option=0) -> None:
    driver = detect_browser_driver()

    if not os.path.exists('output'):
        os.mkdir('output')

    match scrape_option:
        case 1:
            print('Scraping Indeed...')
            __scrape_indeed(driver, job_search_keyword, location_search_keyword)

    print('Scraping Complete.')


def main():
    print("Script started.")
    scrape('cybersecurity', 'Philadelphia', scrape_option=1)


if __name__ == "__main__":
    main()
