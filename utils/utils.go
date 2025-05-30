package utils

import (
	"encoding/json"
	"net/http"
)

func EncodeJSONBody(res http.ResponseWriter, body interface{}) error {
	err := json.NewEncoder(res).Encode(body)
	if err != nil {
		return err
	}
	return nil
}

func DecodeJSONBody(req *http.Request, out interface{}) error {
	err := json.NewDecoder(req.Body).Decode(&out)

	if err != nil {
		return err
	}
	return nil

}
