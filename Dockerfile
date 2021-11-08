FROM golang:1.17 as compiler
COPY . /app
WORKDIR /app
RUN go build -mod=vendor -o /bin/server cmd/server/main.go
RUN go build -mod=vendor -o /bin/swarmupd cmd/cli/main.go

FROM debian:jessie


ENV TOKENS=none
ENV PORT=8000
ENV REGISTRY_USER=none
ENV REGISTRY_PASSWORD=none
ENV SERVICE_PREFIXIES_ONLY=none


ENV DOCKER_HOST unix:///var/run/docker.sock 
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=compiler /bin/server /bin/server
EXPOSE 80 443
CMD ["bin/server"]
