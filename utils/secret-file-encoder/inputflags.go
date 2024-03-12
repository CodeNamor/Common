package main

import (
	"errors"
	"flag"
)

// InputFlags provides a structure for storing all flags input to the executable,
// encapsulating all input data to avoid future parameter changes.
type SecretFileEncoderInputFlags struct {
	apiVersion     *string
	projectName    *string
	configFilePath *string
	outputFilePath *string
}

// validate ensures all required flags, those without default values, have values present
// and panics otherwise.
func (flags *SecretFileEncoderInputFlags) validate() {
	errs := make([]error, 0)

	if *flags.projectName == "" {
		errs = append(errs, errors.New("ERROR: '-p' project name is required (e.g., '-p=omni-membersearch-api')"))
	}

	if *flags.configFilePath == "" {
		errs = append(errs, errors.New("ERROR: '-i' input config file is required (e.g., '-i=dev_config.json')"))
	}

	if len(errs) > 0 {
		panic(errs)
	}

	return
}

// NewValidSecretFileEncoderInputFlags parses and validates the flag values input to the secret file encoder
// executable.
func NewValidSecretFileEncoderInputFlags() SecretFileEncoderInputFlags {
	inputFlags := SecretFileEncoderInputFlags{
		apiVersion:     flag.String("version", "v1", "'v' followed by versioning tuple"),
		projectName:    flag.String("p", "", "project name (e.g., \"omni-membersearch-api\")"),
		configFilePath: flag.String("i", "", "input config file, can be prefixed with relative path (e.g., \"dev_config.json\")"),
		outputFilePath: flag.String("o", "secret.yaml", "output secret file, can be prefixed with relative path"),
	}

	flag.Parse()
	inputFlags.validate()

	return inputFlags
}
