package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validateShouldNotPanicWhenFlagsValid(t *testing.T) {
	version := "v1"
	p := "PROJECT"
	i := "CONFIG.JSON"
	o := "SECRET.YAML"

	mockInputFlags := SecretFileEncoderInputFlags{
		apiVersion:     &version,
		projectName:    &p,
		configFilePath: &i,
		outputFilePath: &o,
	}

	require.NotPanics(t, mockInputFlags.validate)
}

func Test_validateShouldPanicForEmptyProjectName(t *testing.T) {
	version := "v1"
	p := ""
	i := "CONFIG.JSON"
	o := "SECRET.YAML"

	mockInputFlags := SecretFileEncoderInputFlags{
		apiVersion:     &version,
		projectName:    &p,
		configFilePath: &i,
		outputFilePath: &o,
	}

	require.Panics(t, mockInputFlags.validate)
}

func Test_validateShouldPanicForEmptyInputConfigFile(t *testing.T) {
	version := "v1"
	p := "PROJECT"
	i := ""
	o := "SECRET.YAML"

	mockInputFlags := SecretFileEncoderInputFlags{
		apiVersion:     &version,
		projectName:    &p,
		configFilePath: &i,
		outputFilePath: &o,
	}

	require.Panics(t, mockInputFlags.validate)
}
