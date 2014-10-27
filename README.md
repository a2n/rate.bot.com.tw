# 臺灣銀行黃金價格抓取程式


## 序
此程式以 [go](https://golang.org) 語言撰寫，以日期為參數抓取臺灣銀行[匯率利率黃金牌價查詢](http://rate.bot.com.tw/)。最後將結果寫入 records.csv，欄位格式如下：

```Unix timestamp,銀行買入價格,銀行賣出價格```

因為銀行週末休息，所以不更新紀錄。

## API
```GET http://rate.bot.com.tw/Pages/UIP005/UIP00511.aspx?whom=GB0030001000&afterOrNot=0&curcd=TWD&date=2000101```

只需變換 **date** 參數。
