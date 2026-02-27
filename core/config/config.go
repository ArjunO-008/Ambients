package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	Use24Hour        bool   `json:"use24Hour"`
	ActiveSkinID     string `json:"activeSkinID"`
	ActiveSkinCustom bool   `json:"activeSkinCustom"`
	BackgroundPath   string `json:"backgroundPath"`
	BackgroundType   string `json:"backgroundType"`
}

type ConfigService struct {
	path string
}

func NewConfigService() *ConfigService {
	dir := configDir()
	return &ConfigService{
		path: filepath.Join(dir, "settings.json"),
	}
}

func (c *ConfigService) EnsureDir() error {
	return os.MkdirAll(filepath.Dir(c.path), 0755)
}

func (c *ConfigService) Load() Settings {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return Settings{
			Use24Hour:    false,
			ActiveSkinID: "default",
		}
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return Settings{
			Use24Hour:    false,
			ActiveSkinID: "default",
		}
	}
	return s
}

func (c *ConfigService) Save(s Settings) error {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}
