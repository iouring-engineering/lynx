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
