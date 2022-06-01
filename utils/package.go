package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Package struct {
	Name                                       string
	MajorVersion, FeatureVersion, MinorVersion uint64
}

var PACKAGE_STR_REGEX = regexp.MustCompile(`^(?P<name>.+)@(?P<major>\d+)\.(?P<feature>\d+)\.(?P<minor>\d+)$`)

func ParsePackage(packageStr string) (pkg Package, err error) {
	matches := PACKAGE_STR_REGEX.FindStringSubmatch(packageStr)
	if matches == nil {
		err = errors.New("no match found for raw package expression")
		return
	}

	readMatch := func(captureName string) string {
		return matches[PACKAGE_STR_REGEX.SubexpIndex(captureName)]
	}

	pkg.Name = readMatch("name")

	pkg.MajorVersion, err = strconv.ParseUint(readMatch("major"), 10, 64)
	if err != nil {
		return
	}

	pkg.FeatureVersion, err = strconv.ParseUint(readMatch("feature"), 10, 64)
	if err != nil {
		return
	}

	pkg.MinorVersion, err = strconv.ParseUint(readMatch("minor"), 10, 64)
	if err != nil {
		return
	}

	return
}

func (pkg *Package) GetVersion() string {
	if pkg.Name == "" {
		return ""
	}

	return fmt.Sprintf("%d.%d.%d", pkg.MajorVersion, pkg.FeatureVersion, pkg.MinorVersion)
}

func (pkg *Package) VersionRankComparedTo(pkgReq *Package) int {
	if pkg.Name == "" {
		return -4
	}

	if pkg.MajorVersion > pkgReq.MajorVersion {
		return 3
	} else if pkg.MajorVersion < pkgReq.MajorVersion {
		return -3
	}

	if pkg.FeatureVersion > pkgReq.FeatureVersion {
		return 2
	} else if pkg.FeatureVersion < pkgReq.FeatureVersion {
		return -2
	}

	if pkg.MinorVersion > pkgReq.MinorVersion {
		return 1
	} else if pkg.MinorVersion < pkgReq.MinorVersion {
		return -1
	}

	return 0
}
