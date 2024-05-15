package utils

import (
	"encoding/json"
	"fmt"
)

func DecodeResponseJSON(jsonData []byte, target interface{}) (interface{}, error) {
	err := json.Unmarshal(jsonData, &target)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON data: %v", err)
	}
	return target, nil
}
