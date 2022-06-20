
# blockchain


## Background

got the idea from

[Learn Blockchains by Building One](https://github.com/dvf/blockchain)

## Building


```
$ cd cmd
$ go build -o goblockchain
```


## Usage

start instances(nodes) as follows

You can start as many nodes as you want with the following command

```
./goblockchain -port 8001
./goblockchain -port 8002
./goblockchain -port 8003
```

## Json Endpoints


get full blockchain

* `curl 127.0.0.1:8001/chain`

mine a new block

* `curl 127.0.0.1:8001/mine`

Adding a new transaction

* `POST 127.0.0.1:8001/transactions/new`

* __Body__:

```json
  {
    "sender": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    "recipient": "1Ez69SnzzmePmZX3WpEzMKTrcBF2gpNQ55",
    "amount": 1000
  }
```


Register a node in the blockchain network

* `POST 127.0.0.1:8001/nodes/register`

* __Body__:

```json
  {
    "nodes": [
        "http://127.0.0.1:8002",
        "http://127.0.0.1:8003"
    ]
}
```

Resolving Blockchain

* `curl 127.0.0.1:8001/nodes/resolve`