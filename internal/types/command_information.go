package types

type CommandInformation struct {
	CommandName      string
	CommandArguments []string
	ExecutionTime    int64
	WasSuccessful    bool
}
