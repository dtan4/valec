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
	secretFile     string
	dotenvTemplate string
	dryRun         bool
	output         string
	override       bool
)
