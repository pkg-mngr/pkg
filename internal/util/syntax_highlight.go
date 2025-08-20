package util

import (
	"fmt"
	"strconv"
	"strings"
)

func SyntaxHighlight(scriptLine string) string {
	command := "\033[38;2;191;137;255m"
	option := "\033[38;2;153;199;255m"
	text := "\033[38;2;255;159;91m"
	number := "\033[38;2;255;115;46m"
	operator := "\033[38;2;170;170;170m"

	parts := strings.Split(scriptLine, " ")
	highlightedLine := make([]string, len(parts))

	for i, part := range parts {
		if i == 0 || parts[i-1] == "|" {
			highlightedLine[i] = fmt.Sprintf("%s%s\033[0m", command, part)
			continue
		}
		if part[0] == '-' || part[0] == '+' {
			highlightedLine[i] = fmt.Sprintf("%s%s\033[0m", option, part)
			continue
		}
		if part == ">" || part == ">>" {
			highlightedLine[i] = fmt.Sprintf("%s%s\033[0m", operator, part)
			continue
		}
		if _, err := strconv.ParseFloat(part, 64); err == nil {
			highlightedLine[i] = fmt.Sprintf("%s%s\033[0m", number, part)
		}
		highlightedLine[i] = fmt.Sprintf("%s%s\033[0m", text, part)
	}

	return strings.Join(highlightedLine, " ")
}
