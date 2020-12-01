package movies

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeRequest(method string, headers map[string]string, t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, "/movies", nil)
	if err != nil {
		t.Fatalf("Error during request: %s", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rr := httptest.NewRecorder()
	finalHandler := http.HandlerFunc(GetMovies)
	Middleware(finalHandler).ServeHTTP(rr, req)
	return rr
}

func TestGetMovies(t *testing.T) {
	rr := makeRequest("GET", map[string]string{"X-Secret": "1234"}, t)
	expected, err := ioutil.ReadFile("test_data/movies_expected_data.json")
	if err != nil {
		t.Fatalf("Error loading file with expected data: %s", err)
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if strings.TrimSpace(rr.Body.String()) != string(expected) {
		t.Errorf("Handler returned unexpected response body: got %v want %v", rr.Body.String(), string(expected))
	}
}

func TestMiddleware(t *testing.T) {
	rr := makeRequest("GET", map[string]string{}, t)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Header missing: got %v want %v", status, http.StatusForbidden)
	}
}

func TestGetMoviesWrongMethod(t *testing.T) {
	rr := makeRequest("POST", map[string]string{"X-Secret": "1234"}, t)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Wrong method: got %v want %v", status, http.StatusBadRequest)
	}
}
