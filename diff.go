package main

import (
	"github.com/BrandonWade/godash"
	"github.com/BrandonWade/synth"
)

// simpleDiff - returns a slice of Files calculated by using a naive one-way diff
func simpleDiff(files *[]*synth.File, filters *[]string) []*synth.File {
	newFiles := []*synth.File{}

	for _, file := range *files {
		if !godash.IncludesSubstr(filters, file.Path) {
			newFiles = append(newFiles, file)
		}
	}

	return newFiles
}
