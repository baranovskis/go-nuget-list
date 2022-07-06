package app

import (
	"fmt"
	"github.com/pterm/pterm"
	"go-nuget-list/pkg/csproj"
	"go-nuget-list/pkg/nuget"
	"go-nuget-list/pkg/sln"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type PackagesScanner struct {
	sources []nuget.PackageSource
}

func NewPackagesScanner(sources []nuget.PackageSource) *PackagesScanner {
	return &PackagesScanner{sources: sources}
}

func (ps *PackagesScanner) Scan(fileName string) (*Output, error) {
	pterm.Info.Println("Starting packages scanner...")

	output := &Output{}

	if strings.HasSuffix(fileName, ".sln") {
		pterm.Info.Println("Solution file detected...")
		err := ps.scanSolution(output, fileName)
		if err != nil {
			return nil, err
		}
	} else {
		pterm.Info.Println("Project file detected...")
		err := ps.scanProject(output, fileName)
		if err != nil {
			return nil, err
		}
	}

	pterm.Info.Println("Total packages:", output.TotalPackages)

	sort.Slice(output.Packages, func(i, j int) bool {
		return output.Packages[i].Id+output.Packages[i].Version < output.Packages[j].Id+output.Packages[j].Version
	})

	progress, _ := pterm.DefaultProgressbar.WithTotal(len(output.Packages)).WithTitle("Fetching NuGet data..").Start()
	nc := nuget.NewNugetClient()

	for _, p := range output.Packages {
		var found = false

		for _, source := range ps.sources {
			progress.UpdateTitle(fmt.Sprintf("Fetching NuGet package data (%s) from '%s'...", p.Id,
				source.SourceName))

			response, err := nc.Search(source.Path, p.Id)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("Failed to fetch NuGet package (%s) from '%s'\nError: %s", p.Id,
					source.SourceName, err))
			}

			for _, d := range response.Data {
				// check package version?
				if d.Id != p.Id /*|| d.Version != p.Version*/ {
					continue
				}

				p.Name = d.Title
				p.Description = d.Description
				p.Summary = d.Summary
				p.Authors = d.Authors.Values
				p.Tags = d.Tags.Values
				p.LicenseUrl = d.LicenseUrl
				p.ProjectUrl = d.ProjectUrl

				found = true

				pterm.Success.Println(fmt.Sprintf("NuGet package (%s) has been successfully fetched", p.Id))
				break
			}

			if found {
				break
			}
		}

		progress.Increment()

		if !found {
			pterm.Warning.Println(fmt.Sprintf("Failed to fetch NuGet package (%s)", p.Id))
		}
	}

	progress.RemoveWhenDone = true
	progress.Stop()

	pterm.Info.Println("Packages scanning has been completed")
	return output, nil
}

func (ps *PackagesScanner) scanSolution(data *Output, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	fileDir := filepath.Dir(fileName)

	spinner, _ := pterm.DefaultSpinner.Start("Parsing solution file...")

	sp, err := sln.NewSolutionParser(f).Parse()
	if err != nil {
		return err
	}
	spinner.Success()

	spinner, _ = pterm.DefaultSpinner.Start("Parsing project files...")

	for _, p := range sp.Projects {
		if !strings.HasSuffix(p.ProjectFile, ".csproj") {
			continue
		}

		if err := ps.scanProject(data, path.Join(fileDir, p.ProjectFile)); err != nil {
			pterm.Error.Println(fmt.Sprintf("Failed to parse project file (%s)\nError: %s", p.ProjectFile, err))
		}
	}

	spinner.Success()
	pterm.Info.Println("Scanned projects in solution:", data.ScannedProjects)
	return nil
}

func (ps *PackagesScanner) scanProject(data *Output, fileName string) error {
	pr := csproj.NewProjectParser()
	prj, err := pr.Parse(fileName)
	if err != nil {
		return err
	}

	data.ScannedProjects++

	for _, ig := range prj.ItemGroups {
		for _, pr := range ig.PackageReferences {
			found := false

			// check for duplicates
			for _, c := range data.Packages {
				if c.Id == pr.Include && c.Version == pr.Version {
					found = true
					break
				}
			}

			if !found {
				data.TotalPackages++
				data.Packages = append(data.Packages, &OutputPackage{
					Id:      pr.Include,
					Version: pr.Version,
				})
			}
		}
	}

	return nil
}
