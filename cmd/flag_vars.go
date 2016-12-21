package cmd

// global flag variable
var (
	debug     bool
	keyAlias  string
	noColor   bool
	tableName string
	region    string
)

// command-specific flag variable
var (
	keys           string
	secretFile     string
	dotenvTemplate string
	dryRun         bool
	interactive    bool
	output         string
	override       bool
	quote          bool
)
