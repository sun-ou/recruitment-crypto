package wallet

import (
	"time"

	"crypto.com/pkg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	// r.Use(gin.Logger())
	r.Use(gin.Recovery())

	c := NewController()

	api := r.Group("/api")
	api.POST("/deposit", c.Deposit)
	api.POST("/withdraw", c.Withdraw)
	api.POST("/transfer", c.Transfer)
	api.GET("/balance", c.Balance)
	api.GET("/history", c.History)
	api.GET("/health", c.Health)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("positive_decimal", IsPositiveDecimal)
	}

	return r
}

func NewController() *walletController {
	return &walletController{}
}

// Deposit Deposit to specify user wallet
func (w *walletController) Deposit(ctx *gin.Context) {
	p := &ParamDeposit{}
	err := ctx.ShouldBind(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	result := ResponseBalance{
		UserName: p.UserName,
		Balance:  Cent2String(NewUser().Get(p.UserName).Deposit(String2Cent(p.Money))),
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// Withdraw Withdraw from specify user wallet
func (w *walletController) Withdraw(ctx *gin.Context) {
	p := &ParamWithdraw{}
	err := ctx.ShouldBind(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	balance, ok := NewUser().Get(p.UserName).Withdraw(String2Cent(p.Money))
	if !ok {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails("not enougth money"))
		return
	}

	result := ResponseBalance{
		UserName: p.UserName,
		Balance:  Cent2String(balance),
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// Transfer Transfer from one user to another user
func (w *walletController) Transfer(ctx *gin.Context) {
	p := &ParamTransfer{}
	err := ctx.ShouldBind(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	sender := NewUser().Get(p.UserName)
	receiver := NewUser().Get(p.Receiver)
	senderBalance, receiverBalance, ok := sender.Transfer(receiver, String2Cent(p.Money))
	if !ok {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails("not enougth money"))
		return
	}

	result := ResponseTransfer{
		Sender:   ResponseBalance{UserName: sender.UserName, Balance: Cent2String(senderBalance)},
		Receiver: ResponseBalance{UserName: receiver.UserName, Balance: Cent2String(receiverBalance)},
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// Balance Get specify user balance
func (w *walletController) Balance(ctx *gin.Context) {
	p := &ParamUser{}
	err := ctx.ShouldBind(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	result := ResponseBalance{
		UserName: p.UserName,
		Balance:  Cent2String(NewUser().Get(p.UserName).Balance()),
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// History Get specify user transaction history
func (w *walletController) History(ctx *gin.Context) {
	p := &ParamUser{}
	err := ctx.ShouldBind(p)
	if err != nil {
		pkg.NewResponse(ctx).ToErrorResponse(pkg.InvaildParams.WitchDetails(err.Error()))
		return
	}

	result := ResponseHistory{
		List: w.FormatHistory(NewUser().Get(p.UserName).History()),
	}

	pkg.NewResponse(ctx).ToResponse(result)
}

// Health health check
func (w *walletController) Health(ctx *gin.Context) {
	pkg.NewResponse(ctx).ToResponse(`ok`)
}

// FormatHistory format history list
func (w *walletController) FormatHistory(originHistory []Transaction) (formatHistory []FormatTransaction) {
	idMap := make(map[uint]string)
	for _, v := range originHistory {
		if _, ok := idMap[v.SenderId]; !ok {
			idMap[v.SenderId] = ""
		}
		if _, ok := idMap[v.ReceiverId]; !ok {
			idMap[v.ReceiverId] = ""
		}
	}

	NewUser().GetName(idMap)

	for _, v := range originHistory {
		vv := FormatTransaction{
			Id:         v.Id,
			SenderId:   v.SenderId,
			ReceiverId: v.ReceiverId,
			Money:      Cent2String(v.Money),
			Balance:    Cent2String(v.Balance),
			Action:     v.Action,
			CreateDate: time.Unix(v.CreateDate, 0).Format("2006-01-02 15:04:05"),
		}
		if _, ok := idMap[v.SenderId]; ok {
			vv.Sender = idMap[v.SenderId]
		}
		if _, ok := idMap[v.ReceiverId]; ok {
			vv.Receiver = idMap[v.ReceiverId]
		}

		switch v.Action {
		case ActionDeposit:
			vv.Receiver = vv.Sender
			vv.ReceiverId = vv.SenderId
			vv.Sender = ""
			vv.SenderId = 0
		case ActionWithdraw, ActionSend:
			// no need to process
		case ActionReceive:
			vv.Sender, vv.SenderId, vv.Receiver, vv.ReceiverId = vv.Receiver, vv.ReceiverId, vv.Sender, vv.SenderId
		}

		formatHistory = append(formatHistory, vv)
	}
	return
}
