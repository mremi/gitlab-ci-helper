// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
)

type MetaBuild struct {
	Id      string `json:"id"`
	Ref     string `json:"ref"`
	RefName string `json:"ref_name"`
	Tag     string `json:"tag"`
	Stage   string `json:"stage"`
	JobName string `json:"job_name"`
}

type MetaProject struct {
	Id  string `json:"id"`
	Dir string `json:"dir"`
}

type MetaServer struct {
	Name     string `json:"name"`
	Revision string `json:"revision"`
	Version  string `json:"version"`
}

type Meta struct {
	Build   *MetaBuild   `json:"build"`
	Project *MetaProject `json:"project"`
	Server  *MetaServer  `json:"server"`
}

type CiDumpMetaCommand struct {
	Ui       cli.Ui
	Verbose  bool
	MetaFile string
}

func (c *CiDumpMetaCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("ci:meta", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.StringVar(&c.MetaFile, "file", "ci.json", "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	meta := &Meta{
		Build: &MetaBuild{
			Id:      helper.GetEnv("CI_JOB_ID", os.Getenv("CI_BUILD_ID")),
			Ref:     helper.GetEnv("CI_COMMIT_SHA", os.Getenv("CI_BUILD_REF")),
			RefName: helper.GetEnv("CI_COMMIT_REF_NAME", os.Getenv("CI_BUILD_REF_NAME")),
			Tag:     helper.GetEnv("CI_COMMIT_TAG", os.Getenv("CI_BUILD_TAG")),
			Stage:   helper.GetEnv("CI_JOB_STAGE", os.Getenv("CI_BUILD_STAGE")),
			JobName: helper.GetEnv("CI_JOB_NAME", os.Getenv("CI_BUILD_NAME")),
		},
		Project: &MetaProject{
			Id:  os.Getenv("CI_PROJECT_ID"),
			Dir: os.Getenv("CI_PROJECT_DIR"),
		},
		Server: &MetaServer{
			Name:     os.Getenv("CI_SERVER_NAME"),
			Revision: os.Getenv("CI_SERVER_REVISION"),
			Version:  os.Getenv("CI_SERVER_VERSION"),
		},
	}

	fp, _ := os.Create(c.MetaFile)
	defer fp.Close()

	b, _ := json.Marshal(meta)

	var out bytes.Buffer
	json.Indent(&out, b, "", "    ")

	out.WriteTo(fp)

	return 0
}

func (c *CiDumpMetaCommand) Synopsis() string {
	return "Dump a json file with build information."
}

func (c *CiDumpMetaCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper ci:meta [options]

  Dump meta information about ci into a ci.json file

Options:

  -file               Target file (default: ci.json)
  -verbose            Add verbose information to the output
`
	return strings.TrimSpace(helpText)
}
