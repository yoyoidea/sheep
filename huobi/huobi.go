package huobi

import (
	"encoding/json"
	"fmt"
)

type Huobi struct {
	AccessKey string
	SecretKey string
}

func (h *Huobi) GetExchangeName() string {
	return "HuobiPro"
}

// 查询当前用户的所有账户, 根据包含的私钥查询
// return: AccountsReturn对象
func (h *Huobi) GetAccounts() AccountsReturn {
	accountsReturn := AccountsReturn{}

	strRequest := "/v1/account/accounts"
	jsonAccountsReturn := apiKeyGet(make(map[string]string), strRequest, h.AccessKey, h.SecretKey)
	json.Unmarshal([]byte(jsonAccountsReturn), &accountsReturn)

	return accountsReturn
}

// 根据账户ID查询账户余额
// nAccountID: 账户ID, 不知道的话可以通过GetAccounts()获取, 可以只现货账户, C2C账户, 期货账户
// return: BalanceReturn对象
func (h *Huobi) GetAccountBalance(strAccountID string) BalanceReturn {
	balanceReturn := BalanceReturn{}

	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", strAccountID)
	jsonBanlanceReturn := apiKeyGet(make(map[string]string), strRequest, h.AccessKey, h.SecretKey)
	json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)

	return balanceReturn
}

// 下单
// placeRequestParams: 下单信息
// return: PlaceReturn对象
func (h *Huobi) Place(placeRequestParams PlaceRequestParams) PlaceReturn {
	placeReturn := PlaceReturn{}

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
	jsonPlaceReturn := apiKeyPost(mapParams, strRequest, h.AccessKey, h.SecretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	return placeReturn
}

// 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func (h *Huobi) SubmitCancel(strOrderID string) PlaceReturn {
	placeReturn := PlaceReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", strOrderID)
	jsonPlaceReturn := apiKeyPost(make(map[string]string), strRequest, h.AccessKey, h.SecretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	return placeReturn
}

// 查询订单详情
// strOrderID: 订单ID
// return: OrderReturn对象
func (h *Huobi) GetOrderInfo(strOrderID string) OrderReturn {
	orderReturn := OrderReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s", strOrderID)
	jsonPlaceReturn := apiKeyGet(make(map[string]string), strRequest, h.AccessKey, h.SecretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &orderReturn)

	return orderReturn
}

func NewHuobi(accessKey, secretKey string) *Huobi {
	return &Huobi{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}
