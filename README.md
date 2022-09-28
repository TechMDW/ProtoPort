# ProtoPort

**Description**: This CLI tool will allow you to build all of your proto files with one command

## Dependencies

[Protoc](https://developers.google.com/protocol-buffers/docs/downloads)

### Download Pre Built

All the Pre Built binary are available in the [release page](https://github.com/TechMDW/ProtoPort/releases)

### Installation with golang

```bash
go install github.com/TechMDW/ProtoPort/cmd/proto-port@latest
```

## Protoc plugins installation

For more information visit [grpc.io](https://grpc.io/)
| Lang | plugin |
| ---- | --------------------------------------------------------------- |
| Dart | `dart pub global activate protoc_plugin` |
| Go | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2` |
| Go | `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28` |

## Usage

```bash
ProtoPort [command] [options]
```

### Commands

| Syntax | Description                              |
| ------ | ---------------------------------------- |
| basic  | generate proto files from a local folder |
| github | generate proto files from a GitHub repo  |

#### Basic options

| Syntax   | Description                                                                                                                               |
| -------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| -inputs  | Path to proto files (if path not specified, it will scan current folder for proto files and generate them with the same folder structure) |
| -outputs | Path to output files (if not specified, it will store the file in the current folder and it will preserve the input folder structure)     |
| -lang    | Choose language to generate \*[required](go, cpp, csharp, java, kotlin, objc, php, python, pyi, ruby, dart, node)                         |

#### GitHub options

| Syntax   | Description                                                                                                                                      |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| -inputs  | The url to the GitHub repo (if path not specified, it will scan the whole repo for proto files and generate them with the same folder structure) |
| -outputs | Path to output files (if not specified, it will store the file in the current folder and it will preserve the input folder structure)            |
| -pat     | GitHub Personal Access Token (only needed if repo is private)                                                                                    |
| -lang    | Choose language to generate \*[required](go, cpp, csharp, java, kotlin, objc, php, python, pyi, ruby, dart, node)                                |

### TODO

Will work on this list when I got some free time. If you want to contribute, feel free to do so.

- [x] Compiling proto files.
- [x] Persistent folder structure pawn build.
- [x] Using GitHub as source for protofiles.
- [ ] Using Gitlab as source for protofiles.
- [ ] Using Bitbucket as source for protofiles.
- [ ] Arguments passthrough to protoc.
- [ ] Arguments pass through to protoc plugins.
- [ ] Auto installing plugins and protoc by default.
- [x] Add binary for windows.
- [x] Add binary for linux.
- [x] Add binary for mac.

## Getting help

If you have questions, concerns, bug reports, etc, please file an issue in this repository's Issue Tracker.

## Guidance on how to contribute

> All contributions to this project will be released under the TechMDW AB
> dedication. By submitting a pull request, or filing a bug, issue, or
> feature-request you are agreeing to comply with this waiver of copyright interest.
> Details can be found in our [LICENSE](LICENSE).

There are two primary ways to help:

- Using the issue tracker, and
- Changing the code-base.

## Using the issue tracker

Use the issue tracker to suggest feature requests, report bugs, and ask questions.
This is also a great way to connect with the developers of the project as well
as others who are interested in this solution.

Use the issue tracker to find ways to contribute. Find a bug or a feature, mention in
the issue that you will take on that effort, then follow the _Changing the code-base_
guidance below.

## Changing the code-base

Generally speaking, you should fork this repository, make changes in your
own fork, and then submit a pull-request. All new code should have associated unit
tests that validate implemented features and the presence or lack of defects.
Additionally, the code should follow any stylistic and architectural guidelines
prescribed by the project. In the absence of such guidelines, mimic the styles
and patterns in the existing code-base.

---
