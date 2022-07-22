package jsonutils

import (
	"common/util/fileutils"
	"encoding/json"

	"io/ioutil"
)

func WriteJsonPretty(v interface{}, fpath string) error {

	if data, err := json.MarshalIndent(v, "", "\t"); err != nil {
		return err
	} else {

		return fileutils.WriteFile(fpath, data)
	}
}

func UNMarshalFromFile(v interface{}, fpath string) (err error) {

	if data, err := ioutil.ReadFile(fpath); err != nil {
		return err
	} else {

		return json.Unmarshal(data, v)
	}
}

func ToJsonString(v interface{}, pretty bool) string {

	var data []byte
	var err error

	if pretty {

		data, err = json.MarshalIndent(v, "", "\t")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return ""
	}

	return string(data)
}
