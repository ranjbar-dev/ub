#### order placement

what kind of orders we support

kind | description | required fields
---- | ----------- | ---- 
market     | the order without price from users | pair, type, amount 
limit      | the order users choose the price | pair, type, amount, price
stop limit | the order with criteria price | pair, type, amount, price, stop_point_price

##### market order
This kind of order is the order which has no price and should be traded with the best existing price in orderbook
users only provides the amount and we try to fulfill the order with the market current price
parameters of this kind of orders

filed name       | type   | description
---------------- | ------ |-----------
type             | string | buy or sell 
exchange_type    | string | is market
amount           | string | the amount users wants to buy or sell like "2000.0"
pair_currency_id | int64 | the id of pair like 2
price            | string | empty for market orders


##### limit order
This kind of order is the order with price provided by user, and should be traded with the orders in  orderbook
that pass the price criteria the user asks for
parameters of this kind of orders

filed name       | type   | description
---------------- | ------ |-----------
type             | string | buy or sell 
exchange_type    | string | is limit
amount           | string | the amount users wants to buy or sell like "0.2"
pair_currency_id | int64 | the id of pair like 2
price            | string | the price users wants to buy or sell at  like "35000.0"

 ##### stop order
 This kind of order is the order with stop point price  and price both provided by user, this kind of order is the potential 
 order and it means that the user wants if the price reaches the stop point price, an order with the amount and price would be
 submitted to order book which could be fulfilled or not
 parameters of this kind of orders
 
 filed name       | type   | description
 ---------------- | ------ |-----------
 type             | string | buy or sell 
 exchange_type    | string | is limit
 amount           | string | the amount users wants to buy or sell like "0.2"
 pair_currency_id | int64 | the id of pair like 2
 price            | string | the price users wants to buy or sell like "35000.0"
 stop_point_price | string | the price users wants to place order if it reached like "35000.0"
 
 ### order placement process
 the order is created in `OrderCreationService`,in create process we insert order in database and freeze the user balance (the payed by coin).
 if the order is `stop order` then we just add it to redis for the time its `stop point price` reaches. if it is limit or market, 
 then we decide where should be the order placed. we have two options,one is our orderbook, and the other one is external exchange (binance),
 this decision is happened in `DecisionManager` service which will choose based on following criteria
 if the amount of order is exceeding from the `our_exchange_limit` in pair_currencies table  and the order is market then
 the order would be sent to external exchange,if the order is limit it would be sent to our order matching engine.
 all market orders should be fulfilled, the limit ones should be fulfilled if their price is in the threshold defined in `bot_rules`
 
 