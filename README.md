# Assignment Solution

The solution provided here attempts to fulfill all the task requirements. Running the binary should put a local server up that can be
accessed via HTTP calls.

To run the server (in the directory with all the solution files):
```
$ go run .
```

To query the server:
```
$ curl -XPOST  'http://localhost:8080/transaction' --data '{"denominations": [100]}'
{"id":"de502339-fcd3-41da-b622-5da7b5980b1c","total_price":0,"total_inserted":[100],"total_returned":null,"completed":false}

$ curl -XPOST  'http://localhost:8080/transaction' --data '{"denominations": [100], "transaction_id":"de502339-fcd3-41da-b622-5da7b5980b1c"}'
{"id":"de502339-fcd3-41da-b622-5da7b5980b1c","total_price":0,"total_inserted":[100,100],"total_returned":null,"completed":false}

$ curl  -XPOST  'http://localhost:8080/transaction' --data '{"denominations": [200], "transaction_id":"de502339-fcd3-41da-b622-5da7b5980b1c", "total_price": 150, "finish_transaction":true}'
{"id":"de502339-fcd3-41da-b622-5da7b5980b1c","total_price":0,"total_inserted":[100,100,200],"total_returned":[200,50],"completed":true}
```

## High Level Design

Overall, the implementation is quite simple. Since the problem deals with money, there is a `denomination` type that
is modelled on top of a `uint`. There are predefined currency values with minimum denomination being a cent, equal to 1. This
prevents floating point complications.

The primary function that does payment procesing is `processPaymentInCents`. This function takes in a register type `registerer[denomination]`
which provides the two basic interface to allow deposite and withdrawl. This function is not fully transactional (transactions are not specified by assignment requirements).

The entire service is designed to work with `denomination` type.

The interface for register is generic so it can support any type of denomination. The benefit is no dependency on a single money (or denomination) type. A generic interface doesn't help much with modelling generic behaviour, rather it is useful to put general constraint around what a specific behaviour.

The transaction flow is:

1. Client makes a POST call to the server with denominations that they insert into the machine.
2. This returns a transaction with
