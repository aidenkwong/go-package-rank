FROM golang:1.18

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .

ENV PORT=8080

EXPOSE 8080

CMD ["/app/main"]
