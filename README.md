# PersonalCLI

PersonalCLI is a command-line interface (CLI) tool designed to help you manage your daily tasks, notes, calendar events, and weather information directly from your terminal. It's built with Go and aims to be a fast, efficient, and convenient personal assistant for developers and power users.

## Features

PersonalCLI comes with several useful commands to streamline your workflow:

### 📝 To-Do List (`personalcli todo`)

Manage your tasks with ease.
*   **Add tasks:** `personalcli todo add "Buy groceries"`
*   **List tasks:** `personalcli todo list`
*   **Mark tasks as done:** `personalcli todo done <task_id>`
*   **Clear all tasks:** `personalcli todo clear`

### 🗒️ Notes (`personalcli note`)

Keep track of your thoughts and ideas.
*   **Create a new note:** `personalcli note new "Meeting agenda"`
*   **List all notes:** `personalcli note list`
*   **Find notes by keyword:** `personalcli note find "important"`

### 🗓️ Calendar (`personalcli calendar`)

View your upcoming Google Calendar events.
*   **View upcoming events:** `personalcli calendar events`
    *   (Requires initial OAuth 2.0 authentication through your web browser)

**Setting Up Google Calendar API:**
1.  Go to the [Google Cloud Console](https://console.cloud.google.com/).
2.  Create a new project or select an existing one.
3.  Enable the Google Calendar API for your project.
4.  Navigate to "Credentials" and click "Create Credentials" > "OAuth 2.0 Client IDs".
5.  For "Application type", select "Desktop application".
6.  Download the credentials file (JSON) and rename it to `credentials.json`.
7.  Place this file in `~/.config/personalcli/` directory (the tool will prompt you to authenticate the first time).
8.  Run `personalcli calendar events` and follow the authentication flow in your browser.

### ☀️ Weather (`personalcli weather`)

Get current weather information for any location.
*   **Get weather by zip code:** `personalcli weather -z 90210`
*   **Get weather by city name:** `personalcli weather -l "London"`
*   **Get weather by city and state/country:** `personalcli weather -l "Mount Prospect, IL"`
    *   Displays temperature in Fahrenheit and Celsius, along with conditions, humidity, and wind speed.
    *   **Requires an OpenWeatherMap API key.** You can provide it using the `--api-key` flag or by setting the `WEATHER_API_KEY` environment variable. Get a free API key from [OpenWeatherMap](https://openweathermap.org/).

### Comprehensive Usage Examples

#### Todo List Examples:
```bash
# Add a new task
personalcli todo add "Complete project report"

# List all tasks
personalcli todo list

# Mark task #3 as completed
personalcli todo done 3

# Clear all tasks
personalcli todo clear
```

#### Notes Examples:
```bash
# Create a new note
personalcli note new "Meeting notes: discussed project timeline"

# List all notes
personalcli note list

# Find notes containing specific keyword
personalcli note find "meeting"
```

#### Calendar Examples:
```bash
# View upcoming events
personalcli calendar events
```

#### Weather Examples:
```bash
# Get weather by city
personalcli weather -l "New York"

# Get weather by city and state
personalcli weather -l "Chicago, IL"

# Get weather by zip code
personalcli weather -z 60601

# Provide API key via command line
personalcli weather -l "London" --api-key your_api_key_here
```

## Installation

To get PersonalCLI up and running on your system:

1.  **Install Go:** If you don't have Go installed, download and install it from [golang.org/doc/install](https://golang.org/doc/install).
2.  **Clone the Repository:**
    ```bash
    git clone https://github.com/your-username/PersonalCLI.git
    cd PersonalCLI
    ```
3.  **Build the Executable:**
    ```bash
    go build -o personalcli ./cmd/personalcli
    ```
    This will create an executable named `personalcli` in your project's root directory.

4.  **Add to your PATH (Optional, Recommended):**
    To run `personalcli` from any directory, add its location to your system's `PATH`. For example, you can move it to `/usr/local/bin`:
    ```bash
    sudo mv personalcli /usr/local/bin/
    ```
    Or, if you prefer to keep it in your home directory, add this to your shell's configuration file (e.g., `.bashrc`, `.zshrc`):
    ```bash
    export PATH=$PATH:/path/to/your/PersonalCLI
    ```
    (Replace `/path/to/your/PersonalCLI` with the actual path where you built the `personalcli` executable). Remember to `source` your shell config file after editing.

## Usage

Simply run `personalcli` followed by the command you want to use. Refer to the "Features" section for command-specific usage examples.

For help with any command, you can use the `--help` flag:
```bash
personalcli --help
personalcli weather --help
```

## Troubleshooting

### Common Issues

**Google Calendar Authentication Issues:**
*   If the authentication page doesn't open automatically, manually visit the URL shown in the terminal.
*   Make sure your `credentials.json` file is saved in the correct location: `~/.config/personalcli/credentials.json`
*   If you encounter OAuth errors, delete the `token.json` file in `~/.config/personalcli/` and try again.

**Weather API Issues:**
*   Ensure your OpenWeatherMap API key is valid and properly formatted.
*   Set the `WEATHER_API_KEY` environment variable correctly:
  ```bash
  export WEATHER_API_KEY="your_actual_api_key_here"
  ```
*   Check that your API key has the necessary permissions enabled in the OpenWeatherMap dashboard.

**General Setup Issues:**
*   If you get "command not found" errors, make sure you've added PersonalCLI to your PATH or are running the full path to the executable.
*   Ensure you have proper permissions to create files in your home directory's `.config` folder.
*   If configuration files aren't being created, check that you have write permissions to `~/.config/` directory.

### Data Storage
*   PersonalCLI stores tasks in `~/.config/personalcli/tasks.json`
*   Notes are stored in `~/.config/personalcli/notes.json`
*   Google Calendar authentication token is stored in `~/.config/personalcli/token.json`
*   Google Calendar credentials should be in `~/.config/personalcli/credentials.json`

## Contributing

We welcome contributions! If you have ideas for new features, bug fixes, or improvements, please feel free to open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
