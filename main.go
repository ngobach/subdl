package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ngobach/subdl/sub"
	"github.com/urfave/cli/v2"
)

var defaultSubServiceName = "subscene"

func cmdDownload(c *cli.Context) error {
	svcName := c.String("service")
	svc, found := sub.Hub[svcName]
	if !found {
		return fmt.Errorf("could not find service %s", svcName)
	}
	keyword := ""
	survey.AskOne(&survey.Input{
		Message: "Please enter the movie title",
	}, &keyword, nil)
	result := svc.Search(keyword)
	if len(result) == 0 {

		return fmt.Errorf("no result found")
	}
	options := []string{}
	for _, entry := range result {
		options = append(options, entry.DisplayName)
	}
	var answer int
	survey.AskOne(&survey.Select{
		Message: "Please select one",
		Options: options,
	}, &answer, nil)
	fmt.Printf("Downloading %s\n", result[answer].DisplayName)
	svc.Download(result[answer].Id)
	return nil
}

func cmdServices(c *cli.Context) error {
	fmt.Println("Available sub services:")
	for name := range sub.Hub {
		isDefault := ""
		if name == defaultSubServiceName {
			isDefault = " (default)"
		}
		fmt.Printf("- %s%s\n", name, isDefault)
	}
	return nil
}

func main() {
	app := cli.App{
		Name:        "subdl",
		Usage:       "Downloading subtitles done right",
		Description: "A command line tool to download subtitles from various sources.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "service",
				Value: defaultSubServiceName,
				Usage: "Subtitle service",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "download",
				Usage:  "Search and download subtitles",
				Action: cmdDownload,
			},
			{
				Name:   "services",
				Usage:  "Print list of supported subtitle services",
				Action: cmdServices,
			},
		},
	}
	app.Run(os.Args)
}
