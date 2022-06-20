package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    int
}

type block struct {
	Index         int
	timestamp     int64
	Transactions  []Transaction
	Proof         int
	Previous_hash string
}

type Blockchain struct {
	current_transactions []Transaction
	chain                []block
	nodes                []string
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{}

	bc.current_transactions = make([]Transaction, 0)
	bc.chain = make([]block, 0)
	bc.nodes = make([]string, 0)

	bc.New_Block(1, "1")

	return bc
}

func (bc *Blockchain) GetChain() []block {
	return bc.chain
}

func hash(data block) string {

	fmt.Println("hash.....")

	timestamp := []byte(strconv.FormatInt(data.timestamp, 10))
	previous_hash := []byte(data.Previous_hash)

	data_str := bytes.Join([][]byte{previous_hash, timestamp}, []byte{})

	hash := sha256.Sum256(data_str)

	hash_str := string(hash[:])

	fmt.Printf("hash_str:%s \n", hash_str)

	return hash_str

}

func valid_proof(last_proof int, proof int) bool {

	// fmt.Printf("valid proof:%d\n", proof)

	str_last_proorf := []byte(strconv.FormatInt(int64(last_proof), 10))
	str_proof := []byte(strconv.FormatInt(int64(proof), 10))

	str_data := bytes.Join([][]byte{str_last_proorf, str_proof}, []byte{})

	guess_hash := sha256.Sum256(str_data)

	// fmt.Println("guess hash:", guess_hash)

	return bytes.Equal(guess_hash[:2], []byte("00"))
}

func (bc *Blockchain) New_Block(proof int, previous_hash string) *block {

	fmt.Println("new block")

	blockinstance := &block{}

	blockinstance.Index = len(bc.chain) + 1
	blockinstance.timestamp = time.Now().Unix()
	blockinstance.Transactions = bc.current_transactions
	blockinstance.Proof = proof
	blockinstance.Previous_hash = previous_hash

	trans_len := len(bc.current_transactions)

	fmt.Printf("trans_len:%d \n", trans_len)
	fmt.Printf("index:%d \n", blockinstance.Index)

	bc.current_transactions = bc.current_transactions[trans_len:] //clear
	bc.chain = append(bc.chain, *blockinstance)

	return blockinstance
}

func (bc *Blockchain) Proof_of_work(last_proof int) int {

	fmt.Println("proof_of_work")

	proof := 0

	for !(valid_proof(last_proof, proof)) {
		proof++
	}

	return proof
}

func (bc *Blockchain) New_Transaction(sender string, recipient string, amount int) int {

	fmt.Println("new transaction")

	var trans Transaction

	trans.Sender = sender
	trans.Recipient = recipient
	trans.Amount = amount

	bc.current_transactions = append(bc.current_transactions, trans)

	block := bc.Last_block()

	fmt.Printf("index:%d\n", block.Index)

	return block.Index + 1

}

func (bc *Blockchain) Last_block() block {

	fmt.Println("get last block")

	height := len(bc.chain)
	block := bc.chain[height-1]
	return block
}

func (bc *Blockchain) Register_Node(address string) {

	fmt.Println("register node")

	u, err := url.Parse(address)

	if err != nil {
		panic(err)
	}

	bc.nodes = append(bc.nodes, u.Host)

}

func (bc *Blockchain) Valid_chain(chain []block) bool {

	fmt.Println("valid chain")

	last_block := chain[0]
	current_index := 1

	for current_index < len(chain) {

		current_block := chain[current_index]

		fmt.Println("current block:", current_block)
		fmt.Println("last block:", last_block)
		fmt.Println("---------------------")

		if current_block.Previous_hash != hash(last_block) {

			fmt.Println("wo ou, not valid chain")

			return false

		}

		if !(valid_proof(last_block.Proof, current_block.Proof)) {
			return false
		}

		last_block = current_block
		current_index++

	}

	return true

}

func (bc *Blockchain) Resolve_Conflicts() bool {

	fmt.Println("resolve conflicts")

	neighbours := bc.nodes

	var new_chain []block

	// max_length := len(bc.chain)

	for _, node := range neighbours {

		fmt.Println("node:", node)

		url := fmt.Sprintf("http://%s/chain", node)

		fmt.Println("url:", url)

		resp, err := http.Get(url)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("resp:", resp)

		if resp.StatusCode == 200 {

			//TODO
		}

	}

	if new_chain != nil {
		bc.chain = new_chain
		return true
	}

	return false

}
