FROM golang:1.23 as build

WORKDIR /app


# Copy the go.mod and go.sum files of each module to the container
COPY e2ee_api/go.mod e2ee_api/go.sum ./e2ee_api/
COPY e2ee_handlers/go.mod e2ee_handlers/go.sum ./e2ee_handlers/
COPY x3dh_core/go.mod x3dh_core/go.sum ./x3dh_core/
COPY x3dh_server/go.mod x3dh_server/go.sum ./x3dh_server/

# Download the dependencies for each module
RUN cd /app/e2ee_api && go mod tidy
RUN cd /app/e2ee_handlers && go mod tidy
RUN cd /app/x3dh_core && go mod tidy
RUN cd /app/x3dh_server && go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go app (make sure to specify the main module directory)
WORKDIR /app/e2ee_handlers
RUN export CGO_ENABLED=0 && go build -o myapp .


# Step 2: Create a minimal container for the app (uses a smaller image)
FROM alpine:latest  

# Install required dependencies for your app (if any)
RUN apk --no-cache add ca-certificates

# Download the Amazon RDS global certificate bundle
RUN wget https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem -O /usr/local/share/ca-certificates/global-bundle.crt

# Update the CA certificates store
RUN update-ca-certificates

# Copy the compiled Go binary from the build stage
COPY --from=build /app/e2ee_handlers/myapp /usr/local/bin/myapp

# Expose the port the app runs on (modify if necessary)
EXPOSE 8080

# Run the Go app
CMD ["myapp"]