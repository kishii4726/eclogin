//go:build ci

package session

import (
	"encoding/json"
)

// StartSession は CI 環境でのテスト用にモック化されたバージョンです
func StartSession(sessionData, inputData []byte, region string) error {
	// セッションデータとインプットデータが有効な JSON であることを確認
	var session, input interface{}
	if err := json.Unmarshal(sessionData, &session); err != nil {
		return err
	}
	if err := json.Unmarshal(inputData, &input); err != nil {
		return err
	}
	return nil
}
