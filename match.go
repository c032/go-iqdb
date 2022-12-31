package iqdb

// Match is a match found by iqdb.
type Match struct {
	URL string

	// Similarity is how similar the match is to the requested image.
	//
	// It's a value between `0` and `1`.
	Similarity float64
}
