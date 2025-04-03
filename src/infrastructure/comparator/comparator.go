package comparator

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
)

type Comparator struct{}

func (c *Comparator) Compare(expectedPath string, outputPath string) (abstractions.ComparisonResult, error) {
	return c.CompareZoj(expectedPath, outputPath)
}

func (c *Comparator) CompareZoj(expectedPath, outputPath string) (abstractions.ComparisonResult, error) {
	expectedBytes, err1 := os.ReadFile(expectedPath)
	if err1 != nil {
		return abstractions.OJ_RE, fmt.Errorf("error reading expected output: %w", err1)
	}

	outputBytes, err2 := os.ReadFile(outputPath)
	if err2 != nil {
		return abstractions.OJ_RE, fmt.Errorf("error reading user output: %w", err2)
	}

	// Compare directly first
	if bytes.Equal(bytes.TrimRight(expectedBytes, "\r\n \t"), bytes.TrimRight(outputBytes, "\r\n \t")) {
		return abstractions.OJ_AC, nil
	}

	// Normalize whitespace for PE check
	expectedNorm := normalizeWhitespace(expectedBytes)
	outputNorm := normalizeWhitespace(outputBytes)

	if expectedNorm == outputNorm {
		return abstractions.OJ_PE, nil
	}

	return abstractions.OJ_WA, nil
}

func normalizeWhitespace(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanWords)

	var builder strings.Builder
	first := true
	for scanner.Scan() {
		if !first {
			builder.WriteRune(' ')
		}
		builder.WriteString(scanner.Text())
		first = false
	}
	return builder.String()
}
