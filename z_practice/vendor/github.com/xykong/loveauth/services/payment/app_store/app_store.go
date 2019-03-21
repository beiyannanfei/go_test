package app_store

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"io/ioutil"
	"net/http"
	"time"
)

//
// Send the Receipt Data to the App Store
// On your server, create a JSON object with the following keys:
// Binding from JSON
// noinspection ALL
type DoVerifyReceiptReq struct {
	//
	// The base64 encoded receipt data.
	//
	ReceiptData string `json:"receipt-data" structs:"receipt-data"`
	//
	// Only used for receipts that contain auto-renewable subscriptions.
	// Your app’s shared secret (a hexadecimal string).
	//
	Password string `json:"password" structs:"password"`
	//
	// Only used for iOS7 style app receipts that contain auto-renewable or non-renewing subscriptions.
	// If value is true, response includes only the latest renewal transaction for any subscriptions.
	//
	ExcludeOldTransactions string `json:"exclude-old-transactions" structs:"exclude-old-transactions"`
}

//
// in: body
// swagger:parameters profile_query_vip
type DoVerifyReceiptReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoVerifyReceiptReq DoVerifyReceiptReq
}

type InApp struct {
	Quantity                string `json:"quantity"`
	ProductId               string `json:"product_id"`
	TransactionId           string `json:"transaction_id"`
	OriginalTransactionId   string `json:"original_transaction_id"`
	PurchaseDate            string `json:"purchase_date"`
	PurchaseDateMs          string `json:"purchase_date_ms"`
	PurchaseDatePst         string `json:"purchase_date_pst"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
	IsTrialPeriod           string `json:"is_trial_period"`
}

type Receipt struct {
	//OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
	PurchaseDateMs         string `json:"purchase_date_ms"`
	UniqueIdentifier       string `json:"unique_identifier"`
	OriginalTransactionId  string `json:"original_transaction_id"`
	Bvrs                   string `json:"bvrs"`
	TransactionId          string `json:"transaction_id"`
	Quantity               string `json:"quantity"`
	UniqueVendorIdentifier string `json:"unique_vendor_identifier"`
	ItemId                 string `json:"item_id"`
	ProductId              string `json:"product_id"`
	PurchaseDate           string `json:"purchase_date"`
	//OriginalPurchaseDate    string `json:"original_purchase_date"`
	PurchaseDatePst string `json:"purchase_date_pst"`
	Bid             string `json:"bid"`
	//OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`

	ReceiptType                string  `json:"receipt_type"`
	AdamId                     int     `json:"adam_id"`
	AppItemId                  int     `json:"app_item_id"`
	BundleId                   string  `json:"bundle_id"`
	ApplicationVersion         string  `json:"application_version"`
	DownloadId                 int     `json:"download_id"`
	VersionExternalIdentifier  int     `json:"version_external_identifier"`
	ReceiptCreationDate        string  `json:"receipt_creation_date"`
	ReceiptCreationDateMs      string  `json:"receipt_creation_date_ms"`
	ReceiptCreationDatePst     string  `json:"receipt_creation_date_pst"`
	RequestDate                string  `json:"request_date"`
	RequestDateMs              string  `json:"request_date_ms"`
	RequestDatePst             string  `json:"request_date_pst"`
	OriginalPurchaseDate       string  `json:"original_purchase_date"`
	OriginalPurchaseDateMs     string  `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePst    string  `json:"original_purchase_date_pst"`
	OriginalApplicationVersion string  `json:"original_application_version"`
	InApp                      []InApp `json:"in_app"`
}

//
// Parse the Response
// The response’s payload is a JSON object that contains the following keys and values:
// swagger:response DoVerifyReceiptRsp
type DoVerifyReceiptRsp struct {
	// in: body
	Body struct {
		//
		// Either 0 if the receipt is valid, or one of the error codes listed in Table 2-1.
		//
		// For iOS 6 style transaction receipts, the status code reflects the status of the specific transaction’s receipt.
		//
		// For iOS 7 style app receipts, the status code is reflects the status of the app receipt as a whole.
		// For example, if you send a valid app receipt that contains an expired subscription,
		// the response is 0 because the receipt as a whole is valid.
		//
		// Table 2-1  Status codes Description
		//
		// 21000 The App Store could not read the JSON object you provided.
		//
		// 21002 The data in the receipt-data property was malformed or missing.
		//
		// 21003 The receipt could not be authenticated.
		//
		// 21004 The shared secret you provided does not match the shared secret on file for your account.
		//
		// 21005 The receipt server is not currently available.
		//
		// 21006 This receipt is valid but the subscription has expired. When this status code is returned to your server,
		//  the receipt data is also decoded and returned as part of the response.
		//  Only returned for iOS 6 style transaction receipts for auto-renewable subscriptions.
		//
		// 21007 This receipt is from the test environment, but it was sent to the production environment for verification.
		//  Send it to the test environment instead.
		//
		// 21008 This receipt is from the production environment, but it was sent to the test environment for verification.
		//  Send it to the production environment instead.
		//
		// 21010 This receipt could not be authorized. Treat this the same as if a purchase was never made.
		//
		// 21100-21199 Internal data access error.
		//
		Status int `json:"status"`
		//
		// A JSON representation of the receipt that was sent for verification.
		// For information about keys found in a receipt, see Receipt Fields.
		//
		Receipt Receipt `json:"receipt"`

		Environment string `json:"environment"`

		//
		// Only returned for receipts containing auto-renewable subscriptions.
		// For iOS 6 style transaction receipts, this is the base-64 encoded receipt for the most recent renewal.
		// For iOS 7 style app receipts, this is the latest base-64 encoded app receipt.
		//
		LatestReceipt string `json:"latest_receipt"`
		//
		// Only returned for receipts containing auto-renewable subscriptions.
		// For iOS 6 style transaction receipts, this is the JSON representation of the receipt for the most recent renewal.
		// For iOS 7 style app receipts, the value of this key is an array containing all in-app purchase transactions.
		// This excludes transactions for a consumable product that have been marked as finished by your app.
		//
		LatestReceiptInfo string `json:"latest_receipt_info"`
		//
		// Only returned for iOS 6 style transaction receipts, for an auto-renewable subscription.
		// The JSON representation of the receipt for the expired subscription.
		//
		LatestExpiredReceiptInfo string `json:"latest_expired_receipt_info"`
		//
		// Only returned for iOS 7 style app receipts containing auto-renewable subscriptions.
		// In the JSON file, the value of this key is an array where each element contains
		// the pending renewal information for each auto-renewable subscription identified by the Product Identifier.
		// A pending renewal may refer to a renewal that is scheduled in the future or a renewal
		// that failed in the past for some reason.
		//
		PendingRenewalInfo string `json:"pending_renewal_info"`
		//
		// Retry validation for this receipt. Only applicable to status codes 21100-21199 (listed in Table 2-1)
		//
		IsRetryable string `json:"is-retryable"`
	}
}

func errorMessageFromCode(code int) string {
	switch {
	case code == 21000:
		return "The App Store could not read the JSON object you provided."

	case code == 21002:
		return "The data in the receipt-data property was malformed or missing."
	case code == 21003:
		return "The receipt could not be authenticated."
	case code == 21004:
		return "The shared secret you provided does not match the shared secret on file for your account."
	case code == 21005:
		return "The receipt server is not currently available."
	case code == 21006:
		return "This receipt is valid but the subscription has expired. When this status code is returned to your server," +
			"the receipt data is also decoded and returned as part of the response." +
			"Only returned for iOS 6 style transaction receipts for auto-renewable subscriptions."
	case code == 21007:
		return "This receipt is from the test environment, but it was sent to the production environment for verification." +
			" Send it to the test environment instead."
	case code == 21008:
		return "This receipt is from the production environment, but it was sent to the test environment for verification." +
			" Send it to the production environment instead."
	case code == 21010:
		return "This receipt could not be authorized. Treat this the same as if a purchase was never made."
	case code >= 21100 && code <= 21199:
		return "Internal data access error."
	}

	return "Undefined error."
}

func doVerifyReceipt(url string, request DoVerifyReceiptReq) (*DoVerifyReceiptRsp, error) {

	start := time.Now()
	jsonValue, _ := json.Marshal(request)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {

		logrus.WithFields(logrus.Fields{
			"elapsed": time.Since(start),
			"url":     url,
			"body":    string(jsonValue),
		}).Error("VerifyReceipt post request.")

		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"elapsed":  time.Since(start),
		"url":      url,
		"body":     string(jsonValue),
		"response": string(respBody),
	}).Info("VerifyReceipt post request.")

	var data DoVerifyReceiptRsp
	if err := json.Unmarshal(respBody, &data.Body); err != nil {
		return nil, errors.NewCodeError(errors.Failed, err)
	}

	//if data.Body.Status != 0 {
	//	return nil, errors.NewCodeString(errors.Failed, errorMessageFromCode(data.Body.Status))
	//}

	return &data, nil
}

func VerifyReceipt(request DoVerifyReceiptReq) (*DoVerifyReceiptRsp, error) {

	urlSandbox := "https://sandbox.itunes.apple.com/verifyReceipt"
	urlProduction := "https://buy.itunes.apple.com/verifyReceipt"

	// prefer production fist
	data, err := doVerifyReceipt(urlProduction, request)
	if data != nil && data.Body.Status == 21007 {
		data, err = doVerifyReceipt(urlSandbox, request)
	}

	//// prefer sandbox fist
	//data, err := doVerifyReceipt(urlSandbox, request)
	//if data != nil && data.Body.Status == 21008 {
	//	data, err = doVerifyReceipt(urlProduction, request)
	//}

	if err != nil {
		return nil, errors.NewCodeError(errors.Failed, err)
	}

	if data.Body.Status != 0 {
		return nil, errors.NewCodeString(errors.Failed, errorMessageFromCode(data.Body.Status))
	}

	return data, err
}
