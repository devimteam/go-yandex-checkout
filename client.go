package go_yandex_checkout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const ApiEndpoint = "https://payment.yandex.net/api/v3/"

type YandexCheckoutClient interface {
	CreatePayment(request *CreatePaymentRequest, transactionId string) (*CreatePaymentResponse, error)
	GetPaymentInfo(paymentId string) (*GetPaymentResponse, error)
	CapturePayment(paymentId string, reqObj *CapturePaymentRequest) (*CapturePaymentResponse, error)
	CancelPayment(paymentId string, reqObj *CancelPaymentRequest) (*CancelPaymentResponse, error)
	CreateRefund(request *CreateRefundRequest) (*CreateRefundResponse, error)
	GetRefundInfo(refundId string) (*GetRefundResponse, error)
}

type ycc struct {
	url       string
	shopId    string
	secretKey string
}

func (svc *ycc) send(request *http.Request, idempotenceKey string) (*http.Response, error) {
	request.SetBasicAuth(svc.shopId, svc.secretKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotence-Key", idempotenceKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return resp, nil
	} else if resp.StatusCode == 202 {
		time.Sleep(500 * time.Millisecond)
		return svc.send(request, idempotenceKey)
	} else if resp.StatusCode >= 400 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		pr := &ProcessingResponse{}
		err = json.Unmarshal(respBody, pr)
		if resp.StatusCode == 429 {
			if *pr.RetryAfter > 0 {
				time.Sleep(pr.GetSleepTime())
				return svc.send(request, idempotenceKey)
			}
		} else {
			return resp, err
		}
	}

	return resp, nil
}

func (svc *ycc) sendPostJson(url string, reqObj interface{}, resObj interface{}, idempotenceKey string) error {
	jsonStr, err := json.Marshal(reqObj)
	if nil != err {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	resp, err := svc.send(req, idempotenceKey)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, resObj)
	if err != nil {
		return err
	}

	return nil
}

func (svc *ycc) CreatePayment(request *CreatePaymentRequest, transactionId string) (*CreatePaymentResponse, error) {
	resObj := &CreatePaymentResponse{}
	err := svc.sendPostJson(svc.url, request, resObj, transactionId)
	if nil != err {
		return nil, err
	}

	return resObj, nil
}

func (svc *ycc) GetPaymentInfo(paymentId string) (*GetPaymentResponse, error) {
	transactionId := "GetPayment-" + paymentId
	resObj := &GetPaymentResponse{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%spayments/%s", svc.url, paymentId), nil)
	resp, err := svc.send(req, transactionId)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, resObj)
	if nil != err {
		return nil, err
	}

	return resObj, nil
}

func (svc *ycc) CapturePayment(paymentId string, reqObj *CapturePaymentRequest) (*CapturePaymentResponse, error) {
	transactionId := "CapturePayment-" + paymentId
	resObj := &CapturePaymentResponse{}

	url := fmt.Sprintf("%spayments/%s/capture", svc.url, paymentId)
	err := svc.sendPostJson(url, reqObj, resObj, transactionId)

	if nil != err {
		return nil, err
	}

	return resObj, nil
}

func (svc *ycc) CancelPayment(paymentId string, reqObj *CancelPaymentRequest) (*CancelPaymentResponse, error) {
	transactionId := "CancelPayment-" + paymentId
	resObj := &CancelPaymentResponse{}

	url := fmt.Sprintf("%spayments/%s/cancel", svc.url, paymentId)
	err := svc.sendPostJson(url, reqObj, resObj, transactionId)

	if nil != err {
		return nil, err
	}

	return resObj, nil
}

func (svc *ycc) CreateRefund(request *CreateRefundRequest) (*CreateRefundResponse, error) {
	resObj := &CreateRefundResponse{}
	url := fmt.Sprintf("%srefunds", svc.url)
	err := svc.sendPostJson(url, request, resObj, request.PaymentID)
	if nil != err {
		return nil, err
	}

	return resObj, nil
}

func (svc *ycc) GetRefundInfo(refundId string) (*GetRefundResponse, error) {
	transactionId := "get-refund-" + refundId
	resObj := &GetRefundResponse{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%srefunds/%s", svc.url, refundId), nil)
	resp, err := svc.send(req, transactionId)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, resObj)
	if nil != err {
		return nil, err
	}

	return resObj, nil
}

func NewYandexCheckoutClient(url string, shopId string, secretKey string) YandexCheckoutClient {
	return &ycc{
		url:       url,
		shopId:    shopId,
		secretKey: secretKey,
	}
}
