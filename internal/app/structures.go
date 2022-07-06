package app

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
)

type Output struct {
	ScannedProjects int32            `json:"scannedProjects"`
	TotalPackages   int32            `json:"usedPackages"`
	Packages        []*OutputPackage `json:"packages"`
}

type OutputPackage struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Summary     string   `json:"summary,omitempty"`
	Version     string   `json:"version"`
	Authors     []string `json:"authors,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	LicenseUrl  string   `json:"licenseUrl"`
	ProjectUrl  string   `json:"projectUrl"`
}

func (o *Output) Print() {
	fmt.Println()

	td := pterm.TableData{
		{"Id", "Version", "License", "Project"},
	}

	for _, p := range o.Packages {
		tmp := make([]string, 4)
		tmp[0] = p.Id
		tmp[1] = p.Version
		tmp[2] = p.LicenseUrl
		tmp[3] = p.ProjectUrl

		td = append(td, tmp)
	}

	pterm.DefaultTable.WithHasHeader().WithData(td).Render()
	fmt.Println()
}

func (o *Output) SaveToFile(fileName string) error {
	outputFile, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, outputFile, 0644)
}
