FROM golang:1.23
LABEL authors="najwang"

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 8080

CMD ["app"]

#You can then build and run the Docker image:
#$ docker build -t groupie-tracker-1 .
#$ docker run -p 80:8080 -d --restart unless-stopped --name groupie-tracker-1 groupie-tracker-1
