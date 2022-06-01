package tui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TuiApplication struct {
	application      *tview.Application
	pages            *tview.Pages
	repositories     map[string]*utils.RepoWithPackagesFound
	packagesRequired []utils.Package
	lock             *sync.Mutex
}

func NewTuiApplicaition(repositories map[string]utils.Repository, packagesRequired []utils.Package) *TuiApplication {
	reposWithPkgs := make(map[string]*utils.RepoWithPackagesFound, len(repositories))
	for key, repo := range repositories {
		reposWithPkgs[key] = &utils.RepoWithPackagesFound{
			Repository: repo,
			Packages:   make([]utils.Package, len(packagesRequired)),
		}
	}

	return &TuiApplication{
		application:      tview.NewApplication(),
		pages:            tview.NewPages(),
		repositories:     reposWithPkgs,
		packagesRequired: packagesRequired,
		lock:             &sync.Mutex{},
	}
}

func (data *TuiApplication) Start() {
	defer clearDownloadedRepos()

	currentPage := 0

	data.pages.AddPage(
		fmt.Sprint(currentPage),
		data.createProcessingTable(),
		true,
		true,
	)

	data.application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			data.application.Stop()
			return nil
		case tcell.KeyRight:
			currentPage = (currentPage + 1) % data.pages.GetPageCount()
			data.pages.SwitchToPage(fmt.Sprint(currentPage))
			return nil
		case tcell.KeyLeft:
			currentPage = (data.pages.GetPageCount() + currentPage - 1) % data.pages.GetPageCount()
			data.pages.SwitchToPage(fmt.Sprint(currentPage))
			return nil
		}
		return event
	})

	if err := data.application.SetRoot(data.pages, true).SetFocus(data.pages).Run(); err != nil {
		panic(err)
	}
}

type ProgressFlag struct {
	repoIdx  int
	progress int
	status   string
	done     bool
}

func (data *TuiApplication) createProcessingTable() *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorders(true)
	table.SetBorderPadding(0, 0, 20, 20)

	table.SetTitle("Repositories")
	table.SetCell(0, 0, utils.CrateCell("Name", tcell.ColorYellow))
	table.SetCell(0, 1, utils.CrateCell("Repo", tcell.ColorYellow))
	table.SetCell(0, 2, utils.CrateCell("Progress", tcell.ColorYellow))
	table.SetCell(0, 3, utils.CrateCell("Status", tcell.ColorYellow))

	idx := 1
	progressbars := make([]*tview.TableCell, 0, len(data.repositories))
	statusbars := make([]*tview.TableCell, 0, len(data.repositories))

	for _, repo := range data.repositories {
		table.SetCell(idx, 0, utils.CrateNormalCell(repo.Name))
		table.SetCell(idx, 1, utils.CrateCell(repo.Url, tcell.ColorCornflowerBlue))

		progressbar := utils.CrateCell("", tcell.ColorGreen)
		progressbars = append(progressbars, progressbar)
		table.SetCell(idx, 2, progressbar)

		statusbar := utils.CrateCell("", tcell.ColorGreen)
		statusbars = append(statusbars, statusbar)
		table.SetCell(idx, 3, statusbar)

		idx++
	}

	go data.fetchAndProcessRepos(progressbars, statusbars)

	return table
}

func (data *TuiApplication) fetchAndProcessRepos(progressbars, statusbars []*tview.TableCell) {
	progressChannel := make(chan ProgressFlag, len(data.repositories))

	os.Chdir("./downloaded-repos")

	idx := 0
	downloadStepWaiter := &sync.WaitGroup{}
	downloadStepWaiter.Add(len(data.repositories))
	for repoName := range data.repositories {
		go data.processRepo(repoName, idx, progressChannel, downloadStepWaiter)
		idx++
	}

	for totalProcessesLeft := len(data.repositories); totalProcessesLeft > 0; {
		progressFlag := <-progressChannel
		if progressFlag.done {
			totalProcessesLeft--
		}
		data.application.QueueUpdateDraw(func() {
			progressbars[progressFlag.repoIdx].SetText(
				fmt.Sprintf(
					"%s%s",
					strings.Repeat("#", progressFlag.progress),
					strings.Repeat("-", 3-progressFlag.progress),
				),
			)
			statusbars[progressFlag.repoIdx].SetText(progressFlag.status)
		})
	}

	data.addPackagePages()
}

func (data *TuiApplication) addPackagePages() {
	numPages := data.pages.GetPageCount()

	for pkgIdx, pkgReq := range data.packagesRequired {
		table := tview.NewTable()
		table.SetBorder(true)
		table.SetBorders(true)
		table.SetBorderPadding(0, 0, 20, 20)

		table.SetTitle(fmt.Sprintf("%s @ %s", pkgReq.Name, pkgReq.GetVersion()))
		table.SetCell(0, 0, utils.CrateCell("Repo Name", tcell.ColorYellow))
		table.SetCell(0, 1, utils.CrateCell("Repo Url", tcell.ColorYellow))
		table.SetCell(0, 2, utils.CrateCell("Version Found", tcell.ColorYellow))
		table.SetCell(0, 3, utils.CrateCell("Does Satisfy ?", tcell.ColorYellow))

		repoIdx := 1
		for _, repo := range data.repositories {
			table.SetCell(repoIdx, 0, utils.CrateNormalCell(repo.Name))
			table.SetCell(repoIdx, 1, utils.CrateCell(repo.Url, tcell.ColorCornflowerBlue))
			table.SetCell(repoIdx, 2, utils.CrateCell(repo.Packages[pkgIdx].GetVersion(), tcell.ColorGreen))
			doesSatisfyVersion := "no"
			versionComparision := repo.Packages[pkgIdx].VersionRankComparedTo(&pkgReq)

			if -1 <= versionComparision && versionComparision <= 1 {
				doesSatisfyVersion = "yes"
			}
			table.SetCell(repoIdx, 3, utils.CrateNormalCell(doesSatisfyVersion))

			repoIdx++
		}

		data.pages.AddPage(
			fmt.Sprint(numPages),
			table,
			true,
			false,
		)
		numPages++
	}
}

func (data *TuiApplication) processRepo(
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

		data.lock.Lock()
		repository.Packages[idx].Name = pkg.Name
		repository.Packages[idx].MajorVersion, _ = strconv.ParseUint(readMatch("major"), 10, 64)
		repository.Packages[idx].FeatureVersion, _ = strconv.ParseUint(readMatch("feature"), 10, 64)
		repository.Packages[idx].MinorVersion, _ = strconv.ParseUint(readMatch("minor"), 10, 64)
		data.lock.Unlock()
	}

	progressChannel <- ProgressFlag{repoIdx: repoIdx, progress: 3, status: "Completed", done: true}
}

func fetchRepo(repository *utils.RepoWithPackagesFound) {
	cmd := exec.Command("git", "clone", repository.Url, repository.Name)

	if err := cmd.Run(); err != nil {
		return
	}
}

func clearDownloadedRepos() {
	os.Chdir("./downloaded-repos")

	childDirs, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, dir := range childDirs {
		os.RemoveAll(dir.Name())
	}
}
