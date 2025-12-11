package cmd

import (
	"gopkg.in/yaml.v3"
	"os"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

func SaveToYaml(filePath string, data any) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := yaml.Marshal(data)
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
