package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

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

var cobraBlackList = []string{"cert_manager", "stern"}

type Plugin struct {
	Name                          string `yaml:"name"`
	CompletionFunctionName        string `yaml:"completionFunctionName"`
	Description                   string `yaml:"description"`
	PluginSupportsCobraCompletion bool   `yaml:"supportsCobraCompletion"`
}

type PluginConfigImpl struct {
	ConfigFile    string   `yaml:"configfile,omitempty"`
	Shell         string   `yaml:"shell"`
	ShellLocation string   `yaml:"shellLocation"`
	Plugins       []Plugin `yaml:"plugins"`
}

type PluginConfig interface {
	AddPlugin(plugin Plugin) error
	DisplayConfig() error
	DoesPluginExist(plugin string) bool
	DoesPluginUseCobra(plugin string) bool
	EditPlugin(pluginName string, completionFunctionName, description string, completionFunctionSet, descriptionSet bool) error
	GetPlugin(pluginName string) (*Plugin, error)

	GeneratePluginConfig() error
	GetCompletionFunctionName(plugin string) string

	GetConfig() *PluginConfigImpl
	WriteCleanConfig() error
	GenerateCompletionScript() (string, error)

	checkConfigFile() error

	createConfigFolder() error
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
			Description:                   fmt.Sprintf("A kubectl plugin called %s", plugin),
			PluginSupportsCobraCompletion: completionFunctionName != "",
		})

		if err != nil {
			return err
		}
	}

	pc.removeOldPlugins(plugins)
	return nil
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
	if !pc.DoesPluginUseCobra(plugin) {
		return ""
	}

	switch pc.Shell {
	case "zsh":
		return "_" + plugin
	case "bash":
		return "__start_" + plugin
	}
	return "unsupported"
}

func (pc *PluginConfigImpl) DoesPluginUseCobra(plugin string) bool {
	if inCobraBlackList(plugin) {
		return false
	}
	command := "grep ValidArgsFunction $(which kubectl-" + plugin + ") || true"
	output, err := exec.Command(pc.ShellLocation, "-c", command).Output()

	if err != nil {
		fmt.Println(command)
		log.Fatal(string(output[:]))
	}

	if strings.TrimSpace(string(output[:])) == "" {
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
			initPluginCode += fmt.Sprintf("\tsource <(kubectl-%s completion %s)\n", plugin.Name, pc.Shell)

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
		return fmt.Sprintf(scripts.BashCompletionScript, initPluginCode, plugins, plugins), nil
	case "zsh":
		return fmt.Sprintf(scripts.ZshCompletionScript, plugins, initPluginCode), nil
	}
	return "", nil
}

func inCobraBlackList(plugin string) bool {
	for _, pluginName := range cobraBlackList {
		if plugin == pluginName {
			return true
		}
	}
	return false

}
