all: clean build

build:
	go build .

clean:
	rm gragas || :

test:
	go test ./...

serve: all
	./gragas
