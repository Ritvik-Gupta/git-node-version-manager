package parser

import (
	"errors"
	"regexp"

	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
)

type RawParser struct {
	repos []string
}

func NewRawParser(repos []string) *RawParser {
	return &RawParser{repos}
}

var RAW_REPO_REGEX = regexp.MustCompile(`^(?P<name>.+)=(?P<url>.+)$`)

func (rawParser *RawParser) ParseWriteInto(repositories map[string]utils.Repository) error {
	for _, rawRepo := range rawParser.repos {
		matches := RAW_REPO_REGEX.FindStringSubmatch(rawRepo)
		if matches == nil {
			return errors.New("no match found for raw repository expression")
		}

		readMatch := func(captureName string) string {
			return matches[RAW_REPO_REGEX.SubexpIndex(captureName)]
		}

		name, url := readMatch("name"), readMatch("url")
		if _, exists := repositories[name]; exists {
			return errors.New("conflicting repositories found")
		}

		repositories[name] = utils.Repository{Name: name, Url: url}
	}

	return nil
}
