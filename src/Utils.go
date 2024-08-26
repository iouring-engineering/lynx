package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"slices"
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
		startTime := time.Now()

		var audit = &ServiceAudit{}
		audit.ReqT = CurrentTime()

		audit.IP = rq.Header.Get(INTERNAL_X_REAL_IP)
		audit.ReqId = rq.Header.Get(INTERNAL_X_REQ_ID)
		audit.Method = rq.Method
		audit.Url = rq.URL.Path
		audit.Source = rq.Header.Get(INTERNAL_X_SOURCE)

		audit.SetBuild(rq.Header.Get(INTERNAL_X_BUILD), rq.Header.Get("User-Agent"))
		audit.SetFields(rq)

		var context = &IouHttpContext{Request: rq, RespWriter: rw, Audit: audit}

		handler(context)

		audit.SetSC(http.StatusOK)
		audit.LogAudit(startTime)
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
// 18-01-2024 17:57:00 date format
func validateExpiry(exp *string) bool {
	if *exp == "" {
		return true
	}
	t, err := time.ParseInLocation("02-01-2006 15:04:05", *exp, time.Local)
	if err != nil {
		return false
	}
	*exp = t.Format("2006-01-02 15:04:05")
	return true
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
	return strings.Contains(userAgent, "Android")
}

func isIosWeb(cxt *IouHttpContext) bool {
	var userAgent = cxt.Request.Header.Get("User-Agent")
	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") {
		return true
	}
	return false
}

func isMap(value any) bool {
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Map
}

func getDataString(data any) (string, error) {
	if data != nil && isMap(data) {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}
	if data == nil {
		return "", nil
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
	file, err := os.ReadFile(config.AppConfig.WebHtmlFilePath)
	if err != nil {
		return err
	}
	webHtmlCache = string(file)
	Logger.Info("loaded web html")
	file, err = os.ReadFile(config.AppConfig.Android.HtmlFilePath)
	if err != nil {
		return err
	}
	androidHtmlCache = string(file)
	Logger.Info("loaded android html")
	file, err = os.ReadFile(config.AppConfig.Ios.HtmlFilePath)
	if err != nil {
		return err
	}
	iosHtmlCache = string(file)
	Logger.Info("loaded ios html")
	return nil
}

func frameAndroidWebPage(data DbShortLink, link string) string {
	var social SocialInput
	json.Unmarshal([]byte(data.Social), &social)
	var htmlFile = androidHtmlCache
	replacements := map[string]string{
		"{TITLE}":             getValueOrDefault(social.Title, config.AppConfig.SocialMedia.Title),
		"{DESCRIPTION}":       getValueOrDefault(social.Description, config.AppConfig.SocialMedia.Description),
		"{URL_CONTENT}":       link,
		"{IMAGE_CONTENT}":     getValueOrDefault(social.ImgUrl, config.AppConfig.SocialMedia.ThumbNailImg),
		"{REDIRECT_LOCATION}": link,
		"{ICON}":              config.AppConfig.SocialMedia.ShortIcon,
	}
	for key, val := range replacements {
		htmlFile = strings.ReplaceAll(htmlFile, key, val)
	}
	return htmlFile
}

func frameIosWebPage(data DbShortLink, link, shortCode string, utm map[string]string,
	otherParams map[string]string) string {
	var social SocialInput
	json.Unmarshal([]byte(data.Social), &social)
	var htmlFile = iosHtmlCache
	utm["shortcode"] = shortCode
	values := url.Values{}
	for key, value := range utm {
		values.Add(key, value)
	}
	for key, value := range otherParams {
		values.Add(key, value)
	}
	var finalStringToCopy = values.Encode()
	replacements := map[string]string{
		"{TITLE}":             getValueOrDefault(social.Title, config.AppConfig.SocialMedia.Title),
		"{DESCRIPTION}":       getValueOrDefault(social.Description, config.AppConfig.SocialMedia.Description),
		"{URL_CONTENT}":       link,
		"{IMAGE_CONTENT}":     getValueOrDefault(social.ImgUrl, config.AppConfig.SocialMedia.ThumbNailImg),
		"{REDIRECT_LOCATION}": link,
		"{ICON}":              config.AppConfig.SocialMedia.ShortIcon,
		"{SHORT_CODE}":        finalStringToCopy,
	}
	for key, val := range replacements {
		htmlFile = strings.ReplaceAll(htmlFile, key, val)
	}
	return htmlFile
}

func frameWebPage(cxt *IouHttpContext, data DbShortLink, utm map[string]string, otherParams map[string]string) string {
	var social SocialInput
	json.Unmarshal([]byte(data.Social), &social)
	var htmlFile = webHtmlCache
	replacements := map[string]string{
		"{TITLE}":             getValueOrDefault(social.Title, config.AppConfig.SocialMedia.Title),
		"{DESCRIPTION}":       getValueOrDefault(social.Description, config.AppConfig.SocialMedia.Description),
		"{URL_CONTENT}":       config.AppConfig.BaseUrl,
		"{IMAGE_CONTENT}":     getValueOrDefault(social.ImgUrl, config.AppConfig.SocialMedia.ThumbNailImg),
		"{REDIRECT_LOCATION}": frameCompleteUrl(cxt, data, utm, otherParams),
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

func IsFormRequest(req *http.Request) bool {
	return strings.EqualFold(req.Header.Get("content-type"), APPLICATION_FORM)
}

func IsJsonRequest(req *http.Request) bool {
	return strings.EqualFold(req.Header.Get("content-type"), APPLICATION_JSON)
}

func IsOctetStream(req *http.Request) bool {
	return strings.EqualFold(req.Header.Get("content-type"), OCTET_STREAM)
}

func GetDtInTime(inputFrmt string, date string) time.Time {
	parseTime, _ := time.ParseInLocation(inputFrmt, date, time.Local)
	return parseTime
}

func GetQueryParams(req http.Request) map[string]string {
	var result map[string]string = make(map[string]string)
	for key := range req.URL.Query() {
		if !slices.Contains(validUtmParams, key) {
			var value = req.URL.Query().Get(key)
			result[key] = value
		}
	}
	return result
}

func GetUtmParams(req http.Request) map[string]string {
	var result map[string]string = make(map[string]string)
	for _, param := range validUtmParams {
		if req.URL.Query().Has(param) {
			var value = req.URL.Query().Get(param)
			result[param] = value
		}
	}
	return result
}
