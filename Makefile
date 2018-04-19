build:
	glide update
	cp -r ./vendor/* ${GOPATH}/src/
	go build -o proxy

run: build
	./proxy -upstreamUrl http://testurl -secret mysecret #-allowPath /github-webhook -allowPath /project