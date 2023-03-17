export GPT_ERROR_LOGFILE=/tmp/gptresult
export OH_MY_GPT_CORRECTION_PATH=/wherever/you/put/the/executable
export CHATGPT_API_KEY="your_api_key"

gpt-command-line() {
    [[ -z $BUFFER ]] && zle up-history
    # Run and collect possible errors
    GPTRESULTARRAY=("${(@f)$($OH_MY_GPT_CORRECTION_PATH ''$LBUFFER'' 2>> $GPT_ERROR_LOGFILE)}")
    # We're using return code
    GPTRESCODE="$?"
    if [[ $GPTRESCODE = 1 ]]; then
        BUFFER="$BUFFER # ❓ChatGPT response error, check /tmp/gptresult"
    elif [[ $GPTRESCODE = 2 ]]; then
        BUFFER="$BUFFER # 🔑 API key error"
    elif [[ $GPTRESCODE = 3 ]]; then
        BUFFER="$BUFFER # 🤖 it looks fine to me?"
    elif [[ $GPTRESCODE = 4 ]]; then
        BUFFER="$BUFFER # ❓No prompt given"
    else
        # Keep the old buffer for now
        OLDBUFFER="${BUFFER}"
        zle .reset-prompt
        echo
        # Ask user
        GPTPROMPTTXT=" >🤖 ChatGPT suggest following, does this look better? [y/N]\n
        $GPTRESULTARRAY[1]
        "
        echo -e "$GPTPROMPTTXT"
        read -k 1 GPTCHOICEANSWER
        if [[ $GPTCHOICEANSWER = y ]]; then
            # Replace the current prompt with the suggestion
            zle .reset-prompt
            BUFFER="${GPTRESULTARRAY[2]} # 🚀"
        else
            # Reset to old prompt
            zle .reset-prompt
            BUFFER="${OLDBUFFER} # Restored original buffer"
        fi
    fi
}
zle -N gpt-command-line
# Defined shortcut key: [Esc] o
bindkey "^[o" gpt-command-line

