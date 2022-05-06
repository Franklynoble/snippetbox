 #syntax=docker/dockerfile:1
FROM golang:1.16-alpine3.13 AS builder

WORKDIR /app
#COPY go.mod .
#COPY go.sum .
#RUN go mod download 
RUN go get github.com/githubnemo/CompileDaemon
 COPY . . 
 COPY ./docker-entrypoint.sh   /docker-entrypoint.sh
 #/entrypoint.sh
RUN go build -o main ./cmd/web

#Run stage do
FROM alpine:3.13
WORKDIR /app 
COPY --from=builder /app/main . 

EXPOSE 3000
# wait-for-it requires bash, which alpine doesn't ship with by default. Use wait-for instead
#ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
#RUN chmod +x  /docker-entrypoint.sh
#ENTRYPOINT ["chmod", "+x","/frank/go/src/github.com/Franklynoble/snippetbox/docker-entrypoint.sh"] 
CMD ["/app/main"]
