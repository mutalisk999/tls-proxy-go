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

