package prerender

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_BotRequest(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://www.google.com/", nil)
	req.Header.Set("User-Agent", "twitterbot")

	NewOptions().NewPrerender().ServeHTTP(res, req, nil)

	if len(res.Body.Bytes()) == 0 {
		t.Error("Error, prerender.io not called")
	}
}

func Test_NonBotRequest(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://www.google.com/", nil)

	NewOptions().NewPrerender().ServeHTTP(res, req, nil)
	if len(res.Body.Bytes()) > 0 {
		t.Error("Error, prerender.io called for non-proxy request")
	}
}
