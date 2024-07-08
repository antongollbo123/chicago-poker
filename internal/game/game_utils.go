package game

import (
	"strconv"
	"strings"
)

func ParseInput(input string) []int {
	parts := strings.Fields(input)
	indices := []int{}
	for _, part := range parts {
		idx, err := strconv.Atoi(part)
		if err == nil {
			indices = append(indices, idx)
		}
	}
	return indices
}
