order package files and services includes the following:

#### admin order matching
this is where we try to fulfill the orders sent from admin. if the price of orders are touched
since the time the order is created we would fulfill the order using our Bot.

#### bot aggregation service
the orders which are fulfilled using bot would sent into the redis queue and then we put the
order in our external exchange as the sum of all these orders so we can reduce our risk of losing
money because of fulfilling the orders by our Bot (ask Mr Jabbari why are we doing this).

#### decision manager
this is the service which decides how  the user order should be handled, whether it should be sent to
external exchange or be handled by our engine. it makes its decision based on order amount and price and the
field of `our_exchange_limit` in `pair_currencies` table.

#### engine communicator
this is where we talk to our engine we submit order to engine or ask the engine to remove the orders.

#### engine result handler 
this is like the callback that engine would call after matching process. 

#### events handler
this is where we handle after order status change , for example sending data to clients,invalidate cache or even send data
to external exchange.

#### force trader
this service decides if the order returned from engine should be fulfilled by our bot. it is where
the threshold is computed whether the order should be handled or not. it uses current market price
and the data in `bot_rules` of  `pair_currencies` table.

#### in queue order manager
when the prices changed , some orders in orderbook may meet the threshold to be fulfilled. this is
where this happens. it calls engine with the current price and the engine update the orderbook.

#### order create manager
this is where the order is inserted in db and the  user balance freezes. it also checks whether the
order user asks to create is valid.

#### post order matching service
this service would get the process of updating order data, create child data,decrease or increase user balances,
insert transactions and trades in db. and finally inform the users that their order is matched.

#### redis manager
is the responsible for inserting stop orders data in redis

#### service
is the place that clients interacts with our system, asking to create orders, cancel them and list of open  orders, 
order history and trade history.

#### stop order submission manager
when the prices changed , some stop orders may meet their stop point price.this is
where this happens. it submit an order based on the stop order.

#### trade event handle
this is where we put the fulfilled orders by our bot orders in redis (so we can finally send them to external exchange).

 

 

 
