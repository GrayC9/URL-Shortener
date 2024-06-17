.PHONY:

build:
	go build -o ./url cmd/url-shortener/main.go

run:
	./url