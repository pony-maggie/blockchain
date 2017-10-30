package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var goblockchain = NewBlockchain()

// get ip:port/mine
func MineHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("mine handler")

	type Block_ST struct {
		Index         int           `json:index`
		Message       string        `json:message`
		Transactions  []transaction `json:transactions`
		Proof         int           `json:proof`
		Previous_hash string        `json:previous_hash`
	}

	result := Block_ST{}

	var code int

	code = http.StatusBadRequest
	result.Message = "unsupport"

	req.ParseForm()

	if req.Method == "GET" {

		last_block := goblockchain.last_block()
		last_proof := last_block.proof
		previous_hash := last_block.previous_hash

		fmt.Printf("last_proof:%d\n", last_proof)
		fmt.Printf("previous_hash:%s\n", previous_hash)

		proof := goblockchain.proof_of_work(last_proof)

		var trans_reward transaction
		trans_reward.Sender = "0"
		trans_reward.Recipient = "random address"
		trans_reward.Amount = 1

		block := goblockchain.New_Block(proof, previous_hash)

		code = http.StatusOK

		result.Message = "New Block Forged"
		result.Index = block.index
		result.Proof = block.proof
		result.Transactions = block.transactions
		result.Previous_hash = block.previous_hash
	}

	bytes, _ := json.Marshal(result)
	w.WriteHeader(code)
	fmt.Fprintf(w, string(bytes))

}

//post ip:port/transactions/new
func NewTransactionHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("new transactions handler")

	type Transaction_ST struct {
		Sender    string `json:sender`
		Recipient string `json:recipient`
		Amount    int    `json:amount`
	}

	type ResponseJsonBean struct {
		Message string `json:"message"`
	}

	transaction := Transaction_ST{}
	result := ResponseJsonBean{}

	var code int

	code = http.StatusBadRequest
	result.Message = "unsupport"

	req.ParseForm()

	b, _ := ioutil.ReadAll(req.Body)

	fmt.Printf("body:%s\n", b)

	if req.Method == "POST" {

		err := json.Unmarshal([]byte(b), &transaction)

		if err != nil {
			fmt.Println(err.Error())

			code = http.StatusInternalServerError
			result.Message = "json unparse failed"

		} else {

			fmt.Printf("%+v\n", transaction)

			index := goblockchain.New_Transaction(transaction.Sender, transaction.Recipient, transaction.Amount)

			code = http.StatusCreated
			result.Message = fmt.Sprintf("New nodes have been added to Block %d", index)

		}

	}

	bytes, _ := json.Marshal(result)
	w.WriteHeader(code)
	fmt.Fprintf(w, string(bytes))

}

//get ip:port/chain
func ChainHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("chain transaction handler")

}

//post ip:port/nodes/register
func RegisterNodesHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("nodes register handler")

	type Nodes_ST struct {
		Nodes []string
	}

	type ResponseJsonBean struct {
		Message string   `json:"message"`
		Data    []string `json:"total_nodes"`
	}

	nodeGroup := Nodes_ST{}
	result := ResponseJsonBean{}

	var code int

	code = http.StatusBadRequest
	result.Message = "unsupport"

	req.ParseForm()

	b, _ := ioutil.ReadAll(req.Body)

	fmt.Printf("body:%s\n", b)

	if req.Method == "POST" {

		err := json.Unmarshal([]byte(b), &nodeGroup)

		if err != nil {
			fmt.Println(err.Error())

			code = http.StatusInternalServerError
			result.Message = "json unparse failed"

		} else {

			result.Data = make([]string, 0)

			fmt.Printf("%+v\n", nodeGroup)

			for _, node := range nodeGroup.Nodes {

				fmt.Printf("node:%s\n", node)
				goblockchain.Register_Node(node)
				result.Data = append(result.Data, node)
			}

			code = http.StatusCreated
			result.Message = "New nodes have been added"

		}

	}

	bytes, _ := json.Marshal(result)
	w.WriteHeader(code)
	fmt.Fprintf(w, string(bytes))

}

//get ip:port/nodes/resolve
func ConsensusHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("nodes resolve register")

}

func main() {

	port := flag.String("port", "8001", "use -port <port number>")
	flag.Parse()

	fmt.Printf("port is:%s\n", *port)

	http.HandleFunc("/mine", MineHandler)
	http.HandleFunc("/transactions/new", NewTransactionHandler)
	http.HandleFunc("/nodes/register", RegisterNodesHandler)
	http.HandleFunc("/nodes/resolve", ConsensusHandler)
	http.HandleFunc("/chain", ChainHandler)

	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
