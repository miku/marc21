all: marctoxml marctexttoxml

marctoxml: cmd/marctoxml/main.go
	go build -o $@ $<

marctexttoxml: cmd/marctexttoxml/main.go
	go build -o $@ $<

clean:
	rm -f marctoxml marctexttoxml
