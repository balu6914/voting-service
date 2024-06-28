FROM golang:1.18 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o voting-service .
RUN go build -o wait-for-dynamodb ./cmd/wait-for-dynamodb

FROM golang:1.18

WORKDIR /app

COPY --from=builder /app/voting-service /app/voting-service
COPY --from=builder /app/wait-for-dynamodb /app/wait-for-dynamodb

EXPOSE 8000

ENTRYPOINT ["/app/wait-for-dynamodb"]
