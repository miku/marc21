marctoxml: cmd/marctoxml/main.go
	go build -o $@ $<

clean:
	rm -f marctoxml
