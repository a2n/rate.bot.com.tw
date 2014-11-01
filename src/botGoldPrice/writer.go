package botGoldPrice

import (
	"os"
	"errors"
	"log"
	"fmt"
	"code.google.com/p/go-sqlite/go1/sqlite3"
)

type Writer struct {
	connection *sqlite3.Conn
}

func Connect(path string) (*Writer, error) {
	_, err := os.Open(path)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	conn, err := sqlite3.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return &Writer {
		connection: conn,
	}, nil
}

func (w *Writer) Write(records []Record) (error){
	if len(records) == 0 {
		return errors.New("Empty records to write.")
	}

	sql := "BEGIN TRANSACTION;"
	for _, record := range records {
		sql += fmt.Sprintf("INSERT INTO Price VALUES (%d,%f,%f);", record.Date.Unix(), record.Buy, record.Sell)
	}
	sql += "COMMIT TRANSACTION;"
	err := w.connection.Exec(sql)
	if err != nil {
		return err
	}
	log.Println("Writing recods sucessful.")
	return nil
}
