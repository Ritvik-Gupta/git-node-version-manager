package tui

import (
	"fmt"
	"os"
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
	reposWriteLock   *sync.Mutex
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
		reposWriteLock:   &sync.Mutex{},
	}
}

func (data *TuiApplication) Start() {
	crateTempRepoStore()
	defer clearDownloadedRepos()

	currentPage := 0

	data.pages.AddPage(
		fmt.Sprint(currentPage),
		data.createRepoAnalysisTable(),
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

func crateTempRepoStore() {
	os.Mkdir("./downloaded-repos", os.ModeTemporary)
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
