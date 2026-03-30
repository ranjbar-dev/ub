#### What services we have in backend

our services includes following:

##### exchange-server
written in golang, is responsible for end user request to our system.
it has dependency to mysql,redis,rabbitmq, centrifugo (real-time messaging) all our different services run in docker containers

#### ub-server
written in php, is responsible for admin part of the exchange.
it has dependency to mysql,redis,rabbitmq, centrifugo (real-time messaging).

#### wallet
written in golang, is responsible for handling wallet.
it has dependency to mysql,redis.but these services our different one from the  ones exchange is using.
it also has interaction with admin part and call the callback of from ub-server with the result of
blockchain transactions data.

#### ub-communicator
written in golang, is responsible for sending emails and sms.
it has dependency to mongodb,rabbitmq.the rabbitmq is where relay between other services and the ub-communicator.

