package plugin

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/patterninc/heimdall/internal/pkg/aws"
)

const (
	jobDirectoryPermissions fs.FileMode = 0700
	stdoutName                          = `stdout`
	stderrName                          = `stderr`
	s3Prefix                            = `s3://`
	separator                           = `/`
)

type Runtime struct {
	WorkingDirectory string
	ArchiveDirectory string
	ResultDirectory  string
	Version          string
	UserAgent        string
	Stdout           *os.File
	Stderr           *os.File
}

func (r *Runtime) Set() (err error) {

	// create working directory
	if err := os.MkdirAll(r.WorkingDirectory, jobDirectoryPermissions); err != nil {
		return err
	}

	// create result directory
	if !strings.HasPrefix(r.ResultDirectory, s3Prefix) {
		if err := os.MkdirAll(r.ResultDirectory, jobDirectoryPermissions); err != nil {
			return err
		}
	}

	// set stdout and stderr
	if r.Stdout, err = os.Create(path.Join(r.WorkingDirectory, stdoutName)); err != nil {
		return err
	}

	if r.Stderr, err = os.Create(path.Join(r.WorkingDirectory, stderrName)); err != nil {
		return err
	}

	// fmt.Println(r.WorkingDirectory)

	return nil

}

func (r *Runtime) Close() {

	if r.Stdout != nil {
		r.Stdout.Close()
	}

	if r.Stderr != nil {
		r.Stderr.Close()
	}

	// we're making our best effort to move files from working dir to archive
	// we start it as a go routine to minimize impact of sync commands
	go func() {
		if err := copyDir(r.WorkingDirectory, r.ArchiveDirectory); err != nil {
			// FIXME: implement better error handling
			fmt.Println(`archive:`, err)
			return
		}
	}()

	// ...and now we make the best effort to remove working directory
	go func() {
		if err := os.RemoveAll(r.WorkingDirectory); err != nil {
			// FIXME: implement better error handling
			fmt.Println(`cleanup:`, err)
			return
		}
	}()

}

func copyDir(src, dst string) error {

	// Read the source directory
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// setup write file function based on our destination
	// if we have local filesystem, crete directory as appropriate
	writeFileFunc := os.WriteFile
	if strings.HasPrefix(dst, s3Prefix) {
		writeFileFunc = aws.WriteToS3
	} else {
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			os.MkdirAll(dst, jobDirectoryPermissions)
		}
	}

	// Iterate through each file in the source directory
	for _, file := range files {
		srcPath := src + separator + file.Name()
		dstPath := dst + separator + file.Name()

		// Check if the file is a directory
		if file.IsDir() {
			// Recursively copy the subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy the file
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			if err = writeFileFunc(dstPath, data, 0600); err != nil {
				return err
			}
		}
	}

	return nil

}
