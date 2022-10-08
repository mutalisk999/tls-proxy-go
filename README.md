# tls-proxy-go



### How to generate cert file and key file

```
mkdir certs
rm certs/*

# for tls server
openssl req -new -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 36500

# for tls client
openssl req -new -nodes -x509 -out certs/client.pem -keyout certs/client.key -days 36500
```


```
mkdir certs
rm certs/*

# CA key and cert
openssl genrsa -out certs/ca.key 2048
openssl req -new -key certs/ca.key -out certs/ca.csr
openssl x509 -req -in certs/ca.csr -out certs/ca.pem -signkey certs/ca.key -CAcreateserial -days 36500

# for tls server
openssl genrsa -out certs/server.key 2048
openssl req -new -key certs/server.key -out certs/server.csr
openssl x509 -req -in certs/server.csr -out certs/server.pem -signkey certs/server.key -CA certs/ca.pem -CAkey certs/ca.key -CAcreateserial -days 36500

# for tls client
openssl genrsa -out certs/client.key 2048
openssl req -new -key certs/client.key -out certs/client.csr
openssl x509 -req -in certs/client.csr -out certs/client.pem -signkey certs/client.key -CA certs/ca.pem -CAkey certs/ca.key -CAcreateserial -days 36500

```

### Client

* Modify client.json
```
cp client_example.json client.json
```

### Server

* Modify server.json
```
cp server_example.json server.json
```

