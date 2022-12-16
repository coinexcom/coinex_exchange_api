package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

// APIHTTPHOST api host
const APIHTTPHOST = ""

// ACCESSID apply from www.coinex.com
const ACCESSID = ""

// SECRETKEY apply from www.coinex.com
const SECRETKEY = ""

// CONTENTTYPE http content-type
const CONTENTTYPE = "application/json"

// USERAGENT http User-Agent
const USERAGENT = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"

func generateAuthorization(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5Str := fmt.Sprintf("%x", w.Sum(nil))
	md5Str = strings.ToUpper(md5Str)
	return md5Str
}

// CommonResp common resp: {'code': 0, 'message': 'success'}
type CommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AssetInfo asset balance info
type AssetInfo struct {
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
}

//BalanceResp response for api '/v1/balance/'
type BalanceResp struct {
	CommonResp
	AssetBalance map[string]AssetInfo `json:"data"`
}

// OrderInfo order base info
type OrderInfo struct {
	ID           int64  `json:"id"`
	Amount       string `json:"amount"`
	AvgPrice     string `json:"avg_price"`
	CreateTime   int64  `json:"create_time"`
	DealAmount   string `json:"deal_amount"`
	DealFee      string `json:"deal_fee"`
	DealMoney    string `json:"deal_money"`
	FinishedTime int64  `json:"finished_time"`
	MakerFeeRate string `json:"maker_fee_rate"`
	Market       string `json:"market"`
	OrderType    string `json:"order_type"`
	Price        string `json:"price"`
	Status       string `json:"status"`
	TakerFeeRate string `json:"taker_fee_rate"`
	Type         string `json:"type"`
	ClientID     string `json:"client_id"`
}

//OrderResp put limit order response from server
type OrderResp struct {
	CommonResp
	Order OrderInfo `json:"data"`
}

// CancelOrderResp cancel order response from server
type CancelOrderResp struct {
	CommonResp
	Order OrderInfo `json:"data"`
}

//QueryPageInfo response
type QueryPageInfo struct {
	CurrentPage int  `json:"curr_page"`
	Count       int  `json:"count"`
	HasNext     bool `json:"has_next"`
}

// OrderPendingResp query pending order response
type OrderPendingResp struct {
	CommonResp
	Data struct {
		QueryPageInfo
		Orders []OrderInfo `json:"data"`
	} `json:"data"`
}

// OrderFinishedResp query finished order response
type OrderFinishedResp struct {
	CommonResp
	Data struct {
		QueryPageInfo
		Orders []OrderInfo `json:"data"`
	} `json:"data"`
}

func interfaceToString(v interface{}) string {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Bool:
		return strings.Title(fmt.Sprintf("%v", v))
	case reflect.Slice, reflect.Array:
		var items []string
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			items = append(items, "'"+interfaceToString(s.Index(i))+"'")
		}
		return "[" + strings.Join(items, ", ") + "]"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func httpRequest(method, urlHost string, reqParameters map[string]interface{}) ([]byte, error) {
	params := make(map[string]interface{}, len(reqParameters))
	for k, v := range reqParameters {
		params[k] = v
	}
	currentMilliseconds := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	params["access_id"] = ACCESSID
	params["tonce"] = currentMilliseconds
	client := &http.Client{}
	var reqBody io.Reader
	if method == "POST" {
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(paramsJSON)
	}
	req, err := http.NewRequest(method, urlHost, reqBody)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	queryParamsString := ""
	for _, k := range keys {
		queryParamsString += fmt.Sprintf("%s=%s&", k, interfaceToString(params[k]))
	}
	toEncodeparamsString := queryParamsString + "secret_key=" + SECRETKEY
	req.Header.Set("Content-Type", CONTENTTYPE)
	req.Header.Set("User-Agent", USERAGENT)
	req.Header.Set("authorization", generateAuthorization(toEncodeparamsString))
	if method == "GET" || method == "DELETE" {
		req.URL.RawQuery = queryParamsString
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

//HTTPGet http get method
func HTTPGet(urlHost string, parameters map[string]interface{}) ([]byte, error) {
	return httpRequest("GET", urlHost, parameters)
}

//HTTPPost http get method
func HTTPPost(urlHost string, parameters map[string]interface{}) ([]byte, error) {
	return httpRequest("POST", urlHost, parameters)
}

//HTTPDelete http get method
func HTTPDelete(urlHost string, parameters map[string]interface{}) ([]byte, error) {
	return httpRequest("DELETE", urlHost, parameters)
}

// GetAccount Inquire account asset constructure
func GetAccount() ([]byte, error) {
	resp, err := HTTPGet(APIHTTPHOST+"/v1/balance/", nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// QueryOrderPending Acquire Unexecuted Order List.
func QueryOrderPending(market string, accountID, page, limit int) ([]byte, error) {
	parameters := map[string]interface{}{
		"market":     market,
		"page":       page,
		"limit":      limit,
		"account_id": accountID,
	}
	resp, err := HTTPGet(APIHTTPHOST+"/v1/order/pending", parameters)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// QueryOrderFinished Acquire executed order list
func QueryOrderFinished(market string, accountID, page, limit int) ([]byte, error) {
	parameters := map[string]interface{}{
		"market":     market,
		"page":       page,
		"limit":      limit,
		"account_id": accountID,
	}
	resp, err := HTTPGet(APIHTTPHOST+"/v1/order/finished", parameters)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// PutLimitOrder create limit order
func PutLimitOrder(amount, price, orderType, market string) ([]byte, error) {
	parameters := map[string]interface{}{
		"amount": amount,
		"price":  price,
		"type":   orderType,
		"market": market,
	}
	resp, err := HTTPPost(APIHTTPHOST+"/v1/order/limit", parameters)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// PutMarketOrder create market order
func PutMarketOrder(amount, orderType, market string) ([]byte, error) {
	parameters := map[string]interface{}{
		"amount": amount,
		"type":   orderType,
		"market": market,
	}
	resp, err := HTTPPost(APIHTTPHOST+"/v1/order/market", parameters)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateSubAccount create sub account
func CreateSubAccount(name string, allowedIPs []string) ([]byte, error) {
	parameters := map[string]interface{}{
		"allow_trade":   true,
		"allowed_ips":   allowedIPs,
		"sub_user_name": name,
	}
	resp, err := HTTPPost(APIHTTPHOST+"/v1/sub_account/auth/api", parameters)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CancelOrder Cancel unexecuted order
func CancelOrder(orderID int64, market string) ([]byte, error) {
	parameters := map[string]interface{}{
		"id":     orderID,
		"market": market,
	}
	resp, err := HTTPDelete(APIHTTPHOST+"/v1/order/pending", parameters)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	//Inquire account asset constructure
	accooutRespBody, err := GetAccount()
	if err != nil {
		fmt.Printf("GetAccount Error: %s", err)
		return
	}
	balanceResp := BalanceResp{}
	json.Unmarshal(accooutRespBody, &balanceResp)
	fmt.Printf("%v\n", balanceResp)

	//put limit order
	limitOrderRespBody, err := PutLimitOrder("1", "1", "buy", "BTCUSDT")
	if err != nil {
		fmt.Printf("PutLimitOrder Error: %v\n", err)
		return
	}
	putLimitOrderResp := OrderResp{}
	json.Unmarshal(limitOrderRespBody, &putLimitOrderResp)
	fmt.Printf("PutLimitOrder: %v\n", putLimitOrderResp)

	// Cancel order
	cancelOrderRespBody, err := CancelOrder(putLimitOrderResp.Order.ID, "BTCUSDT")
	if err != nil {
		fmt.Printf("CancelOrder Error: %v", err)
		return
	}
	cancleOrderResp := OrderResp{}
	json.Unmarshal(cancelOrderRespBody, &cancleOrderResp)
	fmt.Printf("CancelOrder: %v\n", cancleOrderResp)

	// put market order
	putMarketRespBody, err := PutMarketOrder("100", "buy", "BTCUSDT")
	if err != nil {
		fmt.Printf("PutMarketOrder Error: %v\n", err)
		return
	}
	putMarketOrderResp := OrderResp{}
	json.Unmarshal(putMarketRespBody, &putMarketOrderResp)
	fmt.Printf("PutMarketOrder: %v\n", putMarketOrderResp)

	//Acquire Unexecuted Order List
	queryOrderRespBody, err := QueryOrderPending("BTCUSDT", 0, 1, 10)
	if err != nil {
		fmt.Printf("CancelOrder Error: %v", err)
		return
	}
	queryOrders := OrderPendingResp{}
	json.Unmarshal(queryOrderRespBody, &queryOrders)
	fmt.Printf("QueryOrder: %v\n", queryOrders)

	// Acquire executed order list
	queryFinishedOrderRespBody, err := QueryOrderFinished("BTCUSDT", 0, 1, 10)
	if err != nil {
		fmt.Printf("CancelOrder Error: %v", err)
		return
	}
	finishedOrders := OrderFinishedResp{}
	json.Unmarshal(queryFinishedOrderRespBody, &finishedOrders)
	fmt.Printf("FinishedOrders: %v\n", finishedOrders)
}
