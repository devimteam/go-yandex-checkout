package go_yandex_checkout

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
)

const ApiEndpoint = "https://payment.yandex.net/api/v3/"

type YandexCheckoutClient interface {
	CreatePayment(request *CreatePaymentRequest, transactionId string) (*CreatePaymentResponse, error)
	GetPaymentInfo(paymentId string) (*GetPaymentResponse, error)
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
	return resp, nil
}

func (svc *ycc) sendPostJson(reqObj interface{}, resObj interface{}, idempotenceKey string) error {
	jsonStr, err := json.Marshal(reqObj)
	if nil != err {
		return err
	}

	req, err := http.NewRequest("POST", svc.url, bytes.NewBuffer(jsonStr))
	resp, err := svc.send(req, idempotenceKey)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		err = json.Unmarshal(respBody, resObj)
		if err != nil {
			return err
		}
	} else if resp.StatusCode == 202 {
		time.Sleep(1 * time.Second)
		return svc.sendPostJson(reqObj, resObj, idempotenceKey)
	} else if resp.StatusCode >= 400 {
		pr := &ProcessingResponse{}
		err = json.Unmarshal(respBody, pr)
		if resp.StatusCode == 429 {
			if *pr.RetryAfter > 0 {
				time.Sleep(pr.GetSleepTime())
				return svc.sendPostJson(reqObj, resObj, idempotenceKey)
			}
		} else {
			return err
		}
	}

	return nil
}

func (svc *ycc) CreatePayment(request *CreatePaymentRequest, transactionId string) (*CreatePaymentResponse, error) {
	res := &CreatePaymentResponse{}
	err := svc.sendPostJson(request, res, transactionId)
	return res, err
}

func (svc *ycc) GetPaymentInfo(paymentId string) (*GetPaymentResponse, error) {
	transactionId := "GetPayment-" + paymentId
	resObj := &GetPaymentResponse{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%spayments/%s", svc.url, paymentId), nil)
	resp, err := svc.send(req, transactionId)
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, resObj)

	// TODO: check response (similar to previous method)

	return resObj, err
}

func NewYandexCheckoutClient(url string, shopId string, secretKey string) YandexCheckoutClient {
	return &ycc{
		url:       url,
		shopId:    shopId,
		secretKey: secretKey,
	}
}
