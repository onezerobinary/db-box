# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev
WORKDIR /app
# Set an env var that matches github repo name,
ENV SRC_DIR=/go/src/github.com/onezerobinary/web-box
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN cd $SRC_DIR; go build -o myapp; cp myapp /app/
ENTRYPOINT ["./myapp"]