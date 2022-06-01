package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Ritvik-Gupta/git-node-version-manager/parser"
	"github.com/Ritvik-Gupta/git-node-version-manager/tui"
	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
	"github.com/urfave/cli/v2"
)

func readRepositories(inputFile string, rawReporitories []string) (map[string]utils.Repository, error) {
	repositories := make(map[string]utils.Repository)

	err := parser.NewCsvParser(inputFile).ParseWriteInto(repositories)
	if err != nil {
		return nil, err
	}

	err = parser.NewRawParser(rawReporitories).ParseWriteInto(repositories)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Value:   "example.csv",
				Usage:   "Input `CSV|TXT File` for Repositories",
			},
			&cli.StringSliceFlag{
				Name:    "repositories",
				Aliases: []string{"r"},
				Usage:   "Provide `REPOS` as an array",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return cli.Exit("No Package Names Provided", 86)
			}

			packages := make([]utils.Package, 0, ctx.NArg())
			for i := 0; i < ctx.NArg(); i++ {
				pkg, err := utils.ParsePackage(ctx.Args().Get(i))
				if err != nil {
					return err
				}

				packages = append(packages, pkg)
			}

			fmt.Println(packages)
			fmt.Println(ctx.String("input"))
			fmt.Println(ctx.StringSlice("repositories"))

			repositories, err := readRepositories(ctx.String("input"), ctx.StringSlice("repositories"))
			if err != nil {
				return err
			}

			b, _ := json.MarshalIndent(repositories, "", " ")
			fmt.Println(string(b))

			tui.NewTuiApplicaition(repositories, packages).Start()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
