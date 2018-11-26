# build docker image starting from this base image in the public repo
FROM golang:1.11.0-stretch as builder
# create and cd into this directory. 
WORKDIR /go-modules
# Copy all files in source working directory into Docker working directory set above
COPY . ./
# using shell form (/bin/sh) set env vars, then build the go pacakge in the working directory. Name the output 'app'
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o app 
# build a new docker image from alpine:3.8 base image
FROM alpine:3.8
# set a new working direcory
WORKDIR /root/
# Copy the binary executable from the previous image into this image
COPY --from=builder /go-modules/app .
# run the app
CMD ["./app"]