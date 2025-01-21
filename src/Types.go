package main

import (
	"database/sql"
	"net/http"
)

type ServiceHandler func(context *IouHttpContext)
type AppSearch string
type LinkType string
type ExpiryType string

type CmdArgsType struct {
	ConfigPath string
	ServerPort *int
	ServerAddr string
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SslMode  string `yaml:"sslmode"`
	MaxIdle  int    `yaml:"maxIdle"`
	MaxOpen  int    `yaml:"maxOpen"`
}

const (
	APP_SEARCH AppSearch = "appsearch"
	CUSTOM     AppSearch = "custom"
)

type Resp struct {
	S   string `json:"s"`
	Msg string `json:"msg,omitempty"`
	D   any    `json:"d,omitempty"`
}

type IouMsgResp struct {
	S   string `json:"s"`
	Msg string `json:"msg,omitempty"`
}

type Config struct {
	LynxDb    Database `yaml:"lynx-database"`
	BasePath  string   `yaml:"base-path"`
	AppConfig struct {
		ShortLinkLen        int    `yaml:"short-link-len"`
		DuplicateRetryCount int    `yaml:"duplicate-retry-count"`
		DefaultFallbackUrl  string `yaml:"default-fallback-url"`
		WebHtmlFilePath     string `yaml:"web-html-path"`
		HtmlFilePath404     string `yaml:"404-html-path"`
		Android             struct {
			HtmlFilePath         string    `yaml:"html-path"`
			AndroidDefaultWebUrl string    `yaml:"android-default-web-url"`
			GooglePlaySearchUrl  string    `yaml:"google-play-search-url"`
			Behaviour            AppSearch `yaml:"behavior"`
			PackageName          string    `yaml:"package-name"`
			Certificate          []string  `yaml:"sha-certificates"`
		} `yaml:"android"`
		Ios struct {
			HtmlFilePath      string    `yaml:"html-path"`
			IosDefaultWebUrl  string    `yaml:"ios-default-web-url"`
			AppStoreSearchUrl string    `yaml:"app-store-search-url"`
			Behaviour         AppSearch `yaml:"behavior"`
			TeamId            string    `yaml:"team-id"`
			BundleIdentifier  string    `yaml:"bundle-identifier"`
			AppLinkPath       []string  `yaml:"app-link-path"`
		} `yaml:"ios"`
		SocialMedia struct {
			Title        string `yaml:"title"`
			Description  string `yaml:"description"`
			ThumbNailImg string `yaml:"thumbnail-image"`
			ShortIcon    string `yaml:"short-icon"`
		} `yaml:"social-media"`
		BaseUrl string `yaml:"base-url"`
		Desktop struct {
			DefaultUrl string `yaml:"default-url"`
			WindowsUrl string `yaml:"windows-url"`
			MacUrl     string `yaml:"mac-url"`
		} `yaml:"desktop"`
	} `yaml:"app-config"`
}

type LibConfig struct {
	Env  string `yaml:"env"`
	Http struct {
		WriteTimeout int `yaml:"writeTimeout"`
		ReadTimeout  int `yaml:"readTimeout"`
	}
	Log struct {
		Dir   string   `yaml:"dir"`
		File  string   `yaml:"file"`
		Level LogLevel `yaml:"level"`
	}
}

type DatabaseType struct {
	Db       *sql.DB
	Endpoint string
}

type IouHttpContext struct {
	Request    *http.Request
	RespWriter http.ResponseWriter
	Audit      *ServiceAudit
}

type EndPointContext struct {
	*IouHttpContext
	Db           *sql.DB
	EndpointName string
	LogOnce      bool
}

type MobileInputs struct {
	// MinimumVersion string `json:"mv"`
	// WebUrl         string `json:"webUrl"`
	Fbl string `json:"fbl"`
}

type SocialInput struct {
	Title       string `json:"title" validate:"optional"`
	Description string `json:"description" validate:"optional"`
	ImgUrl      string `json:"imgUrl" validate:"optional"`
	Icon        string `json:"shortIcon"`
}

type CreateShortLinkRequest struct {
	Expiry  string       `json:"expiry"`
	WebUrl  string       `json:"webUrl"`
	Data    any          `json:"data"`
	Android MobileInputs `json:"android"`
	Ios     MobileInputs `json:"ios"`
	Social  SocialInput  `json:"social"`
}

type DbShortLink struct {
	ShortCode string
	Data      string
	WebUrl    string
	Android   string
	Ios       string
	Social    string
	Expiry    sql.NullString // in minutes
}

type CreateShortLinkResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	LongUrl     string `json:"longUrl"`
	Expiry      string `json:"expiry"`
	Og          struct {
		Icon  string `json:"icon"`
		Image string `json:"image"`
	} `json:"og"`
}

type IosAppDetails struct {
	AppId string   `json:"appID"`
	Paths []string `json:"paths"`
}

type IosAppLinks struct {
	Apps    []string        `json:"apps"`
	Details []IosAppDetails `json:"details"`
}
type IosAppVerifyResponse struct {
	AppLinks IosAppLinks `json:"applinks"`
}

type AndroidTarget struct {
	NameSpace   string   `json:"namespace"`
	PackageName string   `json:"package_name"`
	Sha256      []string `json:"sha256_cert_fingerprints"`
}

type AndroidVerifyResponse struct {
	Relation []string      `json:"relation"`
	Target   AndroidTarget `json:"target"`
}

type ShortCodeDataResponse struct {
	Input     any               `json:"input"`
	AddParams map[string]string `json:"addParams"`
	ShortCode string            `json:"shortcode"`
}

type LogException struct {
	Type   string `json:"ty,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Status string `json:"st,omitempty"`
}

type Build struct {
	Id        string `json:"id"`
	Details   string `json:"d"`
	Model     string `json:"m"`
	OsVer     string `json:"osv"`
	AppVer    string `json:"appv"`
	Type      string `json:"t"`
	UserAgent string `json:"usrAg"`
}

type ServiceAudit struct {
	Url          string       `json:"url"`
	Method       string       `json:"method"`
	Msg          string       `json:"msg,omitempty"`
	IP           string       `json:"ip,omitempty"`
	TT           int64        `json:"tt"`
	ReqT         string       `json:"reqT,omitempty"`
	ResT         string       `json:"resT,omitempty"`
	SvSt         string       `json:"svSt,omitempty"`
	Err          LogException `json:"ex,omitempty"`
	ErrList      []string     `json:"exLst,omitempty"`
	Source       string       `json:"source,omitempty"`
	Build        Build        `json:"build,omitempty"`
	Sc           string       `json:"sc,omitempty"`
	ReqId        string       `json:"reqId,omitempty"`
	Res          any          `json:"res,omitempty"`
	Fields       any          `json:"fields,omitempty"`
	SkipResponse bool         `json:"-"`
}
