package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.arsenm.dev/lure-repo-bot/internal/analyze"
	"go.arsenm.dev/lure-repo-bot/internal/shutils"
	"go.arsenm.dev/lure-repo-bot/internal/spdx"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func init() {
	err := spdx.Update()
	if err != nil {
		fatalErr(err)
	}
}

func main() {
	ctx := context.Background()

	var files []*os.File
	for _, arg := range os.Args[1:] {
		file, err := os.Open(arg)
		if err != nil {
			fatalErr(err)
		}
		files = append(files, file)
	}

	wd, err := os.Getwd()
	if err != nil {
		fatalErr(err)
	}

	issuesFound := false
	for _, file := range files {
		fl, err := syntax.NewParser().Parse(file, "lure.sh")
		if err != nil {
			fatalErr(err)
		}

		var nopRWC shutils.NopRWC
		runner, err := interp.New(
			interp.Env(expand.ListEnviron()),
			interp.StdIO(nopRWC, nopRWC, os.Stderr),
			interp.ExecHandler(shutils.NopExec),
			interp.ReadDirHandler(shutils.NopReadDir),
			interp.OpenHandler(shutils.NopOpen),
			interp.StatHandler(shutils.NopStat),
		)
		if err != nil {
			fatalErr(err)
		}

		err = runner.Run(ctx, fl)
		if err != nil {
			fatalErr(err)
		}

		findings, err := analyze.AnalyzeScript(runner, fl)
		if err != nil {
			fatalErr(err)
		}

		flName := strings.TrimPrefix(file.Name(), wd)
		flName = strings.TrimPrefix(flName, "/")

		fmt.Println(flName + ":")
		if len(findings) == 0 {
			fmt.Println("\tNo issues found!")
		} else {
			issuesFound = true
			for _, finding := range findings {
				var name string
				if finding.Index != nil {
					name = fmt.Sprintf(
						"%s[%v] %s",
						finding.ItemName,
						finding.Index,
						finding.ItemType,
					)
				} else {
					name = fmt.Sprintf(
						"%s %s",
						finding.ItemName,
						finding.ItemType,
					)
				}

				msg := fmt.Sprintf(finding.Msg, name)

				if finding.ExtraMsg == "" {
					fmt.Printf("\tLine %d: %s\n", finding.Line, msg)
				} else {
					fmt.Printf("\tLine %d: %s\n\t\t%s\n", finding.Line, msg, finding.ExtraMsg)
				}
			}
		}
	}

	if issuesFound {
		os.Exit(1)
	}
}

func fatalErr(a ...any) {
	fmt.Println(append([]any{"error:"}, a...)...)
	os.Exit(1)
}
