<p align="center">
    <img src="https://raw.githubusercontent.com/goex-top/cryptoassetsweb/master/public/favicon.ico">
</p>
<p align="center">
  <a href="https://github.com/golang/go">
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/goex-top/cryptoassets">
  </a>

  <a href="https://github.com/goex-top/cryptoassets/master/LICENSE">
    <img src="https://img.shields.io/github/license/mashape/apistatus.svg" alt="license">
  </a>
  <a href="https://www.travis-ci.com/goex-top/cryptoassets">
    <img src="https://www.travis-ci.com/goex-top/cryptoassets.svg?branch=master" alt="build status">
  </a>
</p>
<p align="center">
  <a href="https://github.com/vuejs/vue">
    <img src="https://img.shields.io/badge/vue-2.6.11-brightgreen.svg" alt="vue">
  </a>
  <a href="https://github.com/vuejs/vue">
    <img src="https://img.shields.io/badge/vue-2.6.11-brightgreen.svg" alt="vue">
  </a>
  <a href="https://github.com/ElemeFE/element">
    <img src="https://img.shields.io/badge/element--ui-2.13.0-brightgreen.svg" alt="element-ui">
  </a>
</p>

# Crypto Assets
> [中文](https://github.com/goex-top/cryptoassets/blob/master/README-cn.md)

Recording your crypto assets with multi accounts

![image](https://raw.githubusercontent.com/goex-top/cryptoassetsweb/master/assets.gif)

## Let's start
Run this tool, it can record your assets automatic.

### build with source code
> [install `go`](https://golang.org/doc/install)
* `git clone https://github.com/goex-top/cryptoassets.git`
* `git submodule update --init --recursive`
* `go build`

### Run it
* `./cryptoassets`
* open brower [http://localhost:9000](http://localhost:9000)
* input username and password which in `config.toml` file

### Add API KEY
add API KEY in setting view
![image](https://raw.githubusercontent.com/goex-top/cryptoassetsweb/master/settings.png)

## Config
create a config file, `config.toml`, using `cp sample-config.toml config.toml` , then modify it with followed description

```toml
proxy=""                 # socks5://127.0.0.1:1080
freq=60                  # unit: second, 60 for 1min
debug = true             # enable / disable verbase log print
[user]
username="admin"         #  username for login
password="AbcdEfgh"      # password for login and encrypts and decrypts your apiseckey to store in database
```

## API KEY store in database
* When a user creates an exchange, the security key / passhase key will be encrypted by AES (ECB) and stored in the database. Remember the password in the toml configuration file. This password is the only password to decrypt the key in the database.
* Try to create a read-only API KEY

## Database
ORM using [GORM](https://github.com/jinzhu/gorm), support `MySQL`, `PostgreSQL`, `Sqlite3`, `SQL Server` 
Currently `sqlite3` is used, it can create `sqlite3` file automatically, **CONVENIENCE**, you can copy it to everywhere

### Data models
3 tables
* account
  - store all API KEY of all exchanges
* assets history. Total assets of each exchange, stored at regular intervals (based on freq in the configuration)
  - total valuation of BTC, USD, USDT and CNY
* coin assets history. All coin asset history for all exchanges
  - coin valuation of BTC, USD, USDT and CNY

```sql
CREATE TABLE accounts (
    id               INTEGER       PRIMARY KEY AUTOINCREMENT,
    created_at       DATETIME,
    updated_at       DATETIME,
    deleted_at       DATETIME,
    nick_name        VARCHAR (255) UNIQUE,
    exchange_name    VARCHAR (255),
    api_key          VARCHAR (255),
    api_secret_key   VARCHAR (255),
    api_passphrase   VARCHAR (255),
    last_update_time BIGINT
);

CREATE TABLE assets (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    account_id INTEGER,
    btc        REAL,
    usdt       REAL,
    usd        REAL,
    cny        REAL,
    btc_usdt   REAL,
    btc_usd    REAL,
    btc_cny    REAL,
    usdt_usd   REAL,
    usdt_cny   REAL,
    usd_cny    REAL
);

CREATE TABLE coin_assets (
    id            INTEGER       PRIMARY KEY AUTOINCREMENT,
    created_at    DATETIME,
    updated_at    DATETIME,
    deleted_at    DATETIME,
    asset_id      INTEGER,
    coin_name     VARCHAR (255),
    amount        REAL,
    frozen_amount REAL,
    btc           REAL,
    usdt          REAL,
    usd           REAL,
    cny           REAL
);

```

## Exchanges to support
Exchange | Spot | Future(Contract) | Future(Swap) | LOGO
:-: | :-: | :-: | :-: | :-: 
[BitMEX](https://www.bitmex.com/register/tIRSfz) | | ☑️ | ☑️ | [![bitmex](https://user-images.githubusercontent.com/1294454/27766319-f653c6e6-5ed4-11e7-933d-f0bc3699ae8f.jpg)](https://www.bitmex.com/register/tIRSfz) |
[Binance](https://www.binance.com/?ref=10052861) | ☑️|  | ☑️ | [![binance](https://user-images.githubusercontent.com/1294454/29604020-d5483cdc-87ee-11e7-94c7-d1a8d9169293.jpg)](https://www.binance.com/?ref=10052861) |
[OKEx](https://www.okex.com) | ☑️ | ☑️ | ☑️ |[![OKEx](https://user-images.githubusercontent.com/1294454/32552768-0d6dd3c6-c4a6-11e7-90f8-c043b64756a7.jpg)](https://www.okex.com) |
[Huobi](https://www.huobipro.com/zh-cn/topic/invited/?invite_code=n6d33) | ☑️| ☑️ |  | [![huobipro](https://user-images.githubusercontent.com/1294454/27766569-15aa7b9a-5edd-11e7-9e7f-44791f4ee49c.jpg)](https://www.huobipro.com/zh-cn/topic/invited/?invite_code=n6d33) |
[Poloniex](https://www.poloniex.com/?utm_source=goex&utm_medium=web) | ☑️|  |  | [![poloniex](https://user-images.githubusercontent.com/1294454/27766817-e9456312-5ee6-11e7-9b3c-b628ca5626a5.jpg)](https://www.poloniex.com/?utm_source=goex&utm_medium=web)|
[Bitfinex](https://www.bitfinex.com) | ☑️|  |  | [![bitfinex](https://user-images.githubusercontent.com/1294454/27766244-e328a50c-5ed2-11e7-947b-041416579bb3.jpg)](https://www.bitfinex.com)|
[Bitstamp](https://www.bitstamp.net) | ☑️|  |  | [![bitstamp](https://user-images.githubusercontent.com/1294454/27786377-8c8ab57e-5fe9-11e7-8ea4-2b05b6bcceec.jpg)](https://www.bitstamp.net) |
[Bittrex](https://bittrex.com) | ☑️|  |  | [![bittrex](https://user-images.githubusercontent.com/1294454/27766352-cf0b3c26-5ed5-11e7-82b7-f3826b7a97d8.jpg)](https://bittrex.com) |
[Bithumb](https://www.bithumb.com) | ☑️|  |  | [![bithumb](https://user-images.githubusercontent.com/1294454/30597177-ea800172-9d5e-11e7-804c-b9d4fa9b56b0.jpg)](https://www.bithumb.com)|
[GateIO](https://www.gate.io/signup/330917) | ☑️|  |  | [![GateIO](https://user-images.githubusercontent.com/1294454/31784029-0313c702-b509-11e7-9ccc-bc0da6a0e435.jpg)](https://www.gate.io/signup/330917)|
[ZB](https://www.zb.com) | ☑️|  |  | [![zb](https://user-images.githubusercontent.com/1294454/32859187-cd5214f0-ca5e-11e7-967d-96568e2e2bd1.jpg)](https://www.zb.com)  |
[BigONE](https://b1.run/users/new?code=7JDU9ANL) | ☑️|  |  | [![BigONE](https://user-images.githubusercontent.com/1294454/69354403-1d532180-0c91-11ea-88ed-44c06cefdf87.jpg)](https://b1.run/users/new?code=7JDU9ANL)  |
[HitBTC](https://hitbtc.com/) | ☑️|  |  | [![HitBTC](https://user-images.githubusercontent.com/1294454/27766555-8eaec20e-5edc-11e7-9c5b-6dc69fc42f5e.jpg)](https://hitbtc.com/) |

**All assets in difference type of account per exchange will be merged**

## Rate 
* USD/CNY fetch from [finance of yahoo](https://finance.yahoo.com/)
* USDT/USD fetch from[Binance US](https://www.binance.us/en/trade/USDT_USD)
* BTC/USD fetch from[Binance US](https://www.binance.us/en/trade/BTC_USD)

**Update per 2 hours**

## Source code of frontend
If you want to modify frontend, please check out source code [https://github.com/goex-top/cryptoassetsweb.git](https://github.com/goex-top/cryptoassetsweb.git)
