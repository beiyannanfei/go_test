package alipay

const (
	K_TIME_FORMAT = "2006-01-02 15:04:05"

	K_ALI_PAY_SANDBOX_API_URL     = "https://openapi.alipaydev.com/gateway.do"
	K_ALI_PAY_PRODUCTION_API_URL  = "https://openapi.alipay.com/gateway.do"
	K_ALI_PAY_PRODUCTION_MAPI_URL = "https://mapi.alipay.com/gateway.do"

	K_FORMAT  = "JSON"
	K_CHARSET = "utf-8"
	K_VERSION = "1.0"

	// https://doc.open.alipay.com/docs/doc.htm?treeId=291&articleId=105806&docType=1
	K_SUCCESS_CODE = "10000"
)

const (
	k_RESPONSE_SUFFIX = "_response"
	k_ERROR_RESPONSE  = "error_response"
	k_SIGN_NODE_NAME  = "sign"
)

const (
	K_SIGN_TYPE_RSA2 = "RSA2"
	K_SIGN_TYPE_RSA  = "RSA"
)

const (
	K_CONTENT_TYPE_FORM = "application/x-www-form-urlencoded;charset=utf-8"
)
