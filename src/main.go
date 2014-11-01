package main

import (
	"log"
	"time"
	"botGoldPrice"
)

func timer() {
	ch := time.Tick(1 * time.Hour)
	crawler := botGoldPrice.NewCrawler()
	records := make([]botGoldPrice, )
	for {
		<-ch
		htmls := crawler.GetOneDay(time.Now())
		records = append(records, botGoldPrice.NewParser(html).Parse()...)
	}
}

func main() {
	begin := time.Date(2000, time.Month(10), 30, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	crawler := botGoldPrice.NewCrawler()
	htmls := crawler.GetDateRange(begin, end)
	records := make([]botGoldPrice.Record, 0)
	for _, html := range htmls {
		records = append(records, botGoldPrice.NewParser(html).Parse()...)
	}
	for _, record := range records {
		log.Printf("%04d-%02d-%02d, buy: %f, sell: %f\n", record.Date.Year(), record.Date.Month(), record.Date.Day(), record.Buy, record.Sell)
	}
}
