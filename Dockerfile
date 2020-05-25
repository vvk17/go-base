# Dockerfile References: https://docs.docker.com/engine/reference/builder/

#Start form the latest golang base image
FROM golang:latest as builder

#Add Maintainer Info
LABEL maintainer="vvk17 vvk17@mail.ru"

#Set the current Work Dir inside the container
WORKDIR /app

#Copy go mod amd sum files
COPY go.mod go.sum ./

#Download all dependencies. They are caches if go.mod and go.sum are not changes
RUN go mod download

#Copy source from the current directory into Work Dir of the container
COPY . .

#Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o main .


###### Start a new stage with small image from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

#Copy the Pre-built binary file from previous stage
COPY --from=builder /app/main .

#Expose port 
EXPOSE 8080

#Run the executable
CMD ["./main"]