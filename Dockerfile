#######             #######
####### MULTI STAGE #######
#######             #######

####### BUILD STAGE ####### 

# Base image
FROM golang:1.18.2-alpine AS builder
# Working directory inside the image
WORKDIR /app
# Copy everything from the root of our project to workdir
COPY . .
# Build our app to a binary single executable file.
RUN go build -o main main.go
# Install cURL for downloading golang migrate
RUN apk add curl
# Install golang migrate.
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz


####### RUN STAGE ####### 

FROM alpine
WORKDIR /app
# Copy binary from builder stage (/app/main) to run stage (/app)
COPY --from=builder /app/main .
# Copy migrate from builder stage to run stage
COPY --from=builder /app/migrate /usr/bin/migrate
# Copy env file to load configuration from builder stage (/app/app.env) to run stage (/app)
COPY --from=builder /app/app.env .
# Copy start.sh from builder stage (/app/start.sh) to run stage (/app)
COPY --from=builder /app/start.sh .
# Copy wait-for.sh from builder stage (/app/wait-for.sh) to run stage (/app)
COPY --from=builder /app/wait-for.sh .
# Copy migrations files from build stage to run stage
COPY --from=builder /app/db/migration ./migration

# Expose ports used
EXPOSE 8080
# Define default CMD to run when the container starts
CMD [ "/app/main" ]
# Main entry point of docker image. When CMD used together with entry point, 
# the CMD will be passed as an additional param to the entry point script.
# Similar to ENTRYPOINT ["/app/start.sh", "/app/main"]
ENTRYPOINT ["/app/start.sh"]



