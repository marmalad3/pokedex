generate-swagger: 
	./scripts/generate-swagger.sh

test: lint
	go list ./... | grep -v vendor | xargs go test -count 1 -cover -race -v

fmt:
	gofmt -s -w .

lint:
	go list ./... | grep -v vendor | xargs go vet -v

