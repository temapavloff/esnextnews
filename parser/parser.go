package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/DusanKasan/parsemail"
)

// Parse - return markdown formatted representation of email
func Parse(r io.Reader) (string, error) {
	e, err := parsemail.Parse(r)
	if err != nil {
		return "", err
	}
	items, err := collect(e.TextBody)
	if err != nil {
		return "", err
	}
	return render(items), nil
}

func collect(body string) ([][]string, error) {
	start := "THE LATEST IN JAVASCRIPT AND CROSS-PLATFORM TOOLS"
	end := "----------------------------------------------------------------------------------------------------------------------------------"

	scanner := bufio.NewScanner(strings.NewReader(body))
	buckets := [][]string{}
	i := make([]string, 0)

	for scanner.Scan() {
		if strings.Trim(scanner.Text(), " ") == start {
			break
		}
	}

	for scanner.Scan() {
		l := strings.Trim(scanner.Text(), " ")

		if l == end {
			break
		}

		if l == "" {
			if len(i) > 0 {
				buckets = append(buckets, i)
				i = make([]string, 0)
			}
			continue
		}

		i = append(i, l)
	}

	if len(buckets) == 0 {
		return buckets, errors.New("No content found")
	}

	return buckets, nil
}

func render(buckets [][]string) string {
	result := ""

	for _, item := range buckets {

		result += renderItem(item)
	}

	return strings.Trim(result, "\n")
}

func renderItem(item []string) string {
	l := len(item)

	if l == 1 {
		return "\n**" + item[0] + "**\n\n"
	}

	if l > 1 && item[1][0] == '[' && item[1][len(item[1])-1] == ']' {
		il := len(item[1])
		if item[1][0] == '[' && item[1][il-1] == ']' {
			return "[" + item[0] + "](" + item[1][1:il-1] + ")\n\n"
		}
	}

	return strings.Join(item, " ") + "\n\n"
}
