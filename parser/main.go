package parser

import (
	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
)

type Parser interface {
	ParseWriteInto(map[string]utils.Repository) error
}
