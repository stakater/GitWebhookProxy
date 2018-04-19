build:
	glide update
	sudo cp -r ./vendor/* ${GOPATH}/src/
	go build -o stk