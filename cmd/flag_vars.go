package cmd

// global flag variable
var (
	debug     bool
	keyAlias  string
	noColor   bool
	tableName string
)

// command-specific flag variable
var (
	configFile     string
	dotenvTemplate string
	dryRun         bool
	override       bool
)
