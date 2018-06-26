
GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

time=".3s"

lint: gofmt golint govet

t: lint test build

d: docker docker-run

docker:
	docker build -f Dockerfile -t eyedeekay/jumphelper .

docker-run:
	docker run -i -t -d \
		-p 127.0.0.1:7054:7054 \
		--name jumphelper \
		eyedeekay/jumphelper

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
	./bin/ijh -url="http://i2p-projekt.i2p/en/" -addr=true
	sleep "$(time)"

doecho:
	while true; do make echo; done

curl:
	/usr/bin/curl -l 127.0.0.1:7054/check/i2p-projekt.i2p
	/usr/bin/curl -l 127.0.0.1:7054/i2p-projekt.i2p

deps:
	go get -u github.com/eyedeekay/jumphelper/src
	go get -u golang.org/x/time/rate
