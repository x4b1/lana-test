FROM golang:1.17.1-alpine as builder

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api


FROM alpine:latest

ENV PORT 8080

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /api/main .

EXPOSE ${PORT}

CMD ["./main"]
