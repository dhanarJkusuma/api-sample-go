package helper

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnMarshall(body io.ReadCloser, obj interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}
	return nil
}
