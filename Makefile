all: clean build

build:
	go build .

clean:
	rm bacchus || :

test:
	go test ./...

serve: all
	./bacchus
