FROM golang:1.22

WORKDIR /app
COPY ./ ./
RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
RUN go build -o /avito ./cmd/
EXPOSE 8080
