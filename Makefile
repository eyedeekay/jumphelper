
GO_COMPILER_OPTS = -a -tags netgo -ldflags '-w -extldflags "-static"'

time=".3s"

lint: gofmt golint govet

t: lint test build

d: docker docker-run

docker:
	docker build -f Dockerfile -t eyedeekay/jumphelper .

docker-run: docker-clean
	docker run -i -t -d \
		--net host \
		-p 127.0.0.1:7054:7054 \
		--name jumphelper \
		eyedeekay/jumphelper

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
	./bin/ijh -url="i2p-projekt.i2p/" -port="7054" -addr=true
	sleep "$(time)"

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
	docker logs -f jumphelper

diff:
	bash -c "diff -d <(sort -u alive-hosts.txt | sed 's|=.*||g') <(sort -u <(sort -u addresses.csv | sed 's|,.*||g') <(sort -u alive-hosts.txt | sed 's|=.*||g'))" 1> candidates.diff

start:
	while true; do make docker docker-run follow; done

stop:
	docker rm -f jumphelper; true
