package huobi

type AccountsData struct {
	ID     int64  `json:"id"`      // Account ID
	Type   string `json:"type"`    // 账户类型, spot: 现货账户
	State  string `json:"state"`   // 账户状态, working: 正常, lock: 账户被锁定
	UserID int64  `json:"user-id"` // 用户ID
}

type AccountsReturn struct {
	Status  string         `json:"status"` // 请求状态
	Data    []AccountsData `json:"data"`   // 用户数据
	ErrCode string         `json:"err-code"`
	ErrMsg  string         `json:"err-msg"`
}

// 子账户结构
type SubAccount struct {
	Currency string `json:"currency"` // 币种
	Balance  string `json:"balance"`  // 结余
	Type     string `json:"type"`     // 类型, trade: 交易余额, frozen: 冻结余额
}

type Balance struct {
	ID     int64        `json:"id"`    // 账户ID
	State  string       `json:"state"` // 账户状态, working: 正常, lock: 账户被锁定
	Type   string       `json:"type"`  // 账户类型, spot: 现货账户
	List   []SubAccount `json:"list"`  // 子账户数组
	UserID int64        `json:"user-id"`
}

type BalanceReturn struct {
	Status  string  `json:"status"` // 请求状态
	Data    Balance `json:"data"`   // 账户余额
	ErrCode string  `json:"err-code"`
	ErrMsg  string  `json:"err-msg"`
}

type PlaceRequestParams struct {
	AccountID string `json:"account-id"` // 账户ID
	Amount    string `json:"amount"`     // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
	Price     string `json:"price"`      // 下单价格, 市价单不传该参数
	Source    string `json:"source"`     // 订单来源, api: API调用, margin-api: 借贷资产交易
	Symbol    string `json:"symbol"`     // 交易对, btcusdt, bccbtc......
	Type      string `json:"type"`       // 订单类型, buy-market: 市价买, sell-market: 市价卖, buy-limit: 限价买, sell-limit: 限价卖
}

type PlaceReturn struct {
	Status  string `json:"status"`
	Data    string `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type Order struct {
	ID          int64  `json:"id"`
	Symbol      string `json:"symbol"`
	State       string `json:"state"`
	FieldAmount string `json:"field-amount"`
	Price       string `json:"price"`
	Type        string `json:"type"`
}

type OrderReturn struct {
	Status  string `json:"status"`
	Data    Order  `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type OrdersRequestParams struct {
	Symbol string `json:"symbol"`
	States string `json:"states"`
}

type OrdersReturn struct {
	Status  string  `json:"status"`
	Data    []Order `json:"data"`
	ErrCode string  `json:"err-code"`
	ErrMsg  string  `json:"err-msg"`
}
