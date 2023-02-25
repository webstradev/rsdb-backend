FROM golang:1.20.0-alpine as builder

WORKDIR /app

COPY . . 

RUN C_GO_ENABLED=0 go build -o /bin/api .

FROM alpine

COPY --from=builder /app/bin/api /bin/api

EXPOSE 8080

CMD ["./bin/api"]