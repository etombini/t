package main

import (
	"fmt"
	"strings"
)

type help struct {
	verb        map[string]bool
	description string
}

type helpList []help

func (h helpList) printHelp() {
	fmt.Printf("t is a tool for managing a todo list using a CLI\n")
	fmt.Printf("Usage:\n\tt <command> [arguments]\n\n")
	fmt.Printf("The commands are:\n")

	for i := range h {
		v := make([]string, len(h[i].verb), len(h[i].verb))
		j := 0
		for k := range h[i].verb {
			v[j] = k
			j++
		}
		fmt.Printf("\t%s", strings.Join(v, ", "))
		fmt.Printf("\t\t%s\n", h[i].description)
	}
}
