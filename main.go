package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/alecthomas/kingpin"
)

var cwd string

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

	app.Flag("current-directory", "Run as if git was started in <path> instead of the current working directory.").Short('C').PlaceHolder("<path>").ExistingDirVar(&cwd)

	var dryRun bool
	var fetch bool
	var message []string
	var file string

	getCmd := app.Command("get", "Gets the current version tag.").Default()
	getCmd.Flag("fetch", "Fetch tags first").Default("true").BoolVar(&fetch)

	deleteCmd := app.Command("delete", "Deletes a tag for the last version and prints it.")
	deleteCmd.Flag("dry-run", "Without creating a new tag, show git command.").BoolVar(&dryRun)
	deleteCmd.Flag("fetch", "Fetch tags first").Default("true").BoolVar(&fetch)

	majorCmd := app.Command("major", "Creates a tag for the next major version and prints it.")
	minorCmd := app.Command("minor", "Creates a tag for the next minor version and prints it.")
	patchCmd := app.Command("patch", "Creates a tag for the next patch version and prints it.")
	replaceCmd := app.Command("replace", "Replaces a tag for the last version and prints it.")

	for _, c := range []*kingpin.CmdClause{majorCmd, minorCmd, patchCmd, replaceCmd} {
		c.Flag("dry-run", "Without creating a new tag, show git command.").BoolVar(&dryRun)
		c.Flag("fetch", "Fetch tags first").Default("true").BoolVar(&fetch)
		c.Flag("message", `Use the given tag message (instead of prompting). If multiple -m options are given, their values are concatenated as separate paragraphs.`).Short('m').StringsVar(&message)
		c.Flag("file", `Take the tag message from the given file. Use - to read the message from the standard input`).Short('F').StringVar(&file)
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case getCmd.FullCommand():
		if err := getVersion(fetch); err != nil {
			log.Fatal(err)
		}
	case majorCmd.FullCommand():
		if err := incrementMajor(dryRun, fetch, message, file); err != nil {
			log.Fatal(err)
		}
	case minorCmd.FullCommand():
		if err := incrementMinor(dryRun, fetch, message, file); err != nil {
			log.Fatal(err)
		}
	case patchCmd.FullCommand():
		if err := incrementPatch(dryRun, fetch, message, file); err != nil {
			log.Fatal(err)
		}
	case replaceCmd.FullCommand():
		if err := replaceTag(dryRun, fetch, message, file); err != nil {
			log.Fatal(err)
		}
	case deleteCmd.FullCommand():
		if err := deleteTag(dryRun, fetch); err != nil {
			log.Fatal(err)
		}
	}
}

func getVersion(fetch bool) error {
	latest, err := latestVer(fetch)
	if err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

func deleteTag(dryRun bool, fetch bool) error {
	latest, err := latestVer(fetch)
	if err != nil {
		return err
	}
	if err := removeTag(latest, dryRun); err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

func replaceTag(dryRun bool, fetch bool, message []string, file string) error {
	latest, err := latestVer(fetch)
	if err != nil {
		return err
	}
	if err := removeTag(latest, dryRun); err != nil {
		return err
	}
	if err := createTag(latest, dryRun, message, file); err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

func incrementPatch(dryRun bool, fetch bool, message []string, file string) error {
	latest, err := latestVer(fetch)
	if err != nil {
		return err
	}
	latest.Patch++
	if err := createTag(latest, dryRun, message, file); err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

func incrementMinor(dryRun bool, fetch bool, message []string, file string) error {
	latest, err := latestVer(fetch)
	if err != nil {
		return err
	}
	latest.Minor++
	latest.Patch = 0
	if err := createTag(latest, dryRun, message, file); err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

func incrementMajor(dryRun bool, fetch bool, message []string, file string) error {
	latest, err := latestVer(fetch)
	if err != nil {
		return err
	}
	latest.Major++
	latest.Minor = 0
	latest.Patch = 0
	if err := createTag(latest, dryRun, message, file); err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

func gitCmd(args ...string) (*exec.Cmd, *bytes.Buffer) {
	var stdout bytes.Buffer
	var cmd *exec.Cmd
	if cwd != "" {
		cmd = exec.Command("git", append([]string{"-C", cwd}, args...)...)
	} else {
		cmd = exec.Command("git", args...)
	}
	cmd.Stdout = &stdout
	return cmd, &stdout
}

// Semver :
type Semver struct {
	Major int
	Minor int
	Patch int
}

func (s *Semver) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.Major, s.Minor, s.Patch)
}

func createTag(v *Semver, dryRun bool, message []string, file string) error {
	args := []string{"tag"}
	for _, m := range message {
		args = append(args, "--message", m)
	}
	if file != "" {
		args = append(args, "--file", file)
	}
	git, _ := gitCmd(append(args, v.String())...)
	if dryRun {
		w := csv.NewWriter(os.Stdout)
		w.Comma = ' '
		if err := w.Write(git.Args); err != nil {
			return err
		}
		w.Flush()
		return nil
	}

	if err := git.Run(); err != nil {
		return err
	}

	return nil
}

func removeTag(v *Semver, dryRun bool) error {
	git, _ := gitCmd("tag", "-d", v.String())
	if dryRun {
		w := csv.NewWriter(os.Stdout)
		w.Comma = ' '
		if err := w.Write(git.Args); err != nil {
			return err
		}
		w.Flush()
		return nil
	}

	if err := git.Run(); err != nil {
		return err
	}

	return nil
}

func latestVer(fetch bool) (*Semver, error) {
	if fetch {
		git, _ := gitCmd("fetch", "--tags")
		if err := git.Run(); err != nil {
			return nil, err
		}
	}
	git, stdout := gitCmd("tag", "-l")
	if err := git.Run(); err != nil {
		// var status = 1
		// if exit, ok := err.(*exec.ExitError); ok {
		// 	fmt.Fprint(os.Stderr, string(exit.Stderr))
		// 	if s, ok := exit.Sys().(syscall.WaitStatus); ok {
		// 		status = s.ExitStatus()
		// 	}
		// }
		// os.Exit(status)
		// return
		return nil, err
	}

	latest := &Semver{}
	stream := bufio.NewScanner(stdout)
	for stream.Scan() {
		ver, err := parseVer(stream.Text())
		if err != nil {
			continue
		}
		latest = greaterVer(latest, ver)
	}

	return latest, nil
}

func greaterVer(v1, v2 *Semver) *Semver {
	if v1 == nil {
		return v2
	}

	if v1.Major < v2.Major {
		return v2
	}
	if v1.Major > v2.Major {
		return v1
	}
	if v1.Minor < v2.Minor {
		return v2
	}
	if v1.Minor > v2.Minor {
		return v1
	}
	if v1.Patch < v2.Patch {
		return v2
	}
	if v1.Patch > v2.Patch {
		return v1
	}
	return v1
}

var semverRegex = regexp.MustCompile(`^v?(?P<major>\d+)(\.(?P<minor>\d+))?(\.(?P<patch>\d+))?(?:-.*)?$`)

func parseVer(s string) (*Semver, error) {
	match := semverRegex.FindStringSubmatch(s)
	if len(match) == 0 {
		return nil, errors.New("invalid version syntax")
	}
	result := map[string]int{}
	for i, name := range semverRegex.SubexpNames() {
		if i == 0 {
			continue
		}
		if i < len(match) {
			result[name], _ = strconv.Atoi(match[i])
		}
	}
	return &Semver{
		Major: result["major"],
		Minor: result["minor"],
		Patch: result["patch"],
	}, nil
}
