package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func NewClient(ServerAddress string) *wallectClient {
	return &wallectClient{ServerAddress: ServerAddress}
}

func (w *wallectClient) Post(action string, data interface{}) (code int, body []byte, err error) {
	url := "http://" + w.ServerAddress + action
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) // nolint
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	return resp.StatusCode, body, err
}

func (w *wallectClient) Deposit(p *ParamDeposit) (result *ResponseBalance, err error) {
	var body []byte
	var code int
	code, body, err = w.Post("/api/deposit", p)
	if err == nil {
		if code == http.StatusOK {
			_ = json.Unmarshal(body, &result)
		} else {
			respError := ResponseError{}
			_ = json.Unmarshal(body, &respError)
			err = fmt.Errorf("err code: %d, msg: %s, details: %v", respError.Code, respError.Msg, respError.Details)
		}
	}
	return
}

func (w *wallectClient) Withdraw(p *ParamWithdraw) (result *ResponseBalance, err error) {
	var body []byte
	var code int
	code, body, err = w.Post("/api/withdraw", p)
	if err == nil {
		if code == http.StatusOK {
			_ = json.Unmarshal(body, &result)
		} else {
			respError := ResponseError{}
			json.Unmarshal(body, &respError)
			err = fmt.Errorf("err code: %d, msg: %s, details: %v", respError.Code, respError.Msg, respError.Details)
		}
	}
	return
}

func (w *wallectClient) Transfer(p *ParamTransfer) (result *ResponseTransfer, err error) {
	var body []byte
	var code int
	code, body, err = w.Post("/api/transfer", p)
	if err == nil {
		if code == http.StatusOK {
			json.Unmarshal(body, &result)
		} else {
			respError := ResponseError{}
			json.Unmarshal(body, &respError)
			err = fmt.Errorf("err code: %d, msg: %s, details: %v", respError.Code, respError.Msg, respError.Details)
		}
	}
	return
}

func (w *wallectClient) Balance(p *ParamUser) (result *ResponseBalance, err error) {
	var body []byte
	var code int
	code, body, err = w.Post("/api/balance", p)
	if err == nil {
		if code == http.StatusOK {
			json.Unmarshal(body, &result)
		} else {
			respError := ResponseError{}
			json.Unmarshal(body, &respError)
			err = fmt.Errorf("err code: %d, msg: %s, details: %v", respError.Code, respError.Msg, respError.Details)
		}
	}
	return
}

func (w *wallectClient) History(p *ParamUser) (result *ResponseHistory, err error) {
	var body []byte
	var code int
	code, body, err = w.Post("/api/history", p)
	if err == nil {
		if code == http.StatusOK {
			json.Unmarshal(body, &result)
		} else {
			respError := ResponseError{}
			json.Unmarshal(body, &respError)
			err = fmt.Errorf("err code: %d, msg: %s, details: %v", respError.Code, respError.Msg, respError.Details)
		}
	}
	return
}

func String2Cent(param string) uint {
	f, _ := strconv.ParseFloat(param, 64)
	return uint(f * 100)
}

func Cent2String(param uint) string {
	return fmt.Sprintf("%.2f", float64(param)/100)
}
