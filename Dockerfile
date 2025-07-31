# ##
#FROM golang:1.24.5-alpine  AS build
#
#WORKDIR /app
#
#COPY . .
#
#RUN --mount=type=secret,id=ENV \
#    --mount=type=secret,id=CREDENTIALS_JSON \
#    go build -o /main .
#
#
###
### Deploy
###
#FROM alpine:latest AS deploy
#
#WORKDIR /
#
#COPY --from=build /main /
#
#EXPOSE 9000
#
## run tracker on the background
#CMD [ "./main" ]

FROM golang:1.24.5-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 9000
CMD ["./main"]
