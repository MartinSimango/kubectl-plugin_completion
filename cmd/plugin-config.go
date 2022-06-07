package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/MartinSimango/kubectl-plugin_completion/scripts"
	"gopkg.in/yaml.v2"
)

var homeDir, _ = os.UserHomeDir()
var configFolder = homeDir + "/.kube/plugin-completion-config"

type Shells struct {
	Name       string
	ConfigFile string
}

var shells = map[string]string{"zsh": "zsh.yaml", "bash": "bash.yaml"}

//  list will be populated if DoesPluginHaveCobraSupport function fails to work correctly for a particular
//  plugin
var cobraBlackList = []string{"view_secret", "auth_proxy", "blame"}

type Plugin struct {
	Name                          string `yaml:"name"`
	CompletionFunctionName        string `yaml:"completionFunctionName"`
	Description                   string `yaml:"description"`
	PluginSupportsCobraCompletion bool   `yaml:"supportsCobraCompletion"`
}

type PluginConfigImpl struct {
	ConfigFile             string   `yaml:"configfile,omitempty"`
	Shell                  string   `yaml:"shell"`
	ShellLocation          string   `yaml:"shellLocation"`
	Plugins                []Plugin `yaml:"plugins"`
	KubectlOverridePlugins []string `yaml:"kubectlOverridePlugins"`
}

type PluginConfig interface {
	AddPlugin(plugin Plugin) error
	DisplayConfig() error
	DoesPluginExist(plugin string) bool
	DoesPluginHaveCobraSupport(plugin string) bool
	EditPlugin(pluginName string, completionFunctionName, description string, completionFunctionSet, descriptionSet bool) error
	GetPlugin(pluginName string) (*Plugin, error)
	PrintAllPlugins()
	PrintCobraPlugins()

	GeneratePluginConfig() error
	GetCompletionFunctionName(plugin string) string
	GetCompletionScriptSection(plugin string) string

	GetConfig() *PluginConfigImpl
	GenerateCompletionScript() (string, error)
	WriteCleanConfig() error

	checkConfigFile() error
	createConfigFolder() error

	doesPluginUseCobra(plugin string) bool
	doesPluginHaveCobraCompletionCommand(plugin string) bool
	doesPluginCompletionOverrideKubectl(plugin string) bool

	getAllKubectlPlugins() ([]string, error)

	populatePlugins([]string) error
	removeOldPlugins(plugins []string)

	writeConfig() error
}

var _ PluginConfig = &PluginConfigImpl{}

func NewEmptyPluginImpl(shell, shellLocation string) *PluginConfigImpl {
	//validate shell name

	config := PluginConfigImpl{Shell: shell, ShellLocation: shellLocation, ConfigFile: configFolder + "/" + shells[shell]}
	config.checkConfigFile()
	return &config
}

func NewPluginConfigImpl(shell, shellLocation string) *PluginConfigImpl {
	//validate shell name
	config := NewEmptyPluginImpl(shell, shellLocation)
	return config.GetConfig()
}

func NewZshPluginConfigImpl() *PluginConfigImpl {
	return NewPluginConfigImpl("zsh", "/bin/zsh")
}

func NewBashPluginConfigImpl() *PluginConfigImpl {
	return NewPluginConfigImpl("bash", "/bin/bash")

}

func (pc *PluginConfigImpl) AddPlugin(plugin Plugin) error {
	if pc.DoesPluginExist(plugin.Name) {
		return fmt.Errorf("Plugin with name %s already exists", plugin.Name)
	}

	pc.Plugins = append(pc.Plugins, plugin)
	if pc.doesPluginCompletionOverrideKubectl(plugin.Name) {
		pc.KubectlOverridePlugins = append(pc.KubectlOverridePlugins, plugin.Name)
	}
	return nil
}

func (pc *PluginConfigImpl) EditPlugin(pluginName string, completionFunctionName, description string, completionFunctionSet, descriptionSet bool) error {
	if !pc.DoesPluginExist(pluginName) {
		return fmt.Errorf("Plugin with name '%s' does not exist", pluginName)
	}
	plugin, err := pc.GetPlugin(pluginName)

	if err != nil {
		return err
	}

	if completionFunctionSet {
		plugin.CompletionFunctionName = completionFunctionName
	}
	if descriptionSet {
		plugin.Description = description
	}

	pc.writeConfig()
	return nil

}

func (pc *PluginConfigImpl) GetPlugin(pluginName string) (*Plugin, error) {
	for i, plugin := range pc.Plugins {
		if plugin.Name == pluginName {
			return &pc.Plugins[i], nil
		}
	}
	return nil, fmt.Errorf("Plugin with name %s does not exist", pluginName)
}

func (pc *PluginConfigImpl) getAllKubectlPlugins() ([]string, error) {
	command := `kubectl plugin list | grep kubectl | xargs basename | cut -d "-" -f2`
	output, err := exec.Command(pc.ShellLocation, "-c", command).Output()

	if err != nil {
		return nil, fmt.Errorf((string(output[:])))
	}

	plugins := strings.Split(string(output[:]), "\n")
	return plugins[:len(plugins)-1], nil
}

func (pc *PluginConfigImpl) GeneratePluginConfig() error {

	plugins, err := pc.getAllKubectlPlugins()

	if err != nil {
		return nil
	}

	err = pc.populatePlugins(plugins)

	if err != nil {
		return err
	}

	err = pc.writeConfig()

	if err != nil {
		return err
	}

	return nil
}

func (pc *PluginConfigImpl) populatePlugins(plugins []string) error {
	for _, plugin := range plugins {
		if pc.DoesPluginExist(plugin) {
			//check for completion and check for cobra and print
			continue
		}

		completionFunctionName := pc.GetCompletionFunctionName(plugin)

		err := pc.AddPlugin(Plugin{
			Name:                          plugin,
			CompletionFunctionName:        completionFunctionName,
			Description:                   pc.getPluginDescription(plugin),
			PluginSupportsCobraCompletion: completionFunctionName != "",
		})

		if err != nil {
			return err
		}
	}

	pc.removeOldPlugins(plugins)
	return nil
}

func (pc *PluginConfigImpl) getPluginDescription(plugin string) string {
	switch plugin {
	case "krew":
		return "krew is the kubectl plugin manager"
	case "plugin_completion":
		return "A kubectl plugin for allowing shell completions for kubectl plugins"
	case "stern":
		return "Multi pod and container log tailing"
	}
	command := "kubectl krew search " + strings.ReplaceAll(plugin, "_", "-") + " | head -n 2 | awk -F '  ' 'NR>1 {print $2}'" //awk '/DESCRIPTION:/,NF==0' | sed -n '1!p'"

	// command := "kubectl krew info " + strings.ReplaceAll(plugin, "_", "-") + " | awk '/DESCRIPTION:/,NF==0' | sed -n '1!p'"
	output, err := exec.Command(pc.ShellLocation, "-c", command).Output()
	outputString := string(output[:])
	if err != nil {
		fmt.Println(command)
		log.Fatal(string(outputString))
	}

	if strings.TrimSpace(outputString) == "" || strings.HasPrefix(outputString, "error:") {
		return fmt.Sprintf("A kubectl plugin called %s", plugin)
	}
	description := strings.ReplaceAll(outputString, "\n", " ")

	return description
}

func (pc *PluginConfigImpl) removeOldPlugins(currentPlugins []string) {
	for i := 0; i < len(pc.Plugins); i++ {
		matched := false
		for _, plugin := range currentPlugins {
			if pc.Plugins[i].Name == plugin {
				matched = true
				break
			}
		}
		if !matched { // remove plugin TODO: make this more efficient
			pc.Plugins = append(pc.Plugins[:i], pc.Plugins[i+1:]...)
			i--
		}
	}
}
func (pc *PluginConfigImpl) GetConfig() *PluginConfigImpl {

	pc.checkConfigFile()

	yfile, err := ioutil.ReadFile(pc.ConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yfile, pc)

	if err != nil {
		log.Fatal(err)
	}
	return pc
}

func (pc *PluginConfigImpl) DisplayConfig() error {

	yamlData, err := yaml.Marshal(removeConfigFile(*pc))

	if err != nil {
		// return fmt.Errorf("Error while Marshaling. %v", err)
		return err
	}

	fmt.Print(string(yamlData[:]))
	return nil

}

func (pc *PluginConfigImpl) WriteCleanConfig() error {
	return NewEmptyPluginImpl(pc.Shell, pc.ShellLocation).writeConfig()
}

func (pc *PluginConfigImpl) writeConfig() error {

	yamlData, err := yaml.Marshal(removeConfigFile(*pc))
	if err != nil {
		// return fmt.Errorf("Error while Marshaling. %v", err)
		return err
	}

	err = ioutil.WriteFile(pc.ConfigFile, yamlData, 0666)

	if err != nil {
		return err
	}

	return nil
}

func (pc *PluginConfigImpl) checkConfigFile() error {
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		return pc.createConfigFolder()
	}
	if _, err := os.Stat(pc.ConfigFile); os.IsNotExist(err) {
		return pc.writeConfig()
	}
	return nil
}

func (pc *PluginConfigImpl) createConfigFolder() error {
	err := os.MkdirAll(configFolder, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func removeConfigFile(pc PluginConfigImpl) PluginConfigImpl {
	pc.ConfigFile = ""
	return pc
}

func (pc *PluginConfigImpl) DoesPluginExist(plugin string) bool {
	for _, p := range pc.Plugins {
		if p.Name == plugin {
			return true
		}
	}
	return false
}

func (pc *PluginConfigImpl) GetCompletionFunctionName(plugin string) string {
	if !pc.DoesPluginHaveCobraSupport(plugin) {
		return ""
	}

	switch plugin {
	case "support_bundle":
		plugin = "support-bundle"
	}
	switch pc.Shell {
	case "zsh":
		return "_" + plugin
	case "bash":
		return "__start_" + plugin
	}
	return "unsupported shell"
}

func (pc *PluginConfigImpl) DoesPluginHaveCobraSupport(plugin string) bool {
	if isInList(cobraBlackList, plugin) {
		return false
	}
	return pc.doesPluginUseCobra(plugin) && pc.doesPluginHaveCobraCompletionCommand(plugin)

}

func (pc *PluginConfigImpl) doesPluginUseCobra(plugin string) bool {
	// command := "grep ValidArgsFunction $(which kubectl-" + plugin + ") || true"

	// output, err := exec.Command(pc.ShellLocation, "-c", command).Output()
	// outputString := string(output[:])
	// if err != nil {
	// 	fmt.Println(command)
	// 	log.Fatal(string(outputString))
	// }

	// if strings.TrimSpace(outputString) == "" {
	// 	return false
	// } else {
	// 	return true
	// }
	return true

}

func (pc *PluginConfigImpl) doesPluginHaveCobraCompletionCommand(plugin string) bool {
	command := "strings $(which kubectl-" + plugin + ")" +
		` | grep "Generate the autocompletion script for powershell.
Generate the autocompletion script for the fish shell.
Generate the autocompletion script for the zsh shell.
Generate the autocompletion script for the bash shell." | cut -d$'\n' -f1
`

	output, err := exec.Command(pc.ShellLocation, "-c", command).Output()
	outputString := string(output[:])
	if err != nil {
		fmt.Println(command)
		log.Fatal(string(outputString))
	}

	// if the output from the command starts with this string then shell completion
	// is mostly likely not guarenteed
	noCompletionPrefix := "Examples: `/foo` would allow `/foo`"

	if strings.TrimSpace(outputString) == "" || strings.HasPrefix(outputString, noCompletionPrefix) {
		return false
	} else {
		return true
	}
}

func (pc *PluginConfigImpl) GenerateCompletionScript() (string, error) {

	pluginList, err := pc.getAllKubectlPlugins()

	if err != nil {
		return "", err
	}

	plugins := strings.Join(pluginList, " ")

	initPluginCode := ""
	for _, plugin := range pc.Plugins {
		pluginUsesCobra := plugin.PluginSupportsCobraCompletion

		// source section
		if pluginUsesCobra {
			initPluginCode += pc.GetConfig().GetCompletionScriptSection(plugin.Name)
		}

		//completion function section

		if plugin.CompletionFunctionName != "" {
			switch pc.Shell {
			case "bash":
				initPluginCode += fmt.Sprintf("\tpluginCompletionFunction[%s]=\"%s\"\n", plugin.Name, plugin.CompletionFunctionName)
			case "zsh":
				initPluginCode += fmt.Sprintf("\tcompdef %s kubectl-%s\n", plugin.CompletionFunctionName, plugin.Name)
			}
		}

		// Plugin Description
		if pc.Shell == "zsh" {
			initPluginCode += fmt.Sprintf("\tpluginDescription[%s]=\"%s\"\n", plugin.Name, plugin.Description)
		}

		// Does use cobra section
		if pc.Shell == "zsh" {
			if strings.TrimSpace(plugin.Description) != "" {
				if pluginUsesCobra {
					initPluginCode += fmt.Sprintf("\tcobraSupported[%s]=true\n", plugin.Name)
				}
			}
		}

		initPluginCode += "\n"

	}

	switch pc.Shell {
	case "bash":
		return fmt.Sprintf(scripts.BashCompletionScript, initPluginCode, plugins, plugins, strings.Join(pc.KubectlOverridePlugins, " ")), nil
	case "zsh":
		return fmt.Sprintf(scripts.ZshCompletionScript, plugins, initPluginCode, strings.Join(pc.KubectlOverridePlugins, " ")), nil
	}
	return "", nil
}

func (pc PluginConfigImpl) GetCompletionScriptSection(plugin string) string {
	switch plugin {
	case "stern":
		return fmt.Sprintf("\tsource <(kubectl-%s --completion %s)\n", plugin, pc.Shell)
	}

	return fmt.Sprintf("\tsource <(kubectl-%s completion %s)\n", plugin, pc.Shell)

}

func isInList(list []string, searchString string) bool {
	for _, item := range list {
		if item == searchString {
			return true
		}
	}
	return false
}

func (pc PluginConfigImpl) PrintCobraPlugins() {
	var count = 1
	for _, plugin := range pc.Plugins {
		if plugin.PluginSupportsCobraCompletion {
			fmt.Printf("%d\t%s\n", count, plugin.Name)
			count++
		}
	}
}

func (pc PluginConfigImpl) PrintAllPlugins() {

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 8, 8, 0, '\t', tabwriter.AlignRight)

	defer w.Flush()

	fmt.Fprintf(w, "\n%s\t%s\t%s\t", "NAME", "DESCRIPTION", "COMPLETION SUPPORTED")

	for _, plugin := range pc.Plugins {
		supported := "yes"
		if !plugin.PluginSupportsCobraCompletion {
			supported = "no"
		}

		maxLength := 100
		if len(plugin.Description) > maxLength {
			plugin.Description = plugin.Description[0:maxLength] + "..."
		}
		fmt.Fprintf(w, "\n%s\t%s\t%s", plugin.Name, plugin.Description, supported)
	}
	fmt.Fprintf(w, "\n")

}

func (pc PluginConfigImpl) doesPluginCompletionOverrideKubectl(plugin string) bool {
	switch plugin {
	case "allctx":
		return true
	}
	return false
}
