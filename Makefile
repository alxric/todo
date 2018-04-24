-include CONFIG
-include CREDENTIALS

build:
	GOARCH=amd64 go build -o $(APPLICATION)-$(VERSION)-linux-amd64/$(APPLICATION) .
	GOARCH=386 go build -o $(APPLICATION)-$(VERSION)-linux-386/$(APPLICATION) .
	GOOS=windows GOARCH=amd64 go build -o $(APPLICATION)-$(VERSION)-windows-amd64/$(APPLICATION).exe .
	GOOS=darwin GOARCH=amd64 go build -o $(APPLICATION)-$(VERSION)-mac-amd64/$(APPLICATION) .

create_dir: build
	mkdir -p $(APPLICATION)-$(VERSION)-linux-amd64
	mkdir -p $(APPLICATION)-$(VERSION)-linux-386
	mkdir -p $(APPLICATION)-$(VERSION)-windows-amd64
	mkdir -p $(APPLICATION)-$(VERSION)-mac-amd64

create_tar: create_dir
	tar -cvzf $(APPLICATION)-$(VERSION)-linux-amd64.tar.gz $(APPLICATION)-$(VERSION)-linux-amd64/$(APPLICATION)
	tar -cvzf $(APPLICATION)-$(VERSION)-linux-386.tar.gz $(APPLICATION)-$(VERSION)-linux-386/$(APPLICATION)
	tar -cvzf $(APPLICATION)-$(VERSION)-windows-amd64.tar.gz $(APPLICATION)-$(VERSION)-windows-amd64/$(APPLICATION).exe
	tar -cvzf $(APPLICATION)-$(VERSION)-mac-amd64.tar.gz $(APPLICATION)-$(VERSION)-mac-amd64/$(APPLICATION)

release: create_tar
