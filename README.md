# Storage Management CLI

A simple tool to help find and clear out your biggest files

## Usage

Download the executable.
Run the executable in your terminal.

```bash
$./storage-management-cli
```

## Initial Arguments

Upon inital execution of the file, you can pass in these arguments:

`root` - By default the tool uses the current directory. Use the `root` flag to define a custom directory.

```bash
$./storage-management-cli -root /your/path/to/analyze
```

`ignore` - A comma seperated list of words to ignore in the search of files and directories. If you use more than one, make sure to wrap them in quotes

```bash
$./storage-management-cli -i .git
OR
$./storage-management-cli -i ".git, node_modules"
```

`resultCount` - By default the tool shows the top 10 results. Use the `resultCount` flag to define a custom # of results.

```bash
$./storage-management-cli -resultCount 25
```

## Usage

Once you fire up the executable, you can enter any of the commands:

`cd <PATH>` - Changes the directory to the string passed in

`delete <ID>` - Deletes the path indicated by the ID shown in the table above

`open <ID>` - Opens the file explorer of the path ID specified. If the path is a file, it opens the parent directory.

`more <Number>` - Changes the amount of results that are shown in the table

`exit` - Exits the CLI

## TODO

- Unit Tests
- CI automated tests and builds
- Better handling of larger file directories (Memory Crashes)

## Contributing

Open to suggestions and PRs!
