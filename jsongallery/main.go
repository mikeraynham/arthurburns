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

type Root struct {
	Thumb string
	Small string
	Large string
}

func main() {
	var baseDir string;
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	}

	root := Root{
		Thumb: `images/products/thumbs`,
		Small: `images/products/small`,
		Large: `images/products/large`,
	}

	pathRe := regexp.MustCompile(`^(.+)_tn_([0-9]+)\.jpg$`)
	images := make(map[string][]Image)

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

			name := res[1]
			index := res[2]
			title := strings.Title(strings.Replace(name, "-", " ", -1))

			images[title] = append(images[title], Image{
				Thumb:      filepath.Join(root.Thumb, name+"_tn_"+index+".jpg"),
				Small:      filepath.Join(root.Small, name+"_sm_"+index+".jpg"),
				Large:      filepath.Join(root.Large, name+"_lg_"+index+".jpg"),
				Title:      title + " " + index,
				CoverImage: isCoverImage(thumbPath, name, index),
			})

			return nil
		},
	)

	json, _ := json.MarshalIndent(images, "", "  ")
	fmt.Println(string(json))
}

func isCoverImage(thumbPath string, name string, index string) bool {
	coverPath := filepath.Join(filepath.Dir(thumbPath), name+"_tn_"+index)
	_, err := os.Stat(coverPath)
	return !os.IsNotExist(err)
}
