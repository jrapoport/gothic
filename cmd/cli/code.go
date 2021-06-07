package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/cmd/cli/root"
	"github.com/jrapoport/gothic/core/context"
	"github.com/jrapoport/gothic/utils"
	"github.com/spf13/cobra"
)

var codeCmd = &cobra.Command{
	Use:  "code [COUNT]",
	Long: "generate new signup codes ",
	RunE: codeRunE,
	Args: cobra.MaximumNArgs(1),
}

var (
	codeUses   int
	codeOutput string
)

func init() {
	fs := codeCmd.Flags()
	fs.IntVarP(&codeUses, "max-uses", "m", 1, "maximum times a code can be used")
	fs.StringVarP(&codeOutput, "out", "o", "", "output csv to file path")
}

func codeRunE(_ *cobra.Command, args []string) error {
	var count int
	if len(args) > 0 {
		var err error
		count, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
	}
	if count <= 0 {
		count = 1
	}
	client, err := root.NewAdminClient()
	if err != nil {
		return err
	}
	defer func() {
		client.Close()
	}()
	req := &admin.CreateSignupCodesRequest{
		Uses:  int64(codeUses),
		Count: int64(count),
	}
	res, err := client.CreateSignupCodes(context.Background(), req)
	if err != nil {
		return err
	}
	list := res.GetCodes()
	if codeOutput != "" {
		err = utils.WriteCSV(codeOutput, "code", list)
		if err != nil {
			return err
		}
	}
	fmt.Printf("created %d codes\n", len(list))
	utils.PrintGrid(os.Stdout, list, 10)
	return nil
}
