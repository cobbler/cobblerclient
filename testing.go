// SPDX-LLicense-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright SUSE LLC

package cobblerclient

import (
	"os"
	"testing"
)

// FailOnError ...
func FailOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// Fixture implies that from the context where the test is being run a "fixtures" folder exists.
func Fixture(fn string) ([]byte, error) {
	// Disable semgrep (linter in Codacy) since this is testcode
	return os.ReadFile("./fixtures/" + fn) // nosemgrep
}
