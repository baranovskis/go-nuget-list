package csproj

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type ProjectParser struct{}

func NewProjectParser() *ProjectParser {
	return &ProjectParser{}
}

func (pp *ProjectParser) Parse(path string) (Project, error) {
	var prj Project

	xmlFile, err := os.Open(path)
	if err != nil {
		return prj, err
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return prj, err
	}

	err = xml.Unmarshal(byteValue, &prj)
	return prj, err
}
