# oh-my-gpt-correction

A Zsh helper to auto-correct the current command via ChatGPT.

Just compile the go code to executable. Code is written _badly_ but it does its job:

- Looks for an API key in CHATGPT_API_KEY environment variable
- It accepts one argument, intended to get a whole command inside this argument
- Exit codes are:
    - 1: Chat completion error (remote or lib-wise issues)
    - 2: No API key given
    - 3: Can't find an issue with given command
    - 4: No command given
- In case of a suggestion, it prints 2 lines:
    - First line is colored diff from original command
    - Second line is the original suggestion
- We are catching this output on [ZSH part](./example.zsh) and utilise accordingly

### Demo

![A simple demo](./demo.gif)
