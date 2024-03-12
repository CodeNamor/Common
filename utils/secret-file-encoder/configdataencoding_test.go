package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_encodeConfigDataFromReaderSourceShouldCorrectlyEncodeSourceData(t *testing.T) {
	testConfigFilePath := "./testfiles/example_config.json"
	mockConfigData, err := ioutil.ReadFile(testConfigFilePath)

	require.NoErrorf(t, err, "Test config file is missing from '%s'", testConfigFilePath)

	configDataReader := bytes.NewReader(mockConfigData)
	encodedConfigDataResult := encodeConfigDataFromReaderSource(configDataReader)

	result, _ := base64.RawStdEncoding.WithPadding('=').DecodeString(encodedConfigDataResult)

	// This is effectively an encode/decode test where the original file contents are compared to the
	// end result of encoding an decoding the contents. Removed the file comparison because line encodings
	// between systems appeared to upset base64 comparisons
	require.Equal(t, mockConfigData, result)
}
