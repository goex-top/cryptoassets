package main

// API

/* setting

GET - /setting   获取配置中的平台名称与API KEY
req - none
rsp - [{exchange_name:"binance.com", api_key:"xxxxxx"}]

POST - /setting  新增一个平台
req - [{exchange_name:"binance.com", api_key:"xxxxxx", sec_key:"xxxxxx"}]
rsp - none

*/

/* asset

GET - /asset_history 获取历史用户总资金
req - none
rsp - [{date:unix, value:200},{date:unix, value:400}]

GET - /asset 获取当前用户资金明细
req - none
rsp - [{exchange_name:"binance.com", value_usdt:400, value_btc: 0.001}]


GET - /exchange_detail 获取当前用户某平台资金明细
req - {exchange_name:"binance.com"}
rsp - {exchange_name:"binance.com", currencies:[usdt:400, btc: 0.001]}


*/
