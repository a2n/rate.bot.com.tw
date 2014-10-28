package botGoldPrice

import (
	"testing"
	"io/ioutil"
)

func TestParser(t *testing.T) {
	path := "../../data/emptyRecords"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	slice := NewParser(string(b)).Parse()
	if len(slice) == 0 {
		t.Error("Empty record... failed.\n")
	} else {
		t.Log("Empty record... passed.\n")
	}
}
