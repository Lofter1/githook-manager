package prepush

import (
	"fmt"
	"log"
	"slices"

	"github.com/Lofter1/githook-manager/util"
	"github.com/spf13/cobra"
)

type blockOptions struct {
	Branches []string `mapstructure:"branches"`
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Blocks the push to certain branches",
	Run:   runBlock,
}

func runBlock(cmd *cobra.Command, args []string) {
	var opts blockOptions

	if err := util.LoadOptions(cmd, &opts); err != nil {
		log.Fatal(err)
	}

	pushBranch := util.GetBranchNameFromRef(remoteRef)
	if slices.Contains(opts.Branches, pushBranch) {
		fmt.Println("Push to branch " + pushBranch + " has been blocked")
	}
}

func init() {
	PrepushCmd.AddCommand(blockCmd)
	util.RegisterOptions(blockCmd, blockOptions{})
}
