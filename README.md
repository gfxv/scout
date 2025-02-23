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
### Configuration
| CLI Flags | Environment  | Type | Default | Description |
| --------- | ------------ | ---- | ------- | ----------- |
| `-index` | `SCOUT_INDEX`   | `bool`   | `false`        | Index files in the specified directory and exit (`-files` flag required). |
| `-serve` | `SCOUT_SERVE`   | `bool`   | `false`        | Start the web server. |
| `-files` | `SCOUT_FILES`   | `string` | *empty string* | Directory path containing files to index (required with `-index`). |
| `-port`  | `SCOUT_PORT`    | `string` | `"6969"`       | Port to listen on when serving (e.g., 8080). |
| `-db`    | `SCOUT_DB_PATH` | `string` | `"meta.db"`  | Path to the SQLite database file used to store the search index. |

**Note:**
- Either `-index` or `-serve` must be specified.
- `-index` and `-serve` cannot be used together.
- When using `-index`, the `-files` flag is required.
- Flags override environment variables if both are provided.

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
    - Access the search interface at http://localhost:6969.

### Running with Docker Compose
#### Sample docker-compose file
```yml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "6969:6969"          # map host port to container port (adjustable via SCOUT_PORT)
    volumes:
      - ./files:/app/files   # mount local files for indexing
      - ./data:/app          # persist the database (meta.db)
    environment:
      - SCOUT_INDEX=false    # set to true to index files
      - SCOUT_SERVE=true     # set to true to serve the web interface
      - SCOUT_FILES=/app/files
      - SCOUT_PORT=6969
      - SCOUT_DB_PATH=/app/meta.db
    restart: unless-stopped
```

#### Setup
Follow these steps to run Scout with Docker Compose. 
You’ll first index your files (if it’s the first run or you have new files), then serve the indexed files.

1) Index Files (First Run or New Files)
    - Update `docker-compose.yml` to enable indexing
    ```yml
    environment:
      - SCOUT_INDEX=true    # enable indexing
      - SCOUT_SERVE=false   # disable serving
    ```
    - Run the indexing process:
    ```bash
    docker compose up
    ```
    - This indexes files from `./files` into `./data/meta.db` and exits. Check the output for "Indexing took: ..." to confirm completion.
2) Serve Indexed Files
    - Update `docker-compose.yml` to enable serving:
    ```yml
    environment:
      - SCOUT_INDEX=false    # disable indexing
      - SCOUT_SERVE=true     # enable serving
    ```
    - Start the web server:
    ```bash
    docker compose up -d
    ```
    - Access the search interface at http://localhost:6969

**Note:**
- The `meta.db` in `./data` persists across runs. Re-index only when adding new files.
- Use `docker compose down -v` to remove volumes (`./data` and `./files`) if you want a fresh start.
