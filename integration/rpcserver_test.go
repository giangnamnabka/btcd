// Copyright (c) 2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// This file is ignored during the regular tests due to the following build tag.
// +build rpctest

package integration

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"github.com/giangnamnabka/btcd/chaincfg"
	"github.com/giangnamnabka/btcd/chaincfg/chainhash"
	"github.com/giangnamnabka/btcd/integration/rpctest"
	"github.com/giangnamnabka/btcd/rpcclient"
)

func testGetBestBlock(r *rpctest.Harness, t *testing.T) {
	_, prevbestHeight, err := r.Client.GetBestBlock()
	if err != nil {
		t.Fatalf("Call to `getbestblock` failed: %v", err)
	}

	// Create a new block connecting to the current tip.
	generatedBlockHashes, err := r.Client.Generate(1)
	if err != nil {
		t.Fatalf("Unable to generate block: %v", err)
	}

	bestHash, bestHeight, err := r.Client.GetBestBlock()
	if err != nil {
		t.Fatalf("Call to `getbestblock` failed: %v", err)
	}

	// Hash should be the same as the newly submitted block.
	if !bytes.Equal(bestHash[:], generatedBlockHashes[0][:]) {
		t.Fatalf("Block hashes do not match. Returned hash %v, wanted "+
			"hash %v", bestHash, generatedBlockHashes[0][:])
	}

	// Block height should now reflect newest height.
	if bestHeight != prevbestHeight+1 {
		t.Fatalf("Block heights do not match. Got %v, wanted %v",
			bestHeight, prevbestHeight+1)
	}
}

func testGetBlockCount(r *rpctest.Harness, t *testing.T) {
	// Save the current count.
	currentCount, err := r.Client.GetBlockCount()
	if err != nil {
		t.Fatalf("Unable to get block count: %v", err)
	}

	if _, err := r.Client.Generate(1); err != nil {
		t.Fatalf("Unable to generate block: %v", err)
	}

	// Count should have increased by one.
	newCount, err := r.Client.GetBlockCount()
	if err != nil {
		t.Fatalf("Unable to get block count: %v", err)
	}
	if newCount != currentCount+1 {
		t.Fatalf("Block count incorrect. Got %v should be %v",
			newCount, currentCount+1)
	}
}

func testGetBlockHash(r *rpctest.Harness, t *testing.T) {
	// Create a new block connecting to the current tip.
	generatedBlockHashes, err := r.Client.Generate(1)
	if err != nil {
		t.Fatalf("Unable to generate block: %v", err)
	}

	info, err := r.Client.GetInfo()
	if err != nil {
		t.Fatalf("call to getinfo cailed: %v", err)
	}

	blockHash, err := r.Client.GetBlockHash(int64(info.Blocks))
	if err != nil {
		t.Fatalf("Call to `getblockhash` failed: %v", err)
	}

	// Block hashes should match newly created block.
	if !bytes.Equal(generatedBlockHashes[0][:], blockHash[:]) {
		t.Fatalf("Block hashes do not match. Returned hash %v, wanted "+
			"hash %v", blockHash, generatedBlockHashes[0][:])
	}
}

func testBulkClient(r *rpctest.Harness, t *testing.T) {
	// Create a new block connecting to the current tip.
	generatedBlockHashes, err := r.Client.Generate(20)
	if err != nil {
		t.Fatalf("Unable to generate block: %v", err)
	}

	var futureBlockResults []rpcclient.FutureGetBlockResult
	for _, hash := range generatedBlockHashes {
		futureBlockResults = append(futureBlockResults, r.BatchClient.GetBlockAsync(hash))
	}

	err = r.BatchClient.Send()
	if err != nil {
		t.Fatal(err)
	}

	isKnownBlockHash := func(blockHash chainhash.Hash) bool {
		for _, hash := range generatedBlockHashes {
			if blockHash.IsEqual(hash) {
				return true
			}
		}
		return false
	}

	for _, block := range futureBlockResults {
		msgBlock, err := block.Receive()
		if err != nil {
			t.Fatal(err)
		}
		blockHash := msgBlock.Header.BlockHash()
		if !isKnownBlockHash(blockHash) {
			t.Fatalf("expected hash %s  to be in generated hash list", blockHash)
		}
	}

}

var rpcTestCases = []rpctest.HarnessTestCase{
	testGetBestBlock,
	testGetBlockCount,
	testGetBlockHash,
	testBulkClient,
}

var primaryHarness *rpctest.Harness

func TestMain(m *testing.M) {
	var err error

	// In order to properly test scenarios on as if we were on mainnet,
	// ensure that non-standard transactions aren't accepted into the
	// mempool or relayed.
	btcdCfg := []string{"--rejectnonstd"}
	primaryHarness, err = rpctest.New(
		&chaincfg.SimNetParams, nil, btcdCfg, "",
	)
	if err != nil {
		fmt.Println("unable to create primary harness: ", err)
		os.Exit(1)
	}

	// Initialize the primary mining node with a chain of length 125,
	// providing 25 mature coinbases to allow spending from for testing
	// purposes.
	if err := primaryHarness.SetUp(true, 25); err != nil {
		fmt.Println("unable to setup test chain: ", err)

		// Even though the harness was not fully setup, it still needs
		// to be torn down to ensure all resources such as temp
		// directories are cleaned up.  The error is intentionally
		// ignored since this is already an error path and nothing else
		// could be done about it anyways.
		_ = primaryHarness.TearDown()
		os.Exit(1)
	}

	exitCode := m.Run()

	// Clean up any active harnesses that are still currently running.This
	// includes removing all temporary directories, and shutting down any
	// created processes.
	if err := rpctest.TearDownAll(); err != nil {
		fmt.Println("unable to tear down all harnesses: ", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func TestRpcServer(t *testing.T) {
	var currentTestNum int
	defer func() {
		// If one of the integration tests caused a panic within the main
		// goroutine, then tear down all the harnesses in order to avoid
		// any leaked btcd processes.
		if r := recover(); r != nil {
			fmt.Println("recovering from test panic: ", r)
			if err := rpctest.TearDownAll(); err != nil {
				fmt.Println("unable to tear down all harnesses: ", err)
			}
			t.Fatalf("test #%v panicked: %s", currentTestNum, debug.Stack())
		}
	}()

	for _, testCase := range rpcTestCases {
		testCase(primaryHarness, t)

		currentTestNum++
	}
}
