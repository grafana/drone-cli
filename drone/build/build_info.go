package build

import (
	"os"
	"strconv"
	"text/template"

	"github.com/grafana/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var buildInfoCmd = cli.Command{
	Name:      "info",
	Usage:     "show build details",
	ArgsUsage: "<repo/name> [build]",
	Action:    buildInfo,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplBuildInfo,
		},
	},
}

func buildInfo(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}
	buildArg := c.Args().Get(1)

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	var number int
	if buildArg == "last" || len(buildArg) == 0 {
		// Fetch the build number from the last build
		build, err := client.BuildLast(owner, name, "")
		if err != nil {
			return err
		}
		number = int(build.Number)
	} else {
		number, err = strconv.Atoi(buildArg)
		if err != nil {
			return err
		}
	}

	build, err := client.Build(owner, name, number)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, build)
}

// template for build information
var tmplBuildInfo = `Number: {{ .Number }}
Status: {{ .Status }}
Event: {{ .Event }}
Commit: {{ .After }}
Branch: {{ .Target }}
Ref: {{ .Ref }}
Author: {{ .Author }} {{ if .AuthorEmail }}<{{.AuthorEmail}}>{{ end }}
Message: {{ .Message }}
`
