.POSIX:
.SUFFIXES:

GO = go
RM = rm
GOFLAGS =
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin

goflags = $(GOFLAGS)

all: hidrocor

hidrocor:
	$(GO) build $(goflags) -o hidrocor .

clean:
	$(RM) -rf hidrocor

install: all
	mkdir -p $(DESTDIR)$(BINDIR)
	cp -f hidrocor $(DESTDIR)$(BINDIR)

.PHONY: hidrocor clean install
