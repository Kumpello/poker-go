package articles

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var ErrInvalidStatusCode = errors.New("invalid status code (!=200)")

type provider interface {
	provide(what string) (io.ReadCloser, error)
}

// htmlProvider provides a file as "standard" curl request
type htmlProvider struct {
}

func (s htmlProvider) provide(url string) (io.ReadCloser, error) {
	res, err := http.Get(url) // nolint:gosec // G107: url from trusted source
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrInvalidStatusCode
	}

	return res.Body, nil
}

// fileProvider returns the file as io.ReadCloser
// but read the whole file at once
// it's more for test purposes
type fileProvider struct{}

func (f fileProvider) provide(filepath string) (io.ReadCloser, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot read the file")
	}
	return io.NopCloser(strings.NewReader(string(file))), nil
}
