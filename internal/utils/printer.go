package utils

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"pigil/internal/types"
	"time"
)

func PrintInformation(data *[]types.CommandInformation) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Time of Execution", "Command Name",
		"Command Arguments", "Execution Time (in seconds)",
		"Was Successful"})
	t.SetAutoIndex(true)
	for _, items := range *data {
		t.AppendRow(table.Row{time.UnixMicro(items.TimeOfExecution).Format(
			time.Stamp), items.CommandName, items.CommandArguments,
			items.ExecutionTime,
			items.WasSuccessful})
	}
	t.Render()
}
