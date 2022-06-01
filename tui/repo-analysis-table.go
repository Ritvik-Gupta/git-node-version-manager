package tui

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ProgressFlag struct {
	repoIdx  int
	progress int
	status   string
	done     bool
}

func (data *TuiApplication) createRepoAnalysisTable() *tview.Table {
	table := tview.NewTable()
	table.SetBorder(true)
	table.SetBorders(true)
	table.SetBorderPadding(0, 0, 20, 20)

	table.SetTitle("Repository Analysis")
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
		go data.analyseRepo(repoName, idx, progressChannel, downloadStepWaiter)
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
