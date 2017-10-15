package main

import (
	"encoding/json"
	"fmt"
	"github.com/mikeraynham/arthurburns/jsongallery/image"
	"github.com/mikeraynham/arthurburns/jsongallery/pathfmt"
	"log"
	"os"
	"path/filepath"
)

type Section struct {
	SectionTitle string
	SectionDir   string
	Images       []image.Image
}

func AppendSection(sections []Section, section Section) []Section {
	if len(section.Images) > 0 {
		sections = append(sections, section)
	}

	return sections
}

func main() {
	var baseDir string = "."
	if len(os.Args) > 1 {
		baseDir = filepath.Clean(os.Args[1])
	}

	var section Section
	var sections []Section
	var prevSectionDir string

	filepath.Walk(
		filepath.Join(baseDir, "thumbs"),
		func(thumbPath string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}

			if info.IsDir() {
				return nil
			}

			sectionDir := filepath.Base(filepath.Dir(thumbPath))
			sectionTitle := pathfmt.ToTitle(sectionDir)

			if sectionDir != prevSectionDir {
				sections = AppendSection(sections, section)
				section = Section{
					SectionDir:   sectionDir,
					SectionTitle: sectionTitle,
				}
			}

			img := image.New(thumbPath, sectionDir)

			if img == nil {
				return nil
			}

			section.Images = append(section.Images, *img)
			prevSectionDir = sectionDir

			return nil
		},
	)

	sections = AppendSection(sections, section)

	json, _ := json.MarshalIndent(sections, "", "  ")
	fmt.Println(string(json))
}
