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

type ParamUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}

type ParamDeposit struct {
	ParamUser
	Money uint `json:"money" binding:"required,min=1,max=1000000000"`
}

type ParamWithdraw struct {
	ParamDeposit
}

type ParamTransfer struct {
	ParamDeposit
	Receiver string `json:"receiver" binding:"required"`
}

type ResponseTransfer struct {
	Sender   ResponseBalance `json:"sender"`
	Receiver ResponseBalance `json:"receiver"`
}

type ResponseBalance struct {
	UserName string `json:"user_name"`
	Balance  uint   `json:"balance"`
}

type ResponseHistory struct {
	List []Transaction `json:"list"`
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
