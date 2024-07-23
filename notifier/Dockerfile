FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod go.mod
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /consumer_group ./cmd/consumer_group/main.go

FROM scratch
COPY --from=builder consumer_group /bin/consumer_group

ENTRYPOINT ["/bin/consumer_group"]
