# redis-in-go

# Run
1. Run `$ make run` in a terminal
2. Run `$ redis-cli` in another terminal and start running the commands

# Commands supported
## PING
```
> ping
> ping abc
```
## GET
`> GET key`
## SET
`> SET key value`
## DEL
`> DEL key1 key2`
## EXISTS
`> EXISTS key1 key2`
## HSET
`HSET hashKey1 innerKey1 value`
## HGET
`HGET hashKey1 innerKey1`
## HGETALL
`HGETALL hashKey1`