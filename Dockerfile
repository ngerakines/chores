FROM golang:1.14.6-alpine3.12 as chores-build
LABEL maintainer="Nick Gerakines <nick.gerakines@gmail.com>"
RUN apk --no-cache add build-base ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go generate -tags prod .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -ldflags "-w -s -extldflags '-static'" -tags prod github.com/ngerakines/chores/...

FROM alpine:3.12 as chores
RUN apk add --no-cache --update ca-certificates tzdata
RUN mkdir -p /app
RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime
RUN echo "US/Eastern" >  /etc/timezone
WORKDIR /app
COPY --from=chores-build /go/bin/chores /go/bin/
EXPOSE 8080
STOPSIGNAL SIGINT
ENTRYPOINT ["/go/bin/chores"]