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
	packages := make(map[string]string)
	for i := 2; i < len(args); i = i + 1 {
		packages[args[i]] = args[i]
	}

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "//") {
			comments = append(comments, scanner.Text()[2:])
		} else if generate(scanner.Text(), comments, &buf, packages) {
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
		for p, i := range packages {
			if p == i {
				fmt.Fprintf(fs, "import \"%s\"\r\n", p)
			} else {
				fmt.Fprintf(fs, "import %s \"%s\"\r\n", i, p)
			}
		}

		io.Copy(fs, &buf)
	}
}

func generate(line string, comments []string, buf io.Writer, packages map[string]string) bool {
	if len(comments) == 0 {
		return false
	}
	parts := strings.Split(strings.Trim(line, " "), " ")
	if len(parts) < 2 || parts[0] != "type" {
		return false
	}

	docs := ""
	route := ""
	var middleware []string

	for _, c := range comments {
		if strings.HasPrefix(c, "route") {
			route = strings.Trim(c[5:], " ")
		} else if strings.HasPrefix(c, "middleware") {
			middleware = strings.Split(c[10:], ",")
			docs = docs + "\r\n" + c
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
return %s
}
`, parts[1], "`" + route + "`")))
	}

	if len(middleware) != 0 {
		packages["github.com/nirandas/ice"] = "ice"
		for i := 0; i < len(middleware); i = i + 1 {
			middleware[i] = strings.Trim(middleware[i], " ")
			if strings.Index(middleware[i], ".") == -1 {
				middleware[i] = "ice." + middleware[i]
			}
		}
		buf.Write([]byte(fmt.Sprintf(`
func (r *%s) Middlewares()[]ice.Middleware{
 return []ice.Middleware{%s}
}
`, parts[1], strings.Join(middleware, ","))))
	}

	return true
}
