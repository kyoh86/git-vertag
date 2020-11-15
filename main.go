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
	var prefix string
	var ancestors bool
	app.Flag("current-directory", "Run as if git was started in <path> instead of the current working directory.").Short('C').PlaceHolder("<path>").ExistingDirVar(&cwd)
	app.Flag("dry-run", "Without deleting tag, show git command.").Envar("GIT_VERTAG_DRYRUN").BoolVar(&dryRun)
	app.Flag("fetch", "Fetch tags first").Envar("GIT_VERTAG_FETCH").Default("true").BoolVar(&fetch)
	app.Flag("prefix", "Prefix for tag").Envar("GIT_VERTAG_PREFIX").Default("v").StringVar(&prefix)
	app.Flag("ancestors", "With ancestor versions (vN and vN.N)").Envar("GIT_VERTAG_ANCESTORS").BoolVar(&ancestors)

	getCmd := app.Command("get", "Gets the current version tag.").Default()
	majorCmd := app.Command("major", "Creates a tag for the next major version and prints it.")
	minorCmd := app.Command("minor", "Creates a tag for the next minor version and prints it.")
	patchCmd := app.Command("patch", "Creates a tag for the next patch version and prints it.")
	releaseCmd := app.Command("release", "Creates a tag to remove pre-release meta information.")
	preCmd := app.Command("pre", "Creates a tag for the next pre-release version and prints it.")
	buildCmd := app.Command("build", "Creates a tag for the next build version and prints it.")

	var message []string
	var file string
	var pushTo string

	for _, c := range []*kingpin.CmdClause{majorCmd, minorCmd, patchCmd, preCmd, buildCmd, releaseCmd} {
		c.Flag("message", "Use the given tag message (instead of prompting). If multiple -m options are given, their values are concatenated as separate paragraphs.").Short('m').StringsVar(&message)
		c.Flag("file", "Take the tag message from the given file. Use - to read the message from the standard input").Short('F').StringVar(&file)
		c.Flag("push-to", "The remote repository that is destination of a push operation. This parameter can be either a URL or the name of a remote.").PlaceHolder("REPOSITORY").StringVar(&pushTo)
	}

	var pre internal.PreReleaseFlag
	for _, c := range []*kingpin.CmdClause{majorCmd, minorCmd, patchCmd} {
		c.Flag("pre", "Update pre-release notation. It accepts only alphanumeric or numeric identities.").SetValue(&pre)
	}
	preCmd.Arg("pre", "Pre-release notation. It accepts only alphanumeric or numeric identities.").Required().SetValue(&pre)

	var build internal.BuildFlag
	for _, c := range []*kingpin.CmdClause{majorCmd, minorCmd, patchCmd, preCmd, releaseCmd} {
		c.Flag("build", "Update build notation. It accepts only alphanumeric or numeric identities.").SetValue(&build)
	}
	buildCmd.Arg("build", "Update build notation. It accepts only alphanumeric or numeric identities.").Required().SetValue(&build)

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		app.FatalUsage("%s", err)
	}

	tag := internal.NewTagger()
	if dryRun {
		tag.Runner = internal.NewMockRunner()
	}

	tag.Workdir = cwd
	tag.PushTo = pushTo

	mgr := internal.Manager{
		Prefix:    prefix,
		Tagger:    tag,
		Fetch:     fetch,
		Ancestors: ancestors,
	}

	switch cmd {
	case getCmd.FullCommand():
		v, err := mgr.GetVer()
		if err != nil {
			return
		}
		fmt.Println(v)

	case majorCmd.FullCommand():
		printResult(mgr.UpdateMajor(pre, build, message, file))

	case minorCmd.FullCommand():
		printResult(mgr.UpdateMinor(pre, build, message, file))

	case patchCmd.FullCommand():
		printResult(mgr.UpdatePatch(pre, build, message, file))

	case preCmd.FullCommand():
		printResult(mgr.UpdatePre(pre, build, message, file))

	case releaseCmd.FullCommand():
		printResult(mgr.Release(build, message, file))

	case buildCmd.FullCommand():
		printResult(mgr.Build(build, message, file))
	}
}

func printResult(cur, next string, err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update %s to %s\n", cur, next)
}
