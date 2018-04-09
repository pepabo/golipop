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
	os.Exit(cli.run(os.Args[1:]))
}

// CLI struct
type CLI struct {
	outStream, errStream io.Writer
	client               *lolp.Client

	Args          []string
	Command       string
	SubCommand    string
	Kind          string            `long:"kind" short:"k" arg:"(wordpress|php|rails|node)" description:"kind of project"`
	SubDomain     string            `long:"sub-domain" short:"s" description:"your sub-domain"`
	CustomDomains []string          `long:"custom-domain" short:"c" description:"your custom domain"`
	Payload       map[string]string `long:"payload" short:"a" description:"payload for project"`
	Database      map[string]string `long:"database" short:"d" description:"database for project"`
	Username      string            `long:"username" short:"u" description:"username for login"`
	Password      string            `long:"password" short:"p" description:"password for login"`

	OptLogLevel string `long:"loglevel" short:"l" arg:"(debug|info|warn|error)" description:"specify log-level"`
	OptHelp     bool   `long:"help" short:"h" description:"show this help message and exit"`
	OptVersion  bool   `long:"version" short:"v" description:"prints the version number"`
}

const (
	// ExitOK for exit code
	ExitOK int = 0

	// ExitErr for exit code
	ExitErr int = 1
)

// CLI executes for cli
func (c *CLI) run(a []string) int {
	p := flags.NewParser(c, flags.PrintErrors|flags.PassDoubleDash)
	args, err := p.ParseArgs(a)
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
	attrs := strings.Join(c.buildHelp([]string{
		"Username",
		"Password",
		"Payload",
		"Database",
		"CustomDomains",
		"SubDomain",
	}), "\n")

	opts := strings.Join(c.buildHelp([]string{
		"OptLogLevel",
		"OptHelp",
		"OptVersion",
	}), "\n")

	help := `
Usage: lolp [<option>] <command> [<args|attributes>]

Commands: login, project

Attributes:
%s
Options:
%s

Examples:
  login -u <your-email> -p <your-password>
  project create -k <php|rails|node> -d password:<password>
  project create -k wordpress -a username:<wp-user> -a password:<wp-pw> -a email:<wp-email>
  project list
  project delete <project-sub-domain>
`
	fmt.Fprintf(c.outStream, help, attrs, opts)
}

func (c *CLI) buildHelp(names []string) []string {
	var help []string
	t := reflect.TypeOf(CLI{})

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
		help = append(help, fmt.Sprintf("  %-40s %s", o, desc))
	}

	return help
}

// callAPI calls API for cli
func (c *CLI) callAPI() error {
	c.client = lolp.New()
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
			err = c.project()
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
	token, err := c.client.Authenticate(c.Username, c.Password)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "export %s=%s", lolp.TokenEnvVar, token)
	return nil
}

// createProject creates project
func (c *CLI) createProject() error {
	n := new(lolp.ProjectNew)
	n.Kind = c.Kind
	if len(c.SubDomain) > 0 {
		n.SubDomain = c.SubDomain
	}
	if len(c.CustomDomains) > 0 {
		n.CustomDomains = c.CustomDomains
	}
	if len(c.Payload) > 0 {
		payload := make(map[string]interface{})
		for k, v := range c.Payload {
			payload[k] = v
		}
		n.Payload = payload
	}
	if len(c.Database) > 0 {
		n.Database = map[string]interface{}{"password": c.Database["password"]}
	}
	p, err := c.client.CreateProject(n)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.outStream, "%s\n", p.Domain)
	return nil
}

// projects gets project list
func (c *CLI) projects() error {
	projects, err := c.client.Projects()
	if err != nil {
		return err
	}
	fmt.Fprintf(c.outStream, "%-38s  %-36s %s\n", "ID", "Name", "Kind")
	for _, v := range *projects {
		fmt.Fprintf(c.outStream, "%-38s  %-36s %s\n", v.ID, v.Name, v.Kind)
	}
	return nil
}

// project get a project
func (c *CLI) project() error {
	p, err := c.client.Project(c.SubCommand)
	if err != nil {
		return err
	}
	c.showStruct(p)
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

// showStruct shows a struct
func (c *CLI) showStruct(s interface{}) {
	ss := reflect.ValueOf(s).Elem()
	typeOfT := ss.Type()
	for i := 0; i < ss.NumField(); i++ {
		f := ss.Field(i)
		fmt.Fprintf(c.outStream, "%-20s %#v\n", typeOfT.Field(i).Name, f.Interface())
	}
}
