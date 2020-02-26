# Crypto Assets
统计你的个人资产

## 缘由
FCoin 跑路时，提币需邮件回复账户里资金数量，有多少人能知道自己资金具体数目？难。

## 快速上手
打开工具时自动记录其各平台账户资金数目，统计其占比

## 配置
创建一份`config.toml`配置文件，如`cp sample-config.toml config.toml` ，修改其内容

```toml
proxy=""                 # socks5://127.0.0.1:1080
freq=60                   # unit: second, 60 for 1min
[user]
username="admin"         #  username for login
password="AbcdEfgh"      # password for login and encrypts and decrypts your apiseckey to store in database
```

## 密钥存储
用户创建交易所时，密钥会通过AES(ECB)加密后存储至数据库中，切记`toml`配置文件中的`password`，这个`password`是解密数据库中密钥的唯一密码。