package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func Test_NewSecretFileShouldCreateSecretFileWhenRequiredValuesArePresent(t *testing.T) {
	require.NotPanics(t, func() {
		NewSecretFile("V1", "NAME", "CONFIG")
	})

	secretFileResult := NewSecretFile("V1", "NAME", "CONFIG")

	require.Equal(t, secretFileResult.ApiVersion, "V1")
	require.Equal(t, secretFileResult.Metadata.Name, "NAME")
	require.Equal(t, secretFileResult.Data.Config, "CONFIG")
}

func Test_NewSecretFileShouldPanicOnEmptyStringForApiVersion(t *testing.T) {
	require.Panics(t, func() {
		NewSecretFile("", "NAME", "CONFIG")
	})
}

func Test_NewSecretFileShouldPanicOnEmptyStringForName(t *testing.T) {
	require.Panics(t, func() {
		NewSecretFile("V1", "", "CONFIG")
	})
}

func Test_NewSecretFileShouldPanicOnEmptyStringForEncoding(t *testing.T) {
	require.Panics(t, func() {
		NewSecretFile("V1", "NAME", "")
	})
}

func Test_WriteToShouldEncodeAndWriteSecretFileCorrectly(t *testing.T) {
	mockSecretFile := NewSecretFile("v1.2.1", "NAME", "CONFIG")
	stringBuffer := bytes.NewBufferString("")
	mockSecretFile.Write(stringBuffer)

	resultSecretFile := SecretFile{}
	yaml.Unmarshal(stringBuffer.Bytes(), &resultSecretFile)

	if !reflect.DeepEqual(mockSecretFile, resultSecretFile) {
		t.Fatal(fmt.Sprintf("EXPECTED: %v\nGOT: %v", mockSecretFile, resultSecretFile))
	}
}
