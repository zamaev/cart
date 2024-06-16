FROM golang:1.22-alpine as builder

WORKDIR /build

ENV PROTOC_ZIP=protoc-25.1-linux-x86_64.zip
RUN apk add unzip curl
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v25.1/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local \
    && export PATH="$PATH:$HOME/local/bin" \
    && rm -f $PROTOC_ZIP

RUN apk add --no-cache make

COPY Makefile Makefile
RUN make bin-deps

COPY go.mod go.mod
RUN go mod download

COPY . .
RUN make vendor-proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/server.go

FROM scratch
COPY --from=builder server /bin/server

ENTRYPOINT ["/bin/server"]
