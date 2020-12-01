package movies

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMovies(t *testing.T) {
	req, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Secret", "1234")
	expected, err := ioutil.ReadFile("test_data/movies_expected_data.json")
	if err != nil {
		t.Fatalf("Error loading file with expected data: %s", err)
	}
	rr := httptest.NewRecorder()
	finalHandler := http.HandlerFunc(GetMovies)
	Middleware(finalHandler).ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("Handler returned unexpected response body: got %v want %v", rr.Body.String(), string(expected))
	}
}
