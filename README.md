# writeme

**Take notes without breaking your flow.**  
A simple CLI tool to reword your commit messages, ideas, or thoughts using local or cloud-based AI models.

---

## Installation

1. Go to [Releases](https://github.com/yourusername/writeme/releases).
2. Download the binary for your platform.
3. Add the binary to your `PATH`.


## Requirements

If you want to use AI features, you’ll need one of:

* **[Ollama](https://ollama.com)** installed and running locally. (get llama3.1:latest)
* An **OpenAI API Key**.

Set your preferences and credentials in the config file. See [`config.template.yaml`](./config.template.yaml) for reference.

## Flow

1. `writeme create` --> will create a file named `NOTES.md`
2. `writeme config init` --> will create a `writeme` directory and a `config.yaml` file inside your system’s standard config location (e.g. `~/.config` on Linux/macOS, `%APPDATA%` on Windows).
3. `writeme config edit` --> will open up vi (or notepad) so you can edit the file. you can also just open this file with vscode or anyother editor.
4. `writeme note "message"` or `writeme note "message" -a`, the latter for AI rewording