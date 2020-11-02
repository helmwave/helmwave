package main

import "github.com/urfave/cli/v2"

func authors() []*cli.Author {
	return []*cli.Author{
		{
			Name:  "💎 Dmitriy Zhilyaev",
			Email: "helmwave+zhilyaev.dmitriy@gmail.com",
		},
	}
}
