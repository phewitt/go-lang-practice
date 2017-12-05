package funding

import (
	"sync"
	"testing"
)

const WORKERS = 10

func BenchmarkFund(b *testing.B) {
	// Skip N = 1
	if b.N < WORKERS {
		return
	}

	fund := NewFund(b.N)
	dollarsPerWorker := b.N / WORKERS

	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for i := 0; i < dollarsPerWorker; i++ {
				fund.Withdraw(1)
			}
		}()
	}

	wg.Wait()

	if fund.Balance() != 0 {
		b.Error("Balance wasn't zero:", fund.Balance())
	}
}

func BenchmarkWithdrawals(b *testing.B) {
	// Skip N = 1
	if b.N < WORKERS {
		return
	}
	server := NewFundServer(b.N)
	dollarsPerWorker := b.N / WORKERS
	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < dollarsPerWorker; i++ {
				server.Commands <- WithdrawCommand{Amount: 1}
			}
		}()
	}

	wg.Wait()

	balanceResponseChan := make(chan int)
	server.Commands <- BalanceCommand{Response: balanceResponseChan}
	balance := <-balanceResponseChan

	if balance != 0 {
		b.Error("Balance wasn't zero:", balance)
	}
}
