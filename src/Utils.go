package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

func loadConfigFile[S any](configPath string, config *S) {
	buf, err := os.ReadFile(configPath)
	if err != nil {
		log.Println(err)
	}

	err = yaml.Unmarshal(buf, config)
	if err != nil {
		log.Println(err)
	}
}

func LoadLibConfig(configPath string) {
	buf, err := os.ReadFile(configPath)
	if err != nil {
		log.Println(err)
	}

	libConfig = &LibConfig{}
	err = yaml.Unmarshal(buf, libConfig)
	if err != nil {
		log.Println(err)
	}
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		err := errors.New(path + " is a directory, not a normal file")
		return err
	}
	return nil
}

func ParseCmdArgs(CmdArgs *CmdArgsType) error {
	flag.StringVar(&CmdArgs.ConfigPath, "config", "config.yml", "path to config file")
	CmdArgs.ServerPort = flag.Int("server-port", DEFAULT_SERVER_PORT, "server port number")
	flag.Parse()

	CmdArgs.ServerAddr = fmt.Sprintf(":%d", *CmdArgs.ServerPort)

	if err := ValidateConfigPath(CmdArgs.ConfigPath); err != nil {
		return err
	}
	return nil
}

func InitializeConfigs[S any](configRef *S) {
	CmdArgs = &CmdArgsType{}
	err := ParseCmdArgs(CmdArgs)
	if err != nil {
		log.Fatal(err)
	}
	LoadLibConfig(CmdArgs.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	InitLogging()

	loadConfigFile(CmdArgs.ConfigPath, configRef)
}

func BaseMW(handler ServiceHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, rq *http.Request) {
		var context = &IouHttpContext{Request: rq, RespWriter: rw}
		handler(context)
	}
}

func Base62Encode(number uint64) string {
	if number == 0 {
		return string(ALPHABETS[0])
	}

	chars := make([]byte, 0)
	length := uint64(len(ALPHABETS))

	for number > 0 {
		result := number / length
		remainder := number % length
		chars = append([]byte{ALPHABETS[remainder]}, chars...)
		number = result
	}

	return string(chars)
}

func getRandomId() int64 {
	var t RandomIdGenerator = &RandomID{}
	t.IdGenerator()
	id, _ := t.GetId()
	return id
}

func genShortUrl() string {
	iSnowflakeId := getRandomId()
	sEncodedURL := Base62Encode(uint64(iSnowflakeId))
	var startInd = len(sEncodedURL) - config.AppConfig.ShortLinkLen
	return sEncodedURL[startInd:]
}

func JSONMarshal(t any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func CurrentTime() string {
	return time.Now().Format(STD_TIME_FORMAT)
}

// returns in minutes
func calculateExpiry(exp string) int64 {
	var value int64 = 0
	return value
}

func isDuplicateLink(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == MYSQL_DUPLICATE_INDEX {
			return true
		}
	}
	return false
}

func isAndroidWeb(cxt *IouHttpContext) bool {
	var userAgent = cxt.Request.Header.Get("User-Agent")
	if strings.Contains(userAgent, "Android") {
		return false
	}
	return true
}

func isIosWeb(cxt *IouHttpContext) bool {
	var userAgent = cxt.Request.Header.Get("User-Agent")
	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") {
		return false
	}
	return true
}

// func isDesktopWeb(cxt *IouHttpContext) bool {
// 	var userAgent = cxt.Request.Header.Get("User-Agent")
// 	if strings.Contains(userAgent, "Mobile") {
// 		return false
// 	}
// 	return true
// }

func isMap(value any) bool {
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Map
}

func getDataString(data any) (string, error) {
	if isMap(data) {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}
	return fmt.Sprint(data), nil
}

func getValueOrDefault(value string, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func loadHtmlFile() error {
	file, err := os.ReadFile(config.AppConfig.HtmlFilePath)
	if err != nil {
		return err
	}
	htmlCache = string(file)
	file, err = os.ReadFile(config.AppConfig.HtmlFilePath404)
	if err != nil {
		return err
	}
	htmlCache404 = string(file)
	return nil
}

func frame404WebPage() string {
	replacements := map[string]string{
		"{TITLE}":             config.AppConfig.SocialMedia.Title,
		"{DESCRIPTION}":       config.AppConfig.SocialMedia.Description,
		"{URL_CONTENT}":       config.AppConfig.BaseUrl,
		"{IMAGE_CONTENT}":     config.AppConfig.SocialMedia.ThumbNailImg,
		"{REDIRECT_LOCATION}": config.AppConfig.DefaultFallbackUrl,
		"{ICON}":              config.AppConfig.SocialMedia.ShortIcon,
	}
	for key, val := range replacements {
		htmlCache404 = strings.ReplaceAll(htmlCache404, key, val)
	}
	return htmlCache404
}

func frameWebPage(data DbShortLink) string {
	var social SocialInput
	json.Unmarshal([]byte(data.Social), &social)
	var htmlFile = htmlCache
	replacements := map[string]string{
		"{TITLE}":             getValueOrDefault(social.Title, config.AppConfig.SocialMedia.Title),
		"{DESCRIPTION}":       getValueOrDefault(social.Description, config.AppConfig.SocialMedia.Description),
		"{URL_CONTENT}":       config.AppConfig.BaseUrl,
		"{IMAGE_CONTENT}":     getValueOrDefault(social.ImgUrl, config.AppConfig.SocialMedia.ThumbNailImg),
		"{REDIRECT_LOCATION}": frameCompleteUrl(data),
		"{ICON}":              config.AppConfig.SocialMedia.ShortIcon,
	}
	for key, val := range replacements {
		htmlFile = strings.ReplaceAll(htmlFile, key, val)
	}
	return htmlFile
}

func isValidJson(data string) bool {
	var a any
	raw := json.RawMessage(data)
	json.Unmarshal(raw, &a)
	return a != nil
}
