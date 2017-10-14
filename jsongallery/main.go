package main

import (
	"encoding/json"
	"fmt"
	. "github.com/mikeraynham/arthurburns/jsongallery/image"
	"github.com/mikeraynham/arthurburns/jsongallery/pathfmt"
	"log"
	"os"
	"path/filepath"
)

type Section struct {
	SectionTitle string
	SectionDir   string
	Images       []Image
}

type Sections struct {
	Sections []Section
}

func (s *Sections) AppendSection(section Section) {
	if len(section.Images) > 0 {
		s.Sections = append(s.Sections, section)
	}
}

func main() {
	var baseDir string = "."
	if len(os.Args) > 1 {
		baseDir = filepath.Clean(os.Args[1])
	}

	var section Section
	var sections Sections
	var sectionDir string
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

			sectionDir = filepath.Base(filepath.Dir(thumbPath))
			sectionTitle := pathfmt.ToTitle(sectionDir)
			image := NewImage(thumbPath, sectionDir)

			if image == nil {
				return nil
			}

			if sectionDir != prevSectionDir {
				sections.AppendSection(section)
				section = Section{
					SectionDir:   sectionDir,
					SectionTitle: sectionTitle,
				}
			}

			section.Images = append(section.Images, *image)
			prevSectionDir = sectionDir

			return nil
		},
	)

	sections.AppendSection(section)

	json, _ := json.MarshalIndent(sections.Sections, "", "  ")
	fmt.Println(string(json))
}
