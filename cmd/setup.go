package cmd

// TODO: Refactor

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupCmdDescription = `An easy user interface to setup git hooks and automatically 
creates the config file and hook scripts.

Any exisitng hook script that conflicts with the hooks set 
up by the user will be backed up in the hooks folder as <hook-name>.bak`

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Quickly configure git hooks",
	Long:  setupCmdDescription,
	Run:   runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) {
	_, err := os.Stat("./.git")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal(fmt.Errorf("not in a git repository"))
		}
		log.Fatal(err)
	}

	hooks := getHooks(rootCmd)
	functionToConfig, err := promptSelectHookFunctionToConfig(hooks)
	if err != nil {
		log.Fatal(err)
	}

	options, err := promptConfigValues(functionToConfig)
	if err != nil {
		log.Fatal(err)
	}

	setOptions(functionToConfig, options)

	if err := writeHookScript(functionToConfig.Parent()); err != nil {
		log.Fatal(err)
	}

	if err := writeConfigFile(); err != nil {
		log.Fatal(err)
	}

}

func getHooks(root *cobra.Command) []*cobra.Command {
	var hooks []*cobra.Command
	for _, cmd := range root.Commands() {
		if _, ok := cmd.Annotations["hook-file"]; ok {
			hooks = append(hooks, cmd)
		}
	}
	return hooks
}

func getHookFunctions(hook *cobra.Command) []*cobra.Command {
	var functions []*cobra.Command
	for _, cmd := range hook.Commands() {
		functions = append(functions, cmd)
	}
	return functions
}

func writeHookScript(hook *cobra.Command) error {
	script, ok := hook.Annotations["script"]
	if !ok || script == "" {
		return fmt.Errorf("hook %q does not have a script annotation", hook.Use)
	}

	hookFile, ok := hook.Annotations["hook-file"]
	if !ok || hookFile == "" {
		return fmt.Errorf("hook %q does not have a hook-file annotation", hook.Use)
	}

	hookPath := filepath.Join(".git", "hooks", hookFile)
	backupPath := hookPath + ".bak"

	err := os.Rename(hookPath, backupPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	err = os.WriteFile(hookPath, []byte(script), 0755)
	if err != nil {
		return err
	}

	return nil
}

func setOptions(function *cobra.Command, options map[string]any) {
	for optionName, optionValue := range options {
		setOption(function, optionName, optionValue)
	}
}

func setOption(function *cobra.Command, optionName string, value any) {
	optPath := fmt.Sprintf("%s.%s.%s", function.Parent().Use, function.Use, optionName)
	viper.Set(optPath, value)
}

func writeConfigFile() error {
	if viper.ConfigFileUsed() != "" {
		return viper.WriteConfig()
	}
	if cfgFile != "" {
		return viper.WriteConfigAs(cfgFile)
	}

	path := filepath.Join(defaultConfigDir, defaultConfigName+"."+defaultConfigType)
	return viper.WriteConfigAs(path)
}

func getFunctionOptions(function *cobra.Command) (map[string]string, error) {
	if function == nil {
		return nil, fmt.Errorf("no function was provided")
	}
	optJson, hasOptions := function.Annotations["options"]
	if !hasOptions {
		return nil, nil
	}
	var schema map[string]string
	err := json.Unmarshal([]byte(optJson), &schema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// huh helper

func stringListPrompt(title string) ([]string, error) {
	var inputs []string

	for continueAdd := true; continueAdd; {
		input, err := stringPrompt(title)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)

		err = huh.NewConfirm().
			Title("Add another?").
			Affirmative("Yes").
			Negative("No").
			Value(&continueAdd).Run()
		if err != nil {
			return nil, err
		}
	}
	return inputs, nil
}

func stringPrompt(title string) (string, error) {
	var input string
	err := huh.NewInput().Title(title).Value(&input).Run()
	if err != nil {
		return "", err
	}
	return input, nil
}

func promptSelectHookFunctionToConfig(hooks []*cobra.Command) (*cobra.Command, error) {
	var selectedHookIndex int
	var selectedFunction *cobra.Command
	var huhHookOptions []huh.Option[int]

	for i, hook := range hooks {
		huhHookOptions = append(huhHookOptions, huh.NewOption(hook.Use, i))
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().Options(huhHookOptions...).
				Title("Hook selection").
				Description("Select a git hook to configure").
				Value(&selectedHookIndex),

			huh.NewSelect[*cobra.Command]().
				OptionsFunc(func() []huh.Option[*cobra.Command] {
					var huhFunctionOptions []huh.Option[*cobra.Command]
					for _, funcCmd := range getHookFunctions(hooks[selectedHookIndex]) {
						huhOptKey := fmt.Sprintf("%s - %s", funcCmd.Use, funcCmd.Short)
						huhOpt := huh.NewOption(huhOptKey, funcCmd)
						huhFunctionOptions = append(huhFunctionOptions, huhOpt)
					}
					return huhFunctionOptions
				}, &selectedHookIndex).
				Title("Hook function").
				Description("Select a function to register for the selected hook").
				Value(&selectedFunction),
		),
	).Run()
	if err != nil {
		return nil, err
	}
	return selectedFunction, nil
}

func promptConfigValues(fn *cobra.Command) (map[string]any, error) {
	configuredOptions := map[string]any{}
	opts, err := getFunctionOptions(fn)
	if err != nil {
		return nil, err
	}

	for opt, optsType := range opts {
		var input any
		switch optsType {
		case "string":
			input, err = stringPrompt(opt)
			if err != nil {
				return nil, err
			}
		case "[]string":
			input, err = stringListPrompt(opt)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported option type: %s", optsType)
		}
		configuredOptions[opt] = input

	}

	return configuredOptions, nil
}
