package main

const STD_TIME_FORMAT = "2006-01-02 15:04:05.000000"

const (
	DEFAULT_SERVER_PORT = 8080
	ALPHABETS           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

const (
	RESP_OK     = "ok"
	RESP_ERROR  = "error"
	RESP_NODATA = "no-data"
)

// expiry types
const (
	EXPIRY_MINUTES = "minutes"
	EXPIRY_HOURS   = "hours"
	EXPIRY_DAYS    = "days"
)

// seconds
const (
	MIN_PER_HOUR = 60
	MIN_PER_DAY  = 1440
)

// LinkTypes default,deep,web
const (
	LINK_DEFAULT LinkType = "default"
	LINK_DEEP    LinkType = "deep"
	LINK_WEB     LinkType = "web"
)

// error messages
const (
	UnAuthorized     = "Unauthorized"
	Forbidden        = "Forbidden"
	ResourceNotFound = "Resource not found"
	UnKnownErr       = "Unknown error"
	EndpointErr      = "Unable to connect to endpoint. Kindly try again."
	ShortLinkFailed  = "Creating short link failed, please retry after sometime"
	InvalidShortUrl  = "Invalid short URL"
)

// headers
const (
	APPLICATION_JSON = "application/json"
	APPLICATION_FORM = "application/x-www-form-urlencoded"
	HTTP_UPGRADE     = "Upgrade"
	HTTP_WEBSOCKET   = "websocket"
	IF_NONE_MATCH    = "If-None-Match"
	ETAG             = "ETag"
	OCTET_STREAM     = "application/octet-stream"
)

const (
	LOCAL = "LOCAL"
	DEV   = "DEV"
	QA    = "QA"
	UAT   = "UAT"
	PROD  = "PROD"
)

// android constants
const (
	ANDROID_NAMESPACE = "android_app"
	ANDROID_RELATION  = "delegate_permission/common.handle_all_urls"
)
