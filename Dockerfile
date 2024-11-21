FROM golang:1.23.3

RUN mkdir -p /usr/src/app

WORKDIR /usr/src/app

COPY . /usr/src/app

RUN go mod download

EXPOSE 3000

CMD ["go", "run", "./cmd"]

