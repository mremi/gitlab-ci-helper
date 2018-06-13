// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mitchellh/cli"
	helper "github.com/rande/gitlab-ci-helper"
	"github.com/stretchr/testify/assert"
)

func Test_Ci_Dump_Revision(t *testing.T) {
	path := fmt.Sprintf("%s/REVISION", os.TempDir())

	reqs := []*helper.FakeRequest{}
	envs := map[string]string{
		"CI_BUILD_REF": "CI_BUILD_REF",
	}

	helper.WrapperTestCommand(reqs, envs, t, func(ts *httptest.Server) {
		ui := &cli.MockUi{}
		c := &CiDumpRevisionCommand{
			Ui: ui,
		}

		code := c.Run([]string{"-file", path})

		assert.Equal(t, 0, code)

		r, err := os.Open(path)

		assert.NoError(t, err)
		defer r.Close()

		content, err := ioutil.ReadFile(path)

		assert.NoError(t, err)
		assert.Equal(t, []byte("CI_BUILD_REF"), content)

	})

	os.Remove(path)
}

func Test_Ci_Revision_Help(t *testing.T) {
	c := &CiDumpRevisionCommand{
		Ui: &cli.MockUi{},
	}

	assert.True(t, len(c.Help()) > 0)
	assert.True(t, len(c.Synopsis()) > 0)
}

func Test_Ci_Dump_Revision_InvalidRun(t *testing.T) {
	c := &CiDumpRevisionCommand{
		Ui: &cli.MockUi{},
	}

	assert.Equal(t, 1, c.Run([]string{"--foobar"}))
}
