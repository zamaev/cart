build:
	docker build -t cart:latest ./

run:
	docker run -p 8082:8082 cart:latest

test-cover:
	go test -cover ./internal/pkg/service/cart
	go test -cover ./internal/pkg/repository

cyclo:
	gocyclo -ignore "_test|mock/" .

cognit:
	gocognit -ignore "_test|mock/" .
