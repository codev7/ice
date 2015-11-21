package routes

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func Process(args []string) {
	f, err := os.Open(args[0])
	if err != nil {
		log.Fatalf("Failed to open the file %s. %s", args[0], err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var buf bytes.Buffer
	comments := []string{}
	var requests int

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "//") {
			comments = append(comments, scanner.Text()[2:])
		} else if generate(scanner.Text(), comments, &buf) {
			requests = requests + 1
			comments = []string{}
		} else {
			comments = []string{}
		}
	}

	if requests > 0 {
		fs, err := os.Create(args[0] + ".generated.go")
		if err != nil {
			log.Fatalf("Error writing output", err)
		}
		defer fs.Close()
		fs.Write([]byte("package " + args[1] + "\r\n\r\n"))
		io.Copy(fs, &buf)
	}
}

func generate(line string, comments []string, buf io.Writer) bool {
	if len(comments) == 0 {
		return false
	}
	parts := strings.Split(strings.Trim(line, " "), " ")
	if len(parts) < 2 || parts[0] != "type" {
		return false
	}
	docs := ""
	route := ""
	for _, c := range comments {
		if strings.HasPrefix(c, "route") {
			route = strings.Trim(c[5:], " ")
		} else {
			docs = docs + "\r\n" + c
		}
	}

	docs = "`" + docs + "`"

	buf.Write([]byte(fmt.Sprintf(`
func (r *%s) RequestDescription()string{
return %s
}
`, parts[1], docs)))

	if route != "" {
		buf.Write([]byte(fmt.Sprintf(`
func (r *%s) Route()string{
return "%s"
}
`, parts[1], route)))
	}

	return true
}
