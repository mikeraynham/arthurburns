package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Section struct {
	SectionTitle string
	SectionDir   string
	Images       []Image
}

type Image struct {
	Thumb      string
	Small      string
	Large      string
	Title      string
	CoverImage bool
}

type Root struct {
	Thumb string
	Small string
	Large string
}

func main() {
	var baseDir string
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	}

	root := Root{
		Thumb: `images/products/thumbs`,
		Small: `images/products/small`,
		Large: `images/products/large`,
	}

	pathRe := regexp.MustCompile(`^(.+)_tn_([0-9]+)\.jpg$`)
	var section Section
	var sections []Section
	var sectionDir string
	var prevSectionDir string

	filepath.Walk(
		filepath.Join(baseDir, root.Thumb),
		func(thumbPath string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}

			if info.IsDir() {
				return nil
			}

			res := pathRe.FindStringSubmatch(info.Name())

			if res == nil {
				return nil
			}

			imgName := res[1]
			imgIndex := res[2]
			imgTitle := toTitle(imgName)
			sectionDir = filepath.Base(filepath.Dir(thumbPath))
			sectionTitle := toTitle(sectionDir)

			if sectionDir != prevSectionDir {
				if len(section.Images) > 0 {
					sections = append(sections, section)
				}
				section = Section{
					SectionDir: sectionDir,
					SectionTitle: sectionTitle,
				}
			}

			buildPath := func(rootDir string, size string) string {
				return filepath.Join(
					"/",
					rootDir,
					sectionDir,
					imgName+"_"+size+"_"+imgIndex+".jpg",
				)
			}

			image := Image{
				Thumb:      buildPath(root.Thumb, "tn"),
				Small:      buildPath(root.Small, "sm"),
				Large:      buildPath(root.Large, "lg"),
				Title:      imgTitle,
				CoverImage: isCoverImage(thumbPath, imgName, imgIndex),
			}

			section.Images = append(section.Images, image)
			prevSectionDir = sectionDir

			return nil
		},
	)

	if len(section.Images) > 0 {
		sections = append(sections, section)
	}

	json, _ := json.MarshalIndent(sections, "", "  ")
	fmt.Println(string(json))
}

func toTitle(fileName string) string {
	return strings.Title(strings.Replace(fileName, "-", " ", -1))
}

func isCoverImage(thumbPath string, imgName string, imgIndex string) bool {
	coverPath := filepath.Join(filepath.Dir(thumbPath), imgName+"_tn_"+imgIndex)
	_, err := os.Stat(coverPath)
	return !os.IsNotExist(err)
}
