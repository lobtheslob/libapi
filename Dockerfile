FROM golang:1.17 as base

# Install the watcher
RUN go get github.com/codegangsta/gin

# Download deps
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source as late in the process as possible (to speed up local builds)
COPY . .

#################
# BUILDER STAGE #
#################
FROM base AS builder
WORKDIR /app
COPY --from=base /app .
RUN GOOS=linux GOARCH=386\
				go build -v\
				-o app\
				src/*.go

#################
# FINAL STAGE #
#################
FROM alpine:3.7
WORKDIR /app
COPY --from=builder /app/app .
CMD [ "./app"]
