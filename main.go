package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

var Version = "v0.1.0"

const (
	flagIn           = "in"
	flagInShort      = "i"
	flagVerbose      = "verbose"
	flagVerboseShort = "ver"
	flagDebug        = "debug"
	flagDebugShort   = "de"
)

func main() {
	app := &cli.App{
		Name:    "GoFind",
		Usage:   "Find all files and directories with pattern",
		Action:  find,
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagIn,
				Aliases: []string{flagInShort},
				Usage:   "Search in directory",
				Value:   ".",
			},
			&cli.BoolFlag{
				Name:    flagVerbose,
				Aliases: []string{flagVerboseShort},
				Usage:   "print all the files and directories searched for",
			},
			&cli.BoolFlag{
				Name:    flagDebug,
				Aliases: []string{flagDebugShort},
				Usage:   "print all the files and directories searched for",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var (
	logErr  bool
	verbose bool
)

func find(c *cli.Context) error {
	if len(c.Args().First()) == 0 {
		log.Fatal("Please enter file or directory name to search for")
	}

	if len(c.String(flagIn)) == 0 {
		log.Fatal("In parameter can't be empty")
	}

	if _, err := os.Stat(c.String(flagIn)); err != nil {
		log.Fatal("Failed to find the path defined at in parameter")
	}

	fmt.Printf("Searching for: %s in: %s\n", c.Args().First(), c.String(flagIn))

	logErr, verbose = c.Bool(flagVerbose), c.Bool(flagVerbose)

	start := time.Now()

	matches := []string{}

	err := filepath.WalkDir(c.String(flagIn), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logError(err)
			return nil
		}

		logInfo(path)

		ok, err := filepath.Match(c.Args().First(), d.Name())
		if err != nil {
			logError(err)
			return nil
		}

		if ok {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		logError(err)
	}

	elapsed := time.Since(start)
	fmt.Println("Matches:")
	if len(matches) == 0 {
		fmt.Println("File or directory not found")
	} else {
		for _, m := range matches {
			fmt.Println(m)
		}
	}
	fmt.Printf("Elapsed time: %s\n", elapsed)

	return nil
}

func logError(e error) {
	if logErr {
		log.Println(e)
	}
}

func logInfo(s string) {
	if verbose {
		log.Println(s)
	}
}
