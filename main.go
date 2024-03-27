package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// define a data structure to store the scraped data
type ScrapedItem struct {
	url, image, name, price string
}

// it varifies a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func main() {
	// initialize the slice of structs that will hold the scraped data
	var scrapedItems []ScrapedItem

	// initialize the list of pages to scrape with an empty slice
	var pagesToScrape []string

	//  the first pagination URL to scrape
	pageToScrape := "https://scrapeme.live/shop/page/1/"

	// initialize the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// current iteration, starts from first page
	currentIteration := 1

	// maximum number of iterations
	maxIteration := 5

	// initialize the colly instance
	c := colly.NewCollector()

	// setting a valid user agent header
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

	// iterate over the list of pagination links to implement the crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		// if the page discovered is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	// scraping the data
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		scrapedItem := ScrapedItem{}
		scrapedItem.url = e.ChildAttr("a", "href")
		scrapedItem.image = e.ChildAttr("img", "src")
		scrapedItem.name = e.ChildText("h2")
		scrapedItem.price = e.ChildText(".price")

		scrapedItems = append(scrapedItems, scrapedItem)
	})

	c.OnScraped(func(response *colly.Response) {
		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && currentIteration < maxIteration {
			// get the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// increment the current iteration
			currentIteration++

			// visiting a new page
			c.Visit(pageToScrape)
		}
	})

	// visiting the first page
	c.Visit(pageToScrape)

	// opening a CSV file
	file, err := os.Create("scraped.csv")
	if err != nil {
		log.Fatalln("failed to create csv file", err)
	}
	defer file.Close()

	// initializing the CSV writer
	writer := csv.NewWriter(file)

	// defining the CSV	headers
	headers := []string{"url", "image", "name", "price"}

	// writing the header columns
	writer.Write(headers)

	// adding each items to the ".csv" file
	for _, item := range scrapedItems {
		// converting the scraped item into array of strings
		record := []string{
			item.url,
			item.image,
			item.name,
			item.price,
		}

		// writing a new CSV record
		writer.Write(record)
	}
	defer writer.Flush()
}
