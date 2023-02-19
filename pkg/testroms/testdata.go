//nolint:wrapcheck,gomnd,funlen
package testroms

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	binHash  = "5c3eb80066420002bc3dcc7ca4ab6efad7ed4ae5"
	fileMode = 0o600
)

type Data struct {
	Dir       string `exhaustruct:"optional"`
	Files     []File `exhaustruct:"optional"`
	DoCleanup bool   `exhaustruct:"optional"`
}

func (d *Data) Close() error {
	if d.DoCleanup {
		return os.RemoveAll(d.Dir)
	}

	return nil
}

// AddFile adds a file to the Files array.
func (d *Data) AddFile(path string, sha1 string) {
	d.Files = append(d.Files, File{Path: path, SHA1: sha1})
}

// NewBinFiles makes all the rom types with a Bin extension.
func NewBinFiles(data *Data) error {
	binFile := make([]byte, 512)
	binExts := []string{
		".bin", ".a26", ".rom", ".cue", ".gdi", ".gb", ".gba",
		".gbc", ".lyx", ".32x", ".gg", ".pce", ".sms", ".sg",
		".col", ".int", ".ngp", ".ngc", ".vb", ".vec", ".gam",
		".a78", ".j64", ".jag", ".lnx", ".mgw", ".nds", ".fds",
	}

	for _, e := range binExts {
		path := filepath.Join(data.Dir, fmt.Sprintf("test%s", e))

		if err := os.WriteFile(path, binFile, fileMode); err != nil {
			return err
		}

		data.Files = append(data.Files, File{Path: path, SHA1: binHash})
	}

	return nil
}

// NewLynxFiles makes rom files for the Atari Lynx.
func NewLynxFiles(data *Data) error {
	lnxFile := make([]byte, 512+64)
	copy(lnxFile, []byte("LYNX"))

	lnxPath := filepath.Join(data.Dir, "test.lnx")
	lyxPath := filepath.Join(data.Dir, "test.lyx")

	if err := os.WriteFile(lnxPath, lnxFile, fileMode); err != nil {
		return err
	}

	if err := os.WriteFile(lyxPath, lnxFile, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: lnxPath, SHA1: binHash})
	data.Files = append(data.Files, File{Path: lyxPath, SHA1: binHash})

	return nil
}

// NewAtari7800Files makes rom files for the Atari 7800.
func NewAtari7800Files(data *Data) error {
	a78File := make([]byte, 512+128)
	copy(a78File, []byte(" ATARI7800"))

	a78Path := filepath.Join(data.Dir, "a7800.a78")

	if err := os.WriteFile(a78Path, a78File, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: a78Path, SHA1: binHash})

	return nil
}

// NewN64Files makes rom files for the Nintendo 64.
func NewN64Files(data *Data) error {
	v64File := make([]byte, 1024)
	n64File := make([]byte, 1024)
	z64File := make([]byte, 1024)

	copy(v64File, []byte{0, 0x80, 0, 0})
	copy(n64File, []byte{0, 0, 0, 0x80})
	copy(z64File, []byte{0x80, 0, 0, 0})

	for i := 4; i < 1024; i += 4 {
		v64File[i], z64File[i+1], n64File[i+2] = 1, 1, 1
		v64File[i+1], z64File[i], n64File[i+3] = 2, 2, 2
		v64File[i+2], z64File[i+3], n64File[i] = 3, 3, 3
		v64File[i+3], z64File[i+2], n64File[i+1] = 4, 4, 4
	}

	v64Path := filepath.Join(data.Dir, "test-v64.v64")
	n64Path := filepath.Join(data.Dir, "test-n64.v64")
	z64Path := filepath.Join(data.Dir, "test-z64.v64")

	if err := os.WriteFile(v64Path, v64File, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: v64Path, SHA1: "00ba552537f953776b37a05230e9f1c2f6d4c145"})

	if err := os.WriteFile(n64Path, n64File, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: n64Path, SHA1: "00ba552537f953776b37a05230e9f1c2f6d4c145"})

	if err := os.WriteFile(z64Path, z64File, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: z64Path, SHA1: "00ba552537f953776b37a05230e9f1c2f6d4c145"})

	v64Path = filepath.Join(data.Dir, "test-bad-v64.v64")
	n64Path = filepath.Join(data.Dir, "test-bad-n64.v64")
	z64Path = filepath.Join(data.Dir, "test-bad-z64.v64")

	if err := os.WriteFile(v64Path, []byte{0, 0x80, 0, 0, 0, 0}, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: v64Path, SHA1: "a5d06af4902696ab97fd92747bc7886c990dfed5"})

	if err := os.WriteFile(n64Path, []byte{0, 0, 0, 0x80, 0, 0}, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: n64Path, SHA1: "a5d06af4902696ab97fd92747bc7886c990dfed5"})

	if err := os.WriteFile(z64Path, []byte{0x80, 0, 0, 0, 0, 0}, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: z64Path, SHA1: "a5d06af4902696ab97fd92747bc7886c990dfed5"})

	return nil
}

// NewSNESFiles makes rom files for the Super Nintendo.
func NewSNESFiles(data *Data) error {
	snesFile1 := make([]byte, 1024)
	snesFile2 := make([]byte, 1536)
	snesPath1 := filepath.Join(data.Dir, "test.smc")
	snesPath2 := filepath.Join(data.Dir, "test.sfc")

	if err := os.WriteFile(snesPath1, snesFile1, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: snesPath1, SHA1: "60cacbf3d72e1e7834203da608037b1bf83b40e8"})

	if err := os.WriteFile(snesPath2, snesFile2, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: snesPath2, SHA1: "60cacbf3d72e1e7834203da608037b1bf83b40e8"})

	return nil
}

// NewSegaFiles makes rom files for the Sega Genesis.
func NewSegaFiles(data *Data) error {
	// Sega, missing SEGA
	mdFile := make([]byte, 0x10000)
	smdFile := make([]byte, 0x10000)
	mgdFile := make([]byte, 0x10000)

	for idx := range mdFile {
		addr := idx / 0x4000
		halfAddr := (idx - addr*0x4000) / 0x2000
		halfAddr++
		addr++
		mdFile[idx] = byte((idx % 2) * addr)
		smdFile[idx] = byte((halfAddr % 2) * addr)

		if idx >= 0x8000 {
			mgdFile[idx] = 0
		} else {
			mgdFile[idx] = byte(halfAddr + ((addr - 1) * 2))
		}
	}

	mdPath := filepath.Join(data.Dir, "nosega.md")
	smdPath := filepath.Join(data.Dir, "nosega.smd")
	mgdPath := filepath.Join(data.Dir, "nosega.mgd")

	if err := os.WriteFile(mdPath, mdFile, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: mdPath, SHA1: "526289b04144ebb25afcbeaf0febb1c0cd60bf79"})

	if err := os.WriteFile(smdPath, append(make([]byte, 512), smdFile...), fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: smdPath, SHA1: "526289b04144ebb25afcbeaf0febb1c0cd60bf79"})

	if err := os.WriteFile(mgdPath, mgdFile, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: mgdPath, SHA1: "526289b04144ebb25afcbeaf0febb1c0cd60bf79"})

	copy(mdFile[256:273], []byte("SEGA GENESIS    "))
	copy(smdFile[128:136], []byte("EAGNSS  "))
	copy(mgdFile[128:136], []byte("EAGNSS  "))
	copy(smdFile[8320:8328], []byte("SG EEI  "))
	copy(mgdFile[32896:32904], []byte("SG EEI  "))

	mdPath = filepath.Join(data.Dir, "sega-md.md")
	smdPath = filepath.Join(data.Dir, "sega-smd.md")
	mgdPath = filepath.Join(data.Dir, "sega-mgd.md")

	if err := os.WriteFile(mdPath, mdFile, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: mdPath, SHA1: "5d2fa3c5c334d6f5c1c0959b040d4452b983d60f"})

	if err := os.WriteFile(smdPath, append(make([]byte, 512), smdFile...), fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: smdPath, SHA1: "5d2fa3c5c334d6f5c1c0959b040d4452b983d60f"})

	if err := os.WriteFile(mgdPath, mgdFile, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: mgdPath, SHA1: "5d2fa3c5c334d6f5c1c0959b040d4452b983d60f"})

	return nil
}

// NewNESFiles makes rom files for the Nintendo Entertainment System.
func NewNESFiles(data *Data) error {
	nesHeaderv1 := []byte{0x4E, 0x45, 0x53, 0x1A, 0x02, 0x06, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	nesFilev1 := make([]byte, 16+32768+49152)
	nesHeaderv1Trainer := []byte{0x4E, 0x45, 0x53, 0x1A, 0x02, 0x06, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	nesFilev1Trainer := make([]byte, 16+512+32768+49152)
	nesHeaderv2 := []byte{0x4E, 0x45, 0x53, 0x1A, 0x02, 0x06, 0, 0x08, 0, 0x11, 0, 0, 0, 0, 0, 0}
	nesFilev2 := make([]byte, 16+4227072+2146304)
	nesFileNoHeader := make([]byte, 32768+49152)

	copy(nesFilev1, nesHeaderv1)
	copy(nesFilev1Trainer, nesHeaderv1Trainer)
	copy(nesFilev2, nesHeaderv2)

	for idx := 16; idx < len(nesFilev1); idx++ {
		if idx < 32768+16 {
			nesFilev1[idx] = 0
			nesFileNoHeader[idx-16] = 0
			nesFilev1Trainer[idx+512] = 0
		} else {
			nesFilev1[idx] = 1
			nesFileNoHeader[idx-16] = 1
			nesFilev1Trainer[idx+512] = 1
		}
	}

	for idx := 16; idx < len(nesFilev2); idx++ {
		if idx < 4227072+16 {
			nesFilev2[idx] = 0
		} else {
			nesFilev2[idx] = 1
		}
	}

	nesPathv1 := filepath.Join(data.Dir, "nes-v1.nes")
	nesPathv1Trainer := filepath.Join(data.Dir, "nes-v1-trainer.nes")
	nesPathv2 := filepath.Join(data.Dir, "nes-v2.nes")
	nesPathNoHeader := filepath.Join(data.Dir, "nes-noheader.nes")

	if err := os.WriteFile(nesPathv1, nesFilev1, fileMode); err != nil {
		return err
	}

	data.AddFile(nesPathv1, "310127efa1522ee9cb559ec502c0f6bb7fde308c")

	if err := os.WriteFile(nesPathv1Trainer, nesFilev1Trainer, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: nesPathv1Trainer, SHA1: "310127efa1522ee9cb559ec502c0f6bb7fde308c"})

	if err := os.WriteFile(nesPathNoHeader, nesFileNoHeader, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: nesPathNoHeader, SHA1: "310127efa1522ee9cb559ec502c0f6bb7fde308c"})

	if err := os.WriteFile(nesPathv2, nesFilev2, fileMode); err != nil {
		return err
	}

	data.Files = append(data.Files, File{Path: nesPathv2, SHA1: "60afb4f8dcb8e1d1c3ba48a8a836f80a89301f65"})

	return nil
}

// ZipRom takes the file, zips it, and returns a file struct for the zipped file.
func ZipRom(file File) (File, error) {
	var err error

	zipPath := fmt.Sprintf("%s.zip", file.Path)

	// Create the .zip file for writing.
	var fileWriter *os.File

	if fileWriter, err = os.Create(zipPath); err != nil {
		return File{}, err
	}
	defer fileWriter.Close()

	// zip compression writer wrapping the .zip file writer.
	zipWriter := zip.NewWriter(fileWriter)

	// Create the ROM file (in the .zip file) for writing.
	var romInZipWriter io.Writer

	if romInZipWriter, err = zipWriter.Create(filepath.Base(file.Path)); err != nil {
		return File{}, err
	}
	defer zipWriter.Close()

	var romData []byte

	// Read the original ROM file.
	if romData, err = os.ReadFile(file.Path); err != nil {
		return File{}, err
	}

	// Write the original ROM to the ROM file in the .zip file.
	if _, err = romInZipWriter.Write(romData); err != nil {
		return File{}, err
	}

	return File{Path: zipPath, SHA1: file.SHA1}, nil
}

// NewZipFiles makes a zip file with a ROM in it.
func NewZipFiles(data *Data) error {
	for _, romFile := range data.Files {
		switch filepath.Ext(romFile.Path) {
		case ".zip", ".gz":
			continue
		}

		zip, err := ZipRom(romFile)
		if err != nil {
			return err
		}

		data.Files = append(data.Files, zip)
	}

	return nil
}

func GzipRom(file File) (File, error) {
	var err error

	gzipPath := fmt.Sprintf("%s.gz", file.Path)

	var fileWriter *os.File

	if fileWriter, err = os.Create(gzipPath); err != nil {
		return File{}, err
	}

	zipWriter := gzip.NewWriter(fileWriter)
	defer zipWriter.Close()

	zipWriter.Header.Name = filepath.Base(file.Path)

	var fileData []byte

	if fileData, err = os.ReadFile(file.Path); err != nil {
		return File{}, err
	}

	if _, err = zipWriter.Write(fileData); err != nil {
		return File{}, err
	}

	return File{Path: gzipPath, SHA1: file.SHA1}, nil
}

// NewGzipFiles makes a gzip file with a ROM in it.
func NewGzipFiles(data *Data) error {
	for _, romFile := range data.Files {
		switch filepath.Ext(romFile.Path) {
		case ".zip", ".gz":
			continue
		}

		gzip, err := GzipRom(romFile)
		if err != nil {
			return err
		}

		data.Files = append(data.Files, gzip)
	}

	return nil
}

func New() (*Data, error) {
	var err error

	data := &Data{}

	dir, err := os.MkdirTemp("", "roms")
	if err != nil {
		return data, err
	}

	data.Dir = dir

	defer func() {
		data.Close()
	}()

	fileCreators := []func(*Data) error{
		NewBinFiles,
		NewLynxFiles,
		NewAtari7800Files,
		NewN64Files,
		NewSNESFiles,
		NewSegaFiles,
		NewNESFiles,
		// These have to be last
		NewZipFiles,
		NewGzipFiles,
	}

	for _, fileCreator := range fileCreators {
		if err := fileCreator(data); err != nil {
			return data, err
		}
	}

	return data, nil
}

type File struct {
	Path string
	SHA1 string
}
