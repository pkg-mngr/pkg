package util

import (
	"fmt"
	"strings"
)

func WrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}
	output := []string{""}
	for word := range strings.SplitSeq(text, " ") {
		if len(output[len(output)-1]+word) > width {
			output = append(output, "")
		}
		output[len(output)-1] += fmt.Sprintf("%s ", word)
	}
	return strings.Join(output, "\n")
}
