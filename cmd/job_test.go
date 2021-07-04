package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitFolders(t *testing.T) {
	f := splitFolders("a/b/c")
	if diff := cmp.Diff(f, []string{"a", "b", "c"}); diff != "" {
		t.Errorf("Split didn't work as expected\n%s", diff)
	}
}
