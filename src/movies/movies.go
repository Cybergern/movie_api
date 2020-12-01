package movies

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Movie struct {
	Title   string `json:"Title"`
	Year    string `json:"Year"`
	Runtime string `json:"Runtime"`
	Genre   string `json:"Genre"`
	ImdbId  string `json:"imdbID"`
}

var storage *Storage

const SECRET = "1234"
const API_KEY = "xxxxxx"

func init() {
	storage = NewStorage()
}

func HeaderMethodCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("X-Secret")
		if secret != SECRET {
			http.Error(w, "", http.StatusForbidden)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Cache code inspired by https://goenning.net/2017/03/18/server-side-cache-go/

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if content := storage.Get(r.RequestURI); content != nil {
		if _, err := w.Write(content); err != nil {
			log.Fatalln(err)
		}
	} else {
		ids := []string{"tt0116282", "tt0209144", "tt0133093", "tt0118715"}
		var movies = getOmdbInformation(ids)
		sort.Slice(movies, func(i, j int) bool {
			return clean(movies[i].Title) < clean(movies[j].Title)
		})
		if content, err := json.Marshal(movies); err != nil {
			log.Printf("Could not encode response. err:%s\n", err)
		} else {
			if _, err := w.Write(content); err != nil {
				log.Fatalln(err)
			}
			addToCache("60s", r, content)
		}
	}
}

func addToCache(duration string, r *http.Request, content []byte) {
	if d, err := time.ParseDuration(duration); err == nil {
		log.Printf("New page cached: %s for %s\n", r.RequestURI, duration)
		storage.Set(r.RequestURI, content, d)
	} else {
		log.Printf("Page not cached, could not parse cache duration. err: %s\n", err)
	}
}

func getOmdbInformation(imdbIds []string) []Movie {
	var result []Movie
	wg := sync.WaitGroup{}
	for _, id := range imdbIds {
		wg.Add(1)
		go func(id string) {
			resp, err := http.Get(fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&i=%s", API_KEY, id))
			if err != nil {
				log.Fatalln(err)
				return
			}
			var movie Movie

			if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
				log.Fatalln(err)
				return
			}
			result = append(result, movie)
			wg.Done()
		}(id)
	}
	wg.Wait()
	return result
}

func clean(title string) string {
	return strings.TrimPrefix(title, "The ")
}
