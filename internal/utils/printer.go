package utils

import (
	"fmt"
	"github.com/DarthBenro008/pigil/internal/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
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

func GreenPrinter(information string) {
	fmt.Printf("%s%s%s\n", Green, information, Reset)
}

func StatusPrinter(data *[]types.ConfigurationInformation) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Setting Name", "Value"})
	t.SetAutoIndex(true)
	for _, items := range *data {
		if items.Key != UserAT && items.Key != UserRT {
			t.AppendRow(table.Row{items.Key, items.Value})
		}
	}
	t.Render()
}
