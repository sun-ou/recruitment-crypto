package wallet

type User struct {
	Id       uint
	UserName string
	balance  uint
	history  []Transaction
}

type Transaction struct {
	Id         uint
	SenderId   uint
	ReceiverId uint
	Money      uint
	Balance    uint
	Action     string
	CreateDate int64
}

type FormatTransaction struct {
	Id         uint   `json:"id"`
	SenderId   uint   `json:"sender_id"`
	Sender     string `json:"sender"`
	ReceiverId uint   `json:"receiver_id"`
	Receiver   string `json:"receiver"`
	Money      string `json:"money"`
	Balance    string `json:"balance"`
	Action     string `json:"action"`
	CreateDate string `json:"create_date"`
}

type walletController struct{}

type ParamUser struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
}

type ParamDeposit struct {
	ParamUser
	Money string `json:"money" binding:"required,positive_decimal"`
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
	Balance  string `json:"balance"`
}

type ResponseHistory struct {
	List []FormatTransaction `json:"list"`
}

type ResponseError struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"`
}

const (
	ActionDeposit  = "Desposit"
	ActionWithdraw = "Withdraw"
	ActionSend     = "Send"
	ActionReceive  = "Receive"
)
