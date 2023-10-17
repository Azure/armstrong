package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ms-henglu/pal/formatter"
	"github.com/ms-henglu/pal/formatter/azapi"
	"github.com/ms-henglu/pal/trace"
)

const version = "0.4.0"

var showHelp = flag.Bool("help", false, "Show help")
var showVersion = flag.Bool("version", false, "Show version")

func main() {
	input := ""
	output := ""
	mode := ""

	flag.StringVar(&input, "i", "", "Input terraform log file")
	flag.StringVar(&output, "o", "", "Output directory")
	flag.StringVar(&mode, "m", "markdown", "Output format, allowed values are `markdown`, `oav` and `azapi`")

	// backward compatibility, the first argument is the input file
	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); err == nil {
			input = os.Args[1]
			mode = "markdown"
		}
	}
	if input == "" {
		flag.Parse()
		if *showHelp {
			flag.Usage()
			os.Exit(0)
		}
		if *showVersion {
			fmt.Println(version)
			os.Exit(0)
		}
	}
	if input == "" {
		flag.Usage()
		log.Fatalf("[ERROR] input file is required")
	}

	if output == "" {
		output = path.Dir(input)
	}

	log.Printf("[INFO] input file: %s", input)
	log.Printf("[INFO] output directory: %s", output)
	log.Printf("[INFO] output format: %s", mode)

	traces, err := trace.RequestTracesFromFile(input)
	if err != nil {
		log.Fatalf("[ERROR] failed to parse request traces: %v", err)
	}

	for _, t := range traces {
		out := trace.VerifyRequestTrace(t)
		if len(out) > 0 {
			log.Printf("[WARN] verification failed: url %s\n%s", t.Url, strings.Join(out, "\n"))
		}
	}

	switch mode {
	case "oav":
		format := formatter.OavTrafficFormatter{}
		index := 0
		for _, t := range traces {
			out := format.Format(t)
			index = index + 1
			outputPath := path.Join(output, fmt.Sprintf("trace-%d.json", index))
			if err := os.WriteFile(outputPath, []byte(out), 0644); err != nil {
				log.Fatalf("[ERROR] failed to write file: %v", err)
			}
			log.Printf("[INFO] output file: %s", outputPath)
		}
	case "markdown":
		content := markdownPrefix
		format := formatter.MarkdownFormatter{}
		for _, t := range traces {
			content += format.Format(t)
		}
		outputPath := path.Clean(path.Join(output, "output.md"))
		if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
			log.Fatalf("[ERROR] failed to write file: %v", err)
		}
		log.Printf("[INFO] output file: %s", outputPath)
	case "azapi":
		content := azapiPrefix
		format := azapi.AzapiFormatter{}
		for _, t := range traces {
			if res := format.Format(t); res != "" {
				content += res
				content += "\n"
			}
		}
		outputPath := path.Clean(path.Join(output, "pal-main.tf"))
		if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
			log.Fatalf("[ERROR] failed to write file: %v", err)
		}
		log.Printf("[INFO] output file: %s", outputPath)
	default:
		log.Fatalf("[ERROR] unsupported output format: %s", mode)
	}

}

const markdownPrefix = `<!--
Tips:

1. Use Markdown preview mode to get a better reading experience.
2. If you want to select some of the request traces, in VSCode, use shortcut "Ctrl + K, 0" to fold all blocks.

-->

`

const azapiPrefix = `
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

`
