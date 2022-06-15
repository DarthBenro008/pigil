package utils

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"gnoty/internal/types"
	"os"
	"time"
)

func PrintInformation(data *[]types.CommandInformation) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Command Name", "Command Arguments",
		"Time of Execution",
		"Was Successful"})
	t.SetAutoIndex(true)
	for _, items := range *data {
		t.AppendRow(table.Row{items.CommandName, items.CommandArguments,
			time.UnixMicro(items.ExecutionTime).Format(time.Stamp),
			items.WasSuccessful})
	}
	t.Render()
}
