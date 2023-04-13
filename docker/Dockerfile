# Builder
FROM golang:1.19-buster as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . /app

RUN go build -o /app/main /app/cmd/http/main.go

# Final stage
FROM gcr.io/distroless/base-debian11

COPY --from=builder /app/main /app/main

EXPOSE 8080

ENTRYPOINT ["/app/main"]
