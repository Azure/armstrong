# Parsing Azure's Logs - Pal

----

## Introduction

Pal is a simple tool to parse Terraform Azure provider logs.

It can output different formats of the request traces, including:

- Markdown:
  
  It can output the request traces in Markdown format, which can be used to create a GitHub issue or a forum post.
  
  The Markdown output supports expanding/collapsing the request traces. This is useful when the request trace is very long.

- OAV Traffic:

  It can output the request traces in OAV traffic format, which can be used to validate the request traces with [OAV](https://github.com/Azure/oav).

- AzAPI config:

    It can output the request traces in AzAPI config format, which can be used to reproduce the deployed resources with [AzAPI](https://registry.terraform.io/providers/Azure/azapi/latest).

## Usage

```bash
$ pal {path to terraform_log_file}
```

Full usage:

```bash
Usage of pal:
  -help
        Show help
  -i string
        Input terraform log file
  -m markdown
        Output format, allowed values are markdown, `oav` and `azapi` (default "markdown")
  -o string
        Output directory
  -version
        Show version
```

## Example

```bash
$ cd ./testdata
$ pal ./input.txt
```

Above command will generate a [markdown file named "output.md"](https://github.com/ms-henglu/pal/tree/main/testdata/output.md) in the same working directory.

## How to install?

```bash
$ go install github.com/ms-henglu/pal@latest
```