package image

import (
	"path/filepath"
	"os"
	"regexp"
	"github.com/mikeraynham/arthurburns/jsongallery/pathfmt"
)

var pathRe = regexp.MustCompile(`^(.+)_tn_([0-9]+)\.jpg$`)

type Image struct {
	Thumb      string
	Small      string
	Large      string
	Title      string
	CoverImage bool
}

func New(thumbPath string, sectionDir string) *Image {

	thumbBase := filepath.Base(thumbPath)
	res := pathRe.FindStringSubmatch(thumbBase)

	if res == nil {
		return nil
	}

	imgName := res[1]
	imgIndex := res[2]
	imgTitle := pathfmt.ToTitle(imgName)

	buildPath := func(sizeDir string , size string) string {
		return filepath.Join(
			sizeDir,
			sectionDir,
			imgName+"_"+size+"_"+imgIndex+".jpg",
		)
	}

	return &Image{
		Thumb:      buildPath("thumb", "tn"),
		Small:      buildPath("small", "sm"),
		Large:      buildPath("large", "lg"),
		Title:      imgTitle,
		CoverImage: isCoverImage(thumbPath, imgName, imgIndex),

	}
}

func isCoverImage(thumbPath string, imgName string, imgIndex string) bool {
	coverPath := filepath.Join(filepath.Dir(thumbPath), imgName+"_tn_"+imgIndex)
	_, err := os.Stat(coverPath)
	return !os.IsNotExist(err)
}
