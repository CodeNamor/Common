package transform

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// DecodeJSONToStruct fills in a defaulted structure with
// values from decoded JSON. If values are missing from the JSON
// then the existing defaults will be used. targetStructPtr should
// be a pointer to a structure with defaults already applied.
// You would typically use a factory function to create a defaulted
// Struct and then load from JSON to overwrite some of the values.
func DecodeJSONToStruct(bytes []byte, targetStructPtr interface{}) error {
	//deserialize to a map so we can determine what exactly was loaded
	m := make(map[string]interface{})
	err := json.Unmarshal(bytes, &m)

	if err != nil {
		return err
	}

	// apply map to struct fields, make sure all keys are used
	decoderConfig := &mapstructure.DecoderConfig{
		ErrorUnused: true, // error if extra (mistyped) fields are found
		Result:      targetStructPtr,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)

	if err != nil {
		return err
	}

	err = decoder.Decode(m)

	if err != nil {
		return errors.Wrap(err, "errored trying to match JSON data to struct")
	}

	return nil
}
