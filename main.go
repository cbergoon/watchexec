package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/fsnotify/fsnotify"
)

const ConfigFile = "watchexec.yaml"

type Config struct {
	Commands []Command `yaml:"commands"`
	Ignore   []string  `yaml: "ignore,flow`
}

type Command struct {
	Executable string   `yaml:"executable"`
	Arguments  []string `yaml:"arguments,flow"`
	Sequence   int      `yaml:"sequence"`
}

var watcher *fsnotify.Watcher
var lastProcess *exec.Cmd

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
	fmt.Println(1)
	_, err = os.Stat(filepath.Join(wd, ConfigFile))
	fmt.Println(2)
	if os.IsNotExist(err) {
		fmt.Println(3)
		ccf := filepath.Join(wd, ConfigFile)
		err := createConfigFile(ccf)
		if err != nil {
			log.Fatalf("error: %+v", err)
		}
		fmt.Printf("Initialized project with watchexec config file at: %s", ccf)
	}

	fmt.Println(4)
	configBytes, err := ioutil.ReadFile(filepath.Join(wd, ConfigFile))
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
	fmt.Println(5)
	config := &Config{}

	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	sort.Slice(config.Commands, func(i, j int) bool {
		return config.Commands[i].Sequence < config.Commands[j].Sequence
	})

	if len(config.Commands) == 0 {
		return
	}

	var changed bool
	// var changedMutex sync.Mutex

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(wd, watchDir); err != nil {
		log.Fatalf("error: %+v", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				ignored := false
				for _, il := range config.Ignore {
					if strings.HasSuffix(event.Name, il) {
						ignored = true
					}
				}
				if !ignored {
					changed = true
				}
			case <-watcher.Errors:
			}
		}
	}()

	go func() {
		for {
			if changed {
				if lastProcess != nil {
					lastProcess.Process.Kill()
				}
				go ExecCommands(config.Commands)
				changed = false
				time.Sleep(time.Millisecond * 5000)
			} else {
				time.Sleep(time.Millisecond * 5000)
			}
		}
	}()

	<-done
}

func createConfigFile(filePath string) error {
	sampleConfig := `ignore: []
commands:
- executable: go
  arguments: [version]
  sequence: 0`
	err := ioutil.WriteFile(filePath, []byte(sampleConfig), 0644)
	if err != nil {
		return err
	}
	return nil
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}

func ExecCommands(commands []Command) {
	for ci, c := range commands {
		fmt.Printf("$ %s %s \n", c.Executable, strings.Join(c.Arguments, " "))
		cmd := exec.Command(c.Executable, c.Arguments...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
		if len(commands)-1 == ci { // Last Command
			lastProcess = cmd
		}
		cmd.Wait()
	}
}
