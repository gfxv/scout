# Scout

> [!NOTE]
> **Scout** is currently under development, so bugs and breaking changes can happen.
> I highly appreciate any feedback!

**Scout** is a lightweight and easy to use local search engine written in Go, that 
allows to index and search across user-provided documents using the TF-IDF algorithm.

## Building
To build **Scout** from source:
1) If not installed, install **go 1.23.3 or newer** on your machine. You can get **go** from [the official website](https://go.dev/doc/install).
2) If not installed, install **Taskfile**. You can get **Taskfile** from [the official website](https://taskfile.dev/installation/).
3) Clone this repo: `git clone https://github.com/gfxv/scout.git`
4) In the `scout` directory run `task` or `task build`, which will result in a binary file named `scout`.

## Usage
### Command-Line Flags
| Flag | Type | Default | Description |
| ---- | ---- | ------- | ----------- |
| `-index` | `bool` | `false` | Index files in the specified directory and exit (`-files` flag required). |
| `-serve` | `bool` | `false` | Start the web server. |
| `-files` | `string` | *empty string* | Directory path containing files to index (required with `-index`). |
| `-port` | `string` | `"6969"` | Port to listen on when serving (e.g., 8080). |
| `-db` | `string` | `"search.db"` | Path to the SQLite database file used to store the search index. |

**Note:**
- Either `-index` or `-serve` must be specified.
- `-index` and `-serve` cannot be used together.
- When using `-index`, the `-files` flag is required.

### Using Scout
1) Prepare your documents
    - Create a directory (e.g., `docs`) and place the documents you want to index inside it.
    - Supported file types: `.md`, `.txt`, `.pdf`, `.xml`, `.html`, `.xhtml`.
    - Ensure the directory is accessible from where you run **Scout**.
2) Index the documents
    - Run Scout with the -index flag and specify the directory using -files:
    ```bash
    scout -index -files ./docs
    ```
    - This builds a search index in the SQLite database (default: `meta.db`).
3) Start the Web Server
    - Launch the web interface with the `-serve` flag:
    ```bash
    scout -serve
    ```
    - Access it by opening your browser to http://localhost:6969.

