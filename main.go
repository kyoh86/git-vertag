package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/git-vertag/internal"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	// Set command name and description
	app := kingpin.New(
		"git-vertag",
		"A tool to manage version-tag with the semantic versioning specification.",
	).Author("kyoh86").Version(version)

	var cwd string
	var dryRun bool
	var fetch bool

	app.Flag("current-directory", "Run as if git was started in <path> instead of the current working directory.").Short('C').PlaceHolder("<path>").ExistingDirVar(&cwd)
	app.Flag("dry-run", "Without deleting tag, show git command.").BoolVar(&dryRun)
	app.Flag("fetch", "Fetch tags first").Default("true").BoolVar(&fetch)

	var message []string
	var file string
	var push bool

	getCmd := app.Command("get", "Gets the current version tag.").Default()

	deleteCmd := app.Command("delete", "Deletes a tag for the last version and prints it.")
	deleteCmd.Flag("push", "Delete tag from remote.").BoolVar(&push)

	majorCmd := app.Command("major", "Creates a tag for the next major version and prints it.")
	minorCmd := app.Command("minor", "Creates a tag for the next minor version and prints it.")
	patchCmd := app.Command("patch", "Creates a tag for the next patch version and prints it.")
	replaceCmd := app.Command("replace", "Replaces a tag for the last version and prints it.")

	for _, c := range []*kingpin.CmdClause{majorCmd, minorCmd, patchCmd, replaceCmd} {
		c.Flag("message", `Use the given tag message (instead of prompting). If multiple -m options are given, their values are concatenated as separate paragraphs.`).Short('m').StringsVar(&message)
		c.Flag("file", `Take the tag message from the given file. Use - to read the message from the standard input`).Short('F').StringVar(&file)
		c.Flag("push", `Push a new tag to remote`).BoolVar(&push)
	}

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	tag := internal.NewTagger()
	//UNDONE: support setting remote name
	//tag.Remote = "???"
	if dryRun {
		tag.Runner = internal.NewMockRunner()
	}
	tag.Workdir = cwd
	tag.Push = push

	mgr := internal.Manager{
		Tagger: tag,
	}

	v, err := mgr.GetVer(fetch)
	if err != nil {
		return
	}

	switch cmd {
	case getCmd.FullCommand():
		fmt.Println(v)
	case majorCmd.FullCommand():
		ver := v.Increment(internal.LevelMajor)
		if err := mgr.CreateVer(ver, message, file); err != nil {
			log.Fatal(err)
		}
		fmt.Println(ver)
	case minorCmd.FullCommand():
		ver := v.Increment(internal.LevelMinor)
		if err := mgr.CreateVer(ver, message, file); err != nil {
			log.Fatal(err)
		}
		fmt.Println(ver)
	case patchCmd.FullCommand():
		ver := v.Increment(internal.LevelPatch)
		if err := mgr.CreateVer(ver, message, file); err != nil {
			log.Fatal(err)
		}
		fmt.Println(ver)
	case replaceCmd.FullCommand():
		if err := mgr.ReplaceVer(v, message, file); err != nil {
			log.Fatal(err)
		}
	case deleteCmd.FullCommand():
		if err := mgr.DeleteVer(v); err != nil {
			log.Fatal(err)
		}
	}
}
