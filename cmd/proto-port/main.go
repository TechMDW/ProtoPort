package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TechMDW/ProtoPort/internal/github"
	"github.com/TechMDW/ProtoPort/internal/protoc"
	"github.com/TechMDW/ProtoPort/internal/utilities"
)

type Command struct {
	fs *flag.FlagSet

	input  string
	output string
	pat    string
	lang   string
}
type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

var (
	version = "0.0.6"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	configPath, err := os.UserConfigDir()

	if err != nil {
		log.Println(err)
	}

	techMDWPath := filepath.Join(configPath, "TechMDW")

	if _, err := os.Stat(techMDWPath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(techMDWPath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	ProtoPortPath := filepath.Join(techMDWPath, "ProtoPort")

	// clear the previous files
	utilities.DeleteAll(ProtoPortPath)

	if err := root(os.Args[1:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (c *Command) Name() string {
	return c.fs.Name()
}

func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *Command) Run() error {
	var output string
	var input string
	execution, err := os.Executable()
	if err != nil {
		return err
	}
	executionPath := filepath.Dir(execution)

	switch c.Name() {
	case "github":

		if c.input == "" {
			return fmt.Errorf("input is required")
		}

		if c.lang == "" {
			return fmt.Errorf("lang is required")
		}

		if c.output == "" {
			output = filepath.Join(executionPath, "outputs")
			if _, err := os.Stat(output); os.IsNotExist(err) {
				os.Mkdir(output, 0755)
			}
		} else if strings.HasPrefix(c.output, "./") {
			if _, err := os.Stat(c.output); os.IsNotExist(err) {
				return fmt.Errorf("output path does not exist")
			}
			output = c.output
		} else {
			if _, err := os.Stat(c.output); os.IsNotExist(err) {
				return fmt.Errorf("output path does not exist")
			}
			output = c.output
		}

		if c.pat == "" {
			input, err = github.GithubReadAndGenrateProtos(c.input, "", "", true)
			if err != nil {
				return err
			}
		} else {
			input, err = github.GithubReadAndGenrateProtos(c.input, "", c.pat, false)
			if err != nil {
				return err
			}
		}

		err = protoc.ReadDirForProto(input, output, c.lang, true)
		if err != nil {
			return err
		}

	case "basic":
		log.Println("Basic Commands")
		if c.lang == "" {
			return fmt.Errorf("lang is required")
		}

		if c.input == "" {
			input = filepath.Join(executionPath, "inputs")
			if _, err := os.Stat(input); os.IsNotExist(err) {
				return fmt.Errorf("no input folder was found")
			}
		} else {
			input = c.input
		}

		if c.output == "" {
			output = filepath.Join(executionPath, "outputs")
			if _, err := os.Stat(output); os.IsNotExist(err) {
				os.Mkdir(output, 0755)
			}
		} else if strings.HasPrefix(c.output, "./") {
			if _, err := os.Stat(c.output); os.IsNotExist(err) {
				return fmt.Errorf("output path does not exist")
			}
			output = c.output
		} else {
			if _, err := os.Stat(c.output); os.IsNotExist(err) {
				return fmt.Errorf("output path does not exist")
			}
			output = c.output
		}

		err = protoc.ReadDirForProto(input, output, c.lang, false)
		if err != nil {
			return err
		}

	case "help":
		fmt.Println("Usage: ProtoPort [command] [options]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  github 	- generate proto files from a github repo")
		fmt.Println("  basic 	- generate proto files from a local folder")
		fmt.Println("  version 	- print the version of ProtoPort")
		fmt.Println("")

		fmt.Println("Options:")
		fmt.Println("  -inputs 	- Path to proto files (if path not specified, it will scan current folder for proto files and generate them with the same folder structure)")
		fmt.Println("  -outputs 	- Path to output files (if not specified, it will store the file in the current folder and it will preserve the input folder structure)")
		fmt.Println("  -lang 	- Choose language to generate *(required)[go, cpp, csharp, java, kotlin, objc, php, python, pyi, ruby]")
		fmt.Println("  -pat 		- Github Personal Access Token (only needed if repo is private)")
		fmt.Println("")

	case "version":
		fmt.Println("ProtoProt Version " + version)

	case "v":
		fmt.Println("ProtoProt Version " + version)
	}

	return nil
}

func GithubCommands() *Command {
	gc := &Command{
		fs: flag.NewFlagSet("github", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.input, "inputs", "", "The url to the github repo (if path not specified, it will scan the whole repo for proto files and generate them with the same folder structure)")
	gc.fs.StringVar(&gc.output, "outputs", "", "Path to output files (if not specified, it will store the file in the current folder and it will preserve the input folder structure)")
	gc.fs.StringVar(&gc.pat, "pat", "", "Github Personal Access Token (only needed if repo is private)")
	gc.fs.StringVar(&gc.lang, "lang", "", "Choose language to generate *(required)[go, cpp, csharp, java, kotlin, objc, php, python, pyi, ruby]")

	return gc
}

func BasicCommands() *Command {
	bc := &Command{
		fs: flag.NewFlagSet("basic", flag.ContinueOnError),
	}

	bc.fs.StringVar(&bc.input, "inputs", "", "Path to proto files (if path not specified, it will scan current folder for proto files and generate them with the same folder structure)")
	bc.fs.StringVar(&bc.output, "outputs", "", "Path to output files (if not specified, it will store the file in the current folder and it will preserve the input folder structure)")
	bc.fs.StringVar(&bc.lang, "lang", "", "Choose language to generate *(required)[go, cpp, csharp, java, kotlin, objc, php, python, pyi, ruby]")

	return bc
}

func HelpCommands() *Command {
	bc := &Command{
		fs: flag.NewFlagSet("help", flag.ContinueOnError),
	}
	return bc
}

func VersionCommands() *Command {
	bc := &Command{
		fs: flag.NewFlagSet("version", flag.ContinueOnError),
	}
	return bc
}

func VCommands() *Command {
	bc := &Command{
		fs: flag.NewFlagSet("v", flag.ContinueOnError),
	}
	return bc
}

func root(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: ProtoPort [command] [options]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  github 	- generate proto files from a github repo")
		fmt.Println("  basic 	- generate proto files from a local folder")
		fmt.Println("  version 	- print the version of ProtoPort")
		fmt.Println("  help 	- print the help menu")
		fmt.Println("")

		fmt.Println("Options:")
		fmt.Println("  -inputs 	- Path to proto files (if path not specified, it will scan current folder for proto files and generate them with the same folder structure)")
		fmt.Println("  -outputs 	- Path to output files (if not specified, it will store the file in the current folder and it will preserve the input folder structure)")
		fmt.Println("  -lang 	- Choose language to generate *(required)[go, cpp, csharp, java, kotlin, objc, php, python, pyi, ruby]")
		fmt.Println("  -pat 		- Github Personal Access Token (only needed if repo is private)")
		fmt.Println("")
		return nil
	}

	cmds := []Runner{
		GithubCommands(),
		BasicCommands(),
		HelpCommands(),
		VersionCommands(),
		VCommands(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}
