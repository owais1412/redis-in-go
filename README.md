# redis-in-go

# Run
`$ make run`

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