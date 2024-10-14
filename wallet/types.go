package wallet

import "sync"

type Bank struct {
	users map[string]*User
	mu    sync.Mutex
}

type User struct {
	UserName string
	balance  uint
	history  []Transaction
	mu       sync.RWMutex
}

type Transaction struct {
	Sender     string
	Receiver   string
	Money      uint
	Balance    uint
	Action     string
	CreateDate int64
}

type walletController struct{}

type wallectClient struct {
	ServerAddress string
}

type ParamUser struct {
	UserName string `form:"user_name" binding:"required"`
}

type ParamDeposit struct {
	ParamUser
	Money uint `form:"money" binding:"required,min=1,max=1000000000"`
}

type ParamWithdraw struct {
	ParamDeposit
}

type ParamTransfer struct {
	ParamDeposit
	Receiver ParamUser
}

type ResponseTransfer struct {
	Sender   ResponseBalance
	Receiver ResponseBalance
}

type ResponseBalance struct {
	UserName string
	Balance  uint
}

type ResponseHistory struct {
	List []Transaction
}

type ResponseError struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"`
}

const (
	ActionDeposit  = "Desposit"
	ActionWithdraw = "Withdraw"
	ActionTransfer = "Transfer"
)
