package main

import (
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"runtime"
	"sort"
	"code.google.com/p/go.net/html"
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

type Record struct {
	Date time.Time
	Buy float32
	Sell float32
}

type ByDate []Record

func (d ByDate) Len() int {
	return len(d)
}

func (d ByDate) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d ByDate) Less(i, j int) bool {
	return d[i].Date.Before(d[j].Date)
}

type Parser struct {
	Content string
}

func NewParser(content string) *Parser {
	if len(content) == 0 {
		return nil
	}

	return &Parser {
		Content: content,
	}
}

func (p *Parser) parse() []Record {
	doc, err := html.Parse(strings.NewReader(p.Content))
	if err != nil {
		return nil
	}

	slice := make([]Record, 0)
	var date time.Time
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode {
			// Date
			if n.Data == "th" && n.FirstChild != nil{
				if n.FirstChild.Data == "資料日期" {
					content := n.NextSibling.NextSibling.FirstChild.Data
					date = p.getDate(content)
				}
			}

			// Records
			if n.Data == "tr" {
				for _, attr := range n.Attr {
					if attr.Key == "class" {
						if attr.Val == "color0" || attr.Val == "color1" {
						record := p.getRecord(n, date)
						slice = append(slice, record)
						}
					}
				}
			}
		}

        for node := n.FirstChild; node != nil; node = node.NextSibling {
            f(node)
        }
    }
    f(doc)

	return slice
}

func (p *Parser) getDate(content string) time.Time {
	if len(content) == 0 {
		return time.Date(0, time.Month(0), 0, 0, 0, 0, 0, time.UTC)
	}

	var year, month, day int
	_, err := fmt.Sscanf(content, "%d/%d/%d", &year, &month, &day)
	if err != nil {
		log.Printf("Parser.getDate has error, sscanf error with %s\n", content)
		return time.Date(0, time.Month(0), 0, 0, 0, 0, 0, time.UTC)
	}

	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Printf("Parser.getDate has error, %s\n", err)
		return time.Date(0, time.Month(0), 0, 0, 0, 0, 0, time.UTC)
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
}

func (p *Parser) getRecord(node *html.Node, date time.Time) Record {
	// Time
	newNode := node.FirstChild
	if newNode != nil {
		var hour, minute int
		str := newNode.FirstChild.Data
		_, err := fmt.Sscanf(str, "%02d:%02d", &hour, &minute)
		if err == nil {
			duration := time.Duration(hour) * time.Hour
			duration += time.Duration(minute) * time.Minute
			date = date.Add(duration)
		}
	}

	// Buy price
	newNode = newNode.NextSibling.NextSibling.NextSibling
	var buy float32
	if newNode != nil {
		if newNode.FirstChild.Data != "" {
			str := newNode.FirstChild.Data
			var price float32
			_, err := fmt.Sscanf(str, "%f", &price)
			if err == nil {
				buy = price
			}
		}
	}

	// Sell price
	newNode = newNode.NextSibling
	var sell float32
	if newNode != nil {
		if newNode.FirstChild.Data != "" {
			str := newNode.FirstChild.Data
			var price float32
			_, err := fmt.Sscanf(str, "%f", &price)
			if err == nil {
				sell = price
			}
		}
	}
	return Record {
		Date: date,
		Buy: buy,
		Sell: sell,
	}
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

	slice := make([]Record, 0)
	for i := 0; i < validDate; i++ {
		str := <-ch
		for _, record := range NewParser(str).parse() {
			slice = append(slice, record)
		}
	}
	sort.Sort(ByDate(slice))

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
