package parser

import (
	"encoding/csv"
	"errors"
	"os"

	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
)

type CsvParser struct {
	fileName string
}

func NewCsvParser(fileName string) *CsvParser {
	return &CsvParser{fileName}
}

func (csvParser *CsvParser) ParseWriteInto(repositories map[string]utils.Repository) error {
	file, err := os.Open("./inputs/" + csvParser.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		name := line[0]
		if _, exists := repositories[name]; exists {
			return errors.New("conflicting repositories found")
		}

		repositories[name] = utils.Repository{Name: name, Url: line[1]}
	}
	return nil
}
