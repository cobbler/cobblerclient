// SPDX-LLicense-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright SUSE LLC

package cobblerclient

import "testing"

func FailOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
