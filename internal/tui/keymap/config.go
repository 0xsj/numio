// internal/tui/keymap/config.go

package keymap

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the keybinding configuration file structure.
type Config struct {
	Normal   map[string]string `toml:"normal"`
	Insert   map[string]string `toml:"insert"`
	Visual   map[string]string `toml:"visual"`
	Operator map[string]string `toml:"operator"`
}

// DefaultConfigPath returns the default config file path.
func DefaultConfigPath() string {
	// Try XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "numio", "keybindings.toml")
	}

	// Fall back to ~/.config
	home, err := os.UserHomeDir()
	if err != nil {
		return "keybindings.toml"
	}

	return filepath.Join(home, ".config", "numio", "keybindings.toml")
}

// LoadConfig loads keybindings from a TOML file.
func LoadConfig(path string) (*Config, error) {
	var config Config

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves keybindings to a TOML file.
func SaveConfig(path string, config *Config) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header comment
	header := `# Numio Keybindings Configuration
# 
# Format: key = "action"
#
# Special keys:
#   ctrl+x, alt+x, shift+x - Modifier combinations
#   enter, tab, backspace, delete, esc - Special keys
#   up, down, left, right - Arrow keys
#   home, end, pgup, pgdown - Navigation keys
#   f1-f12 - Function keys
#   space - Space bar
#
# Multi-key sequences:
#   "gg" = "goto_top"
#   "dd" = "delete_line"
#
# Available actions:
#   Mode switching: normal_mode, insert_mode, append_mode, visual_mode
#   Movement: move_up, move_down, move_left, move_right
#            move_word_next, move_word_prev
#            goto_line_start, goto_line_end, goto_top, goto_bottom
#            page_up, page_down
#   Editing: delete_char, delete_line, delete_to_end
#           yank_line, yank, paste, paste_above
#           undo, redo, join_lines
#           open_below, open_above
#   Insert: insert_char, insert_newline, backspace, delete, insert_tab
#   Operators: operator_delete, operator_yank, operator_change
#   General: quit, force_quit, save, save_quit
#           toggle_help, refresh_rate
#           toggle_line_numbers, toggle_wrap

`
	if _, err := file.WriteString(header); err != nil {
		return err
	}

	// Encode TOML
	encoder := toml.NewEncoder(file)
	return encoder.Encode(config)
}

// ApplyConfig applies a config to a KeyMap.
func (km *KeyMap) ApplyConfig(config *Config) {
	if config.Normal != nil {
		for key, actionStr := range config.Normal {
			action := ParseAction(actionStr)
			if action != ActionNone {
				km.Normal.Bind(key, action)
			}
		}
	}

	if config.Insert != nil {
		for key, actionStr := range config.Insert {
			action := ParseAction(actionStr)
			if action != ActionNone {
				km.Insert.Bind(key, action)
			}
		}
	}

	if config.Visual != nil {
		for key, actionStr := range config.Visual {
			action := ParseAction(actionStr)
			if action != ActionNone {
				km.Visual.Bind(key, action)
			}
		}
	}

	if config.Operator != nil {
		for key, actionStr := range config.Operator {
			action := ParseAction(actionStr)
			if action != ActionNone {
				km.Operator.Bind(key, action)
			}
		}
	}
}

// ToConfig converts a KeyMap to a Config.
func (km *KeyMap) ToConfig() *Config {
	config := &Config{
		Normal:   make(map[string]string),
		Insert:   make(map[string]string),
		Visual:   make(map[string]string),
		Operator: make(map[string]string),
	}

	for _, b := range km.Normal.AllBindings() {
		config.Normal[b.Key] = string(b.Action)
	}

	for _, b := range km.Insert.AllBindings() {
		config.Insert[b.Key] = string(b.Action)
	}

	for _, b := range km.Visual.AllBindings() {
		config.Visual[b.Key] = string(b.Action)
	}

	for _, b := range km.Operator.AllBindings() {
		config.Operator[b.Key] = string(b.Action)
	}

	return config
}

// LoadFromFile loads keybindings from a file, applying them on top of defaults.
func (km *KeyMap) LoadFromFile(path string) error {
	config, err := LoadConfig(path)
	if err != nil {
		return err
	}

	km.ApplyConfig(config)
	return nil
}

// SaveToFile saves current keybindings to a file.
func (km *KeyMap) SaveToFile(path string) error {
	config := km.ToConfig()
	return SaveConfig(path, config)
}

// LoadOrCreate loads keybindings from file, or creates default config if not exists.
func LoadOrCreate(path string) (*KeyMap, error) {
	km := Default()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config file
		if err := km.SaveToFile(path); err != nil {
			// Non-fatal - just use defaults
			return km, nil
		}
		return km, nil
	}

	// Load existing config
	if err := km.LoadFromFile(path); err != nil {
		// Non-fatal - just use defaults
		return km, nil
	}

	return km, nil
}

// DefaultConfig returns a Config with default bindings.
func DefaultConfig() *Config {
	km := Default()
	return km.ToConfig()
}

// MergeConfig merges a config into another, overwriting existing bindings.
func MergeConfig(base, overlay *Config) *Config {
	result := &Config{
		Normal:   make(map[string]string),
		Insert:   make(map[string]string),
		Visual:   make(map[string]string),
		Operator: make(map[string]string),
	}

	// Copy base
	for k, v := range base.Normal {
		result.Normal[k] = v
	}
	for k, v := range base.Insert {
		result.Insert[k] = v
	}
	for k, v := range base.Visual {
		result.Visual[k] = v
	}
	for k, v := range base.Operator {
		result.Operator[k] = v
	}

	// Apply overlay
	for k, v := range overlay.Normal {
		result.Normal[k] = v
	}
	for k, v := range overlay.Insert {
		result.Insert[k] = v
	}
	for k, v := range overlay.Visual {
		result.Visual[k] = v
	}
	for k, v := range overlay.Operator {
		result.Operator[k] = v
	}

	return result
}

// ValidateConfig checks if all actions in a config are valid.
func ValidateConfig(config *Config) []string {
	var errors []string

	validateMode := func(name string, bindings map[string]string) {
		for key, actionStr := range bindings {
			action := ParseAction(actionStr)
			if action == ActionNone && actionStr != "" {
				errors = append(errors, name+": unknown action '"+actionStr+"' for key '"+key+"'")
			}
		}
	}

	validateMode("normal", config.Normal)
	validateMode("insert", config.Insert)
	validateMode("visual", config.Visual)
	validateMode("operator", config.Operator)

	return errors
}