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
	Command              string
	SubCommand           string
	OptLogLevel          string            `long:"loglevel" short:"l" arg:"(debug|info|warn|error)" description:"specify log-level"`
	OptHelp              bool              `long:"help" short:"h" description:"show this help message and exit"`
	OptVersion           bool              `long:"version" short:"v" description:"prints the version number"`
	Kind                 string            `long:"kind" arg:"(wordpress|php|rails|node)" description:"kind for project"`
	Payload              map[string]string `long:"payload" description:"payload for resource"`
	Username             string            `long:"username" description:""`
	Password             string            `long:"password" description:""`
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
Usage: lolp [<option>] <command> [<args|attributes>]

Commands:
  login --username <id> --password <pw>
  project create --kind <php|rails|node> --database pw:<pw>
  project create --kind wordpress --payload username:<wp-user> --payload password:<wp-pw> --payload email:<wp-email>
  project list
  project delete <name>

Options:
`)

	t := reflect.TypeOf(CLI{})
	names := []string{
		"OptLogLevel",
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

		var o, a string
		if a = tag.Get("arg"); a != "" {
			a = fmt.Sprintf("=%s", a)
		}
		if s := tag.Get("short"); s != "" {
			o = fmt.Sprintf("-%s, --%s%s", tag.Get("short"), tag.Get("long"), a)
		} else {
			o = fmt.Sprintf("--%s%s", tag.Get("long"), a)
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
		fmt.Fprintf(c.outStream, "  %-40s %s\n", o, desc)
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

// login logins to lolipop
func (c *CLI) login() error {
	token, err := c.client.Login(c.Username, c.Password)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "export %s=%s", lolp.TokenEnvVar, token)
	return nil
}

// createProject creates project
func (c *CLI) createProject() error {
	attrs := make(map[string]interface{})
	attrs["Kind"] = c.Kind
	if len(c.Payload) > 0 {
		payload := make(map[string]interface{})
		for k, v := range c.Payload {
			payload[k] = v
		}
		attrs["Payload"] = payload
	}
	p, err := c.client.CreateProject(attrs)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "%#v\n", p)
	return nil
}

// projects gets project list
func (c *CLI) projects() error {
	projects, err := c.client.Projects(make(map[string]interface{}))
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "%-38s  %-36s %s\n", "ID", "Name", "Kind")
	for _, v := range *projects {
		fmt.Fprintf(c.outStream, "%-38s  %-36s %s\n", v.ID, v.Name, v.Kind)
	}
	return nil
}

// deleteProject deletes project
func (c *CLI) deleteProject() error {
	err := c.client.DeleteProject(c.Args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "delete successfuly\n")
	return nil
}
