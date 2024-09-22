// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"fmt"
	"testing"

	"gopl.io/ch9/bank1"
)

func TestBank(t *testing.T) {
	aliceBalances := make(map[int]int)
	done := make(chan struct{})
	for i := 0; i < 1000; i++ {
		aliceBalance, finalBalance := runTest(done)
		if aliceBalance != 200 && aliceBalance != 300 {
			t.Errorf("Unexpected balance %d for Alice", aliceBalance)
		}
		if finalBalance != 300 {
			t.Errorf("Unexpected final balance %d", finalBalance)
		}
		aliceBalances[aliceBalance]++
		bank.Deposit(-300)	// Reset balance to 0
	}
	for aliceBalance, count := range aliceBalances {
		fmt.Println("Alice see balance", aliceBalance, count, "times")
	}
}

func runTest(done chan struct{}) (int, int) {
	var aliceBalance int

	// Alice
	go func() {
		bank.Deposit(200)
		aliceBalance = bank.Balance()
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	return aliceBalance, bank.Balance()
}
