FROM golang:1.25-alpine AS builder

# install dependencies
RUN apk add --no-cache curl

# set working directory
WORKDIR /app

# copy the project files into the container
COPY go.mod go.sum ./
RUN go mod download

# copy the working directory
COPY . .

# build the application
RUN CGO_ENABLE=0 GOOS=linux go build -o /app/auth-service cmd/auth-service/main.go

FROM alpine:latest 

# install dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata

# copy the binary from the builder stage to the final
COPY --from=builder /app/auth-service /app/auth-service
RUN chmod +x /app/auth-service

# set working directory
WORKDIR /app

# expose port
EXPOSE 8080

# run the application
CMD ["/app/auth-service"]








