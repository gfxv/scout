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
Follow these steps to set up and use **Scout**:
1) Prepare the `files` Directory
- Create a directory named `files` near binary executable `scout`.
- Place the documents you want to index inside this directory.
- Supported file types: `.md`, `.txt`, `.pdf`, `.xml`, `.html`, `.xhtml`.
2) Run **Scout**
- Run the binary executable like so `./scout`
- After, the indexing process should begin
3) Use **Scout** 
- When **Scout** finish indexing files, you will be able to access the web interface.
- Open your web browser and navigate to http://localhost:6969
4) Search Your Documents
- Enter your query in the search bar and press `Enter` or Search button.
- The results will display relevant documents from the `files` directory.

