all: clean build

build:
	go build .

clean:
	rm -f bacchus

test:
	go test `glide novendor`

serve: all
	./bacchus

genproto:
	./proto/gen_go.sh

init:
	git submodule update --init

install:
	glide install

syncproto:
	cd proto && git pull origin master

docker-build:
	docker build -t bacchus .

docker-push:
	docker tag bacchus:latest 096202052535.dkr.ecr.us-west-2.amazonaws.com/bacchus:latest
	docker push 096202052535.dkr.ecr.us-west-2.amazonaws.com/bacchus:latest
