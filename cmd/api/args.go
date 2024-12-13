package main

import "github.com/urfave/cli/v2"

const (
	Name    = "CROU API"
	Usage   = "CROU API 서비스입니다"
	Version = "1.0.0"
)

var Flags = []cli.Flag{
	EnvFlag,
}

var EnvFlag = &cli.StringFlag{
	Name:    "env",
	Aliases: []string{"e"},
	Value:   "env.yaml",
	Usage:   "env 파일 경로를 호출합니다",
}
