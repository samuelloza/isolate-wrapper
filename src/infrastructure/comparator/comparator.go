package comparator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/samuelloza/isolate-wrapper/src/application/abstractions"
)

type Comparator struct{}

func (c *Comparator) Compare(expectedPath string, outputPath string) (abstractions.ComparisonResult, error) {
	return c.CompareZoj(expectedPath, outputPath, "/box/input.txt")
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}

func findNextNonSpace(reader1 *bufio.Reader, reader2 *bufio.Reader, c1 *rune, c2 *rune) (abstractions.ComparisonResult, error) {
	var err error
	for unicode.IsSpace(*c1) || unicode.IsSpace(*c2) {
		if *c1 != *c2 {
			if *c2 == -1 {
				for unicode.IsSpace(*c1) {
					r, _, e := reader1.ReadRune()
					if e == io.EOF {
						*c1 = -1
						break
					}
					*c1 = r
				}
				continue
			} else if *c1 == -1 {
				for unicode.IsSpace(*c2) {
					r, _, e := reader2.ReadRune()
					if e == io.EOF {
						*c2 = -1
						break
					}
					*c2 = r
				}
				continue
			} else if *c1 == '\r' && *c2 == '\n' {
				*c1, _, err = reader1.ReadRune()
			} else if *c2 == '\r' && *c1 == '\n' {
				*c2, _, err = reader2.ReadRune()
			} else {
				return abstractions.OJ_PE, nil
			}
		}
		if unicode.IsSpace(*c1) {
			*c1, _, err = reader1.ReadRune()
			if err == io.EOF {
				*c1 = -1
			}
		}
		if unicode.IsSpace(*c2) {
			*c2, _, err = reader2.ReadRune()
			if err == io.EOF {
				*c2 = -1
			}
		}
	}
	return abstractions.OJ_AC, nil
}

func (c *Comparator) CompareZoj(expectedPath, outputPath, inputPath string) (abstractions.ComparisonResult, error) {
	f1, err1 := os.Open(expectedPath)
	if err1 != nil {
		return abstractions.OJ_RE, fmt.Errorf("error opening expected output: %w", err1)
	}
	defer f1.Close()

	f2, err2 := os.Open(outputPath)
	if err2 != nil {
		return abstractions.OJ_RE, fmt.Errorf("error opening user output: %w", err2)
	}
	defer f2.Close()

	r1 := bufio.NewReader(f1)
	r2 := bufio.NewReader(f2)

	c1, _, _ := r1.ReadRune()
	c2, _, _ := r2.ReadRune()
	ret, _ := findNextNonSpace(r1, r2, &c1, &c2)

	for {
		for (!isSpace(c1) && c1 != -1) || (!isSpace(c2) && c2 != -1) {
			if c1 == -1 && c2 == -1 {
				return ret, nil
			}
			if c1 == -1 || c2 == -1 {
				return abstractions.OJ_WA, nil
			}
			if c1 != c2 {
				return abstractions.OJ_WA, nil
			}
			c1, _, _ = r1.ReadRune()
			c2, _, _ = r2.ReadRune()
		}
		ret, _ = findNextNonSpace(r1, r2, &c1, &c2)
		if c1 == -1 && c2 == -1 {
			return ret, nil
		}
		if c1 == -1 || c2 == -1 {
			return abstractions.OJ_WA, nil
		}
		if (c1 == '\n' || c1 == -1) && (c2 == '\n' || c2 == -1) {
			continue
		}
	}
}

func copyN(dst io.Writer, src io.Reader, n int64) error {
	_, err := io.CopyN(dst, src, n)
	return err
}
