# go-linux

A CLI tool that translates natural language into Linux shell commands, helps users understand them, and optionally runs them.

## Project Explanation

`go-linux` (also known as `lihelp`) is designed to assist new Linux users by converting plain English descriptions into executable shell commands. It leverages AI (via Google's Generative AI) to generate accurate commands based on user queries, provides clear explanations, and offers the option to execute the commands directly. This tool is particularly useful for beginners who are learning Linux command-line operations.

### Features

- **Natural Language to Command Translation**: Input a description like "list all files in the current directory" and get the corresponding `ls` command.
- **Command Explanation**: Each generated command comes with a detailed explanation to help users learn.
- **Optional Execution**: Choose to run the command or just view it (with dry-run mode available).
- **Command Logging**: Keeps a history of generated commands for reference.
- **Monitoring Mode**: Run in monitor mode to continuously generate commands.

## Installation

1. Ensure you have Go installed (version 1.24 or later). You can download it from [golang.org](https://golang.org/dl/).
2. Clone the repository:
   ```bash
   git clone https://github.com/kairavb/go-linux.git
   cd go-linux
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Project

### Build the Project

To build the executable:

```bash
go build -o lihelp .
```

### Run Directly

You can run the tool directly using Go:

```bash
go run . "list all files in the current directory"
```

Or after building:

```bash
./lihelp "list all files in the current directory"
```

### Usage Examples

- Generate a command without running it:

  ```bash
  ./lihelp "show disk usage"
  ```

- Use dry-run mode (command not executed):

  ```bash
  ./lihelp --dry-run "create a new directory called test"
  ```

- Run in monitor mode (for continuous command generation):
  ```bash
  ./lihelp --run-monitor
  ```

### Flags

- `--dry-run` or `-d`: Generate and explain the command but do not execute it.
- `--run-monitor` or `-m`: Enable monitoring mode for ongoing command assistance.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License.
