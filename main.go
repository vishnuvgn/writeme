/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"writeme/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {

	cmd.SetVersionInfo(version, commit, date)

	cmd.Execute()
}
