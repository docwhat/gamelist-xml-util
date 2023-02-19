package romhash_test

import (
	"crypto/sha1" //nolint:gosec // We are using this for identification of ROMs, not for security.
	"path/filepath"
	"testing"

	"docwhat.org/gamelist-xml-util/pkg/romhash"
	testdata "docwhat.org/gamelist-xml-util/pkg/testroms"
)

func TestSHA1(t *testing.T) {
	t.Parallel()

	data, err := testdata.New()
	if err != nil {
		t.Fatal(err)
	}
	defer data.Close()

	for _, file := range data.Files {
		if e := filepath.Ext(file.Path); !romhash.KnownExt(e) {
			t.Errorf("KnownExt(%q) => false; want true", e)
		}

		buf := make([]byte, 4*1024*1024)

		//nolint:gosec // We aren't interested in the strength of the hash; we are conforming to an existing standard.
		if got, err := romhash.Hash(file.Path, sha1.New(), buf); err != nil {
			t.Errorf("Hash(%q, sha1.New()) => err = %v; want nil", file.Path, err)
		} else if got != file.SHA1 {
			t.Errorf("Hash(%q, sha1.New()) => %q; want %q", file.Path, got, file.SHA1)
		}
	}
}
