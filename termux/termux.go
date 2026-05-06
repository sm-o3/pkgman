package termux

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

type Package struct {
	Name         string
	Category     string
	Version      string
	Architecture string
	Installed    bool
}

// SearchPackages searches for packages matching the query.
func SearchPackages(query string) ([]byte, error) {
	cmd := exec.Command("pkg", "search", query)
	return cmd.CombinedOutput()
}

// GetCommand returns an exec.Cmd for the given operation.
func GetCommand(operation string, args ...string) *exec.Cmd {
	var cmdArgs []string
	switch operation {
	case "install":
		cmdArgs = append([]string{"install", "-y"}, args...)
	case "remove":
		cmdArgs = append([]string{"uninstall", "-y"}, args...)
	case "reinstall":
		cmdArgs = append([]string{"reinstall", "-y"}, args...)
	case "update":
		cmdArgs = []string{"upgrade", "-y"}
	case "clean":
		cmdArgs = []string{"clean"}
	case "autoremove":
		// apt autoremove is better in termux sometimes, but pkg autoremove might work in newer versions
		cmdArgs = []string{"autoremove", "-y"}
		return exec.Command("apt", cmdArgs...)
	case "repo":
		return exec.Command("termux-change-repo")
	default:
		return nil
	}
	return exec.Command("pkg", cmdArgs...)
}

// ListAllPackages returns a list of all packages available and their categories.
func ListAllPackages() ([]Package, error) {
	cmd := exec.Command("apt", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []Package
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Listing...") || strings.HasPrefix(line, "WARNING:") || strings.TrimSpace(line) == "" {
			continue
		}
		// Format: 0verkill/stable 1:0.16-1 aarch64 [installed]
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		nameCat := strings.SplitN(parts[0], "/", 2)
		name := nameCat[0]
		category := "unknown"
		if len(nameCat) > 1 {
			catParts := strings.Split(nameCat[1], ",")
			category = catParts[0]
		}

		version := parts[1]
		arch := parts[2]
		installed := strings.Contains(line, "[installed")

		packages = append(packages, Package{
			Name:         name,
			Category:     category,
			Version:      version,
			Architecture: arch,
			Installed:    installed,
		})
	}
	return packages, nil
}
