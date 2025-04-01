package comparator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

const (
	OJ_AC  = 0
	OJ_WA  = 1
	OJ_TLE = 2
	OJ_MLE = 3
	OJ_RE  = 4
	OJ_CE  = 5
	OJ_PE  = 6
)

func isSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\r' || r == '\t'
}

func findNextNonSpace(reader1 *bufio.Reader, reader2 *bufio.Reader, c1 *rune, c2 *rune) (int, error) {
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
				return OJ_PE, nil
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
	return OJ_AC, nil
}

func CompareZOJ(fileExpected, fileOutput, inputFile string) (int, error) {
	f1, err1 := os.Open(fileExpected)
	f2, err2 := os.Open(fileOutput)
	if err1 != nil || err2 != nil {
		return OJ_RE, fmt.Errorf("error opening files: %v, %v", err1, err2)
	}
	defer f1.Close()
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
				MakeDiffOut(fileExpected, fileOutput, inputFile, c1, c2)
				return OJ_WA, nil
			}
			if c1 != c2 {
				MakeDiffOut(fileExpected, fileOutput, inputFile, c1, c2)
				return OJ_WA, nil
			}
			c1, _, _ = r1.ReadRune()
			c2, _, _ = r2.ReadRune()
		}
		ret, _ = findNextNonSpace(r1, r2, &c1, &c2)
		if c1 == -1 && c2 == -1 {
			return ret, nil
		}
		if c1 == -1 || c2 == -1 {
			MakeDiffOut(fileExpected, fileOutput, inputFile, c1, c2)
			return OJ_WA, nil
		}
		if (c1 == '\n' || c1 == -1) && (c2 == '\n' || c2 == -1) {
			continue
		}
	}
}

func MakeDiffOut(file1, file2, inputFile string, c1, c2 rune) {
	f1, err := os.Open(file1)
	if err != nil {
		fmt.Println("Error abriendo output esperado:", err)
		return
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		fmt.Println("Error abriendo output del usuario:", err)
		return
	}
	defer f2.Close()

	out, err := os.OpenFile("diff.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error abriendo diff.out:", err)
		return
	}
	defer out.Close()

	fmt.Fprintln(out, "Entrada")
	fmt.Fprintln(out, "=================")

	if fin, err := os.Open(inputFile); err == nil {
		copyN(out, fin, 512)
		fin.Close()
	} else {
		fmt.Fprintln(out, "[No se pudo abrir input de entrada]")
	}

	fmt.Fprintln(out, "\n=================")
	fmt.Fprintln(out, "Respuesta Correcta:")
	copyN(out, f1, 900)
	fmt.Fprintln(out, "\n-----------------")
	fmt.Fprintln(out, "Tu respuesta:")
	copyN(out, f2, 900)
	fmt.Fprintln(out, "\n=================")
	fmt.Fprintln(out, "Este modulo esta en modo beta. No se confié.")
	fmt.Fprintf(out, "Dato esperado '%c', Tu salida '%c'.\n", c1, c2)
}

func copyN(dst io.Writer, src io.Reader, n int64) error {
	_, err := io.CopyN(dst, src, n)
	return err
}
