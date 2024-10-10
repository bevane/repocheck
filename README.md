# Repo Check

Repo Check is a CLI tool that lists out all local git repos in a directory along with
additional information such as the last modified date of the repo, whether the
repo is synced with remote etc.

![repocheck cli output](docs/demo.png)

## Installation

Install go v1.22 or later - [Official installation instructions](https://go.dev/doc/install)

To download and install repocheck, run:

`go install github.com/bevane/repocheck@latest`

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

# Discussions

Interested in this project and have requests for specific features or want to discuss about this repo?
Let me know in [discussions](https://github.com/bevane/repocheck/discussions)


# License
Repo Check is released under the MIT License
