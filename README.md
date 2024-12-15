# Repo Check

Repo Check is a CLI tool that lists out all local git repos in a directory along with
additional information such as the last modified date of the repo, whether the
repo is synced with remote etc.

![repocheck cli output](docs/demo.gif)

```
Usage:
  repocheck [flags]

Flags:
  -h, --help                  help for repocheck
  -L, --lastmodified string   Filter by last modified date of repo
                              options: yyyy-mm-dd | ">yyyy-mm-dd" | ">=yyyy-mm-dd"
                              note: surround any filters containing < or > with quotes
      --no-fetch              Run without doing a git fetch for each repo
  -s, --sort string           Sort results
                              options: lastmodified | name | path | synced (default "lastmodified")
  -S, --synced string         Filter by synced status of repo
                              options: y | n
  -t, --tsv                   Output as tab separated values
```
For more detailed usage instructions see [Usage](#usage)

## Installation

Pre-built packages for Windows, macOS, and Linux are found on the [Releases](https://github.com/bevane/repocheck/releases) page.

### Homebrew on macOS or Linux
```
brew tap bevane/tap
brew install repocheck
```

### Go install - multiplatform
Requires [go v1.22](https://go.dev/doc/install) or later

`go install github.com/bevane/repocheck@latest`

### Shell completions
Installing pre-built packages and brew package will install [completions](https://en.wikipedia.org/wiki/Command-line_completion) automatically.

For other cases including `go install`, if you want shell completions, you need to manually install it.

Manual Instructions:
1. Clone this repo `git clone https://github.com/bevane/repocheck`
2. Run `./scripts/completions.sh` to generate shell fragments for bash, zsh & fish in `/completions`
3. Copy the shell fragment for your shell into the appropriate directory for your shell.
   - Example:
    `mv completions/repocheck.bash /usr/share/bash-completion/completions/repocheck`
## Usage

### Basic Usage

Command help

`repocheck -h` or `repocheck --help`

Running repocheck without any args will list the repos in current directory

`repocheck`

Target directory can be passed in as an arg either as a relative path or an absolute path

`repocheck projects`

`repocheck /home/user/projects`

### Additional flags

#### No fetch
By default, repocheck will run a git fetch for each repo before determining
its sync status to ensure that the information is up to date.

Repocheck can be run without doing git fetch for each repo with the no fetch flag:

`repocheck --no-fetch`

#### Sort
Sort flag `-s` or `--sort` can be used to sort the results by a specific key

`repocheck --sort name` to sort by repo name

`repocheck -s synced` to sort by sync status of the repo - unsynced repos will be at the top

#### Filters
Supported filter flags:
- `-L` or `--lastmodified` - filter results by repos that were last modified on, before or after a certain date
- `-S` or `--synced` - filter results by synced status of repo

**Examples**

`repocheck --synced y` to only show repos that are synced

`repocheck --lastmodified 2024-01-01` to only show repos that were last modified on 2024-01-01

`repocheck --lastmodified ">=2024-01-01"` to only show repos that were last modified on or later than 2024-01-01

Multiple filters can be combined:

`repocheck -L "<2024-01-01" -S n` to only show unsynced repos that were last modified before 2024-01-01

*Note: for options containing '<' or '>' surround the entire query with quotes to prevent them from being interpreted as operators by bash*

#### Output formatting
By default, repocheck will output the results in a pretty human-readable table.
Repocheck also supports output flags to change the output format

Supported output flags:
- `-t` or `--tsv` - to output results as tab separated values that are machine-readable

**Examples**

Machine-readable output can be piped to other command line utilities:

`repocheck --tsv | cut -f2` to show only the second column of the results i.e the path data for each repo

`repocheck --tsv | grep exercises` to only show lines containing "exercises"

# Contact

Submit an [issue](https://github.com/bevane/repocheck/issues/new)

[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/bevane50)](https://x.com/bevane50)

# Contributions

Contributions are welcome. Please fork the repo and open pull requests to contribute. Ensure that your code passes the tests and include tests for any changes where applicable.

# License
Repo Check is released under the MIT License
