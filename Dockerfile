FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . . 

RUN C_GO_ENABLED=0 go build -o /app/bin/api .

FROM alpine

COPY --from=builder /app/bin/api /bin/api

COPY --from=builder /app/migrations /migrations

EXPOSE 8080

CMD ["./bin/api"]