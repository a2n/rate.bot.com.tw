package main

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"time"
	"runtime"
	"sort"
	"botGoldPrice"
)

func get(date string) string {
	if len(date) == 0 {
		log.Println("get url has error, empty date.")
		return ""
	}
	url := fmt.Sprintf("http://rate.bot.com.tw/Pages/UIP005/UIP00511.aspx?whom=GB0030001000&afterOrNot=0&curcd=TWD&date=%s", date)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Printf("get url has error, %s", err)
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("get url has error, %s", err)
		return ""
	}

	return string(b)
}

func getLocal(date string) string {
	if len(date) == 0 {
		return ""
	}

	path := fmt.Sprintf("../data/%s", date)
	fh, err := os.Open(path)
	defer fh.Close()

	if os.IsNotExist(err) {
		log.Fatalf("getLocal() has error, the file %s is not found.\n", date)
		return ""
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func main() {
	NCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(NCPU)
	ch := make(chan string, NCPU)
	date := time.Date(2014, time.Month(10), 1, 0, 0, 0, 0, time.UTC)
	validDate := 0
	for i := 0; i < 9; i++ {
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			str := fmt.Sprintf("%04d%02d%02d", date.Year(), date.Month(), date.Day())
			go func(date string) {
				ch <- getLocal(str)
			}(str)
			validDate++
		}
		date = date.Add(time.Duration(24) * time.Hour)
	}

	slice := make([]botGoldPrice.Record, 0)
	for i := 0; i < validDate; i++ {
		str := <-ch
		for _, record := range botGoldPrice.NewParser(str).Parse() {
			slice = append(slice, record)
		}
	}
	sort.Sort(botGoldPrice.ByDate(slice))

	// Write
	str := ""
	for _, record := range slice {
		str = fmt.Sprintf("%s%d,%.0f,%.0f\n", str, record.Date.Unix(), record.Buy, record.Sell)
	}

	fh, err := os.Create("records.csv")
	if err != nil {
		log.Fatalf("create recrods.csv has error, %s\n", err)
	}
	_, err = fh.WriteString(str)
	if err != nil {
		log.Fatalf("write file has error, %s", err)
	}
}
