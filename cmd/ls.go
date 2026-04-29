package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const storageFileName = "storage.json"

type storageJSON struct {
	LastKnownMenubarData menubarData `json:"lastKnownMenubarData"`
}

type menubarData struct {
	Menus menus `json:"menus"`
}

type menus struct {
	File fileMenu `json:"File"`
}

type fileMenu struct {
	Items []menuItem `json:"items"`
}

type menuItem struct {
	ID      string    `json:"id"`
	Label   string    `json:"label"`
	Submenu *submenu  `json:"submenu,omitempty"`
	URI     *uriEntry `json:"uri,omitempty"`
}

type submenu struct {
	Items []menuItem `json:"items"`
}

type uriEntry struct {
	Path string `json:"path"`
}

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List directories opened in VS Code",
		RunE:  runLs,
	}
	return cmd
}

func runLs(cmd *cobra.Command, args []string) error {
	storagePath, err := getStoragePath()
	if err != nil {
		return fmt.Errorf("failed to get storage path: %w", err)
	}

	data, err := os.ReadFile(storagePath)
	if err != nil {
		return fmt.Errorf("failed to read storage file: %w", err)
	}

	var storage storageJSON
	if err := json.Unmarshal(data, &storage); err != nil {
		return fmt.Errorf("failed to parse storage.json: %w", err)
	}

	dirs := extractRecentFolders(storage)

	if len(dirs) == 0 {
		return fmt.Errorf("no recent folders found in VS Code storage")
	}

	for _, dir := range dirs {
		fmt.Println(dir)
	}

	return nil
}

func extractRecentFolders(storage storageJSON) []string {
	var dirs []string

	for _, item := range storage.LastKnownMenubarData.Menus.File.Items {
		if item.ID == "submenuitem.MenubarRecentMenu" && item.Submenu != nil {
			for _, subItem := range item.Submenu.Items {
				if subItem.ID == "openRecentFolder" && subItem.URI != nil && subItem.URI.Path != "" {
					dir := filepath.Clean(subItem.URI.Path)
					dirs = append(dirs, dir)
				}
			}
		}
	}

	return dirs
}

func getStoragePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	storagePath := filepath.Join(
		home,
		"Library",
		"Application Support",
		"Code",
		"User",
		"globalStorage",
		storageFileName,
	)

	return storagePath, nil
}