all: clean build

build:
	go build .

clean:
	rm gragas || :
