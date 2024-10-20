package wallet

import (
	"testing"

	"crypto.com/pkg"
	"go.uber.org/goleak"
	"gotest.tools/v3/assert"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestUserReset(t *testing.T) {
	pkg.SetupTestDBEngine()
	NewUser().Get("Alice").Reset()
	NewUser().Get("Bruce").Reset()
	pkg.DBEngine.Close()
}

func TestGetUser(t *testing.T) {
	pkg.SetupTestDBEngine()

	userA := NewUser().Get("")
	assert.Assert(t, userA == nil, "blank user should return nil")

	userB := NewUser().Get("b")
	assert.Assert(t, userB != nil, "cannot find user b")

	pkg.DBEngine.Close()
}

func TestDeposit(t *testing.T) {
	pkg.SetupTestDBEngine()

	var money uint = 999
	userA := NewUser().Get("a")
	userA.Reset()
	assert.Equal(t, money, userA.Deposit(money), "unexpected result for deposit")
	assert.Equal(t, money, userA.Balance(), "unexpected result for balance")

	pkg.DBEngine.Close()
}

func TestWithdrawSuccess(t *testing.T) {
	pkg.SetupTestDBEngine()

	var money uint = 999
	var minuend uint = 888
	userA := NewUser().Get("a")
	userA.Reset()

	// user not exists
	_, ok := userA.Withdraw(minuend)
	assert.Equal(t, false, ok, "unexpected operation")

	userA.Deposit(money) // create user automatically
	balance, ok := userA.Withdraw(minuend)

	assert.Equal(t, true, ok, "operaton failed")
	assert.Equal(t, money-minuend, balance, "incorrect balance")
	assert.Equal(t, money-minuend, userA.Balance(), "incorrect balance")

	pkg.DBEngine.Close()
}

func TestWithdrawFail(t *testing.T) {
	pkg.SetupTestDBEngine()

	var money uint = 999
	var minuend uint = 8999
	userA := NewUser().Get("a")
	userA.Reset()
	userA.Deposit(money)
	balance, ok := userA.Withdraw(minuend)

	assert.Equal(t, false, ok, "operaton failed")
	assert.Equal(t, money, balance, "incorrect balance")
	assert.Equal(t, money, userA.Balance(), "incorrect balance")

	pkg.DBEngine.Close()
}

func TestTransferSuccess(t *testing.T) {
	pkg.SetupTestDBEngine()

	var amount uint = 999
	var money uint = 888
	userA := NewUser().Get("a")
	userA.Reset()
	userB := NewUser().Get("b")
	userB.Reset()

	userA.Deposit(amount)
	balanceA, balanceB, ok := userA.Transfer(userB, money)

	assert.Equal(t, true, ok, "unexpected operation")
	assert.Equal(t, amount-money, balanceA, "incorrect balance in user a")
	assert.Equal(t, amount-money, userA.Balance(), "incorrect balance in user a")
	assert.Equal(t, money, balanceB, "incorrect balance in user b")
	assert.Equal(t, money, userB.Balance(), "incorrect balance in user b")

	// test update statemtnt
	userA.Deposit(amount)
	balanceA, balanceB, ok = userA.Transfer(userB, money)
	assert.Equal(t, true, ok, "unexpected operation")
	assert.Equal(t, amount*2-money*2, balanceA, "incorrect balance in user a")
	assert.Equal(t, amount*2-money*2, userA.Balance(), "incorrect balance in user a")
	assert.Equal(t, money*2, balanceB, "incorrect balance in user b")
	assert.Equal(t, money*2, userB.Balance(), "incorrect balance in user b")

	pkg.DBEngine.Close()
}

func TestTransferFail(t *testing.T) {
	pkg.SetupTestDBEngine()

	var amount uint = 999
	var money uint = 8999
	var balanceA, balanceB uint
	var ok bool
	userA := NewUser().Get("a")
	userB := NewUser().Get("b")
	NewUser().Reset(userA, userB)

	// user not exists
	_, _, ok = userA.Transfer(userB, money)
	assert.Equal(t, false, ok, "unexpected operation")

	userA.Deposit(amount) // create user automatically
	balanceA, balanceB, ok = userA.Transfer(userB, money)
	assert.Equal(t, false, ok, "unexpected operation")
	assert.Assert(t, amount >= balanceA, "incorrect balance in user a")
	assert.Assert(t, amount >= userA.Balance(), "incorrect balance in user a")
	assert.Assert(t, money > balanceB, "incorrect balance in user b")
	assert.Assert(t, money > userB.Balance(), "incorrect balance in user b")

	pkg.DBEngine.Close()
}

func TestHistoryEmpty(t *testing.T) {
	pkg.SetupTestDBEngine()

	userA := NewUser().Get("a")
	userA.Reset()
	result := userA.History()
	assert.Equal(t, 0, len(result), "transition list should be empty")

	pkg.DBEngine.Close()
}

func TestHistorySuccess(t *testing.T) {
	pkg.SetupTestDBEngine()

	var money uint = 999
	userA := NewUser().Get("a")
	userA.Reset()
	userA.Deposit(money)
	userA.Withdraw(money)
	result := userA.History()
	assert.Equal(t, 2, len(result), "wrong transition list length")

	pkg.DBEngine.Close()
}

func TestGetNameSuccess(t *testing.T) {
	pkg.SetupTestDBEngine()

	var money uint = 999
	userA := NewUser().Get("a")
	userB := NewUser().Get("b")
	NewUser().Reset(userA, userB)
	userA.Deposit(money)
	userB.Deposit(money)

	idMap := map[uint]string{userA.Id: "", userB.Id: ""}
	NewUser().GetName(idMap)

	assert.Equal(t, userA.UserName, idMap[userA.Id], "cannot get the name of user a")
	assert.Equal(t, userB.UserName, idMap[userB.Id], "cannot get the name of user b")
	assert.Assert(t, len(idMap) == 2, "wrong name list length")

	// empty map
	idMap = make(map[uint]string)
	NewUser().GetName(idMap)
	assert.Assert(t, len(idMap) == 0, "wrong name list length")

	pkg.DBEngine.Close()
}
