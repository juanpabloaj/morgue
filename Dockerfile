FROM golang:1.12 AS builder
ADD . /app/backend
WORKDIR /app/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /morgue ./...

# final stage
FROM alpine:latest
COPY --from=builder /morgue ./

ENV PORT 8080

RUN chmod +x ./morgue
ENTRYPOINT ["./morgue"]
