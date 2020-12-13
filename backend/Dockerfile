FROM golang:1.15

WORKDIR /usr/src/app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build
RUN ls

ENTRYPOINT [ "./backend" ]