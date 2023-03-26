package builder

import (
	"fmt"
	"io"
	"os"
)

func BuildGoExecutor(verbose bool) (string, error) {
	entries, err := os.ReadDir("./")
	if err != nil {
		return "Error while listing directories", err
	}

	//var builderr error
	var result string

	for _, e := range entries {
		if e.IsDir() == true && e.Name() != "repository" {
			output := "bin/executor"
			buildCmd := []string{"go", "build"}
			buildCmd = append(buildCmd, "-o", output, ".")

			if verbose {
				if err := ExecCommandPipe("./"+e.Name(), buildCmd, os.Stdout, os.Stderr); err != nil {
					return result, err
				}
			} else {
				if result, err = ExecCommand("./"+e.Name(), buildCmd); err != nil {
					return result, err
				}
			}

			if err := os.Mkdir("./repository/"+e.Name(), os.ModePerm); err != nil {
				return "", err
			}

			copyExecutorErr := copyBinaries("./"+e.Name()+"/"+output, "./repository/"+e.Name()+"/executor")
			if copyExecutorErr != nil {
				return "", copyExecutorErr
			}

			permissionErr := os.Chmod("./repository/"+e.Name()+"/executor", 755)
			if permissionErr != nil {
				return "", permissionErr
			}
		}
	}

	fmt.Printf("Executor built in local environment.\n")
	return result, nil
}

func BuildGoHandler(verbose bool) (string, error) {
	entries, err := os.ReadDir("./")
	if err != nil {
		return "Error while listing directories", err
	}

	//var builderr error
	var result string

	for _, e := range entries {
		if e.IsDir() == true && e.Name() != "repository" {
			output := "bin/handler"
			buildCmd := []string{"go", "build"}
			buildCmd = append(buildCmd, "-o", output, "-buildmode=plugin", "./src")

			if verbose {
				if err := ExecCommandPipe("./"+e.Name(), buildCmd, os.Stdout, os.Stderr); err != nil {
					return result, err
				}
			} else {
				if result, err = ExecCommand("./"+e.Name(), buildCmd); err != nil {
					return result, err
				}
			}

			copyHandlerErr := copyBinaries("./"+e.Name()+"/"+output, "./repository/"+e.Name()+"/handler")
			if copyHandlerErr != nil {
				return "", copyHandlerErr
			}
		}
	}

	fmt.Printf("Handler built in local environment.\n")
	return result, nil
}

func copyBinaries(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
