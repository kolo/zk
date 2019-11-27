package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var zt = template.Must(template.New("zettel").Parse(`# {{.Timestamp}} {{.Title}}`))

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

		filename, err := createZettelFromTemplate(title)
		if err != nil {
			return err
		}

		if err = exec.Command("code", filename).Run(); err != nil {
			return errors.Wrap(err, "can't start editor")
		}

		return nil
	},
}

type zettelTemplate struct {
	Timestamp string
	Title     string
}

func createZettelFromTemplate(title string) (string, error) {
	timestamp := currentTimestamp()
	filename := fmt.Sprintf("%[1]s-%[2]s.md", timestamp, strings.ToLower(title))

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return "", errors.Wrap(err, "can't create zettel file")
	}
	defer f.Close()

	if err = zt.Execute(f, zettelTemplate{timestamp, title}); err != nil {
		return "", errors.Wrap(err, "can't write template to the zettel file")
	}

	return filename, err
}

func currentTimestamp() string {
	return time.Now().Format("20060102_1504")
}

// Execute runs root command.
func Execute() {
	zk.Execute()
}
