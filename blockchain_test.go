package main

import (
	"fmt"
	"testing"
)

func TestValidNodes(t *testing.T) {

	blockchainIns := NewBlockchain()

	blockchainIns.Register_Node("http://127.0.0.1:8002")

	neighbours := blockchainIns.nodes

	found := false

	for _, node := range neighbours {

		fmt.Println("node:", node)

		if node == "127.0.0.1:8002" {
			found = true
		}
	}

	if !found {
		t.Error("not found added node")
	}

}

func TestCreateBlock(t *testing.T) {

	blockchainIns := NewBlockchain()
	blockchainIns.New_Block(123, "abc")

	last_block := blockchainIns.last_block()

	if len(blockchainIns.chain) != 2 {
		t.Error("error length at genesis block")
	}

	if last_block.index != 2 {
		t.Error("genesis block index must be 2")
	}

	if last_block.proof != 123 {
		t.Error("wrong proof")
	}

	if last_block.previous_hash != "abc" {
		t.Error("wrong previous hash")
	}

}

func TestCreateTransaction(t *testing.T) {

	blockchainIns := NewBlockchain()

	blockchainIns.New_Transaction("aa", "bb", 10)

	height := len(blockchainIns.current_transactions)

	transaction := blockchainIns.current_transactions[height-1]

	if transaction.Sender != "aa" {
		t.Error("wrong sender")
	}

	if transaction.Recipient != "bb" {
		t.Error("wrong recipient")
	}

	if transaction.Amount != 10 {
		t.Error("wrong amount")
	}

}
