all: clean build

build:
	go build .

clean:
	rm bacchus || :

test:
	go test `glide novendor`

serve: all
	./bacchus

genproto:
	./proto/gen_go.sh

syncproto:
	cd proto && git pull origin master

docker-build:
	docker build -t bacchus .

docker-push:
	docker tag bacchus:latest 096202052535.dkr.ecr.us-west-2.amazonaws.com/bacchus:latest
	docker push 096202052535.dkr.ecr.us-west-2.amazonaws.com/bacchus:latest
