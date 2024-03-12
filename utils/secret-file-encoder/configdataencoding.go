package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// NewConfigDataEncodingFromFile reads a config file from the given path and
// returns the base64 encoding of that config file
func NewConfigDataEncodingFromFile(configFilePath string) string {
	configFile := getConfigFile(configFilePath)
	return encodeConfigDataFromReaderSource(configFile)
}

func getConfigFile(configFilePath string) *os.File {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		panic(fmt.Errorf("ERROR: fail to read input config file '%s':'%s'", configFilePath, err))
	}
	return configFile
}

func encodeConfigDataFromReaderSource(configDataSource io.Reader) string {
	configData := getConfigData(configDataSource)
	encodedConfigFileContents := encodeConfigData(configData)

	return encodedConfigFileContents
}

// getConfigData abstracts the reading of config data from files for testing purposes
func getConfigData(configDataSource io.Reader) []byte {
	configData, err := ioutil.ReadAll(configDataSource)

	if err != nil {
		panic(err)
	}
	return configData
}

// encodeConfigData returns the input byte array encoded as a padded base64 string.
// Note the padding will append '=' to the end of the encoding as necessary to line
// up the bytes. Padding '=' replicates the default behavior of base64 command
// line tool originally used in secret file generation.
func encodeConfigData(contents []byte) string {
	return base64.RawStdEncoding.WithPadding('=').EncodeToString(contents)
}
