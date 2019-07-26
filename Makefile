CGO_ENABLED=0

all: marctoxml marctexttoxml

marctoxml: cmd/marctoxml/main.go
	go get -v ./...
	CGO_ENABLED=$(CGO_ENABLED) go build -o $@ $<

marctexttoxml: cmd/marctexttoxml/main.go
	go get -v ./...
	CGO_ENABLED=$(CGO_ENABLED) go build -o $@ $<

clean:
	rm -f marctoxml marctexttoxml
