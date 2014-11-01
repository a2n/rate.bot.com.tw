package botGoldPrice

import (
	"fmt"
	"strings"
	"time"
	"log"
	"errors"
	"code.google.com/p/go.net/html"
)

type Record struct {
	Date time.Time
	Buy float32
	Sell float32
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

func (p *Parser) Parse() []Record {
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
							record, err := p.getRecord(n, date)
							if err == nil {
								slice = append(slice, record)
							}
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

func (p *Parser) getRecord(node *html.Node, date time.Time) (Record, error) {
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
		} else {
			return Record{}, err
		}
	}

	// Buy price
	if newNode.NextSibling == nil {
		return Record{}, errors.New("Empty records.")
	}

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
	}, nil
}
