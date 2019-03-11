package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var zk = &cobra.Command{
	Use: "zk",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a title")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args, "_")
		createZettelFromTemplate(title)
	},
}

type zettelTemplate struct {
	Timestamp string
	Title     string
}

func createZettelFromTemplate(title string) error {
	timestamp := currentTimestamp()
	filename := fmt.Sprintf("%[1]s-%[2]s.md", timestamp, strings.ToLower(title))

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrap(err, "can't create zettel file")
	}
	defer f.Close()

	t, err := template.ParseFiles("templates/zettel.md")
	if err != nil {
		return errors.Wrap(err, "can't open zettel template")
	}

	if err = t.Execute(f, zettelTemplate{timestamp, title}); err != nil {
		return errors.Wrap(err, "can't write template to the zettel file")
	}

	fmt.Println(filename)

	return nil
}

func currentTimestamp() string {
	return time.Now().Format("20060102_1504")
}

// Execute runs root command.
func Execute() {
	zk.Execute()
}
