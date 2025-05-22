# Using git hooks made easy

## Why githook manager
1. Easy to use
	- A setup command will configure everything you need
2. Extensible 
	- Easy to add more functionality due to a simple architecture
3. Configurable 
	- The git hook functions are configured using a yaml file 
	- The config file lives inside the .git directory. No unnecesary cluttering of the work directory
4. Powerful 
	- Due to the functions being written in go, you have all the access to features high-level programming languages like go offer you.
	- Beautiful prompts using [huh](https://github.com/charmbracelet/huh)
	- Easy git interaction using [go-git](https://github.com/go-git/go-git)
5. Even more extensible!
	- githook manager is called from a normal shell script in the .git/hooks/ directory
	- You can simply edit the shell script to add whatever else you want to do in the hook

## Usage

### Set up a git hook

This program provides an easy way to set up git hooks. Just call the setup program from inside the git repository you want to configure hooks in

```sh
githook-manager setup
```

A UI will show you all functions and prompt you to configure the necessary options.

It will then write your configuration to the file `.git/.githook-manager.yaml`, back up any existing hook that and add a shell script that executes the functionality.

### Installation

Currently, the only supported way is using `go install`. Make sure you have [golang](https://go.dev/doc/install) installed and the the `$GOPATH/bin` is in your `$PATH`.

```sh
go install github.com/Lofter1/githook-manager
```

## Development
### Creating a new function

Each function has it's own file in the hook directory it will be registered under.

Functions are defined using cobra commands and will be registered in the file it's defined in during `init()`

```go
var myNewFunctionCmd = &cobra.Command{
	Use: "myfunction",
	Short: "Shortly describe what the function does"
	Run: func(cmd *cobra.Command, args []string) {},
    ...
}

func init() {
	HookCmd.AddCommand(myNewFunctionCmd)
}
```

#### Handling options

##### Adding and registering options
In order to register options for a function first define an option struct and then register it in the `init()`. This will make sure that the setup command knows what to prompt the user.

```go
type myfunctionOptions struct {
	MyOption  string `mapstructure:"MyOption"`
	MyOption2 []string `mapstructure:"MyOption2"`
}

func init() {
    util.RegisterOptions(myNewFunctionCmd, myfunctionOptions{})
}
```
The mapstructure tag is important, as these options will be saved in a yaml file.

##### Reading options

You can easily load all options into an object using the `LoadOptions()` helper

```go
var opts myfunctionOptions

if err := util.LoadOptions(cmd, &opts); err != nil {
    log.Fatal(err)
}
```