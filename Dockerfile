FROM golang:1.17.1 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o blockchainapp main.go

##########################################

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/blockchainapp ./blockchainapp

EXPOSE 1323

CMD ["./blockchainapp"]