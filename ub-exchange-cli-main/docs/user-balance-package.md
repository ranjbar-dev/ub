#### how we handle user balance

every user has many user balance for coins we support. each user balance has an address for itself,
the main data in user balance table are the amount,frozen amount and the address.

our system also supports the coins that could be deposit and withdraw for more than one blockchain network.
for example for USDT we have both the Ethereum and Tron. for this user balances we have a json field in
table named `other_addresses` contains the address of other network for coin. for example the USDT user balance
`address` field contains the ETH network address and the Tron network address is saved in field `other_addresses`.

user balance package files and services includes the following:
#### service
this is where create balance. get the address for each coin for any user.
list of balances and all related to user balance. 
 