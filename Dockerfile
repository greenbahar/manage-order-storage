# syntax=docker/dockerfile:1
FROM golang:1.17

#ADD . /go/src/manage-order-storage
WORKDIR /go/src/github.com/greenbahar/manage-order-storage
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy
COPY . ./
RUN go build -o /manage-order-storage
EXPOSE 3001
CMD [ "/manage-order-storage" ]

#ENTRYPOINT [ "/manage-order-storage" ]