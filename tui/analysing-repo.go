package tui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"

	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
)

func (data *TuiApplication) analyseRepo(
	repoName string,
	repoIdx int,
	progressChannel chan<- ProgressFlag,
	downloadStepWaiter *sync.WaitGroup,
) {
	repository := data.repositories[repoName]

	progressChannel <- ProgressFlag{repoIdx: repoIdx, progress: 0, status: "Downloading Repo"}

	fetchRepo(repository)

	progressChannel <- ProgressFlag{repoIdx: repoIdx, progress: 1, status: "Downloaded Repo"}

	downloadStepWaiter.Done()
	downloadStepWaiter.Wait()

	progressChannel <- ProgressFlag{repoIdx: repoIdx, progress: 1, status: "Reading package.json"}

	file, err := os.Open("./" + repository.Name + "/package.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	jsonContent := string(bytes)

	progressChannel <- ProgressFlag{repoIdx: repoIdx, progress: 2, status: "Finding Packages"}

	for idx, pkg := range data.packagesRequired {
		packageRegex := regexp.MustCompile(
			fmt.Sprintf(`"%s":\s*"\^(?P<major>\d+)\.(?P<feature>\d+)\.(?P<minor>\d+)"`, pkg.Name),
		)

		matches := packageRegex.FindStringSubmatch(jsonContent)
		if matches == nil {
			continue
		}

		readMatch := func(captureName string) string {
			return matches[packageRegex.SubexpIndex(captureName)]
		}

		data.reposWriteLock.Lock()
		repository.Packages[idx].Name = pkg.Name
		repository.Packages[idx].MajorVersion, _ = strconv.ParseUint(readMatch("major"), 10, 64)
		repository.Packages[idx].FeatureVersion, _ = strconv.ParseUint(readMatch("feature"), 10, 64)
		repository.Packages[idx].MinorVersion, _ = strconv.ParseUint(readMatch("minor"), 10, 64)
		data.reposWriteLock.Unlock()
	}

	progressChannel <- ProgressFlag{repoIdx: repoIdx, progress: 3, status: "Completed", done: true}
}

func fetchRepo(repository *utils.RepoWithPackagesFound) {
	cmd := exec.Command("git", "clone", repository.Url, repository.Name)

	if err := cmd.Run(); err != nil {
		return
	}
}
