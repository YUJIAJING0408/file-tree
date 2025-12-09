package cmd

import (
	"encoding/csv"
	"os"
)

func NewCSV(filePath string, header []string, data [][]string) {
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(header)
	writer.WriteAll(data)

}
