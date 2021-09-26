.POSIX:
.SUFFIXES:

GO = go
RM = rm
GOFLAGS =
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man
SYSCONFDIR = /etc
SHAREDSTATEDIR = /var/lib

goflags = $(GOFLAGS)

all: hidrocor

hidrocor:
	$(GO) build $(goflags) -o hidrocor .

clean:
	$(RM) -rf hidrocor

install: all
	mkdir -p $(DESTDIR)$(BINDIR)
	mkdir -p $(DESTDIR)$(MANDIR)/man1
	mkdir -p $(DESTDIR)$(SYSCONFDIR)/hidrocor
	cp -f hidrocor $(DESTDIR)$(BINDIR)

.PHONY: hidrocor clean install
