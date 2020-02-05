package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	zettelTemplateStr = `---
id: {{.Timestamp}}
title: {{.Title}}
tags:
---

`
)

var zt = template.Must(template.New("zettel").Parse(zettelTemplateStr))

var zk = &cobra.Command{
	Use: "zk",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a title")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, "_")

		name, err := createZettelFromTemplate(title)
		if err != nil {
			return err
		}

		fmt.Println(name)

		return nil
	},
}

type zettelTemplate struct {
	Timestamp string
	Title     string
}

func createZettelFromTemplate(title string) (string, error) {
	timestamp := currentTimestamp()
	name := fmt.Sprintf("%[1]s-%[2]s.md", timestamp, strings.ToLower(title))
	path, err := filepath.Abs(name)
	if err != nil {
		return "", errors.Wrapf(err, "can't get absolute path for %s", name)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return "", errors.Wrapf(err, "can't create %s", path)
	}
	defer f.Close()

	if err = zt.Execute(f, zettelTemplate{timestamp, title}); err != nil {
		return "", errors.Wrapf(err, "can't write template to %s", path)
	}

	return path, err
}

func currentTimestamp() string {
	return time.Now().Format("200601021504")
}

// Execute runs root command.
func Execute() {
	zk.Execute()
}
