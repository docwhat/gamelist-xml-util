package romhash

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/kjk/lzmadec"
)

// memoizedHas7z returns a memoized haz7z function.
func memoizedHas7z() func() bool {
	hasRun := false
	has7z := false

	return func() bool {
		if !hasRun {
			if runtime.GOOS == "windows" {
				os.Setenv("PATH", fmt.Sprintf("C:\\Program Files\\7-zip;%s", os.Getenv("PATH")))
			}

			_, err := exec.LookPath("7z")

			has7z = err == nil
			hasRun = true
		}

		return has7z
	}
}

// has7z returns true if 7z is installed.
var has7z = memoizedHas7z() //nolint:gochecknoglobals // This is a memoized function.

func decode7Zip(f string) (io.ReadCloser, error) {
	reader, err := lzmadec.NewArchive(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read 7zip file %q: %w", f, err)
	}

	for _, entry := range reader.Entries {
		ext := filepath.Ext(entry.Path)
		if decoder, ok := getDecoder(ext); ok {
			rf, err := reader.GetFileReader(entry.Path)
			if err != nil {
				continue
			}

			rom, err := decoder(rf, entry.Size)
			if err != nil {
				continue
			}

			return rom, nil
		}
	}

	return nil, ErrNoRoms
}
