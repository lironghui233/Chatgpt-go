package cmd

import "flag"

type CommandArgs struct {
	Dict     string
	Config   string
	InitDict bool
}

var Args *CommandArgs

func init() {
	dict := flag.String("dict", "dict.txt", "敏感词汇词库")
	config := flag.String("config", "config.yaml", "配置文件")
	initDict := flag.Bool("init-dict", false, "是否初始化词库")
	flag.Parse()
	Args = &CommandArgs{}
	Args.Dict = *dict
	Args.Config = *config
	Args.InitDict = *initDict
}
