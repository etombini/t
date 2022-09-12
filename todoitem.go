package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TodoItem struct {
	Content  string
	Priority int
	Hashtags []string // #words in Content
	Ats      []string // @words in Content
	Plus     []string // +words in Content
	Creation time.Time
}

type TodoList struct {
	TodoList []TodoItem
	config   Config
}

func NewTodoList(cfg Config) (*TodoList, error) {
	tl := TodoList{
		TodoList: make([]TodoItem, 0),
		config:   cfg,
	}

	f, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return &tl, fmt.Errorf("Can not open file %s: %s", cfg.Path, err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return &tl, fmt.Errorf("Can not read todo.go file at %s: %s", cfg.Path, err)
	}
	if len(content) == 0 {
		return &tl, nil
	}
	if err := json.Unmarshal(content, &tl); err != nil {
		return &tl, fmt.Errorf("Can not unmarshal json content: %s", err)
	}
	tl.config = cfg
	return &tl, nil
}

func (tl *TodoList) write() error {
	content, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return fmt.Errorf("Can not marshal content to json: %s", err)
	}
	f, err := os.OpenFile(tl.config.Path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Can not open file %s: %s", tl.config.Path, err)
	}
	_, err = f.Write(content)
	if err != nil {
		e := fmt.Errorf("Can not write content to todo file %s, see the content below: \n%s\n", tl.config.Path, err)
		return e
	}
	return nil
}

func (tl *TodoList) filter(filter string) []TodoItem {
	f := func(content []string, word string) bool {
		for _, v := range content {
			if v == word {
				return true
			}
		}
		return false
	}
	items := make([]TodoItem, 0)
	for i, t := range tl.TodoList {
		if f(t.Ats, filter) || f(t.Hashtags, filter) || f(t.Plus, filter) {
			items = append(items, tl.TodoList[i])
		} else {
			items = append(items, TodoItem{})
		}
	}
	return items
}

func (tl *TodoList) List(args []string) error {
	todolist := make([]TodoItem, 0)
	if len(args) > 0 {
		todolist = tl.filter(args[0])
	} else {
		todolist = tl.TodoList
	}
	for i := range todolist {
		if todolist[i].Content == "" {
			continue
		}
		content := todolist[i].Content

		//search for AtWords only if there is a @
		for j := range todolist[i].Ats {
			content = strings.Replace(content, todolist[i].Ats[j],
				tl.config.colorMap[tl.config.AtColor]("%s", todolist[i].Ats[j]), -1)
		}

		//search for HashWord only if there is a #
		for j := range todolist[i].Hashtags {
			content = strings.Replace(content, todolist[i].Hashtags[j],
				tl.config.colorMap[tl.config.HashColor]("%s", todolist[i].Hashtags[j]), -1)
		}

		//search for PlusWord only if there is a +
		for j := range todolist[i].Plus {
			content = strings.Replace(content, todolist[i].Plus[j],
				tl.config.colorMap[tl.config.PlusColor]("%s", todolist[i].Plus[j]), -1)
		}

		var p string
		switch todolist[i].Priority {
		case 1:
			p = tl.config.colorMap["p1"](" ")
		case 2:
			p = tl.config.colorMap["p2"](" ")
		case 3:
			p = tl.config.colorMap["p3"](" ")
		case 4:
			p = " "
		}

		fmt.Printf("[%2d] %s %s\n", i, p, content)
		// if i%2 == 1 {
		// 	fmt.Printf("[%3d] %s\n", i, content)
		// } else {
		// 	fmt.Printf(color.New(color.BgBlack).SprintfFunc()("[%3d] %s\n", i, content))
		// }
	}
	return nil
}

func (tl *TodoList) ListAt(args []string) error {
	return tl.listWord("at", args)
}

func (tl *TodoList) ListHash(args []string) error {
	return tl.listWord("hash", args)
}

func (tl *TodoList) ListPlus(args []string) error {
	return tl.listWord("plus", args)
}

func (tl *TodoList) listWord(kind string, args []string) error {
	uniqueWords := make(map[string]bool)
	switch kind {
	case "at":
		for i := range tl.TodoList {
			for j := range tl.TodoList[i].Ats {
				uniqueWords[tl.TodoList[i].Ats[j]] = true
			}
		}
	case "hash":
		for i := range tl.TodoList {
			for j := range tl.TodoList[i].Hashtags {
				uniqueWords[tl.TodoList[i].Hashtags[j]] = true
			}
		}
	case "plus":
		for i := range tl.TodoList {
			for j := range tl.TodoList[i].Plus {
				uniqueWords[tl.TodoList[i].Plus[j]] = true
			}
		}
	}
	if len(uniqueWords) == 0 {
		return nil
	}
	sortedWords := make([]string, len(uniqueWords))
	i := 0
	for w := range uniqueWords {
		sortedWords[i] = w
		i++
	}
	sort.Strings(sortedWords)
	for k := range sortedWords {
		fmt.Printf("%s \t", sortedWords[k])
	}
	fmt.Println("")
	return nil
}

func (tl *TodoList) Delete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("You must provide a Todo number to delete")
	}
	idx, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%s is not a valid index: %s", args[0], err)
	}
	if idx > len(tl.TodoList)-1 || idx < 0 {
		return fmt.Errorf("%s is not an index in available range", args[0])
	}
	for i := idx; i+1 < len(tl.TodoList); i++ {
		tl.TodoList[i] = tl.TodoList[i+1]
	}
	tl.TodoList = tl.TodoList[:len(tl.TodoList)-1]
	return tl.write()
}

func (tl *TodoList) Edit(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Todo index missing")
	}
	idx, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("%s is not a valid index: %s", args[0], err)
	}
	if idx > len(tl.TodoList)-1 {
		return fmt.Errorf("%s is not an index in available range", args[0])
	}
	var content string
	if idx >= 0 {
		content = tl.TodoList[idx].Content
	}

	edited, err := tmpEdit(content, "todo-", tl.config.Editor)
	if err != nil {
		return err
	}
	edited = strings.TrimSpace(edited)

	todo := TodoItem{
		Priority: 4,
		Content:  edited,
		Creation: time.Now().UTC(),
		Hashtags: make([]string, 0),
		Ats:      make([]string, 0),
		Plus:     make([]string, 0),
	}
	if strings.Contains(todo.Content, "p1") {
		todo.Priority = 1
	}
	if strings.Contains(todo.Content, "p2") {
		todo.Priority = 2
	}
	if strings.Contains(todo.Content, "p3") {
		todo.Priority = 3
	}

	//search for AtWord only if there is a @
	if strings.Contains(todo.Content, "@") {
		words := strings.Split(todo.Content, " ")
		uniqueAts := make(map[string]bool)
		for i := range words {
			if strings.HasPrefix(words[i], "@") {
				uniqueAts[words[i]] = true
			}
		}
		for w := range uniqueAts {
			todo.Ats = append(todo.Ats, w)
		}
	}

	//search for HashWord only if there is a #
	if strings.Contains(todo.Content, "#") {
		words := strings.Split(todo.Content, " ")
		uniqueHash := make(map[string]bool)
		for i := range words {
			if strings.HasPrefix(words[i], "#") {
				uniqueHash[words[i]] = true
			}
		}
		for w := range uniqueHash {
			todo.Hashtags = append(todo.Hashtags, w)
		}
	}

	//search for PlusWord only of there is a +
	if strings.Contains(todo.Content, "+") {
		words := strings.Split(todo.Content, " ")
		uniquePlus := make(map[string]bool)
		for i := range words {
			if strings.HasPrefix(words[i], "+") {
				uniquePlus[words[i]] = true
			}
		}
		for w := range uniquePlus {
			todo.Plus = append(todo.Plus, w)
		}
	}

	if idx < 0 { //This is not a replacement/edit, so it is appended
		tl.TodoList = append(tl.TodoList, todo)
	} else { //This is a replacement
		tl.TodoList[idx] = todo
	}
	return tl.write()
}

func tmpEdit(content string, prefix string, editor string) (string, error) {
	//get tmp file
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	filename := f.Name()

	//write content to tmp file
	if content != "" {
		if _, err := f.Write([]byte(content)); err != nil {
			return "", err
		}
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	defer os.Remove(filename)

	//open file with $EDITOR
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	//get content from file
	edited, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(edited), nil
	//delete file
	// see defer

}
