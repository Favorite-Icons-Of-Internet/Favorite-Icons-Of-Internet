GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i

GOPATH:=$(CURDIR)/_vendor:$(CURDIR)/

APPNAME = gofavicon
DISTDIR = dist

all: dist

clean:
	rm -rf bin
	rm -rf $(DISTDIR)

env:
	$(GO) env

test:
	$(GOTEST) gofavicon

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOINSTALL) $(APPNAME)/...
	CGO_ENABLED=0 $(GOINSTALL) $(APPNAME)/...

dist: clean build
	mkdir $(DISTDIR)
	cp -R bin/linux_386/ $(DISTDIR)
	echo "Dist done"
