package scripts

var BashCompletionScript = `
declare -A pluginCompletionFunction

init_plugins(){

    source <(kubectl completion bash)
        
%s    
}

init_plugins 

# The __start_kubectl function is the completion function for the kubectl.
# This function overrides the __start_kubectl function and adds the necessary code to allow 
# for autocompletion of kubectl plugin commands. The plugin autocompletion code can 
# be found between the 'START PLUGIN CODE' and 'END PLUGIN CODE' comments.

__start_kubectl()
{
    local cur prev words cword split
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __kubectl_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("kubectl")
    local command_aliases=()
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function
    local last_command
    local nouns=()
    local noun_aliases=()

    ### START PLUGIN CODE ###

    
    if [[ "$COMP_CWORD" == "1" ]]; then
        COMPREPLY+=( $(compgen -W "%s" -- ${COMP_WORDS[1]}) )
    fi

    if [[ "$COMP_CWORD" > "1" ]]; then
        __is_kubectl_plugin_command ${COMP_WORDS[1]}
        if [[ "$?" == "0" ]]; then 
            __does_kubectl_plugin_override_kubectl ${COMP_WORDS[1]}
            if [[  "$?" == "0" ]]; then 
                del_element=1; 
                words=( "${words[@]:0:$((del_element))}" "${words[@]:$((del_element+1))}" )
                cword=$((cword-1))
                __kubectl_handle_word
                return
            fi
            local plugin
            plugin=${COMP_WORDS[1]}
            if [[ ! -z ${pluginCompletionFunction[$plugin]} ]]; then # make sure description function was specified
                if [[ $(type -t ${pluginCompletionFunction[$plugin]}) == "function" ]]; then # make sure completion function exists
                    ${pluginCompletionFunction[$plugin]}
                    return
                else
                    return
                fi
            
            else
                return
            fi
        fi
    fi

    ### END PLUGIN CODE ###
    
    __kubectl_handle_word
}

__is_kubectl_plugin_command() {

    local kubectl_plugins=(%s)

    current_word=$1
    local plugin
    for plugin in "${kubectl_plugins[@]}"
        do
            if [[ "$current_word" == "$plugin" ]]; then
                return 0
            fi
        done
      return 1
}

__does_kubectl_plugin_override_kubectl() {

    local kubectl_override_plugins=(%s)

    current_word=$1
    local plugin
    for plugin in "${kubectl_override_plugins[@]}"
        do
            if [[ "$current_word" == "$plugin" ]]; then
                return 0
            fi
        done
      return 1
}

`
