# Start from golang base image
FROM golang:1.20.1-alpine

# add maintainer info
LABEL maintainer="Elisio Sa"

# setup folders
RUN mkdir /app
WORKDIR /app

# copy the source from current dir dir inside the container
COPY . .
COPY .env .

# downloads all dependecies
RUN go get -d -v ./...

#install the package
RUN go install -v ./...

# build the go app
RUN go build -o api .

# Expose port 8080 to the outside world
EXPOSE 8080

# run the executable
CMD ["./api"]