FROM golang:1.12.7 AS builder

RUN apt-get -q update && \
    apt-get -q install -y \
    ca-certificates \
    && rm -r /var/lib/apt/lists/*

WORKDIR /calendar
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app ./

EXPOSE 5997

ENTRYPOINT ["./app"]