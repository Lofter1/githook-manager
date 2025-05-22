package prepush

import (
	"github.com/Lofter1/githook-manager/util"
	"github.com/spf13/cobra"

	_ "embed"
)

//go:embed scriptTemplate.txt
var scriptTemplate string

// prepushCmd represents the prepush command
var PrepushCmd = &cobra.Command{
	Use: "prepush",

	Annotations: map[string]string{
		"hook-file": "pre-push",
		"script":    scriptTemplate,
	},
	Run: util.RunConfiguredFunctionsForHook,
}

var (
	localRef   string
	localSha   string
	remoteRef  string
	remoteSha  string
	remoteName string
	remoteUrl  string
)

func init() {
	PrepushCmd.PersistentFlags().StringVar(&localRef, "localRef", "", "")
	PrepushCmd.PersistentFlags().StringVar(&localSha, "localSha", "", "")
	PrepushCmd.PersistentFlags().StringVar(&remoteRef, "remoteRef", "", "")
	PrepushCmd.PersistentFlags().StringVar(&remoteSha, "remoteSha", "", "")
	PrepushCmd.PersistentFlags().StringVar(&remoteName, "remoteName", "", "")
	PrepushCmd.PersistentFlags().StringVar(&remoteUrl, "remoteUrl", "", "")
}
