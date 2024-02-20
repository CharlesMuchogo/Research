FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 9000
USER charles
CMD ["./main"]