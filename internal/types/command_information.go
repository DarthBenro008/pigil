package types

type CommandInformation struct {
	CommandName      string
	CommandArguments []string
	TimeOfExecution  int64
	ExecutionTime    float64
	WasSuccessful    bool
}
