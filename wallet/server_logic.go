package wallet

import (
	"strings"
	"time"
)

var bank = Bank{}

func init() {
	bank.users = make(map[string]*User)
}

// Get get user by name
func (b *Bank) Get(UserName string) *User {
	if strings.TrimSpace(UserName) == "" {
		return nil
	}

	b.mu.Lock()
	u, ok := b.users[UserName]
	if !ok {
		b.users[UserName] = &User{UserName: UserName}
		u = b.users[UserName]
	}
	b.mu.Unlock()

	return u
}

// Balance Get specify user balance
func (u *User) Balance() uint {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.balance
}

// History Get specify user transaction history
func (u *User) History() []Transaction {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.history
}

// Deposit Deposit to specify user wallet
func (u *User) Deposit(money uint) uint {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.balance += money
	u.history = append(u.history, Transaction{
		Receiver:   u.UserName,
		Action:     ActionDeposit,
		Money:      money,
		Balance:    u.balance,
		CreateDate: time.Now().Local().Unix(),
	})
	return u.balance
}

// Withdraw Withdraw from specify user wallet
func (u *User) Withdraw(money uint) (uint, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.balance < money {
		return u.balance, false
	}

	u.balance -= money

	u.history = append(u.history, Transaction{
		Sender:     u.UserName,
		Action:     ActionWithdraw,
		Money:      money,
		Balance:    u.balance,
		CreateDate: time.Now().Local().Unix(),
	})

	return u.balance, true
}

// Transfer Transfer from one user to another user
func (u *User) Transfer(receiver *User, money uint) (uint, uint, bool) {
	bank.mu.Lock()
	defer bank.mu.Unlock()

	if u.balance < money {
		return u.balance, receiver.balance, false
	}

	u.balance -= money
	receiver.balance += money

	now := time.Now().Local().Unix()
	u.history = append(u.history, Transaction{
		Sender:     u.UserName,
		Receiver:   receiver.UserName,
		Action:     ActionTransfer,
		Money:      money,
		Balance:    u.balance,
		CreateDate: now,
	})
	receiver.history = append(receiver.history, Transaction{
		Sender:     u.UserName,
		Receiver:   receiver.UserName,
		Action:     ActionTransfer,
		Money:      money,
		Balance:    receiver.balance,
		CreateDate: now,
	})

	return u.balance, receiver.balance, true
}

// Reset only for testing to clear data
func (u *User) Reset() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.balance = 0
	u.history = []Transaction{}
}
