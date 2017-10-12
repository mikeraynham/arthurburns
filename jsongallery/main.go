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

type Image struct {
	Thumb      string
	Small      string
	Large      string
	Title      string
	CoverImage bool
}

type Section struct {
	SectionDir string
	Images     []Image
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
	sections := make(map[string]*Section)

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
			sectionDir := filepath.Base(filepath.Dir(thumbPath))
			sectionTitle := toTitle(sectionDir)

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

			if _, exists := sections[sectionTitle]; exists {
				sections[sectionTitle].Images = append(
					sections[sectionTitle].Images,
					image,
				)
			} else {
				sections[sectionTitle] = &Section{
					SectionDir: sectionDir,
					Images: []Image{
						image,
					},
				}
			}

			return nil
		},
	)

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
