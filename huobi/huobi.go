package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Account struct {
	ID     int64
	Type   string
	State  string
	UserID int64
}

type Huobi struct {
	accessKey    string
	secretKey    string
	tradeAccount Account
}

func (h *Huobi) GetExchangeName() string {
	return "HuobiPro"
}

// 查询当前用户的所有账户, 根据包含的私钥查询
// return: AccountsReturn对象
func (h *Huobi) GetAccounts() AccountsReturn {
	accountsReturn := AccountsReturn{}

	strRequest := "/v1/account/accounts"
	jsonAccountsReturn := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonAccountsReturn), &accountsReturn)

	return accountsReturn
}

// 根据账户ID查询账户余额
// return: BalanceReturn对象
func (h *Huobi) GetAccountBalance() (*Balance, error) {
	balanceReturn := BalanceReturn{}

	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", h.tradeAccount.ID)
	jsonBanlanceReturn := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)
	if balanceReturn.Status != "ok" {
		return nil, errors.New(balanceReturn.ErrMsg)
	}

	return &balanceReturn.Data, nil
}

// 下单
// placeRequestParams: 下单信息
// return: PlaceReturn对象
func (h *Huobi) Place(amount, price float64, symbol, typ string) (string, error) {
	placeReturn := PlaceReturn{}
	var placeRequestParams PlaceRequestParams
	placeRequestParams.AccountID = strconv.FormatInt(h.tradeAccount.ID, 10)
	placeRequestParams.Amount = strconv.FormatFloat(amount, 'f', -1, 64)
	placeRequestParams.Price = strconv.FormatFloat(price, 'f', -1, 64)
	placeRequestParams.Source = "api"
	placeRequestParams.Symbol = symbol
	placeRequestParams.Type = typ

	mapParams := make(map[string]string)
	mapParams["account-id"] = placeRequestParams.AccountID
	mapParams["amount"] = placeRequestParams.Amount
	if 0 < len(placeRequestParams.Price) {
		mapParams["price"] = placeRequestParams.Price
	}
	if 0 < len(placeRequestParams.Source) {
		mapParams["source"] = placeRequestParams.Source
	}
	mapParams["symbol"] = placeRequestParams.Symbol
	mapParams["type"] = placeRequestParams.Type

	strRequest := "/v1/order/orders/place"
	jsonPlaceReturn := apiKeyPost(mapParams, strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	if placeReturn.Status != "ok" {
		return "", errors.New(placeReturn.ErrMsg)
	}

	return placeReturn.Data, nil

}

// 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func (h *Huobi) SubmitCancel(strOrderID string) error {
	placeReturn := PlaceReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", strOrderID)
	jsonPlaceReturn := apiKeyPost(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	if placeReturn.Status != "ok" {
		return errors.New(placeReturn.ErrMsg)
	}

	return nil
}

// 查询订单详情
// strOrderID: 订单ID
// return: OrderReturn对象
func (h *Huobi) GetOrderInfo(strOrderID string) (*Order, error) {
	orderReturn := OrderReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s", strOrderID)
	jsonPlaceReturn := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &orderReturn)

	if orderReturn.Status != "ok" {
		return nil, errors.New(orderReturn.ErrMsg)
	}

	return &orderReturn.Data, nil

}

func (h *Huobi) GetOrders(params OrdersRequestParams) ([]Order, error) {
	ordersReturn := OrdersReturn{}

	jsonP, _ := json.Marshal(params)

	var paramMap = make(map[string]string)
	json.Unmarshal(jsonP, &paramMap)

	strRequest := "/v1/order/orders"
	ret := apiKeyGet(paramMap, strRequest, h.accessKey, h.secretKey)
	json.Unmarshal([]byte(ret), &ordersReturn)
	if ordersReturn.Status != "ok" {
		return nil, errors.New(ordersReturn.ErrMsg)
	}

	return ordersReturn.Data, nil

}

func NewHuobi(accesskey, secretkey string) (*Huobi, error) {
	h := &Huobi{
		accessKey: accesskey,
		secretKey: secretkey,
	}

	fmt.Println("init huobi.")
	ret := h.GetAccounts()
	if ret.Status != "ok" {
		return nil, errors.New(ret.ErrMsg)
	}

	for _, account := range ret.Data {
		if account.Type == "spot" {
			fmt.Println("account id:", account.ID)
			h.tradeAccount.ID = account.ID
			h.tradeAccount.Type = account.Type
			h.tradeAccount.State = account.State
			h.tradeAccount.UserID = account.UserID
			break
		}

	}

	fmt.Println("init huobi success.")

	return h, nil
}
