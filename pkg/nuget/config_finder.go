package nuget

import (
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/pterm/pterm"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var WindowsLocations = []string{
	"%ProgramFiles(x86)%\\NuGet\\Config",
	"%AppData%\\NuGet",
}

var LinuxLocations = []string{
	"~/.nuget/NuGet",
	"~/.config/NuGet",
}

type ConfigFinder struct{}

// NewNugetConfigFinder returns a new instance of Parser.
func NewNugetConfigFinder() *ConfigFinder {
	return &ConfigFinder{}
}

func (cf *ConfigFinder) Search(paths ...string) ([]PackageSource, error) {
	pterm.Info.Println("Searching NuGet package sources...")

	var locations []string

	if runtime.GOOS == "windows" {
		re := regexp.MustCompile(`\%([^\%\%]*)\%`)

		for _, location := range WindowsLocations {
			subMatchAll := re.FindAllString(location, -1)

			for _, value := range subMatchAll {
				environment := os.Getenv(strings.Trim(value, "%"))
				locations = append(locations, strings.Replace(location, value, environment, -1))
			}
		}
	} else {
		locations = LinuxLocations
	}

	// add custom nuget paths
	for _, path := range paths {
		dir, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if !dir.IsDir() {
			path = filepath.Dir(path)
		}

		locations = append(locations, path)
	}

	var packageSources []PackageSource

	// search all configuration files
	for _, location := range locations {
		err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(strings.ToLower(filepath.Base(path)), ".config") {
				reader, err := os.Open(path)
				if err != nil {
					return err
				}
				defer reader.Close()

				doc, err := xmlquery.Parse(reader)
				if err != nil {
					return err
				}

				result, err := xmlquery.QueryAll(doc, "//configuration/packageSources/add")

				for _, source := range result {
					val := source.SelectAttr("value")

					// TODO: rewrite me :(
					// skip system paths for this moment
					_, err := os.Stat(val)
					if err == nil {
						continue
					}

					// only valid HTTP path is allowed
					if _, err := url.ParseRequestURI(val); err != nil {
						continue
					}

					packageSources = append(packageSources, PackageSource{
						FileName:        path,
						SourceName:      source.SelectAttr("key"),
						Path:            val,
						ProtocolVersion: source.SelectAttr("protocolVersion"),
					})
				}
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}

	if len(packageSources) < 1 {
		return nil, errors.New("no package sources found")
	}

	var displayPaths string
	for _, ps := range packageSources {
		displayPaths += fmt.Sprintf("> %s - %s\n", ps.SourceName, ps.FileName)
	}

	pterm.Info.Printf("Found %d package sources:\n%s", len(packageSources), displayPaths)

	return packageSources, nil
}
