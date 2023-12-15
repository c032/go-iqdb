package iqdb_test

import (
	"testing"

	"github.com/c032/go-iqdb"
)

func TestSearchURL(t *testing.T) {
	const imageURL = "https://files.yande.re/image/a17776e207fd81d3bb946ba530f17e3c/yande.re%20994641%20amiya_%28arknights%29%20animal_ears%20arknights%20bunny_ears%20jazztaki%20selfie.png"

	matches, err := iqdb.SearchURL(imageURL)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(matches); got == 0 {
		t.Fatalf("len(matches) = %d; want > 0", got)
	}

	found := false
	for _, match := range matches {
		if match.URL == "https://yande.re/post/show/994641" && match.Similarity >= 0.9 {
			found = true

			break
		}
	}
	if !found {
		t.Fatal("missing expected match")
	}
}
