package protoc

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TechMDW/ProtoPort/internal/utilities"
)

func ReadDirForProto(path string, output string, lang string, github bool) error {
	directory, err := os.Open(path)
	
	if err != nil {
		return err
	}

	files, err := directory.Readdir(-1)

	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			subdir := filepath.Join(output, file.Name())

			if _, err := os.Stat(subdir); os.IsNotExist(err) {
				os.Mkdir(subdir, 0755)
			}

			ReadDirForProto(filepath.Join(path, file.Name()), subdir, lang, github)
		}

		if utilities.CheckForFileExtension(file.Name(), ".proto") {
			
			err := BuildProto(path, output, file.Name(), lang)
			
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func BuildProto(protoPath, output, protoFile string, lang string) error {
	var command *exec.Cmd

	fullProtoFilePath := filepath.Join(protoPath, protoFile)

	if lang == "" {
		return fmt.Errorf("language not specified")
	}

	switch lang {
	case "go":
		command = exec.Command("protoc", "--go_out="+output, "--go_opt=paths=source_relative", "--go-grpc_out="+output, "--go-grpc_opt=paths=source_relative", "--proto_path="+protoPath, fullProtoFilePath)
	case "cpp":
		command = exec.Command("protoc", "--cpp_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "csharp":
		command = exec.Command("protoc", "--csharp_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "java":
		command = exec.Command("protoc", "--java_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "python":
		command = exec.Command("protoc", "--python_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "ruby":
		command = exec.Command("protoc", "--ruby_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "pyi":
		command = exec.Command("protoc", "--pyi_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "php":
		command = exec.Command("protoc", "--php_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "objc":
		command = exec.Command("protoc", "--objc_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "kotlin":
		command = exec.Command("protoc", "--kotlin_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "node":
		command = exec.Command("protoc", "--node_out="+output, "--proto_path="+protoPath, fullProtoFilePath)
	case "dart":
		command = exec.Command("protoc", "--dart_out="+output, "--proto_path="+protoPath, fullProtoFilePath)

	default:
		return errors.New("language not supported")
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return fmt.Errorf("error building proto file")
	}

	if out.String() == "" {
		fmt.Println("Build successfull:", protoFile)
	}

	return nil
}
