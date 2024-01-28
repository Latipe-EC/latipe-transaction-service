package createPurchase

import (
	"encoding/json"
)

func ParseOrderToByte(request interface{}) ([]byte, error) {
	jsonObj, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}
	return jsonObj, err
}
