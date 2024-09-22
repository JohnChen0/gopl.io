// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"sync"
	"testing"

	"gopl.io/ch9/bank3"
)

func TestBank(t *testing.T) {
	// A gorountine to check the balance.
	go func() {
		prevBalance := 0
		for {
			balance := bank.Balance()
			if balance < 0 || balance > (1000+1)*1000/2 {
				t.Errorf("Bad balance %d", balance)
			}
			if balance < prevBalance {
				t.Errorf("Balance went down from %d to %d",
				         prevBalance, balance)
			}
			prevBalance = balance
		}
	}()

	// Deposit [1..1000] concurrently.
	var n sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		n.Add(1)
		go func(amount int) {
			bank.Deposit(amount)
			n.Done()
		}(i)
	}
	n.Wait()

	if got, want := bank.Balance(), (1000+1)*1000/2; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
