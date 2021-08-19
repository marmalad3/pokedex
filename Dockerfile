FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build ./cmd/pokedex

ENV PORT 5000
EXPOSE 5000

CMD [ "./pokedex" ]