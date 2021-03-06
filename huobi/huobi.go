package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/leizongmin/huobiapi"
)

type MarketTradeDetail struct {
	Ch   string `json:"ch"`
	Tick struct {
		Data []struct {
			Amount    float64 `json:"amount"`
			Direction string  `json:"direction"`
			Price     float64 `json:"price"`
			TS        int64   `json:"ts"`
		} `json:"data"`
	} `json:"tick"`
}

func (m *MarketTradeDetail) String() string {
	return fmt.Sprintln(m.Ch, "实时价格推送  价格:", m.Tick.Data[0].Price, " 数量:", m.Tick.Data[0].Amount, " 买卖：", m.Tick.Data[0].Direction)
}

type MarketDepth struct {
	Ch   string `json:"ch"`
	Tick struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
		TS   int64       `json:"ts"`
	} `json:"tick"`
}

type Account struct {
	ID     int64
	Type   string
	State  string
	UserID int64
}

type Huobi struct {
	accessKey      string
	secretKey      string
	tradeAccount   Account
	market         *Market
	depthListener  DepthlListener
	detailListener DetailListener
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
	strRequest := fmt.Sprintf("/v1/account/accounts/%d/balance", h.tradeAccount.ID)
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

func (h *Huobi) SetDetailListener(listener DetailListener) {
	h.detailListener = listener
}

func (h *Huobi) SetDepthlListener(listener DepthlListener) {
	h.depthListener = listener
}

// Listener 订阅事件监听器
type DetailListener = func(symbol string, detail *MarketTradeDetail)

func (h *Huobi) SubscribeDetail(symbols ...string) {
	for _, symbol := range symbols {
		h.market.Subscribe("market."+symbol+".trade.detail", func(topic string, j *huobiapi.JSON) {
			js, _ := j.MarshalJSON()
			var mtd MarketTradeDetail
			err := json.Unmarshal(js, &mtd)
			if err != nil {
				fmt.Println(err.Error())
			}

			ts := strings.Split(topic, ".")
			if h.detailListener != nil {
				h.detailListener(ts[1], &mtd)
			}

		})
	}

}

// Listener 订阅事件监听器
type DepthlListener = func(symbol string, depth *MarketDepth)

func (h *Huobi) SubscribeDepth(symbols ...string) {
	for _, symbol := range symbols {
		h.market.Subscribe("market."+symbol+".depth.step0", func(topic string, j *huobiapi.JSON) {
			js, _ := j.MarshalJSON()
			var md = MarketDepth{}
			err := json.Unmarshal(js, &md)
			if err != nil {
				fmt.Println(err.Error())
			}

			ts := strings.Split(topic, ".")
			if h.depthListener != nil {
				h.depthListener(ts[1], &md)
			}

		})
	}
}

func NewHuobi(accesskey, secretkey string) (*Huobi, error) {
	h := &Huobi{
		accessKey: accesskey,
		secretKey: secretkey,
	}

	if accesskey != "" {
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
	}

	var err error
	h.market, err = NewMarket()
	if err != nil {
		return nil, err
	}

	go h.market.Loop()

	fmt.Println("init huobi success.")

	return h, nil
}
