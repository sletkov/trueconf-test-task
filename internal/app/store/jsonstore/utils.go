package jsonstore

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
)

func getStoreFromJSONFile(path string) (*UserStore, error) {
	var s *UserStore

	f, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(f, &s)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func writeStoreIntoJSONFile(s *UserStore, path string) error {
	b, err := json.Marshal(&s)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile("users.json", b, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
