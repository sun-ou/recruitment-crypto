package wallet

import (
	"testing"

	"go.uber.org/goleak"
	"gotest.tools/v3/assert"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestGetUser(t *testing.T) {
	userA := bank.Get("")
	assert.Assert(t, userA == nil, "blank user should return nil")

	userB := bank.Get("b")
	assert.Assert(t, userB != nil, "cannot find user b")
}

func TestDeposit(t *testing.T) {
	var money uint = 999
	userA := bank.Get("a")
	userA.Reset()
	assert.Equal(t, money, userA.Deposit(money), "unexpected result for deposit")
	assert.Equal(t, money, userA.Balance(), "unexpected result for balance")
}

func TestWithdrawSuccess(t *testing.T) {
	var money uint = 999
	var minuend uint = 888
	userA := bank.Get("a")
	userA.Reset()
	userA.Deposit(money)
	balance, ok := userA.Withdraw(minuend)

	assert.Equal(t, true, ok, "operaton failed")
	assert.Equal(t, money-minuend, balance, "incorrect balance")
	assert.Equal(t, money-minuend, userA.Balance(), "incorrect balance")
}

func TestWithdrawFail(t *testing.T) {
	var money uint = 999
	var minuend uint = 1888
	userA := bank.Get("a")
	userA.Reset()
	userA.Deposit(money)
	balance, ok := userA.Withdraw(minuend)

	assert.Equal(t, false, ok, "operaton failed")
	assert.Equal(t, money, balance, "incorrect balance")
	assert.Equal(t, money, userA.Balance(), "incorrect balance")
}

func TestTransferSuccess(t *testing.T) {
	var amount uint = 999
	var money uint = 888
	userA := bank.Get("a")
	userA.Reset()
	userB := bank.Get("b")
	userB.Reset()

	userA.Deposit(amount)
	balanceA, balanceB, ok := userA.Transfer(userB, money)

	assert.Equal(t, true, ok, "unexpected operation")
	assert.Equal(t, amount-money, balanceA, "incorrect balance in user a")
	assert.Equal(t, amount-money, userA.Balance(), "incorrect balance in user a")
	assert.Equal(t, money, balanceB, "incorrect balance in user b")
	assert.Equal(t, money, userB.Balance(), "incorrect balance in user b")
}

func TestTransferFail(t *testing.T) {
	var amount uint = 999
	var money uint = 1888
	userA := bank.Get("a")
	userA.Reset()
	userB := bank.Get("b")
	userB.Reset()

	userA.Deposit(amount)
	balanceA, balanceB, ok := userA.Transfer(userB, money)

	assert.Equal(t, false, ok, "unexpected operation")
	assert.Equal(t, amount, balanceA, "incorrect balance in user a")
	assert.Equal(t, amount, userA.Balance(), "incorrect balance in user a")
	assert.Equal(t, uint(0), balanceB, "incorrect balance in user b")
	assert.Equal(t, uint(0), userB.Balance(), "incorrect balance in user b")
}

func TestHistoryEmpty(t *testing.T) {
	userA := bank.Get("a")
	userA.Reset()
	result := userA.History()
	assert.Equal(t, 0, len(result), "transition list should be empty")
}

func TestHistorySuccess(t *testing.T) {
	var money uint = 999
	userA := bank.Get("a")
	userA.Reset()
	userA.Deposit(money)
	userA.Withdraw(money)
	result := userA.History()
	assert.Equal(t, 2, len(result), "wrong transition list length")
}
