package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

const (
	VERSION = "0.0.1"
)

var commands = []cli.Command{
	{
		Name:   "add",
		Usage:  "add alias",
		Action: cmdAdd,
	},
	{
		Name:   "delete",
		Usage:  "delete alias",
		Action: cmdDelete,
	},
}

type config struct {
	Aliases map[string]string `toml:"aliases"`
}

func configDir() string {
	var dir string

	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "goto")
		}
		dir = filepath.Join(dir, "goto")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "goto")
	}
	return dir
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func (cfg *config) load() error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}
	file := filepath.Join(dir, "config.toml")

	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	_, err = os.Create(file)
	if err != nil {
		return err
	}

	return nil
}

func cmdAdd(c *cli.Context) error {
	fmt.Println("Add alias: ", c.Args().First())
	return nil
}

func cmdDelete(c *cli.Context) error {
	fmt.Println("Delete alias: ", c.Args().First())
	return nil
}

func appRun(c *cli.Context) error {
	fmt.Println("dummy\n")
	return nil
}

func run() int {
	app := cli.NewApp()
	app.Name = "goto"
	app.Usage = "Goto "
	app.Version = VERSION
	app.Commands = commands
	app.Action = appRun

	return msg(app.Run(os.Args))

	return 0
}

func main() {
	os.Exit(run())
}
