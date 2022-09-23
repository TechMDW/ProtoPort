package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TechMDW/ProtoPort/internal/github"
	"github.com/TechMDW/ProtoPort/internal/protoc"
)

type Command struct {
	fs *flag.FlagSet

	input  string
	output string
	pat    string
	lang   string

	protoc struct {
		proto_path               string
		encode                   string
		deterministic_output     string
		decode                   string
		decode_raw               string
		descriptor_set_in        string
		descriptor_set_out       string
		include_imports          string
		include_source_info      string
		dependency_out           string
		error_format             string
		fatal_warnings           string
		print_free_field_numbers string
		plugin                   string
	}
}
type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
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

	case "":
		if c.lang == "" {
			return fmt.Errorf("lang is required")
		}

		if c.input == "" {
			input = filepath.Join(executionPath, "inputs")
			if _, err := os.Stat(input); os.IsNotExist(err) {
				return fmt.Errorf("no input folder was found")
			}
		}

		if c.output == "" {
			output = filepath.Join(executionPath, "outputs")
			if _, err := os.Stat(output); os.IsNotExist(err) {
				os.Mkdir(output, 0755)
			}
		}

	}

	err = protoc.ReadDirForProto(input, output, c.lang)
	if err != nil {
		return err
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

func root(args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: ProtoPort [command] [options]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  github 	- generate proto files from a github repo")
		fmt.Println("  basic 	- generate proto files from a local folder")
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
