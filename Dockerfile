FROM golang:1.23-alpine AS builder

WORKDIR /app

# Needed for Templ
RUN apk add --no-cache --update git build-base && \
    go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN templ generate ./views/
RUN CGO_ENABLED=1 go build -o scout ./cmd/scout/main.go


FROM alpine:latest

WORKDIR /app

RUN mkdir /app/files

COPY --from=builder /app/scout .

EXPOSE 6969

VOLUME ["/app/db"]
VOLUME ["/app/files"]

CMD ["./scout"]
