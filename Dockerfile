FROM golang:1.25-alpine

RUN apk add make

COPY . .

RUN go install

RUN go build -o /tunellm .

CMD ["/tunellm"]
