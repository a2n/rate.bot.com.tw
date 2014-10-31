package main

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"runtime"
	"time"
	"math"
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

func getLocal(fileName string) string {
	if len(fileName) == 0 {
		return ""
	}

	path := fmt.Sprintf("../data/%s", fileName)
	fh, err := os.Open(path)
	defer fh.Close()

	if os.IsNotExist(err) {
		log.Fatalf("getLocal() has error, %s\n", fileName, err)
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
	ch := make(chan string, math.MaxUint16)

	date := time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	max := math.MaxUint16
	for i := 0; i < max; i++ {
		str := fmt.Sprintf("%04d%02d%02d", date.Year(), date.Month(), date.Day())
		log.Printf("Getting %s\n", str)
		ch <- get(str)
		date = date.Add(time.Duration(24) * time.Hour)
		if date.After(time.Now()) {
		    close(ch)
		    break
		}
		//ch <- getLocal("validRecords")
		//ch <- getLocal("emptyRecords")
	}

	slice := make([]botGoldPrice.Record, 0)
	for i := 0; i < max; i++ {
		str := <-ch
		if len(str) == 0 {
		    break
		}
		records := botGoldPrice.NewParser(str).Parse()
		slice = append(slice, records...)
	}
	log.Printf("count: %d\n", len(slice))

	// Writing
	str := ""
	for _, record := range slice {
		str = fmt.Sprintf("%s%d,%.0f,%.0f\n", str, record.Date.Unix(), record.Buy, record.Sell)
	}

	file, err := os.Create("records.csv")
	if err != nil {
		log.Fatal(err)
	}
	n, err := file.Write([]byte(str))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Writing %d bytes.", n)

	file.Close()
}
