FROM golang:1.20.3-alpine

# Set the working directory to /app

WORKDIR /app

# Copy the go.mod and go.sum files to the working directory

COPY go.mod go.sum ./

# Download and install any required Go dependencies

RUN go mod download

# Copy the entire source code to the working directory 

COPY . .

# Build the go application

RUN go build -o main . 

# Expose the port specified by the port env variable

EXPOSE 3333

# Set the entry point of the container to the executable
CMD [ "./main" ]