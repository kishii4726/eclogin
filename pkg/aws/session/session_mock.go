//go:build ci

package session

import (
	"encoding/json"
)

func StartSession(sessionData, inputData []byte, region string) error {
	var session, input interface{}
	if err := json.Unmarshal(sessionData, &session); err != nil {
		return err
	}
	if err := json.Unmarshal(inputData, &input); err != nil {
		return err
	}
	return nil
}
