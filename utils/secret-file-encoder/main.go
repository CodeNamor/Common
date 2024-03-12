package main

// go install
// secret-file-encoder is the installed executable name

// example usage:
// secret-file-encoder -p omni-membersearch-api -i ./testfiles/example_config.json -o ./testfiles/example_secret.yaml -version v1.22.3

// for help
// secret-file-encoder -h

import (
	"fmt"
)

func main() {
	defer handlePanics()

	inputs := NewValidSecretFileEncoderInputFlags()

	encodedConfigData := NewConfigDataEncodingFromFile(*inputs.configFilePath)

	secretFile := NewSecretFile(*inputs.apiVersion, *inputs.projectName, encodedConfigData)
	secretFile.WriteToFile(*inputs.outputFilePath)

	fmt.Printf("Secret successfully written to '%s'\n", *inputs.outputFilePath)
}

// handlePanic is a general error top-level panic handler to end the program
// when it cannot proceed while providing failure feedback.
func handlePanics() {
	if r := recover(); r != nil {
		fmt.Println("FAILED TO WRITE SECRET FILE")
		switch rValue := r.(type) {
		case []error:
			for _, err := range rValue {
				fmt.Println(err)
			}
		case error:
			fmt.Println(rValue)
		}
	}
}
