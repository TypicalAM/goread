FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goread

FROM gcr.io/distroless/static
LABEL maintainer="TypicalAM"
COPY --from=builder /app/goread /goread

ENTRYPOINT ["/goread"]
