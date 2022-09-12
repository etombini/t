package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

// Config holds the configuration of the application.
// The corresponding configuration file can be found here:
// $HOME/.config/todo.go.yaml
type Config struct {
	Path      string // where to find tne todo file
	AtColor   string
	HashColor string
	PlusColor string
	Editor    string
	colorMap  map[string]func(format string, a ...interface{}) string
}

func getConfig() (Config, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/.config/todo.go.yaml", os.Getenv("HOME")))
	if err != nil {
		return Config{}, err
	}
	c := &Config{}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return Config{}, err
	}
	if c.AtColor == "" {
		c.AtColor = "blue"
	}
	if c.HashColor == "" {
		c.HashColor = "green"
	}
	if c.PlusColor == "" {
		c.PlusColor = "red"
	}
	c.colorMap = make(map[string]func(format string, a ...interface{}) string)
	c.colorMap["black"] = color.New(color.FgBlack).SprintfFunc()
	c.colorMap["blue"] = color.New(color.FgBlue).SprintfFunc()
	c.colorMap["cyan"] = color.New(color.FgCyan).SprintfFunc()
	c.colorMap["green"] = color.New(color.FgGreen).SprintfFunc()
	c.colorMap["hiblack"] = color.New(color.FgHiBlack).SprintfFunc()
	c.colorMap["hiblue"] = color.New(color.FgHiBlue).SprintfFunc()
	c.colorMap["hicyan"] = color.New(color.FgHiCyan).SprintfFunc()
	c.colorMap["higreen"] = color.New(color.FgHiGreen).SprintfFunc()
	c.colorMap["himagenta"] = color.New(color.FgHiMagenta).SprintfFunc()
	c.colorMap["hired"] = color.New(color.FgHiRed).SprintfFunc()
	c.colorMap["hiwhite"] = color.New(color.FgHiWhite).SprintfFunc()
	c.colorMap["hiyellow"] = color.New(color.FgHiYellow).SprintfFunc()
	c.colorMap["magenta"] = color.New(color.FgMagenta).SprintfFunc()
	c.colorMap["red"] = color.New(color.FgRed).SprintfFunc()
	c.colorMap["white"] = color.New(color.FgWhite).SprintfFunc()
	c.colorMap["yellow"] = color.New(color.FgYellow).SprintfFunc()

	c.colorMap["p1"] = color.New(color.BgRed).SprintfFunc()
	c.colorMap["p2"] = color.New(color.BgYellow).SprintfFunc()
	c.colorMap["p3"] = color.New(color.BgBlue).SprintfFunc()
	c.colorMap["p4"] = color.New(color.BgBlack).SprintfFunc()

	if _, ok := c.colorMap[c.AtColor]; !ok {
		return Config{}, fmt.Errorf("Unknown color %s from configuration", c.AtColor)
	}
	if _, ok := c.colorMap[c.HashColor]; !ok {
		return Config{}, fmt.Errorf("Unknown color %s from configuration", c.HashColor)
	}
	if _, ok := c.colorMap[c.PlusColor]; !ok {
		return Config{}, fmt.Errorf("Unknown color %s from configuration", c.PlusColor)
	}

	if c.Editor == "" {
		c.Editor = "vim"
	}
	return *c, nil
}
