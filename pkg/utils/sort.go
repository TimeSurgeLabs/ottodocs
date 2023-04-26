package utils

import (
	"sort"

	"github.com/chand1012/memory"
)

func SortByScore(fragments []memory.MemoryFragment) []memory.MemoryFragment {
	sort.Slice(fragments, func(i, j int) bool {
		return fragments[i].Score > fragments[j].Score
	})
	return fragments
}

func SortByAverage(fragments []memory.MemoryFragment) []memory.MemoryFragment {
	sort.Slice(fragments, func(i, j int) bool {
		return fragments[i].Avg > fragments[j].Avg
	})
	return fragments
}
