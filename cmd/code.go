package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/storage"
	"github.com/spf13/cobra"
	"text/tabwriter"
)

var codeCmd = &cobra.Command{
	Use:  "code [count]",
	Long: "generate new access codes",
	RunE: codeRunE,
	Args: cobra.MaximumNArgs(1),
}

var (
	codeMulti  bool
	codeOutput string
	codeFormat string
)

func init() {
	fs := codeCmd.Flags()
	fs.BoolVarP(&codeMulti, "multiuse", "m", false, "generate multiuse codes")
	fs.StringVarP(&codeOutput, "out", "o", "", "output csv to file path")
}

func codeRunE(cmd *cobra.Command, args []string) error {
	count := 1
	if len(args) > 0 {
		var err error
		count, err = strconv.Atoi(args[0])
		if err != err {
			return err
		}
	}
	t := api.CodeTypeSingleUse
	if codeMulti {
		t = api.CodeTypeMultiUse
	}
	c := cmdConfig
	db, err := storage.Dial(c, c.Log)
	if err != nil {
		return err
	}
	a := api.NewAPI(c, db)
	codes, err := a.NewRandomAccessCodes(api.CodeFormatPIN, t, count)
	if err != nil {
		err = fmt.Errorf("error generating codes: %w", err)
		return err
	}
	cnt := len(codes)
	list := make([]string, cnt)
	err = db.Transaction(func(tx *storage.Connection) error {
		for i := 0; i < cnt; i++ {
			code := codes[i]
			if err = tx.Create(code).Error; err != nil {
				return err
			}
			list[i] = code.Code
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("error generating codes: %w", err)
		return err
	}
	if codeOutput != "" {
		file, _ := os.Create(codeOutput)
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		csvHeaders := []string{"code"}
		err = writer.Write(csvHeaders)
		if err != nil {
			return err
		}
		for _, ac := range codes {
			err = writer.Write([]string{ac.Code})
			if err != nil {
				return err
			}
			writer.Flush()
		}
	}
	c.Log.Infof("generated %d codes", cnt)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	nCols := 8
	for i := 0; i < cnt; i += nCols {
		cols := make([]interface{}, nCols)
		var n int
		line := ""
		for ; n < nCols; n++ {
			line += "%s\t"
			cols[n] = ""
			if n+i < cnt {
				cols[n] = list[i+n]
			}
		}
		line += "\n"
		_, _ = fmt.Fprintf(w, line, cols...)
	}
	defer w.Flush()
	return nil
}
