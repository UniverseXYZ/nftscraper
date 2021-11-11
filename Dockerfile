FROM golang:1-buster as builder

WORKDIR /usr/local/src/nftscraper

# disable CGO
ENV CGO_ENABLED=0

# copy only module information to take advantage of cache and layers
COPY go.mod go.sum ./

# download dependencies
RUN go mod download

# copy source files
COPY . ./

# build executable
RUN go build -o nftscraper

# make a fresh start to final image
FROM gcr.io/distroless/base-debian10

# copy executable
COPY --from=builder /usr/local/src/nftscraper/nftscraper /usr/local/bin/nftscraper

# use the executable as the main program for image
ENTRYPOINT [ "/usr/local/bin/nftscraper" ]