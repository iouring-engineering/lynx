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
	APP_SEARCH AppSearch = "APP_SEARCH"
	CUSTOM     AppSearch = "CUSTOM"
)

type Config struct {
	LynxDb    Database `yaml:"lynx-database"`
	AppConfig struct {
		ShortLinkLen int    `yaml:"short-link-len"`
		DefaultUrl   string `yaml:"default-url"`
		Android      struct {
			HasApp              bool      `yaml:"has-app"`
			AndroidUrl          string    `yaml:"android-url"`
			AppUriScheme        string    `yaml:"app-uri-scheme"`
			Behaviour           AppSearch `yaml:"behavior"`
			GooglePlaySearchUrl string    `yaml:"google-play-search-url"`
			CustomUrl           struct {
				FallbackUrl string `yaml:"fallback-url"`
				PackageName string `yaml:"package-name"`
			} `yaml:"custom-url"`
			AppLinks struct {
				Enable      bool     `yaml:"enable"`
				Certificate []string `yaml:"certificate"`
			} `yaml:"app-links"`
		} `yaml:"android"`
		Ios struct {
			HasApp            bool      `yaml:"has-app"`
			IosUrl            string    `yaml:"ios-url"`
			AppUriScheme      string    `yaml:"app-uri-scheme"`
			Behaviour         AppSearch `yaml:"behavior"`
			AppStoreSearchUrl string    `yaml:"app-store-search-url"`
			CustomUrl         struct {
				FallbackUrl string `yaml:"fallback-url"`
				PackageName string `yaml:"app-store-id"`
			} `yaml:"custom-url"`
			UniversalLinks struct {
				Enable            bool     `yaml:"enable"`
				BundleIdentifiers []string `yaml:"bundle-identifiers"`
				AppPrefix         string   `yaml:"app-prefix"`
			}
		} `yaml:"ios"`
		SocialMedia struct {
			Title        string `yaml:"title"`
			Description  string `yaml:"description"`
			ThumbNailImg string `yaml:"thumbnail-image"`
		} `yaml:"social-media"`
		Domain  string `yaml:"domain"`
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
		Dir   string  `yaml:"dir"`
		File  string  `yaml:"file"`
		Level LogType `yaml:"level"`
	}
}

type DatabaseType struct {
	Db       *sql.DB
	Endpoint string
}

type IouHttpContext struct {
	Request    *http.Request
	RespWriter http.ResponseWriter
}

type EndPointContext struct {
	*IouHttpContext
	Db           *sql.DB
	EndpointName string
}

type CreateShortLinkRequest struct {
	Expiry struct {
		Type  ExpiryType `json:"type" enums:"minutes,hours,days" validate:"required"`
		Value int        `json:"value" validate:"required"`
	} `json:"expiry"`
	WebUrl  string `json:"webUrl"`
	Data    string `json:"data"`
	Android struct {
		Type LinkType `json:"type" enums:"default,deep,web" validate:"required"`
		Url  LinkType `json:"url"`
		Fbl  LinkType `json:"fbl"`
	} `json:"android"`
	Ios struct {
		Type LinkType `json:"type" enums:"default,deep,web" validate:"required"`
		Url  LinkType `json:"url"`
		Fbl  LinkType `json:"fbl"`
	} `json:"ios"`
	Desktop struct {
		Type LinkType `json:"type" enums:"default,web" validate:"required"`
		Url  LinkType `json:"url"`
	} `json:"desktop"`
	Social struct {
		Title       string `json:"title" validate:"optional"`
		Description string `json:"description" validate:"optional"`
		ImgUrl      string `json:"imgUrl" validate:"optional"`
	} `json:"social"`
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
