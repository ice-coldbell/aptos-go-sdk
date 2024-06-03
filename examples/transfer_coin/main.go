package main

import (
	"fmt"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
)

const FundAmount = 100_000_000
const TransferAmount = 1_000

// main This example shows you how to make an APT transfer transaction in the simplest possible way
func main() {
	// Create a client for Aptos
	client, err := aptos.NewClient(aptos.DevnetConfig)
	if err != nil {
		panic("Failed to create client:" + err.Error())
	}

	// Create accounts locally for alice and bob
	alice, err := aptos.NewEd25519Account()
	if err != nil {
		panic("Failed to create alice:" + err.Error())
	}
	bob, err := aptos.NewEd25519Account()
	if err != nil {
		panic("Failed to create bob:" + err.Error())
	}

	fmt.Printf("\n=== Addresses ===\n")
	fmt.Printf("Alice: %s\n", alice.Address.String())
	fmt.Printf("Bob:%s\n", bob.Address.String())

	// Fund the sender with the faucet to create it on-chain
	err = client.Fund(alice.Address, FundAmount)
	if err != nil {
		panic("Failed to fund alice:" + err.Error())
	}

	aliceBalance, err := client.AccountAPTBalance(&alice.Address)
	if err != nil {
		panic("Failed to retrieve alice balance:" + err.Error())
	}
	bobBalance, err := client.AccountAPTBalance(&bob.Address)
	if err != nil {
		panic("Failed to retrieve bob balance:" + err.Error())
	}
	fmt.Printf("\n=== Initial Balances ===\n")
	fmt.Printf("Alice: %d\n", aliceBalance)
	fmt.Printf("Bob:%d\n", bobBalance)

	// Sign transaction
	signedTxn, err := aptos.APTTransferTransaction(client, alice, bob.Address, TransferAmount)
	if err != nil {
		panic("Failed to sign transaction:" + err.Error())
	}

	// Submit and wait for it to complete
	submitResult, err := client.SubmitTransaction(signedTxn)
	if err != nil {
		panic("Failed to submit transaction:" + err.Error())
	}
	txnHash := submitResult.Hash

	// Wait for the transaction
	_, err = client.WaitForTransaction(txnHash)
	if err != nil {
		panic("Failed to wait for transaction:" + err.Error())
	}
	aliceBalance, err = client.AccountAPTBalance(&alice.Address)
	if err != nil {
		panic("Failed to retrieve alice balance:" + err.Error())
	}
	bobBalance, err = client.AccountAPTBalance(&bob.Address)
	if err != nil {
		panic("Failed to retrieve bob balance:" + err.Error())
	}
	fmt.Printf("\n=== Intermediate Balances ===\n")
	fmt.Printf("Alice: %d\n", aliceBalance)
	fmt.Printf("Bob:%d\n", bobBalance)

	// Now submit as a single call, with a custom payload
	accountBytes, err := bcs.Serialize(&bob.Address)
	if err != nil {
		panic("Failed to serialize bob's address:" + err.Error())
	}

	serializer := &bcs.Serializer{}
	serializer.U64(TransferAmount)
	err = serializer.Error()
	if err != nil {
		panic("Failed to serialize transfer amount:" + err.Error())
	}
	amountBytes := serializer.ToBytes()

	resp, err := client.BuildSignAndSubmitTransaction(alice, aptos.TransactionPayload{
		Payload: &aptos.EntryFunction{
			Module: aptos.ModuleId{
				Address: aptos.AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer",
			ArgTypes: []aptos.TypeTag{},
			Args: [][]byte{
				accountBytes,
				amountBytes,
			},
		}},
	)
	if err != nil {
		panic("Failed to sign transaction:" + err.Error())
	}

	_, err = client.WaitForTransaction(resp.Hash)
	if err != nil {
		panic("Failed to wait for transaction:" + err.Error())
	}

	aliceBalance, err = client.AccountAPTBalance(&alice.Address)
	if err != nil {
		panic("Failed to retrieve alice balance:" + err.Error())
	}
	bobBalance, err = client.AccountAPTBalance(&bob.Address)
	if err != nil {
		panic("Failed to retrieve bob balance:" + err.Error())
	}
	fmt.Printf("\n=== Final Balances ===\n")
	fmt.Printf("Alice: %d\n", aliceBalance)
	fmt.Printf("Bob:%d\n", bobBalance)
}