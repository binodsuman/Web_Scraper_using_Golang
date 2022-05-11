package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Book struct {
	Title string
	Price string
}

func main() {
	start := time.Now()
	file, err := os.Create("books_2.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headers := []string{"Title", "Price"}
	writer.Write(headers)

	c := colly.NewCollector(
		colly.AllowedDomains("books.toscrape.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	//Pagination
	c.OnHTML(".next > a", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	//Each product details
	c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
		book := Book{}
		book.Title = e.ChildAttr(".image_container img", "alt")
		book.Price = e.ChildText(".price_color")
		row := []string{book.Title, book.Price}
		writer.Write(row)
	})

	startUrl := fmt.Sprintf("https://books.toscrape.com/")
	c.Visit(startUrl)
	fmt.Println("Time taken : ", time.Since(start))
}
