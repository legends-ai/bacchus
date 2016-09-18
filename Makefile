all: clean build

build:
	go build .

clean:
	rm bacchus || :

test:
	go test `glide novendor`

serve: all
	./bacchus

docker-build:
	docker build -t simplyianm/bacchus .
