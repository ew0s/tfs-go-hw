# GoTrader

<img src="https://raw.githubusercontent.com/egonelbre/gophers/10cc13c5e29555ec23f689dc985c157a8d4692ab/vector/friends/crash-dummy.svg" alt="gopher" width="30%"/>

[![Build Status](https://img.shields.io/badge/CI-passing-brightgreen)](https://github.com/ew0s/tfs-go-hw/actions)

A cryptocurrency trading bot supporting kraken futures written in Golang.

---

## Current Features

* Support for sending any order on kraken futures (mkt, lmt, etc...)
* Support trading on kraken futures using stop loss & take profit indicator
* REST API support for kraken futures
* Websocket API support for kraken futures
* JWT Token auth support with deleting token on logout from device
* Telegram bot 
* Swagger documentation

---

## Planned Features

* Support multiple kraken api tokens

---

## Exchange support table

| Exchange            | REST API | Streaming API | 
|---------------------|----------|---------------|
| Kraken futures demo | Yes      |  Yes          |
| Kraken futures      | Yes      |  Yes          |

---

## Tech stack

* [Go](https://github.com/golang/go)
* [Postgres](https://www.postgresql.org)
* [Redis](https://redis.io)

---

## Swagger

__When server started:__ ```url: http://{host}:{port}/swagger/index.html```

---

## Installation

### Linux/OSX

```shell
git clone {this repo}
cd {this repo}/course_project/trade-bot
```

* #### Assume you have ```config.yml``` or ```config.yaml``` file in configs folder of type: 

```yaml
server:
  port: (int) 
  websocket:
    readBufferSize: (int)
    writeBufferSize: (int)
    checkOrigin: (true | false)

client:
  # url of server
  url: (string)

postgreDatabase:
  host: (string)
  port: (string)
  username: (string)
  dbname:  (string)
  sslmode: (string)

redisDatabase:
  port: (string)

kraken:
  apiurl: (string)

krakenWS:
  requests:
    writeWaitInSeconds: (int)
    pongWaitInSeconds: (int)
    pingPeriodInSeconds: (int)
    maxMessageSize: (int)
  kraken:
    wsapiurl: (string)
```

* #### Assume you have ```.env``` file on top of project of type:

```.dotenv
DB_PASSWORD = (your postgres db password)

JWT_ACCESS_SIGNING_KEY = (key for signing jwt tokens)

PUBLIC_API_KEY = (public key from kraken futures)
PRIVATE_API_KEY = (private key from kraken futures)

TELEGRAM_APITOKEN = (telegram api token)
WEBHOOK_URL = (your webhook url for telegram bot)
```

* #### Run postgres with settings from your config file
* #### Run redis with settings from your config file
* #### Run migrate files for postgres using ```migrate```
```shell
migrate -path ./schema -database 'postgres://{postgres_username}:{postgres_password}@{host}:{port}/postgres?sslmode={sslmode}' up
```

* #### Then run server and telegram bot

```shell
go run cmd/api/main.go
go run pkg/telegramBot/cmd/api/main.go
```

---



