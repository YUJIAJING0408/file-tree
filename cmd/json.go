package cmd

import (
	"encoding/json"
	"os"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

func SaveToJson(filePath string, data any) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := json.Marshal(data)
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
