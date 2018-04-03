package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/pepabo/golipop"
)

type options struct {
	OptArgs     []string
	OptCommand  string
	OptID       string `long:"id" short:"i" description:"resource id"`
	OptUsername string `long:"username" short:"u" description:"username for authentication"`
	OptPassword string `long:"password" short:"p" description:"password for authentication"`
	OptEmail    string `long:"email" short:"e" description:"password for authentication"`
	OptTemplate string `long:"template" short:"t" arg:"(wordpress|php|rails|node)" description:"project template"`
	OptHelp     bool   `long:"help" short:"h" description:"show this help message and exit"`
	OptVersion  bool   `long:"version" short:"v" description:"prints the version number"`
}

func showHelp() {
	os.Stderr.WriteString(`
Usage: lolp [options] [CMD]

Commands:
  login -u your@example.com -p
  project create -t <wordpress|php|rails|node>
  project list
  project delete -i <id>

Options:
`)
	t := reflect.TypeOf(options{})
	names := []string{
		"OptID",
		"OptUsername",
		"OptPassword",
		"OptEmail",
		"OptTemplate",
		"OptHelp",
		"OptVersion",
	}

	for _, name := range names {
		f, ok := t.FieldByName(name)
		if !ok {
			continue
		}

		tag := f.Tag
		if tag == "" {
			continue
		}
		var o string
		if s := tag.Get("short"); s != "" {
			o = fmt.Sprintf("-%s, --%s", tag.Get("short"), tag.Get("long"))
		} else {
			o = fmt.Sprintf("--%s", tag.Get("long"))
		}

		desc := tag.Get("description")
		if i := strings.Index(desc, "\n"); i >= 0 {
			var buf bytes.Buffer
			buf.WriteString(desc[:i+1])
			desc = desc[i+1:]
			const indent = "                        "
			for {
				if i = strings.Index(desc, "\n"); i >= 0 {
					buf.WriteString(indent)
					buf.WriteString(desc[:i+1])
					desc = desc[i+1:]
					continue
				}
				break
			}
			if len(desc) > 0 {
				buf.WriteString(indent)
				buf.WriteString(desc)
			}
			desc = buf.String()
		}
		fmt.Fprintf(os.Stderr, "  %-21s %s\n", o, desc)
	}
}

func main() {
	os.Exit(NewCLI())
}

func NewCLI() int {
	opts := &options{}
	p := flags.NewParser(opts, flags.PrintErrors|flags.PassDoubleDash)
	args, err := p.Parse()
	if err != nil || opts.OptHelp {
		showHelp()
		return 1
	}

	if opts.OptVersion {
		fmt.Printf("%s\n", lolp.Version)
		return 0
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "command not specified\n")
		return 1
	}

	opts.OptCommand = args[0]
	if len(args) > 1 {
		opts.OptArgs = args[1:]
	}

	Run(opts)
	return 0
}

func Run(o *options) {
	c := lolp.DefaultClient()

	if o.OptCommand == "login" {
		token, err := c.Login(o.OptUsername, o.OptPassword)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "export %s=%s", lolp.TokenEnvVar, token)
		return
	}

	if o.OptCommand == "project" {
		payload := map[string]interface{}{
			"username": o.OptUsername,
			"password": o.OptPassword,
			"email":    o.OptEmail,
		}

		p, err := c.CreateProject(lolp.GetProjectTemplate(o.OptTemplate), payload)
		if err == nil {
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "%#v\n", p)
		return
	}
}
