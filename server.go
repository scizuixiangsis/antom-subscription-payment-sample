package main

import (
	"encoding/json" // 类似 JavaScript 的 JSON.parse/stringify
	"fmt"
	"io"
	"log"
	"net/http" // 处理 HTTP 请求，类似 Express.js
	"strconv"
	"time"

	"github.com/Rhymond/go-money" // 第三方包，类似 npm install 的包
	"github.com/alipay/global-open-sdk-go/com/alipay/api/model"
	"github.com/alipay/global-open-sdk-go/com/alipay/api/request/notify"
	"github.com/alipay/global-open-sdk-go/com/alipay/api/request/subscription"
	"github.com/alipay/global-open-sdk-go/com/alipay/api/response"
	"github.com/alipay/global-open-sdk-go/com/alipay/api/tools"

	defaultAlipayClient "github.com/alipay/global-open-sdk-go/com/alipay/api"

	"github.com/google/uuid"
)

const (
	/*
	  replace with your client id <br>
	  find your client id here: <a href="https://dashboard.alipay.com/global-payments/developers/quickStart">quickStart</a>
	*/
	ClientID = "5YEZ50305NN403523"

	/*
	  replace with your antom public key (used to verify signature) <br>
	  find your antom public key here: <a href="https://dashboard.alipay.com/global-payments/developers/quickStart">quickStart</a>
	*/
	AntomPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkWZnOpL6uX6UT05KuT2GDUiB2bSM1bTq5rFRrbpXy7TihPX7dYc6lJIF4cEBD/DeAU9g9CMGJL/X7d3DgF++4y9tIb+yZ3ihJmjkVMdVwWhuqY1NIqvQKFwlNC+LzLrZiI0qH3SyoEEZtfQLOtBALLbdvKmQLzXcbhJ5uJJqgi0W4CLOJhQXIrPxwlBWHjCL44/BCeqqSCJq1oiTjwmu6CSUCyauOTXs4JAPBw6673OsdJMZq+Cn0m7dE3nb98XisCE2NJqS00JGloZynAeoaVPI9OCTrXt/m6+zsQJF6jibqnok4tiYBfQRxzdWkzT7tmeZI4Jc2HVK7LOS4mJNSwIDAQAB"

	/*
	  replace with your private key (used to sign) <br>
	  please ensure the secure storage of your private key to prevent leakage
	*/
	MerchantPrivateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCEyw05m5tLir11oi9NfnWYvPnz+MTsh+FxQFiahaKZcrzpoZEeWcA66nbC6r4Pivn8KvOP/TQSlxs7vB8iQ2WJ26EF9dk+IS87M0VDKNTxY0xwnolhMWPZxAacbmPZY0QFdKSveVnGkJ6/9yR7AkVzHNSkFEzpyQpZL260bMFl45P7PkMlu4aTZ+mpR0mS505n3+jKyG+yyWOdFwlEvDgw040yKlKEmUixALwfMJcN68mJgObKCwU49lcNfodz/7836bDDzSlfQYtU9PfQAW9/Apb8R1ywcxfAZVTOG6Eru4LjVQfvuJiAtNAM56wlZcxzqDYYNtBeVaYXjhw9POJPAgMBAAECggEAB/ouvDQ12w7lnMh2cQBUzi0kEtzx74emKmiEKhw3DWLQKHTVQ+5+Vsu5ALKbYlGl/NkTsqWyWB+NukPt1uAXnHV+Md54A3x42uSUl5k/WZTuhaFwfU87QVy+TO1wwCFvd6hvMD9o/j23265cGaukQmsL9yNlD3JNVdg8nUgfQiCXv1dBoMBahurEho6KQo8H2dNQAiB+oody38LHM5i+BT92a1cZkO+3atMC7yko7Y4vloi3sHGWOy7YSRxGCgm0PuAr3n1dMDd0LUORtZ9R6nXspALnET9cAwnQfMHCp4Y3rftuBU5nEXyq3tY3kTHDERv/Vasw98ZPgHpKhLkaeQKBgQDz9n+JtZdip74hD1DYRfvD0WPskLaK7yK+9CxtZ426ijakhUgI7jO8rW+hOWPzslYVZwbpXwAavl75VQiEw4HZ84aKgDGn1tiK/5svXVPQdAnkCHEqPQp/6nP3i7M7nlYMUCUiHBSD0oewg3VNFJVtS9uuQ+hqKmjpmvxlPvvqDQKBgQCLWF2lNw0aihucsGRHvl5IRcdb6e1uOUYvqvRuEAi8mEa1YlUgS3hyKgiBTKlTAaYlOtc4Ucuag5/rN0lk664wqE1iCxYY77EWZRpi3ilDHO20D+RNXG1jVdtJdHXUdwlPjChKLRYsfIPXlmGA+0UCDJ1b9bHHnHcnmFkUPKLyywKBgQCrio4PRLKX2h8km+Ja0IrBHADJHNBeTNv/rS14GDJeEkVt1ZHbRbL3XnR5xyLy/ljtX65KdlRaebXKV/JPeDFcEZJu3MkNnVJSGn0CBvuiPZWe1BjOfHFflHnKfF6g8yrKKaiSnXAHaQekJCtc8bZITejAVlucGwn+CM6kWm9EGQKBgGhZlldBMKjtP9xJI++uGgDZcH/eYJWogmz0AvPhQgmpp1nx93mlyt8Dpzbc5/hnRbqfo8hjSKu/YiTNVEMlU17QypJfZv7pkJ4KvIXJhPDjWwb616cvTiOTihIqCos/UVOmzA0wUmiiHkF2NjJW+MieFcFl7upiu8CFEEBdYFGdAoGBAJvJTA7gZTad0RbCQT53HFJDwiU/sg7FQzDYnOs/4hneZ2oZwPyoisfHFTVomgLzcw6MIMm/udjDirgVRKy/v83Np701eIbjiZ4pXjB6LOfOa9seEpTxb7iPg7ONGAfwglv8JrUlIcnqyTycPaeTTEV5VGGnGJk97bYRqMaJT2g1"
)

// Initialize global client
var client *defaultAlipayClient.DefaultAlipayClient

func init() { //init() - 初始化函数，会在main函数之前执行
	client = defaultAlipayClient.NewDefaultAlipayClient(
		"https://open-sea-global.alipay.com",
		ClientID,
		MerchantPrivateKey,
		AntomPublicKey)
}

func main() {
	// Register routes - 注册路由，类似 Express.js 的 app.use('/api', apiRouter)

	// 1. 创建订阅接口 - 前端调用此接口发起订阅支付，返回 Antom 支付页面链接
	http.HandleFunc("/subscriptions/create", enableCORS(handleSubscriptionCreate))

	// 2. 支付通知回调 - Antom 服务器在每次扣款完成后调用此接口通知支付结果（服务器对服务器）
	http.HandleFunc("/subscriptions/receivePaymentNotify", enableCORS(handleReceivePaymentNotify))

	// 3. 订阅状态通知回调 - Antom 服务器在订阅状态变化时调用此接口（创建/取消/到期等）（服务器对服务器）
	http.HandleFunc("/subscriptions/receiveSubscriptionNotify", enableCORS(handleReceiveSubscriptionNotify))

	fmt.Println("Open your browser and visit: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // 启动 HTTP 服务器，监听 8080 端口
}

func handleSubscriptionCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subscriptionVO SubscriptionVO
	json.NewDecoder(r.Body).Decode(&subscriptionVO)

	// Set default terminal type if not provided
	if subscriptionVO.TerminalType == "" {
		subscriptionVO.TerminalType = "WEB"
	}

	request, subscriptionCreateRequest := subscription.NewAlipaySubscriptionCreateRequest()

	// replace with your subscriptionRequestId
	// You can save the relationship between the subscriptionRequestId and the user ID for later information retrieval.
	subscriptionRequestId := uuid.New().String()
	subscriptionCreateRequest.SubscriptionRequestId = subscriptionRequestId
	subscriptionCreateRequest.SubscriptionDescription = "Subscription Description"

	// set subscription start time and end time. you might want to consider time zones
	// If the start time is earlier than the authorization time, the subscription is successful.
	// If the start time is later than the authorization time, the payment is made after the successful authorization, which is the pre-sale
	// For details, please refer to: <a href="https://docs.antom.com/ac/subscriptionpay/activation#uiqBb">Samples</a>
	subscriptionCreateRequest.SubscriptionStartTime = time.Now().Format("2006-01-02T15:04:05-07:00")
	subscriptionCreateRequest.SubscriptionEndTime = time.Now().AddDate(3, 0, 0).Format("2006-01-02T15:04:05-07:00")

	// set periodRule
	subscriptionCreateRequest.PeriodRule = &model.PeriodRule{
		PeriodType:  model.PeriodType(subscriptionVO.PeriodType),
		PeriodCount: subscriptionVO.PeriodCount,
	}

	// set paymentMethod
	subscriptionCreateRequest.PaymentMethod = &model.PaymentMethod{
		PaymentMethodType: subscriptionVO.PaymentMethodType,
	}

	// convert amount unit(in practice, amount should be calculated on your serverside)
	// For details, please refer to: <a href="https://docs.antom.com/ac/ref/cc">Usage rules of the Amount object</a>
	parseFloat, _ := strconv.ParseFloat(subscriptionVO.AmountValue, 64)
	amountValue := money.NewFromFloat(parseFloat, subscriptionVO.Currency).Amount()
	amount := &model.Amount{
		Currency: subscriptionVO.Currency,
		Value:    strconv.FormatInt(amountValue, 10),
	}

	// set payment amount
	subscriptionCreateRequest.PaymentAmount = amount

	// set order info
	subscriptionCreateRequest.OrderInfo = &model.OrderInfo{
		OrderAmount: amount,
	}

	// set settlement strategy
	// replace with your existing settlement currency
	subscriptionCreateRequest.SettlementStrategy = &model.SettlementStrategy{
		SettlementCurrency: "USD",
	}

	// set env info
	terminal := model.TerminalType(subscriptionVO.TerminalType)
	env := &model.Env{TerminalType: terminal}
	if terminal != model.WEB && subscriptionVO.OsType != "" {
		env.OsType = model.OsType(subscriptionVO.OsType)
	}
	subscriptionCreateRequest.Env = env

	// replace with your notify url
	// or configure your notify url here: <a href="https://dashboard.antom.com/global-payments/developers/iNotify">Notification URL</a>
	subscriptionCreateRequest.PaymentNotificationUrl = "http://www.yourNotifyUrl.com/subscriptions/receivePaymentNotify"
	subscriptionCreateRequest.SubscriptionNotificationUrl = "http://www.yourNotifyUrl.com/subscriptions/receiveSubscriptionNotify"

	// replace with your subscription redirect url
	subscriptionCreateRequest.SubscriptionRedirectUrl = "http://localhost:5173/index.html?subscriptionRequestId=" + subscriptionRequestId

	startTime := time.Now()
	subscriptionRequestJson, _ := json.Marshal(subscriptionCreateRequest)
	log.Printf("subscription create request: %s", subscriptionRequestJson)
	response, err := client.Execute(request)
	if err != nil {
		sendJSONResponse(w, http.StatusOK, ApiResponse{
			Status:                "error",
			SubscriptionRequestID: subscriptionRequestId,
			Message:               err.Error(),
		})
		return
	}

	log.Printf("subscription create response: %+v", response)
	log.Printf("subscription create request cost time: %v ms\n", time.Since(startTime).Milliseconds())

	sendJSONResponse(w, http.StatusOK, ApiResponse{
		Status:                "success",
		SubscriptionRequestID: subscriptionRequestId,
		Data:                  response,
	})
}

/*
receive payment notify
*/
func handleReceivePaymentNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "FAIL",
			ResultMessage: "fail.",
			ResultStatus:  "F",
		})
		return
	}
	notifyBody := string(rawBody)
	// retrieve the required parameters from http request
	requestURI := r.RequestURI
	requestMethod := r.Method
	// retrieve the required parameters from request header
	requestTime := r.Header.Get("request-time")
	clientID := r.Header.Get("client-id")
	signature := r.Header.Get("signature")

	// verify the signature of notification
	checkSignature, err := tools.CheckSignature(requestURI, requestMethod, clientID, requestTime, notifyBody, signature, AntomPublicKey)
	if err != nil || !checkSignature {
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "FAIL",
			ResultMessage: "fail.",
			ResultStatus:  "F",
		})
		return
	}

	// deserialize the notification body
	var paymentNotify notify.AlipaySubscriptionPayNotify
	if err := json.Unmarshal(rawBody, &paymentNotify); err != nil {
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "FAIL",
			ResultMessage: "fail.",
			ResultStatus:  "F",
		})
		return
	}

	if paymentNotify.Result.ResultCode == "SUCCESS" {
		// handle your own business logic.
		// e.g. The payment information of the user is saved through the relationship between the subscriptionRequestId and the user ID.
		log.Printf("receive payment notify: %s", notifyBody)
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "SUCCESS",
			ResultMessage: "success.",
			ResultStatus:  "S",
		})
		return
	}

	sendJSONResponse(w, http.StatusOK, response.Result{
		ResultCode:    "SYSTEM_ERROR",
		ResultMessage: "system error.",
		ResultStatus:  "F",
	})
}

/*
receive subscription notify
*/
func handleReceiveSubscriptionNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "FAIL",
			ResultMessage: "fail.",
			ResultStatus:  "F",
		})
		return
	}
	notifyBody := string(rawBody)
	// retrieve the required parameters from http request
	requestURI := r.RequestURI
	requestMethod := r.Method
	// retrieve the required parameters from request header
	requestTime := r.Header.Get("request-time")
	clientID := r.Header.Get("client-id")
	signature := r.Header.Get("signature")

	// verify the signature of notification
	checkSignature, err := tools.CheckSignature(requestURI, requestMethod, clientID, requestTime, notifyBody, signature, AntomPublicKey)
	if err != nil || !checkSignature {
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "FAIL",
			ResultMessage: "fail.",
			ResultStatus:  "F",
		})
		return
	}

	// deserialize the notification body
	var subscriptionNotify notify.AlipaySubscriptionNotify
	if err := json.Unmarshal(rawBody, &subscriptionNotify); err != nil {
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "FAIL",
			ResultMessage: "fail.",
			ResultStatus:  "F",
		})
		return
	}

	if subscriptionNotify.SubscriptionNotificationType == model.SubscriptionNotificationType_CREATE {
		// handle your own business logic.
		// e.g. The subscription information of the user is saved through the relationship between the subscriptionRequestId and the user ID.
		log.Printf("receive subscription notify: %s", notifyBody)
		sendJSONResponse(w, http.StatusOK, response.Result{
			ResultCode:    "SUCCESS",
			ResultMessage: "success.",
			ResultStatus:  "S",
		})
		return
	}

	sendJSONResponse(w, http.StatusOK, response.Result{
		ResultCode:    "SYSTEM_ERROR",
		ResultMessage: "system error.",
		ResultStatus:  "F",
	})
}

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SubscriptionVO represents subscription request data
type SubscriptionVO struct {
	PeriodType        string `json:"periodType"`
	PeriodCount       int    `json:"periodCount"`
	AmountValue       string `json:"amountValue"`
	Currency          string `json:"currency"`
	PaymentMethodType string `json:"paymentMethodType"`
	TerminalType      string `json:"terminalType"`
	OsType            string `json:"osType"`
}

// ApiResponse represents API response structure
type ApiResponse struct {
	Status                string      `json:"status"`
	SubscriptionRequestID string      `json:"subscriptionRequestId,omitempty"`
	Message               string      `json:"message,omitempty"`
	Data                  interface{} `json:"data,omitempty"`
}
