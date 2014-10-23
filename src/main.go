package main

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	//"code.google.com/p/go.net/html"
)

func get(date string) string {
	if len(date) == 0 {
		log.Println("get url error, empty date.")
		return ""
	}
	url := fmt.Sprintf("http://rate.bot.com.tw/Pages/UIP005/UIP00511.aspx?whom=GB0030001000&date=%s&afterOrNot=0&curcd=TWD", date)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Printf("get url err, %s", err)
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("get url err, %s", err)
		return ""
	}

	return string(b)
}

func testGet() string {
	path := "../data/history.html"
	fh, err := os.Open(path)
	defer fh.Close()

	if os.IsNotExist(err) {
		log.Fatalln("testGet has eror, testing data not found.\n")
		return ""
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

type Record struct {
	Date uint
	Price float32
}

func parse(content string) []Record {
	if len(content) == 0 {
		return nil
	}

	slice := make([]Record, 0)


	return slice
}

func main() {
	str := testGet()
	log.Println(str)
}
