To run code directly, set GOPATH to project root and run `go run src/api-test/main.go` in the project root

To run as a docker image, build with `docker build -t go-docker-prod .` and run with `docker run --rm -it -p 8080:8080 go-docker-prod`

Access api at http://localhost:8080/movies

To run the tests, run "go test movies" in the project root.