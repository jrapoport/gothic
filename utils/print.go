package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/segmentio/encoding/json"
)

// PrettyString formats the contents of the obj as a pretty string
func PrettyString(data interface{}) string {
	b, _ := json.MarshalIndent(data, "", "\t")
	return string(b)
}

// PrettyPrint the contents of the obj
func PrettyPrint(data interface{}) {
	fmt.Printf("%s \n", PrettyString(data))
}

// PrintGrid prints a slice of strings as a grid
func PrintGrid(w io.Writer, list []string, cols int) {
	if cols == 0 {
		cols = 1
	}
	cnt := len(list)
	t := new(tabwriter.Writer)
	t.Init(w, 8, 8, 0, '\t', 0)
	for i := 0; i < cnt; i += cols {
		row := make([]interface{}, cols)
		var n int
		line := ""
		for ; n < cols; n++ {
			line += "%s\t"
			row[n] = ""
			if n+i < cnt {
				row[n] = list[i+n]
			}
		}
		line += "\n"
		_, _ = fmt.Fprintf(t, line, row...)
	}
	defer t.Flush()
}

// WriteCSV prints a slice of strings as a csv
func WriteCSV(name, col string, list []string) error {
	file, _ := os.Create(name)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	csvHeaders := []string{col}
	_ = writer.Write(csvHeaders)
	for _, s := range list {
		_ = writer.Write([]string{s})
		writer.Flush()
	}
	return nil
}
