package builder

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func BuildGoExecutor(verbose bool) (string, error) {
	entries, err := os.ReadDir("./apps")
	if err != nil {
		return "Error while listing directories", err
	}

	// Extract exist working dir in userinfo.yaml
	var data map[string]interface{}
	yamlFile, err := ioutil.ReadFile("./config/userinfo.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	gitInfo, ok := data["git_info"].([]interface{})
	if !ok {
		log.Fatal("Failed to get git_info from YAML data")
	}

	var wd = gitInfo[0].(map[interface{}]interface{})["exist_workingdir"]

	var result string

	if wd == nil {
		for _, e := range entries {
			if e.IsDir() == true {
				var tidy_result string
				tidyCmd := []string{"go", "mod"}
				tidyCmd = append(tidyCmd, "tidy")
				if tidy_result, err = ExecCommand("./apps/"+e.Name(), tidyCmd); err != nil {
					return tidy_result, err
				}

				output := "bin/executor"
				buildCmd := []string{"go", "build"}
				buildCmd = append(buildCmd, "-o", output, ".")

				if verbose {
					if err := ExecCommandPipe("./apps/"+e.Name(), buildCmd, os.Stdout, os.Stderr); err != nil {
						return result, err
					}
				} else {
					if result, err = ExecCommand("./apps/"+e.Name(), buildCmd); err != nil {
						return result, err
					}
				}

				if _, err := os.Stat("./repository/" + e.Name()); errors.Is(err, os.ErrNotExist) {
					if err := os.Mkdir("./repository/"+e.Name(), os.ModePerm); err != nil {
						return "", err
					}
				}

				copyExecutorErr := copyBinaries("./apps/"+e.Name()+"/"+output, "./repository/"+e.Name()+"/executor")
				if copyExecutorErr != nil {
					return "", copyExecutorErr
				}

				permissionErr := os.Chmod("./repository/"+e.Name()+"/executor", 755)
				if permissionErr != nil {
					return "", permissionErr
				}
			}
		}
	} else {
		for _, e := range entries {
			if e.IsDir() == true {
				output := "bin/executor"
				buildCmd := []string{"go", "build"}
				buildCmd = append(buildCmd, "-o", output, ".")

				if verbose {
					if err := ExecCommandPipe("./apps/"+e.Name(), buildCmd, os.Stdout, os.Stderr); err != nil {
						return result, err
					}
				} else {
					if result, err = ExecCommand("./apps/"+e.Name(), buildCmd); err != nil {
						return result, err
					}
				}

				destination := strings.ReplaceAll(fmt.Sprintf("%s %v", wd, "/"), " ", "")

				if _, err := os.Stat(destination + e.Name()); errors.Is(err, os.ErrNotExist) {
					if err := os.Mkdir(destination+e.Name(), os.ModePerm); err != nil {
						return "", err
					}
				}

				copyExecutorErr := copyBinaries("./apps/"+e.Name()+"/"+output, destination+e.Name()+"/executor")
				if copyExecutorErr != nil {
					return "", copyExecutorErr
				}

				permissionErr := os.Chmod(destination+e.Name()+"/executor", 755)
				if permissionErr != nil {
					return "", permissionErr
				}
			}
		}
	}

	fmt.Printf("Executor built in local environment.\n")
	return result, nil
}

func BuildGoHandler(verbose bool) (string, error) {
	entries, err := os.ReadDir("./apps")
	if err != nil {
		return "Error while listing directories", err
	}

	// Extract exist working dir in userinfo.yaml
	var data map[string]interface{}
	yamlFile, err := ioutil.ReadFile("./config/userinfo.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	gitInfo, ok := data["git_info"].([]interface{})
	if !ok {
		log.Fatal("Failed to get git_info from YAML data")
	}

	var wd = gitInfo[0].(map[interface{}]interface{})["exist_workingdir"]

	var result string

	if wd == nil {
		for _, e := range entries {
			if e.IsDir() == true {
				output := "bin/handler"
				buildCmd := []string{"go", "build"}
				buildCmd = append(buildCmd, "-o", output, "-buildmode=plugin", "./src")

				if verbose {
					if err := ExecCommandPipe("./apps/"+e.Name(), buildCmd, os.Stdout, os.Stderr); err != nil {
						return result, err
					}
				} else {
					if result, err = ExecCommand("./apps/"+e.Name(), buildCmd); err != nil {
						return result, err
					}
				}

				copyHandlerErr := copyBinaries("./apps/"+e.Name()+"/"+output, "./repository/"+e.Name()+"/handler")
				if copyHandlerErr != nil {
					return "", copyHandlerErr
				}
			}
		}

	} else {
		for _, e := range entries {
			if e.IsDir() == true {
				output := "bin/handler"
				buildCmd := []string{"go", "build"}
				buildCmd = append(buildCmd, "-o", output, "-buildmode=plugin", "./src")

				if verbose {
					if err := ExecCommandPipe("./apps/"+e.Name(), buildCmd, os.Stdout, os.Stderr); err != nil {
						return result, err
					}
				} else {
					if result, err = ExecCommand("./apps/"+e.Name(), buildCmd); err != nil {
						return result, err
					}
				}

				destination := strings.ReplaceAll(fmt.Sprintf("%s %v", wd, "/"), " ", "")

				copyHandlerErr := copyBinaries("./apps/"+e.Name()+"/"+output, destination+e.Name()+"/handler")
				if copyHandlerErr != nil {
					return "", copyHandlerErr
				}
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
