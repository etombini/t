package main

import (
	"fmt"
	"os"
)

var prefix = os.Getenv("HOME") + "todos"
var Version string

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

	var helps helpList = make([]help, 0)

	//verbs
	list := make(map[string]bool)
	list["list"] = true
	list["l"] = true
	helps = append(helps, help{
		verb:        list,
		description: "List todos. An argument can be an @word, @word or +word",
	})

	listAt := make(map[string]bool)
	listAt["list@"] = true
	listAt["l@"] = true
	helps = append(helps, help{
		verb:        listAt,
		description: "List existing @words.",
	})

	listHash := make(map[string]bool)
	listHash["list#"] = true
	listHash["l#"] = true
	helps = append(helps, help{
		verb:        listHash,
		description: "List existing #words.",
	})

	listPlus := make(map[string]bool)
	listPlus["list+"] = true
	listPlus["l+"] = true
	helps = append(helps, help{
		verb:        listPlus,
		description: "List existing +words.",
	})

	delete_ := make(map[string]bool)
	delete_["delete"] = true
	delete_["d"] = true
	helps = append(helps, help{
		verb:        delete_,
		description: "Delete a todo. The argument is a todo id.",
	})

	edit := make(map[string]bool)
	edit["edit"] = true
	edit["e"] = true
	helps = append(helps, help{
		verb:        edit,
		description: "Edit a todo. Argument is a todo id. Editor can be set in the configuration file",
	})

	help_ := make(map[string]bool)
	help_["help"] = true
	help_["h"] = true
	helps = append(helps, help{
		verb:        help_,
		description: "Print this help message.",
	})

	version := make(map[string]bool)
	version["version"] = true
	version["v"] = true
	helps = append(helps, help{
		verb:        version,
		description: "Print current version.",
	})

	alias := make(map[string]bool)
	alias["alias"] = true
	helps = append(helps, help{
		verb:        alias,
		description: "Print alias commands for your shell",
	})

	switch {
	case list[os.Args[1]]:
		err := tl.List(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list todo: %s\n", err)
			os.Exit(1)
		}
		return
	case listAt[os.Args[1]]:
		err := tl.ListAt(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list @words: %s\n", err)
			os.Exit(1)
		}
		return
	case listHash[os.Args[1]]:
		err := tl.ListHash(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list #words: %s\n", err)
			os.Exit(1)
		}
		return
	case listPlus[os.Args[1]]:
		err := tl.ListPlus(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not list #words: %s\n", err)
			os.Exit(1)
		}
		return
	case delete_[os.Args[1]]:
		if err := tl.Delete(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Can not delete todo item: %s\n", err)
			os.Exit(1)
		}
		if err := tl.List([]string{}); err != nil {
			fmt.Fprintf(os.Stderr, "Can not list todo: %s\n", err)
			os.Exit(1)
		}
		return
	case edit[os.Args[1]]:
		if err := tl.Edit(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Can not edit todo item: %s\n", err)
			os.Exit(1)
		}
		if err := tl.List([]string{}); err != nil {
			fmt.Fprintf(os.Stderr, "Can not list todo: %s\n", err)
			os.Exit(1)
		}
		return
	case help_[os.Args[1]]:
		printVersion()
		helps.printHelp()
		return
	case version[os.Args[1]]:
		printVersion()
		return
	case alias[os.Args[1]]:
		if err := printAlias(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Can not print aliases: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknonw command\n")
		helps.printHelp()
	}
}

func printVersion() {
	fmt.Printf("%s\n", Version)
}

func printAlias(args []string) error {
	shell := "bash"
	if len(args) >= 1 {
		shell = args[0]
	}

	switch shell {
	case "bash":
		fmt.Println("alias tl=\"t l\"")
		fmt.Println("alias tl@=\"t l@\"")
		fmt.Println("alias tl#=\"t l#\"")
		fmt.Println("alias tl+=\"t l+\"")
		fmt.Println("alias te=\"t e\"")
		fmt.Println("alias td=\"t d\"")
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	return nil
}
