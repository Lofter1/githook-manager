package prepush

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Lofter1/githook-manager/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type checklistOptions struct {
	Branches []string `mapstructure:"branches"`
	Checks   []string `mapstructure:"checks"`
}

var checklistCmd = &cobra.Command{
	Use:   "checklist",
	Short: "Displays an interactive checklist before pushing",
	Run:   runChecklist,
}

func runChecklist(cmd *cobra.Command, args []string) {
	var opts checklistOptions

	if err := util.LoadOptions(cmd, &opts); err != nil {
		log.Fatal(err)
	}

	var checklistOptions []huh.Option[string]

	if slices.Contains(opts.Branches, util.GetBranchNameFromRef(remoteRef)) {
		for _, item := range opts.Checks {
			checklistOptions = append(checklistOptions, huh.NewOption(item, item))
		}

		var checked []string
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Ensure you have done the following:").
					Options(checklistOptions...).
					Value(&checked),
			),
		).Run()

		if err != nil {
			log.Fatal(err)
		}
		if len(checked) < len(opts.Checks) {
			fmt.Println("Not all checklist items are confirmed. Aborting push.")
			os.Exit(1)
		}
	}
}

func init() {
	PrepushCmd.AddCommand(checklistCmd)
	util.RegisterOptions(checklistCmd, checklistOptions{})
}
