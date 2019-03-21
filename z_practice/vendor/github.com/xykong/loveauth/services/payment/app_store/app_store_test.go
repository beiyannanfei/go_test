package app_store

import (
	"reflect"
	"testing"
)

func TestVerifyReceipt(t *testing.T) {
	type args struct {
		request DoVerifyReceiptReq
	}
	tests := []struct {
		name    string
		args    args
		want    *DoVerifyReceiptRsp
		wantErr bool
	}{
		{
			"",
			args{
				DoVerifyReceiptReq{
					"eyJzaWduYXR1cmUiID0gIkFoaHEwbmFxaG01ci8wSUNuYk9hVThCeVBLNkRja2ZsanRCMDNnZUh4dk0ybEVjVkdqK2NVM1lnWGZ0RkZCZ2lFYkR1NGdoYXFVWFRqRzlpc25Zeit6VWFhTXRUZDJ6YnNLbFhIMitzYm1ZaC9tY2M2eWt3dVFkaFZsOWZKYkxNaFVXWHNTUlIxVlVjQlE4TkxjVGg5ZHNXOTVTREcxdG9DQk45cWM0L01CZlVBQUFEVnpDQ0ExTXdnZ0k3b0FNQ0FRSUNDQnVwNCtQQWhtL0xNQTBHQ1NxR1NJYjNEUUVCQlFVQU1IOHhDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFLREFwQmNIQnNaU0JKYm1NdU1TWXdKQVlEVlFRTERCMUJjSEJzWlNCRFpYSjBhV1pwWTJGMGFXOXVJRUYxZEdodmNtbDBlVEV6TURFR0ExVUVBd3dxUVhCd2JHVWdhVlIxYm1WeklGTjBiM0psSUVObGNuUnBabWxqWVhScGIyNGdRWFYwYUc5eWFYUjVNQjRYRFRFME1EWXdOekF3TURJeU1Wb1hEVEUyTURVeE9ERTRNekV6TUZvd1pERWpNQ0VHQTFVRUF3d2FVSFZ5WTJoaGMyVlNaV05sYVhCMFEyVnlkR2xtYVdOaGRHVXhHekFaQmdOVkJBc01Fa0Z3Y0d4bElHbFVkVzVsY3lCVGRHOXlaVEVUTUJFR0ExVUVDZ3dLUVhCd2JHVWdTVzVqTGpFTE1Ba0dBMVVFQmhNQ1ZWTXdnWjh3RFFZSktvWklodmNOQVFFQkJRQURnWTBBTUlHSkFvR0JBTW1URXVMZ2ppbUx3Ukp4eTFvRWYwZXNVTkRWRUllNndEc25uYWwxNGhOQnQxdjE5NVg2bjkzWU83Z2kzb3JQU3V4OUQ1NTRTa01wK1NheWc4NGxUYzM2MlV0bVlMcFduYjM0bnF5R3g5S0JWVHk1T0dWNGxqRTFPd0Mrb1RuUk0rUUxSQ21lTnhNYlBaaFM0N1QrZVp0REVoVkI5dXNrMytKTTJDb2dmd283QWdNQkFBR2pjakJ3TUIwR0ExVWREZ1FXQkJTSmFFZU51cTlEZjZaZk42OEZlK0kydTIyc3NEQU1CZ05WSFJNQkFmOEVBakFBTUI4R0ExVWRJd1FZTUJhQUZEWWQ2T0tkZ3RJQkdMVXlhdzdYUXd1UldFTTZNQTRHQTFVZER3RUIvd1FFQXdJSGdEQVFCZ29xaGtpRzkyTmtCZ1VCQkFJRkFEQU5CZ2txaGtpRzl3MEJBUVVGQUFPQ0FRRUFlYUpWMlU1MXJ4ZmNxQUFlNUMyL2ZFVzhLVWw0aU80bE11dGE3TjZYelAxcFpJejFOa2tDdElJd2V5Tmo1VVJZSEsrSGpSS1NVOVJMZ3VObDBua2Z4cU9iaU1ja3dSdWRLU3E2OU5JbnJaeUNENjZSNEs3N25iOWxNVEFCU1NZbHNLdDhvTnRsaGdSLzFralNTUlFjSGt0c0RjU2lRR0tNZGtTbHA0QXlYZjd2bkhQQmU0eUN3WVYyUHBTTjA0a2JvaUozcEJseHNHd1YvWmxMMjZNMnVlWUhLWUN1WGhkcUZ3eFZnbTUyaDNvZUpPT3Qvdlk0RWNRcTdlcUhtNm0wM1o5YjdQUnpZTTJLR1hIRG1PTWs3dkRwZU1WbExEUFNHWXoxK1Uzc0R4SnplYlNwYmFKbVQ3aW16VUtmZ2dFWTd4eGY0Y3pmSDB5ajV3TnpTR1RPdlE9PSI7ICJwdXJjaGFzZS1pbmZvIiA9ICJld29KSW05eWFXZHBibUZzTFhCMWNtTm9ZWE5sTFdSaGRHVXRjSE4wSWlBOUlDSXlNREUxTFRFeExUSXpJREE0T2pNeU9qSTNJRUZ0WlhKcFkyRXZURzl6WDBGdVoyVnNaWE1pT3dvSkluVnVhWEYxWlMxcFpHVnVkR2xtYVdWeUlpQTlJQ0psTkRjME1qVmtaamRtWWpaaFlqaGpNMk0zWlRVM016QmpNMkUzWldJMVlqWTFOak00TWpOa0lqc0tDU0p2Y21sbmFXNWhiQzEwY21GdWMyRmpkR2x2YmkxcFpDSWdQU0FpTVRBd01EQXdNREU0TVRRNU1qVTFOeUk3Q2draVluWnljeUlnUFNBaU15STdDZ2tpZEhKaGJuTmhZM1JwYjI0dGFXUWlJRDBnSWpFd01EQXdNREF4T0RFME9USTFOVGNpT3dvSkluRjFZVzUwYVhSNUlpQTlJQ0l4SWpzS0NTSnZjbWxuYVc1aGJDMXdkWEpqYUdGelpTMWtZWFJsTFcxeklpQTlJQ0l4TkRRNE1qazJNelEzTVRNNElqc0tDU0oxYm1seGRXVXRkbVZ1Wkc5eUxXbGtaVzUwYVdacFpYSWlJRDBnSWpWRE1qYzFRalZCTFRORk5qWXROREF5TmkwNU5USXpMVGcwTlVJeU1qVTNOVFZHUXlJN0Nna2ljSEp2WkhWamRDMXBaQ0lnUFNBaWNIVnlZMmhoYzJWeUxtTnZibk4xYldGaWJHVkdaV0YwZFhKbElqc0tDU0pwZEdWdExXbGtJaUE5SUNJeE1EWXhOVFUzTkRnMElqc0tDU0ppYVdRaUlEMGdJbU52YlM1bGN5NVFkWEpqYUdGelpYSWlPd29KSW5CMWNtTm9ZWE5sTFdSaGRHVXRiWE1pSUQwZ0lqRTBORGd5T1RZek5EY3hNemdpT3dvSkluQjFjbU5vWVhObExXUmhkR1VpSUQwZ0lqSXdNVFV0TVRFdE1qTWdNVFk2TXpJNk1qY2dSWFJqTDBkTlZDSTdDZ2tpY0hWeVkyaGhjMlV0WkdGMFpTMXdjM1FpSUQwZ0lqSXdNVFV0TVRFdE1qTWdNRGc2TXpJNk1qY2dRVzFsY21sallTOU1iM05mUVc1blpXeGxjeUk3Q2draWIzSnBaMmx1WVd3dGNIVnlZMmhoYzJVdFpHRjBaU0lnUFNBaU1qQXhOUzB4TVMweU15QXhOam96TWpveU55QkZkR012UjAxVUlqc0tmUT09IjsiZW52aXJvbm1lbnQiID0gIlNhbmRib3giOyJwb2QiID0gIjEwMCI7InNpZ25pbmctc3RhdHVzIiA9ICIwIjt9",
					"",
					"",
				},
			},
			&DoVerifyReceiptRsp{
				Body: struct {
					Status                   int     `json:"status"`
					Receipt                  Receipt `json:"receipt"`
					Environment              string  `json:"environment"`
					LatestReceipt            string  `json:"latest_receipt"`
					LatestReceiptInfo        string  `json:"latest_receipt_info"`
					LatestExpiredReceiptInfo string  `json:"latest_expired_receipt_info"`
					PendingRenewalInfo       string  `json:"pending_renewal_info"`
					IsRetryable              string  `json:"is-retryable"`
				}{
					Status: 0,
					Receipt: Receipt{
						OriginalPurchaseDatePst: "2015-11-23 08:32:27 America/Los_Angeles",
						PurchaseDateMs:          "1448296347138",
						UniqueIdentifier:        "e47425df7fb6ab8c3c7e5730c3a7eb5b6563823d",
						OriginalTransactionId:   "1000000181492557",
						Bvrs:                   "3",
						TransactionId:          "1000000181492557",
						Quantity:               "1",
						UniqueVendorIdentifier: "5C275B5A-3E66-4026-9523-845B225755FC",
						ItemId:                 "1061557484",
						ProductId:              "purchaser.consumableFeature",
						PurchaseDate:           "2015-11-23 16:32:27 Etc/GMT",
						OriginalPurchaseDate:   "2015-11-23 16:32:27 Etc/GMT",
						PurchaseDatePst:        "2015-11-23 08:32:27 America/Los_Angeles",
						Bid:                    "com.es.Purchaser",
						OriginalPurchaseDateMs: "1448296347138",

						ReceiptType:                "",
						AdamId:                     0,
						AppItemId:                  0,
						BundleId:                   "",
						ApplicationVersion:         "",
						DownloadId:                 0,
						VersionExternalIdentifier:  0,
						ReceiptCreationDate:        "",
						ReceiptCreationDateMs:      "",
						ReceiptCreationDatePst:     "",
						RequestDate:                "",
						RequestDateMs:              "",
						RequestDatePst:             "",
						OriginalApplicationVersion: "",
						InApp: nil,
					},
					Environment:              "",
					LatestReceipt:            "",
					LatestReceiptInfo:        "",
					LatestExpiredReceiptInfo: "",
					PendingRenewalInfo:       "",
					IsRetryable:              "",
				},
			},
			false},
		{
			name: "",
			args: args{
				request: DoVerifyReceiptReq{
					ReceiptData:            "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",
					Password:               "",
					ExcludeOldTransactions: "",
				},
			},
			want: &DoVerifyReceiptRsp{
				Body: struct {
					Status                   int     `json:"status"`
					Receipt                  Receipt `json:"receipt"`
					Environment              string  `json:"environment"`
					LatestReceipt            string  `json:"latest_receipt"`
					LatestReceiptInfo        string  `json:"latest_receipt_info"`
					LatestExpiredReceiptInfo string  `json:"latest_expired_receipt_info"`
					PendingRenewalInfo       string  `json:"pending_renewal_info"`
					IsRetryable              string  `json:"is-retryable"`
				}{
					Status: 0,
					Receipt: Receipt{
						PurchaseDateMs:        "",
						UniqueIdentifier:      "",
						OriginalTransactionId: "",
						Bvrs:                       "",
						TransactionId:              "",
						Quantity:                   "",
						UniqueVendorIdentifier:     "",
						ItemId:                     "",
						ProductId:                  "",
						PurchaseDate:               "",
						PurchaseDatePst:            "",
						Bid:                        "",
						ReceiptType:                "ProductionSandbox",
						AdamId:                     0,
						AppItemId:                  0,
						BundleId:                   "cn.lemongrassmedia.loveworld.online",
						ApplicationVersion:         "20290",
						DownloadId:                 0,
						VersionExternalIdentifier:  0,
						ReceiptCreationDate:        "2018-11-23 10:06:27 Etc/GMT",
						ReceiptCreationDateMs:      "1542967587000",
						ReceiptCreationDatePst:     "2018-11-23 02:06:27 America/Los_Angeles",
						RequestDate:                "",
						RequestDateMs:              "",
						RequestDatePst:             "",
						OriginalPurchaseDate:       "2013-08-01 07:00:00 Etc/GMT",
						OriginalPurchaseDateMs:     "1375340400000",
						OriginalPurchaseDatePst:    "2013-08-01 00:00:00 America/Los_Angeles",
						OriginalApplicationVersion: "1.0",
						InApp: []InApp{{
							Quantity:                "1",
							ProductId:               "online.diamond.develop.6",
							TransactionId:           "1000000477111354",
							OriginalTransactionId:   "1000000477111354",
							PurchaseDate:            "2018-11-23 10:06:27 Etc/GMT",
							PurchaseDateMs:          "1542967587000",
							PurchaseDatePst:         "2018-11-23 02:06:27 America/Los_Angeles",
							OriginalPurchaseDate:    "2018-11-23 10:06:27 Etc/GMT",
							OriginalPurchaseDateMs:  "1542967587000",
							OriginalPurchaseDatePst: "2018-11-23 02:06:27 America/Los_Angeles",
							IsTrialPeriod:           "false",
						}},
					},
					Environment:              "Sandbox",
					LatestReceipt:            "",
					LatestReceiptInfo:        "",
					LatestExpiredReceiptInfo: "",
					PendingRenewalInfo:       "",
					IsRetryable:              "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyReceipt(tt.args.request)

			got.Body.Receipt.RequestDate = ""
			got.Body.Receipt.RequestDateMs = ""
			got.Body.Receipt.RequestDatePst = ""

			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyReceipt(%v) error = %v, wantErr %v", tt.args.request, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyReceipt(%v) = %v, want %v", tt.args.request, got, tt.want)
			}
		})
	}
}
