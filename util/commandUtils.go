package util

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunConfiguredFunctionsForHook(hookCmd *cobra.Command, args []string) {
	for _, function := range hookCmd.Commands() {
		if viper.Get(hookCmd.Use+"."+function.Use) != nil {
			function.Run(function, args)
		}
	}
}

func RegisterOptions(cmd *cobra.Command, opts any) {
	t := reflect.TypeOf(opts)

	schema := map[string]string{}
	for i := range t.NumField() {
		field := t.Field(i)
		key := field.Tag.Get("mapstructure")
		if key == "" {
			key = strings.ToLower(field.Name)
		}
		schema[key] = field.Type.String()
	}

	data, err := json.Marshal(schema)
	if err != nil {
		log.Fatalf("failed to marshal schema: %v", err)
	}
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations["options"] = string(data)
}

func LoadOptions(cmd *cobra.Command, opts any) error {
	return viper.UnmarshalKey(cmd.Parent().Use+"."+cmd.Use, &opts)
}
