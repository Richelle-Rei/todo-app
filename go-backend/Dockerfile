FROM golang:1.24.4

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o ./run-app .

CMD ["./run-app"]