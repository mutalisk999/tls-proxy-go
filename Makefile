GOCC=go
GOFLAGS=-ldflags '-w -s'

###############################

all: client server
.PHONY: all

clean: client_clean server_clean
.PHONY: clean

###############################

client: client_clean
	$(GOCC) build $(GOFLAGS) ./cmd/tls-proxy-client
.PHONY: client

client_clean:
	rm -f tls-proxy-client
.PHONY: client_clean

server: server_clean
	$(GOCC) build $(GOFLAGS) ./cmd/tls-proxy-server
.PHONY: server

server_clean:
	rm -f tls-proxy-server
.PHONY: server_clean