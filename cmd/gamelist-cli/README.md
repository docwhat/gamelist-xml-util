# GameList CLI utility

This is a CLI for managing and manipulating game list data written in Go.

It should use the `docwhat.org/gamelist-xml-util/pkg/gamelist` and `docwhat.org/gamelist-xml-util/pkg/miyoogamelist` packages to do the heavy lifting.

The CLI should:

- Require the path to a `gamelist.xml` file. If the path is `-` then it will read from stdin.
- Accept the optional flag `--output` with a path to an output file. The default is `-` which means stdout.
- Accept the optional flag `--roms` with a path to a ROMs directory. The default is the directory of the `gamelist.xml` or `.` if the `gamelist.xml` file is stdin.
- Accept the optional flag `--add-roms` which will add games from the ROMs directory to the `gamelist.xml` file that are not already there.
- Accept the optional flag `--add-images` which will add images from the ROMs directory to games in the `gamelist.xml` that are missing an image but the image exists in the ROMs directory.
