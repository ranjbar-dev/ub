## order matching engine 

Engine uses redis as storage and the data structure that represents our orderbook is redis sorted set.
redis sorted set has some unique features that make it the excellent choice for orderbook management,
the most important feature of redis sorted set is the rank of members which is price of orders in our cases.
engine runs as separate process and also provide interface for submitting orders and cancel them from our code.


### the process of order matching 
when our code decides that the order should be handled by our own engine. the submit method of of engine is called.
the engine push the order in simple queue in redis. then the separate process of engine pops the order from the queue and
load the related orderbook from redis to find best match(es) for the current order.the process of loading orderbook 
happening in goroutine which helps us to handle many matching at the same time. the process of order matching after 
loading the orderbook is based on the fact that the current order is market or limit. if the order is market, we try to 
find the best price order in order book. but we also consider the point that the price should be in a reasonable threshold
in comparison to market price. the matching process continue unless the order amount would fulfilled completely, finally
the result of matching and potential remaining of an order is sent back az callback to the `PostOrderMatching` service. 
this service then update orders databases, increase or decrease the user balances, insert transactions and insert trades.
the data after the matching consist two type of data, done orders which are the orders that are matched and partial order
that is remained because there was no more matching order, the `PostOrderMatching` tries to fulfill the remaining order
if its price is in the threshold, for example all the remaining market orders should be fulfilled, the remaining limit
orders should be fulfilled if their price is in the threshold,this fulfillment is considered as Bot fulfillment since we 
forcefully trade the orders although there was no matching candid for them. if the limit orders are not in threshold 
then the orders would return to the engine and the engine put them in orderbook again.

engine package files and services includes the following:

#### engine 
this service export an interface for other services so they can submit or cancel the order using this interface.
this services starts the worker pool and also run queue and dispatch the orders from the queue to workers 
which are responsible for orderbook and order matching.

#### order
the order struct represent the order which is used in our order matching process, this struct is not the same as the one 
we have for representing order database order placed in order package.

#### orderbook
this service is responsible for loading orderbook and matching orders. it needs an order provider which is redis in our case.

#### pool
pool is the place where workers are created and sending the orders to the workers.

#### queue 
is the redis queue , the orders submitted to engine first will be put in queue, this service is the place where
implements this feature, also it pop the orders too so they can be matched.

#### redis orderbook provider
this is implement the orderbook interface using redis. think about it this way if could have orderbook in mysql
then we would have mysql orderbook provider.

#### worker
is representation of our workers which load the orderbook, give them the data, take the result from them and
callback result handler. 





 
 