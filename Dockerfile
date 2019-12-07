FROM golang:1.13

RUN mkdir /app
WORKDIR /app

COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY main.go .
COPY api/ ./api
COPY models/ ./models

# Copy the toml config file
COPY config.toml .

RUN go build -o leisten .

# VOLUME log
EXPOSE 8000

ENTRYPOINT ["./leisten"]
