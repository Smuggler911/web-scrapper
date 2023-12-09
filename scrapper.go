package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
)

//тип продукта

type Products struct {
	url, image, name, price string
}

// функция проверяющая содержит ли лист строку/строки  или нет
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func main() {

	var products []Products

	///лист
	scrapePage := []string{"https://www.gatorade.com/holiday"}

	i := 1
	limit := 5

	//обозначение нового коллеуктора
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 "

	///проверка на содержание
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {

		newPaginationLink := e.Attr("href")

		if !contains(scrapePage, newPaginationLink) {
			scrapePage = append(scrapePage, newPaginationLink)
		}
	})

	///назначение отребутов по которым производится парсинг
	c.OnHTML("a.product", func(e *colly.HTMLElement) {
		prod := Products{}
		prod.url = e.ChildAttr("a", "href")
		prod.image = e.ChildAttr("img", "src")
		prod.name = e.ChildText("h2")
		prod.price = e.ChildText(".Price_group__XQJOF")

		products = append(products, prod)

	})
	//роутинг по страницам
	c.OnScraped(func(response *colly.Response) {
		if len(scrapePage) != 0 && i < limit {
			scrapePage := scrapePage[0]
			scrapePage = scrapePage[1:]
			i++
			c.Visit(scrapePage)
		}
	})
	///создание файла с парсинг данными
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	wr := csv.NewWriter(file)
	for _, product := range products {

		record := []string{
			product.url,
			product.image,
			product.name,
			product.price,
		}

		wr.Write(record)
	}
	defer wr.Flush()
}
