FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Set destination for COPY
WORKDIR /app

COPY main ./main

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
#EXPOSE 8080

# Run
CMD ["/app/main"]