FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o opgl-cortex-engine main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/opgl-cortex-engine .

EXPOSE 8082

CMD ["./opgl-cortex-engine"]
