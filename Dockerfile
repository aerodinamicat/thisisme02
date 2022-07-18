FROM golang:latest as build-env

ENV APP_NAME thisisme_userauth
ENV CMD_PATH cmd/service/main.go

WORKDIR $GOPATH/src/$APP_NAME

COPY . .
COPY .env .

RUN go get ./...
RUN go install ./...

RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

FROM golang:latest 

ENV APP_NAME thisisme_userauth

COPY --from=build-env /$APP_NAME .

EXPOSE 5070

CMD ./$APP_NAME