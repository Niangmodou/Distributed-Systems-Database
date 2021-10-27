FROM golang:latest

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go get github.com/kataras/iris/v12@master
RUN go get go.mongodb.org/mongo-driver/mongo

RUN go build -o main .

EXPOSE 8080

CMD [ "/app/main" ]