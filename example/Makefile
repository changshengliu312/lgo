  target=test
  ($target):
    go vet ./
    go build -o ${target} -ldflags '-linkmode "external" -extldflags "-static"' *.go
  clean:
    rm -f ${target}
