FROM golang:1.25-alpine

RUN apk add make

COPY . .

RUN go install

# release build
# -h: halt on error
# -s: disable symbol table
# -w: disable DWARF generation
# -trimpath strips computers file paths from error messages
RUN go build \
      -trimpath \
      -ldflags "-s -w -h -X main.version=$(VERSION)" \
      -o ../target/tunnellm

CMD ["/tunnellm"]
