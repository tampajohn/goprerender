package prerender

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_NewPrerender(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://www.stihlusa.com/", nil)

	NewPrerender().ServeHTTP(res, req, nil)
	fmt.Println(len(res.Body.Bytes()))
}

func Test_NewPrerender_WithGzip(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://www.stihlusa.com/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	NewPrerender().ServeHTTP(res, req, nil)
	fmt.Println(len(res.Body.Bytes()))
}

func Test_NewPrerender_NoGzip_NoGzipContent(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://www.stihlusa.com/", nil)

	NewPrerender().ServeHTTP(res, req, nil)
	fmt.Println(len(res.Body.Bytes()))
}

func Test_NewPrerender_Gzip_NoGzipContent(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://www.stihlusa.com/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	NewPrerender().ServeHTTP(res, req, nil)
	fmt.Println(len(res.Body.Bytes()))
}
