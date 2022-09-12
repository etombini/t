package main

import (
	"fmt"
	"os"
)

var prefix = os.Getenv("HOME") + "todos"
var Version string

func printHelp() {

}

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not read configuration file at %s/.config/todo.go.yaml: %s", os.Getenv("HOME"), err)
		os.Exit(1)
	}

	tl, err := NewTodoList(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not get content from file: %s", err)
	}

	if len(os.Args) == 1 {
		err := tl.Edit([]string{"-1"})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not add todo: %s", err)
			os.Exit(1)
		}
		return
	}

	//verbs
	list := make(map[string]bool)
	list["list"] = true
	list["l"] = true
	list["-l"] = true

	listAt := make(map[string]bool)
	listAt["list@"] = true
	listAt["l@"] = true
	listAt["-l@"] = true
	listAt["@"] = true

	listHash := make(map[string]bool)
	listHash["list#"] = true
	listHash["l#"] = true
	listHash["-l#"] = true
	listHash["#"] = true

	listPlus := make(map[string]bool)
	listPlus["list+"] = true
	listPlus["l+"] = true
	listPlus["-l+"] = true
	listPlus["+"] = true

	delete_ := make(map[string]bool)
	delete_["delete"] = true
	delete_["-d"] = true
	delete_["d"] = true

	edit := make(map[string]bool)
	edit["edit"] = true
	edit["-e"] = true
	edit["e"] = true

	help := make(map[string]bool)
	help["help"] = true
	help["-h"] = true
	help["h"] = true

	version := make(map[string]bool)
	version["version"] = true
	version["-v"] = true
	version["v"] = true

	switch {
	case list[os.Args[1]]:
		err := tl.List(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list todo: %s", err)
			os.Exit(1)
		}
		return
	case listAt[os.Args[1]]:
		err := tl.ListAt(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list @words: %s", err)
			os.Exit(1)
		}
		return
	case listHash[os.Args[1]]:
		err := tl.ListHash(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list #words: %s", err)
			os.Exit(1)
		}
		return
	case listPlus[os.Args[1]]:
		err := tl.ListPlus(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list #words: %s", err)
			os.Exit(1)
		}
		return
	case delete_[os.Args[1]]:
		if err := tl.Delete(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Can not delete todo item: %s", err)
			os.Exit(1)
		}
		if err := tl.List([]string{}); err != nil {
			fmt.Fprintf(os.Stderr, "Can not list todo: %s", err)
			os.Exit(1)
		}
		return
	case edit[os.Args[1]]:
		if err := tl.Edit(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Can not edit todo item: %s", err)
			os.Exit(1)
		}
		if err := tl.List([]string{}); err != nil {
			fmt.Fprintf(os.Stderr, "Can not list todo: %s", err)
			os.Exit(1)
		}
		return
	case version[os.Args[1]]:
		printVersion()
		return
	default:
		fmt.Fprintf(os.Stderr, "Unknonw command\n")
		printHelp()
	}
}

func printVersion() {
	fmt.Printf("%s\n", Version)
}
