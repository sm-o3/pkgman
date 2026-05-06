# Pkgman: Termux TUI Package Manager

Pkgman is a terminal user interface (TUI) application built specifically for Termux on Android. It acts as an interactive wrapper around Termux's standard `pkg` and `apt` commands, providing a sleek, categorized, and searchable interface for managing your installed software.

## Features

- **Categorized Package Browser**: View all available and installed packages neatly organized by their repository source (`Stable`, `X11`, `Root`).
- **Fuzzy Search**: Instantly find packages by typing (`/`) within the browser.
- **Context-Aware Actions**: Selecting a package intelligently offers context-relevant actions:
  - If installed: `Update`, `Reinstall`, `Remove`
  - If not installed: `Install`
- **Multi-Select Bulk Operations**: Press `Space` to select multiple packages. Press `Enter` to execute bulk actions (Install, Update, Reinstall, Remove) across all selected items.
- **System Maintenance**: Includes easy shortcuts for full system updates (`pkg upgrade`), cache cleaning (`pkg clean`), and removing orphaned dependencies (`apt autoremove`).

## Installation / Building

To build the project yourself, you will need [Go](https://golang.org/) installed either on your host machine or directly inside Termux.

### Building directly in Termux (Recommended)
1. Ensure Go is installed: `pkg install golang`
2. Clone this repository and navigate to its folder.
3. Build the binary:
   ```bash
   go build -o pkgman
   ```
4. Make it executable and (optionally) move it to your path:
   ```bash
   chmod +x pkgman
   mv pkgman $PREFIX/bin/
   ```

### Cross-Compiling for Termux (From a Linux/macOS Host)
If you are building from a standard computer to send to an Android phone:
1. Navigate to the project root.
2. Compile targeting Android and ARM64:
   ```bash
   GOOS=android GOARCH=arm64 go build -o pkgman
   ```
3. Transfer the `pkgman` binary to your device via `scp`, `adb`, or USB.
4. Execute `./pkgman` within your Termux shell.

## Usage

Simply run the compiled binary inside Termux:

```bash
./pkgman
```
- **Arrows (`↑`/`↓` or `k`/`j`)**: Navigate lists and menus.
- **`Enter`**: Select an option or trigger actions.
- **`Tab` / `Shift+Tab`**: Switch category tabs in the Package Browser.
- **`Space`**: Select multiple items for Bulk Operations.
- **`/`**: Begin fuzzy searching in the Package Browser.
- **`Esc` / `q`**: Go back or quit the application.

## Requirements
- Designed exclusively for the Termux environment.
- Requires standard Termux tools (`pkg`, `apt`) to function correctly.
