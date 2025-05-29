FROM golang:1.24-alpine

RUN apk add --no-cache make git curl

WORKDIR /app

COPY services/warehouse/go.mod ./
COPY services/warehouse/go.sum ./

RUN sed -i 's/^go [0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?$/go 1.23/' go.mod

RUN go mod download

COPY services/warehouse/ ./
COPY services/warehouse/.env.example .env

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseDependency --parseInternal

RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /bin/warehouse .

EXPOSE 8080
ENTRYPOINT ["/bin/warehouse"]