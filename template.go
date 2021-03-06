package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// given services and a template, do the templating
func Template(services ServiceList, templateFile string) (string, error) {
	funcMap := template.FuncMap{
		"getenv": os.Getenv,
	}

	tname := filepath.Base(templateFile)
	tmpl, err := template.New(tname).Funcs(funcMap).ParseFiles(templateFile)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, services)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func WriteConfig(cfg string, outputFile string) error {
	return ioutil.WriteFile(outputFile, []byte(cfg), 0644)
}

func RunCommand(cmdstr string) error {
	cmd := exec.Command("sh", "-c", cmdstr)
	log.Printf("Running: %s\n", cmdstr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	se := bufio.NewScanner(stderr)
	so := bufio.NewScanner(stdout)
	logIO := func(s *bufio.Scanner, prefix string) {
		for s.Scan() {
			log.Printf("%s%s\n", prefix, s.Bytes())
		}
	}
	go logIO(se, "! ")
	go logIO(so, "> ")
	return cmd.Wait()
}

func LoadTasksAndEmitConfig() error {
	log.Printf("Loading tasks from marathon\n")
	err := services.LoadAllAppTasks(client)
	if err != nil {
		return fmt.Errorf("Unable to load tasks from marathon: %s", err.Error())
	}

	log.Printf("Templating %s with %d services\n", config.TemplateFile, len(services))
	output, err := Template(services, config.TemplateFile)
	if err != nil {
		return fmt.Errorf("Unable to compile template: %s", err.Error())
	}

	// compute hash of previous and current file
	if _, err := os.Stat(config.TargetFile); err == nil {
		old_contents, err := ioutil.ReadFile(config.TargetFile)
		if err != nil {
			return err
		}
		h := sha256.New()
		h.Write(old_contents)
		old_hash := hex.EncodeToString(h.Sum(nil))
		h.Reset()
		h.Write([]byte(output))
		new_hash := hex.EncodeToString(h.Sum(nil))
		if old_hash == new_hash {
			log.Printf("New file is the same as whats on disk\n")
			return nil
		}
		//log.Printf("Computed old %s and new %s\n",old_hash, new_hash)
	} else {
		log.Printf("Skipping checksum, old file doesnt exist\n")
	}

	log.Printf("Writing config to %s\n", config.TargetFile)
	err = WriteConfig(output, config.TargetFile)
	if err != nil {
		return fmt.Errorf("Unable to write config to %s: %s", config.TargetFile, err.Error())
	}
	log.Printf("Wrote %s\n", config.TargetFile)

	if config.Command != "" {
		return RunCommand(config.Command)
	}

	return nil

}
