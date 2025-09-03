# writeme

**Take notes without breaking your flow.**  
A simple CLI tool to reword your ideas and thoughts using local or cloud-based AI models.


## Example:
https://drive.google.com/file/d/1nVdbNO9QeiyPoisZ-hF4D4j4oIOAcUIL/view?usp=share_link


---

## Installation

1. **Download the latest release**  
   Visit the [Releases page](https://github.com/vishnuvgn/writeme/releases) and grab the binary for your operating system and CPU architecture.

2. **Install the binary**  
   - Move the downloaded file into a directory on your `PATH` (for example `/usr/local/bin`). 

3. **Verify the installation**  
   ```bash
   writeme --version
   # writeme version x.y.z
---

## Requirements

Writeme can reword text using a local LLM (Ollama) or via the OpenAI API. Choose one backend:

### 1. Ollama (free, local)

1. **Install Ollama**
   Follow the instructions at [ollama.com](https://ollama.com).

2. **Pull the model**

   ```bash
   ollama pull llama3.1:latest
   ```

3. **Start the Ollama server**

   ```bash
   ollama serve
   ```

### 2. OpenAI (cloud-based)

* **API key**
  Create or retrieve your API key on the [OpenAI dashboard](https://platform.openai.com/account/api-keys).

---

## Configuration (optional)

By default, writeme uses Ollama with the `llama3.1:latest` model—no changes required. If you’d rather use OpenAI:


1. Open your `config.yaml` (at `~/.config/writeme/config.yaml` on macOS/Linux—unless you’ve set `XDG_CONFIG_HOME`; on Windows: `%APPDATA%\writeme\config.yaml`). Or run:

   ```bash
   writeme config edit
to open it in your default editor (e.g., `vi`, Notepad, etc.).

2. Change the backend and fill in your OpenAI details:

   ```yaml
   llm:
     backend: openai

   openai:
     model: gpt-4o-mini
     api_key: YOUR_OPENAI_API_KEY
     system_prompt: "Your system prompt here"
## Flow

1. `writeme create`: will create a file named `NOTES.md`
2. `writeme config init`: will create a `writeme` directory and a `config.yaml` file inside your system’s standard config location (e.g. `~/.config` on Linux/macOS, `%APPDATA%` on Windows).
3. `writeme config edit`: will open up vi (or notepad) so you can edit the file. you can also just open this file with vscode or anyother editor.
4. `writeme note "message"` or `writeme note "message" -a`: the latter for AI rewording