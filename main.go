package main

import (
	"github.com/ahelmy/trunk-based-dev/cmd"
	"github.com/ahelmy/trunk-based-dev/app"
)

func main() {
	app.InitApp()
	cmd.Execute()
}