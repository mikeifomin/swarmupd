FROM golang:1.12 as compiler
COPY . /app
WORKDIR /app
RUN go build -mod=vendor -o /bin/server main.go

FROM debian:jessie

ENV DOCKER_HOST unix:///var/run/docker.sock 
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=compiler /bin/server /bin/server
EXPOSE 80 443
ENTRYPOINT ["bin/server"]
