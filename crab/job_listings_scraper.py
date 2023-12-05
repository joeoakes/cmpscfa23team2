import os
from multiprocessing import Pool
from urllib.robotparser import RobotFileParser
import et as et
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import json
import random
import datetime
import re
from bs4 import BeautifulSoup
from lxml import etree as et
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options as ChromeOptions
from selenium.webdriver.firefox.options import Options as FirefoxOptions
from selenium.webdriver.edge.options import Options as EdgeOptions



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
    return et.HTML(page_content)



def __scrape_indeed_page(job_url):


    # Open a new driver instance for each page
    driver = get_web_driver()
    try:
        driver.get(job_url)
        wait = WebDriverWait(driver, 5)  # Reduced wait time

        job_title = __extract_element_text(wait, By.XPATH, '//h1[contains(@class, "jobsearch-JobInfoHeader-title")]', "Title Not Found")
        company = __extract_element_text(wait, By.CLASS_NAME, "css-1saizt3.e1wnkr790", "Company Not Found")
        # First try to extract location with the first selector
        location = __extract_element_text(wait, By.CLASS_NAME, "css-9yl11a.eu4oa1w0", "Location Not Found")

        # If not found, try the second selector
        if location == "Location Not Found":
            location = __extract_element_text(wait, By.CLASS_NAME, "css-ks9svk.eu4oa1w0", "Location Not Found")
        salary = __extract_element_text(wait, By.CLASS_NAME, "css-2iqe2o.eu4oa1w0", "Salary Not Found")
        description = __extract_job_description(wait)

        return {
            "url": job_url,
            "title": job_title,
            "company": company,
            "location": location,
            "salary": salary,
            "description": description
        }
    finally:
        driver.quit()


def __scrape_domain(driver, domain, location, num_pages=1):
    print(f"Starting to scrape jobs for domain: {domain}")
    indeed_pagination_url = f"https://www.indeed.com/jobs?q={domain}&l={location}"
    job_urls = []

    for page in range(num_pages):
        driver.get(indeed_pagination_url + f"&start={page * 10}")
        print(f"Fetching URL: {indeed_pagination_url}")

        wait = WebDriverWait(driver, 10)
        jobs = wait.until(EC.presence_of_all_elements_located((By.CLASS_NAME, "job_seen_beacon")))

        if not jobs:
            print("No job cards found on the page.")
            break

        print(f"Found {len(jobs)} job cards on page {page + 1}.")

        for job_card in jobs:
            job_url_element = job_card.find_element(By.XPATH, './/a[@data-jk]')
            job_urls.append(job_url_element.get_attribute('href'))

    with Pool(60) as p:
        job_listings = p.map(__scrape_indeed_page, job_urls)

    scraped_data = {
        "domain": domain,
        "url": indeed_pagination_url,
        "data": job_listings,
        "metadata": {
            "source": indeed_pagination_url,
            "timestamp": datetime.datetime.now().isoformat()
        }
    }

    return scraped_data



def __extract_element_text(wait, by, locator, default_text):
    try:
        element = wait.until(EC.presence_of_element_located((by, locator)))
        return element.text
    except Exception:
        return default_text



def __extract_job_description(wait):
    try:
        description_container = wait.until(EC.presence_of_element_located((By.ID, "jobDescriptionText")))
        html_content = description_container.get_attribute('outerHTML')
        soup = BeautifulSoup(html_content, 'html.parser')

        description_texts = []
        capture = False

        # Define the keywords and bullet point indicators
        keywords = [r'qualifications', r'experience', r'requirements', r'responsibilities',
                    r'desired skills', r'minimum qualifications', r'recruitment requirements', r'experience',
                    r'skills', r'Recruitment Requirements', r'required experience',
                    r'position requirements', r'Skills and Abilities', r'Abilities',
                    r'Qualifications', r'RECRUITMENT REQUIREMENTS', r'Certification']
        bullet_indicators = ['-', '•', '*']  # Add more indicators if needed

        for element in soup.find_all(['p', 'li', 'ul', 'b']):
            text = element.get_text(strip=True)
            lower_text = text.lower()

            # Check if the element is a keyword in any format
            if any(re.search(keyword, lower_text, re.IGNORECASE) for keyword in keywords):
                capture = True
                continue  # Skip the keyword itself

            if capture:
                if element.name == 'li' or (element.name == 'p' and any(text.lstrip().startswith(indicator) for indicator in bullet_indicators)):
                    description_texts.append(text)
                elif element.name == 'ul':
                    # Capture all list items within the ul
                    for li in element.find_all('li'):
                        description_texts.append(li.get_text(strip=True))
                elif element.name == 'p' and not any(text.lstrip().startswith(indicator) for indicator in bullet_indicators):
                    capture = False  # Stop capturing when encountering a non-bullet point paragraph

        return '\n'.join(description_texts).strip()
    except Exception as e:
        print(f"Error in extracting job description: {e}")
        return "Not Available"



def get_web_driver():
    # Try initializing Firefox WebDriver
    try:
        firefox_options = FirefoxOptions()
        firefox_options.add_argument("-headless")
        return webdriver.Firefox(options=firefox_options)
    except Exception as e:
        print("Firefox WebDriver not found. Trying Edge...", e)

    # Try initializing Edge WebDriver
    try:
        edge_options = EdgeOptions()
        edge_options.add_argument("-headless")
        return webdriver.Edge(options=edge_options)
    except Exception as e:
        print("Edge WebDriver not found. Trying Chrome...", e)

    # Try initializing Chrome WebDriver
    try:
        chrome_options = ChromeOptions()
        chrome_options.add_argument("-headless")
        return webdriver.Chrome(options=chrome_options)
    except Exception as e:
        print("Chrome WebDriver not found.", e)
        raise

def scrape(location_search_keyword='', num_pages=1, scrape_option=0) -> None:
    driver = get_web_driver()

    domains = ['Law', 'Business', 'Software Engineer']

    if not os.path.exists('output'):
        os.makedirs('output')

    for domain in domains:
        domain_data = __scrape_domain(driver, domain, location_search_keyword, num_pages)

        # Create a separate file for each domain
        file_path = os.path.join('output', f'{domain}_jobs.json')
        with open(file_path, 'w', encoding='utf-8') as f:
            json.dump(domain_data, f, ensure_ascii=False, indent=4)
        print(f'Scraped data for {domain} domain saved to {file_path}.')

    driver.quit()

    print('Scraping Complete.')



def main():
    print("Script started.")
    scrape('Pennsylvania', num_pages=3, scrape_option=1)
          
if __name__ == "__main__":
    main()
