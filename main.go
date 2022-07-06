package main

import (
	"errors"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"go-nuget-list/internal/app"
	"go-nuget-list/pkg/nuget"
	"os"
	"strings"
)

func main() {
	cli := &cli.App{
		Name:  "go-nuget-list",
		Usage: "go-nuget-list -o",
		Action: func(c *cli.Context) error {
			fileName := c.Args().Get(0)

			if !strings.HasSuffix(fileName, ".sln") && !strings.HasSuffix(fileName, ".csproj") {
				return errors.New("unknown input file format")
			}

			pterm.Info.Println("Input file:", fileName)

			packageSources, err := nuget.NewNugetConfigFinder().Search(fileName)
			if err != nil {
				return err
			}

			result, err := app.NewPackagesScanner(packageSources).Scan(fileName)
			if err != nil {
				return err
			}

			outputFile := c.String("output")

			if outputFile != "" {
				if err = result.SaveToFile(outputFile); err != nil {
					return err
				}
				pterm.Info.Println("Results successfully saved saved to:", outputFile)
			} else {
				pterm.Warning.Println("No output file has been provided")
				result.Print()

			}

			pterm.Info.Println("DONE!")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Usage:   "output file",
				Aliases: []string{"o"},
			},
		},
	}

	if err := cli.Run(os.Args); err != nil {
		pterm.Error.Println(err)
	}
}
