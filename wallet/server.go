package wallet

import (
	"crypto.com/pkg"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	// r.Use(gin.Logger())
	r.Use(gin.Recovery())

	c := NewController()

	api := r.Group("/api")
	api.POST("/deposit", c.Deposit)
	api.POST("/withdraw", c.Withdraw)
	api.POST("/transfer", c.Transfer)
	api.POST("/balance", c.Balance)
	api.POST("/history", c.History)

	return r
}

func NewController() *walletController {
	return &walletController{}
}

// Deposit Deposit to specify user wallet
func (w *walletController) Deposit(ctx *gin.Context) {
	p := &ParamDeposit{}
	err := ctx.ShouldBindJSON(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	result := ResponseBalance{
		UserName: p.UserName,
		Balance:  bank.Get(p.UserName).Deposit(p.Money),
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// Withdraw Withdraw from specify user wallet
func (w *walletController) Withdraw(ctx *gin.Context) {
	p := &ParamWithdraw{}
	err := ctx.ShouldBindJSON(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	balance, ok := bank.Get(p.UserName).Withdraw(p.Money)
	if !ok {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails("not enougth money"))
		return
	}

	result := ResponseBalance{
		UserName: p.UserName,
		Balance:  balance,
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// Transfer Transfer from one user to another user
func (w *walletController) Transfer(ctx *gin.Context) {
	p := &ParamTransfer{}
	err := ctx.ShouldBindJSON(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	sender := bank.Get(p.UserName)
	receiver := bank.Get(p.Receiver.UserName)
	senderBalance, receiverBalance, ok := sender.Transfer(receiver, p.Money)
	if !ok {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails("not enougth money"))
		return
	}

	result := ResponseTransfer{
		Sender:   ResponseBalance{UserName: sender.UserName, Balance: senderBalance},
		Receiver: ResponseBalance{UserName: receiver.UserName, Balance: receiverBalance},
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// Balance Get specify user balance
func (w *walletController) Balance(ctx *gin.Context) {
	p := &ParamUser{}
	err := ctx.ShouldBindJSON(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	result := ResponseBalance{
		UserName: p.UserName,
		Balance:  bank.Get(p.UserName).Balance(),
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// History Get specify user transaction history
func (w *walletController) History(ctx *gin.Context) {
	p := &ParamUser{}
	err := ctx.ShouldBindJSON(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	result := ResponseHistory{
		List: bank.Get(p.UserName).History(),
	}

	pkg.NewResponse(ctx).ToResponse(result)
}
