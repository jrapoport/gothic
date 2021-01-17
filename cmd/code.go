package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/jrapoport/gothic/api"
	"github.com/jrapoport/gothic/models"
	"github.com/jrapoport/gothic/storage"
	"github.com/spf13/cobra"
	"text/tabwriter"
)

var codeCmd = &cobra.Command{
	Use:  "code [count]",
	Long: "generate new signup codes ",
	RunE: codeRunE,
	Args: cobra.MaximumNArgs(1),
}

var (
	codeMulti  bool
	codeOutput string
)

func init() {
	fs := codeCmd.Flags()
	fs.BoolVarP(&codeMulti, "multiuse", "m", false, "generate multiuse codes")
	fs.StringVarP(&codeOutput, "out", "o", "", "output csv to file path")
}

func codeRunE(_ *cobra.Command, args []string) error {
	count := 1
	if len(args) > 0 {
		var err error
		count, err = strconv.Atoi(args[0])
		if err != err {
			return err
		}
	}
	t := models.SingleUse
	if codeMulti {
		t = models.MultiUse
	}
	c := cmdConfig
	db, err := storage.Dial(c, c.Log)
	if err != nil {
		return err
	}
	a := api.NewAPI(c, db)
	codes, err := a.NewSignupCodes(models.PINFormat, t, count)
	if err != nil {
		err = fmt.Errorf("error generating codes: %w", err)
		return err
	}
	cnt := len(codes)
	list := make([]string, cnt)
	for i := 0; i < cnt; i++ {
		list[i] = codes[i].Code
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
	nCols := 10
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
