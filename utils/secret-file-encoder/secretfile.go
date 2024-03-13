package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	// Default values for currently assumed secret file keys
	defaultKind      = "Secret"
	defaultNamespace = "CodeNamor"
	defaultType      = "Opaque"

	outputFileWritePermissions os.FileMode = 0644
)

// SecretFile models the fields required in a yaml secret file
// as specified for CAPI deployments
type SecretFile struct {
	// ApiVerison is used in an unknown way, but reasoning from example secret files
	// it should be prefix with a 'v' and have a version tuple following
	ApiVersion string `yaml:"apiVersion"`
	// Kind describes the type of information in the yaml file, "secret" is assumed
	// to be the only appropriate value for this application
	Kind string `yaml:"kind"`
	// Metadata is a structure given below
	Metadata SecretMetadata
	// TypeField is used in an unknown way, but reasoning from example secret files
	// "Opaque" is assumed to be the only appropriate value for this application
	TypeField string `yaml:"type"`
	// Data is a structure given below
	Data SecretData
}

// SecretMetadata describes the identifiers for a the Rancher project associated
// with this secret file
type SecretMetadata struct {
	// Namespace is the owning space in Rancher in which you will find the project
	// this secret file is for
	Namespace string `yaml:"namespace"`
	// Name is the name of the project in Rancher
	Name string `yaml:"name"`
}

// SecretData provides the configuration file data encoded in base64
type SecretData struct {
	// Config stores the configuration file for an environment as a base64 encoded string
	Config string `yaml:"config.json"`
}

// SecretFileContents is alias for contents of file, format agnostic
type SecretFileContents []byte

// Write enables abstracted writing of secret file contents
func (contents *SecretFileContents) Write(w io.Writer) {
	w.Write(*contents)
}

// NewSecretFile constructs and returns a SecretFile object providing any assumed/default
// values as well as errors if any necessary values are not provided
func NewSecretFile(apiVersion string, name string, configEncoding string) SecretFile {
	errs := areNewSecretFileInputsValid(apiVersion, name, configEncoding)

	if len(errs) > 0 {
		panic(errs)
	}

	return SecretFile{
		ApiVersion: apiVersion,
		Kind:       defaultKind,
		Metadata: SecretMetadata{
			Name:      name,
			Namespace: defaultNamespace,
		},
		TypeField: defaultType,
		Data: SecretData{
			Config: configEncoding,
		},
	}
}

// WriteToFile writes the the secret file to the specified secretFilePath
func (secret *SecretFile) WriteToFile(secretFilePath string) {
	outputFile := getOutputFile(secretFilePath)
	defer outputFile.Close()

	secret.Write(outputFile)
}

// Write writes the secret to the output writer: could be file, string, buffer, etc
func (secret *SecretFile) Write(writer io.Writer) {
	secretFileContents := secret.marshalSecretFileContents()
	secretFileContents.Write(writer)
}

func areNewSecretFileInputsValid(apiVersion string, name string, configEncoding string) []error {
	errs := make([]error, 0)
	if apiVersion == "" {
		errs = append(errs, fmt.Errorf("ERROR: cannot create SecretFile, 'apiVersion' cannot be empty"))
	}

	if name == "" {
		errs = append(errs, fmt.Errorf("ERROR: cannot create SecretFile, 'name' cannot be empty"))
	}

	if configEncoding == "" {
		errs = append(errs, fmt.Errorf("ERROR: cannot create SecretFile, 'configEncoding' cannot be empty"))
	}

	return errs
}

func (secret *SecretFile) marshalSecretFileContents() *SecretFileContents {
	var secretFileContents SecretFileContents
	secretFileContents, err := yaml.Marshal(secret)
	if err != nil {
		panic(fmt.Errorf("ERROR: could not create yaml file structure: %s", err))
	}
	return &secretFileContents
}

func getOutputFile(path string) *os.File {
	outputFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, outputFileWritePermissions)
	if err != nil {
		panic(err)
	}

	return outputFile
}
