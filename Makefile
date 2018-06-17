
GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

t: gofmt golint govet
	sudo -u i2pd make test
	make build

d: docker docker-run

docker:
	docker build -f Dockerfile -t eyedeekay/jumphelper .

docker-run:
	docker run -i -t -d \
		-p 127.0.0.1:7054:7054 \
		-n jumphelper \
		eyedeekay/jumphelper

test:
	cd src && go test

build:
	GOOS=linux GOARCH=amd64 go build \
		$(GO_COMPILER_OPTS) \
		-o bin/jumphelper \
		./src/server/main.go
	@echo 'built.'

gofmt:
	cd src && gofmt -w *.go

golint:
	cd src && golint *.go

govet:
	cd src && go vet *.go

run:
	sudo -u i2pd ./bin/jumphelper

echo:
	./bin/ijh

curl:
	/usr/bin/curl 127.0.0.1:7054/check/i2p-projekt.i2p
	/usr/bin/curl -l 127.0.0.1:7054/i2p-projekt.i2p
