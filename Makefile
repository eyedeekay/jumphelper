
GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

time="2s"

lint: gofmt golint govet

t: lint test build

d: docker docker-run

docker: docker-host
	docker build -f Dockerfile -t eyedeekay/sam-jumphelper .

docker-network:
	docker network create --subnet 172.81.81.0/29 jumphelper; true

docker-host: docker-network
	docker run \
		-d \
		--name jumphelper-sam-host \
		--network jumphelper \
		--network-alias jumphelper-sam-host \
		--hostname jumphelper-sam-host \
		--link sam-jumphelper \
		--restart always \
		--ip 172.81.81.2 \
		-p :4567 \
		-p 127.0.0.1:7073:7073 \
		--volume jumphelper-sam-host:/var/lib/i2pd:rw \
		-t eyedeekay/sam-host; true

docker-run: docker-network
	docker rm -f sam-jumphelper; true
	docker run \
		-d \
		--name sam-jumphelper \
		--network jumphelper \
		--network-alias sam-jumphelper \
		--hostname sam-jumphelper \
		--link jumphelper-sam-host \
		--restart always \
		--ip 172.81.81.3 \
		-p 127.0.0.1:7054:7054 \
		-t eyedeekay/sam-jumphelper

docker-clean:
	docker rm -f jumphelper; true

install:
	install -m755 bin/ijh /usr/bin
	install -m755 bin/jumphelper /usr/bin

clean:
	rm -f bin/*

test:
	cd src && go test

build: server client

server:
	GOOS=linux GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o bin/jumphelper \
		./src/server/main.go
	@echo 'built.'

client:
	GOOS=linux GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o bin/ijh \
		./src/client/main.go
	@echo 'built.'

gofmt:
	cd src && gofmt -w *.go */*.go

golint:
	cd src && golint *.go

govet:
	cd src && go vet *.go

run:
	./bin/jumphelper

echo:
	@sleep "$(time)" && echo
	ijh -url="i2p-projekt.i2p" -addr=true
	@sleep "$(time)" && echo
	ijh -url="i2p-projekt.i2p"
	@sleep "$(time)" && echo
	ijh -url="fireaxe.i2p" -addr=true
	@sleep "$(time)" && echo
	ijh -url="fireaxe.i2p"

forum:
	./bin/ijh -url="forum.i2p/" -port="7054" -addr=true

doecho:
	while true; do make echo; done

curl:
	/usr/bin/curl -l 127.0.0.1:7054/check/i2p-projekt.i2p
	/usr/bin/curl -l 127.0.0.1:7054/i2p-projekt.i2p

deps:
	go get -u github.com/eyedeekay/jumphelper/src
	go get -u golang.org/x/time/rate

follow:
	docker logs -f sam-jumphelper

diff:
	cd misc && grep -vf addresses.csv helped.csv > test.csv

start:
	while true; do make docker docker-run follow; done

stop:
	docker rm -f sam-jumphelper; true
