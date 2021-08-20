generate-swagger: 
	./scripts/generate-swagger.sh

test: lint
	go list ./... | grep -v vendor | xargs go test -count 1 -cover -race -v

format:
	go list ./... | grep -v vendor | xargs goimports -w

lint:
	go list ./... | grep -v vendor | xargs go vet -v

clean:
	rm ./pokedex || true

build: clean
	go build ./cmd/pokedex

run: build
	./pokedex
	
docker-build:
	docker build --rm -t pokedex .

docker-run: docker-build
	docker run -p 5000:5000 pokedex

docker-run-tests:
	docker run -t -i -v ${PWD}:/pokedex -w /pokedex golang:1.17 sh -c 'go list ./... | grep -v vendor | xargs go test -count 1 -cover -race -v'