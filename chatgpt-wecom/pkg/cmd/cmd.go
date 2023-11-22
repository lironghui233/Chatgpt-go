package cmd

import "flag"

type CommandArgs struct {
	Config string
}

var Args *CommandArgs

func init() {
	config := flag.String("config", "config.yaml", "应用程序配置文件")
	flag.Parse()
	Args = &CommandArgs{}
	Args.Config = *config
}
