FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go build -o app .
EXPOSE 8080
CMD ["./app"]
