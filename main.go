package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var cwd string
var version = "snapshot"

func main() {
	// Set command name and description
	app := kingpin.New(
		"git-vertag",
		"A tool to manage version-tag with the semantic versioning specification.",
	).Author("kyoh86").Version(version)

	app.Flag("current-directory", "Run as if git was started in <path> instead of the current working directory.").Short('C').PlaceHolder("<path>").ExistingDirVar(&cwd)

	var dryRun bool
	getCmd := app.Command("get", "Gets the current version tag.").Default()
	majorCmd := app.Command("major", "Creates a tag for the next major version and prints it.")
	minorCmd := app.Command("minor", "Creates a tag for the next minor version and prints it.")
	patchCmd := app.Command("patch", "Creates a tag for the next patch version and prints it.")

	for _, c := range []*kingpin.CmdClause{majorCmd, minorCmd, patchCmd} {
		c.Flag("dry-run", "Without creating a new tag, show git command.").BoolVar(&dryRun)
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case getCmd.FullCommand():
		latest, err := latestVer()
		if err != nil {
			panic(err)
		}
		fmt.Println(latest)
	case majorCmd.FullCommand():
		latest, err := latestVer()
		if err != nil {
			panic(err)
		}
		latest.Major++
		latest.Minor = 0
		latest.Patch = 0
		if err := createTag(latest, dryRun); err != nil {
			panic(err)
		}
	case minorCmd.FullCommand():
		latest, err := latestVer()
		if err != nil {
			panic(err)
		}
		latest.Minor++
		latest.Patch = 0
		if err := createTag(latest, dryRun); err != nil {
			panic(err)
		}
	case patchCmd.FullCommand():
		latest, err := latestVer()
		if err != nil {
			panic(err)
		}
		latest.Patch++
		if err := createTag(latest, dryRun); err != nil {
			panic(err)
		}
	}
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

func createTag(v *Semver, dryRun bool) error {
	git, _ := gitCmd("tag", v.String())
	if dryRun {
		w := csv.NewWriter(os.Stdout)
		w.Comma = ' '
		w.Write(git.Args)
		w.Flush()
		return nil
	}

	fmt.Println(v)
	if err := git.Run(); err != nil {
		return err
	}

	return nil
}

func latestVer() (*Semver, error) {
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
