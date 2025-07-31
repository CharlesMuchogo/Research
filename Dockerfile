FROM golang:1.24.5-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 9000
CMD ["./main"]
