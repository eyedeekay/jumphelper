
GOPATH = $(PWD)/.go

i2pd_dat?=i2pd_dat

GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

time="2s"

basic: deps test noopts

lint: gofmt golint govet

t: lint test build

d: docker docker-run

docker: docker-build

docker-build:
	docker build --no-cache -f Dockerfile -t eyedeekay/jumphelper .

docker-network:
	docker network create --subnet 172.80.80.0/24 si; true

docker-run: docker-network docker-clean
	docker run \
		-d \
		--name jumphelper \
		--network si \
		--network-alias jumphelper \
		--hostname jumphelper \
		--link sam-host \
		--restart always \
		--ip 172.80.80.13 \
		-p 127.0.0.1:7855:7855 \
		--volume forwarded-jumphelper:/opt/work \
		-t eyedeekay/jumphelper; true

docker-clean:
	docker rm -f jumphelper; true

install-client:
	install -m755 bin/ijh /usr/bin

install-server:
	install -m755 bin/jumphelper /usr/bin

install: install-client install-server

clean:
	rm -f bin/*

test:
	cd src && go test

build: server client

noopts:
	go build \
		-o bin/jumphelper \
		./src/server/main.go
	go build \
		-o bin/ijh \
		./src/client/main.go

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
	./bin/ijh -url="forum.i2p/" -port="7854" -addr=true

doecho:
	while true; do make echo; done

curl:
	/usr/bin/curl -l 127.0.0.1:7854/check/i2p-projekt.i2p
	/usr/bin/curl -l 127.0.0.1:7854/i2p-projekt.i2p

deps:
	go get -u github.com/LarryBattle/nonce-golang
	go get -u github.com/bwesterb/go-pow
	go get -u github.com/eyedeekay/jumphelper/src
	go get -u github.com/eyedeekay/sam-forwarder
	go get -u github.com/eyedeekay/gosam
	go get -u golang.org/x/time/rate
	go get -u github.com/eyedeekay/gosam
	go get -u github.com/kpetku/sam3
	go get -u github.com/eyedeekay/i2pasta/convert
	go get -u golang.org/x/time/rate

follow:
	docker logs -f jumphelper

diff:
	cd misc && grep -vf addresses.csv helped.csv > test.csv

start:
	while true; do make docker docker-run follow; done

