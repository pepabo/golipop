package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/hashicorp/logutils"
	flags "github.com/jessevdk/go-flags"
	"github.com/pepabo/golipop"
)

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.run())
}

// CLI struct
type CLI struct {
	outStream, errStream io.Writer
	client               *lolp.Client
	Args                 []string
	ParsedArgs           map[string]interface{}
	Command              string
	SubCommand           string
	OptLogLevel          string `long:"loglevel" short:"l" arg:"debug|info|warn|error" description:"specify log-level"`
	OptHelp              bool   `long:"help" short:"h" description:"show this help message and exit"`
	OptVersion           bool   `long:"version" short:"v" description:"prints the version number"`
}

const (
	// ExitOK for exit code
	ExitOK int = 0

	// ExitErr for exit code
	ExitErr int = 1

	ArgSplitKey string = "="
)

// CLI executes for cli
func (c *CLI) run() int {
	p := flags.NewParser(c, flags.PrintErrors|flags.PassDoubleDash)
	args, err := p.Parse()
	if err != nil || c.OptHelp {
		c.showHelp()
		return ExitErr
	}

	if c.OptVersion {
		fmt.Fprintf(c.errStream, "%s\n", lolp.Version)
		return ExitOK
	}

	if len(args) == 0 {
		fmt.Fprintf(c.errStream, "command not specified\n")
		return ExitErr
	}

	c.Command = args[0]
	if len(args) > 1 {
		c.SubCommand = args[1]
	}
	if len(args) > 2 {
		c.Args = args[2:]
	}

	if c.OptLogLevel != "" {
		c.OptLogLevel = strings.ToUpper(c.OptLogLevel)
	} else {
		c.OptLogLevel = "ERROR"
	}

	c.ParsedArgs = make(map[string]interface{})
	fmt.Printf("%#v\n", c.Args)
	for _, arg := range c.Args {
		parsed := strings.Split(arg, ArgSplitKey)
		fmt.Printf("%#v\n", parsed)
		if parsed[0] != "" && parsed[1] != "" {
			c.ParsedArgs[strings.Title(parsed[0])] = parsed[1]
		}
	}
	fmt.Printf("%#v\n", c.ParsedArgs)

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(c.OptLogLevel),
		Writer:   c.errStream,
	}
	log.SetOutput(filter)

	if err := c.callAPI(); err != nil {
		fmt.Fprintf(c.errStream, "%s\n", err)
		return ExitErr
	}

	return ExitOK
}

// help shows help
func (c *CLI) showHelp() {
	fmt.Fprintf(c.outStream, `
Usage: lolp [options] [CMD]

Commands:
  login username=your@example.com password=******
  project create kind=<wordpress|php|rails|node> username=wpuser password=wp*** email=wp@example.com
  project list
  project delete id=<ID>

Options:
`)

	t := reflect.TypeOf(CLI{})
	names := []string{
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
		fmt.Fprintf(c.outStream, "  %-21s %s\n", o, desc)
	}
}

// callAPI calls API for cli
func (c *CLI) callAPI() error {
	c.client = lolp.DefaultClient()
	var err error

	switch c.Command {
	case "login":
		err = c.login()
	case "project":
		switch c.SubCommand {
		case "create":
			err = c.createProject()
		case "list":
			err = c.projects()
		case "delete":
			err = c.deleteProject()
		default:
			err = errors.New("unknown sub command")
		}
	default:
		err = errors.New("unknown command")
	}

	if err != nil {
		return err
	}

	return nil
}

type Login struct {
	username string
	password string
}

// login logins to lolipop
func (c *CLI) login() error {
	l := new(Login)
	for f, v := range c.ParsedArgs {
		lolp.SetField(l, f, v)
	}
	token, err := c.client.Login(l.username, l.password)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "export %s=%s", lolp.TokenEnvVar, token)
	return nil
}

// createProject creates project
func (c *CLI) createProject() error {
	p, err := c.client.CreateProject(c.ParsedArgs)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "%#v\n", p)
	return nil
}

// projects gets project list
func (c *CLI) projects() error {
	projects, err := c.client.Projects(map[string]interface{}{})
	if err != nil {
		return err
	}
	for _, v := range *projects {
		fmt.Fprintf(c.outStream, "%#v\n", v)
	}
	return nil
}

// deleteProject deletes project
func (c *CLI) deleteProject() error {
	p := new(lolp.Project)
	for f, v := range c.ParsedArgs {
		lolp.SetField(p, f, v)
	}
	err := c.client.DeleteProject(p.ID)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "delete successfuly\n")
	return nil
}
