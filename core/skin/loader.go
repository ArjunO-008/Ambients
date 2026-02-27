package skin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Skin struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsCustom bool   `json:"isCustom"`
	Path     string `json:"path"`
}

type SkinService struct {
	customDir string
}

func NewSkinService() *SkinService {
	return &SkinService{
		customDir: customSkinsDir(),
	}
}

func (s *SkinService) EnsureCustomDir() error {
	return os.MkdirAll(s.customDir, 0755)
}
func (s *SkinService) GetCustomDir() string {
	return s.customDir
}

func (s *SkinService) ListCustomSkin() string {
	entries, err := os.ReadDir(s.customDir)
	if err != nil {
		return "[]"
	}
	var skins []Skin

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".html") {
			continue
		}

		displayName := strings.TrimSuffix(name, filepath.Ext(name))

		skins = append(skins, Skin{
			ID:       displayName,
			Name:     displayName,
			IsCustom: true,
			Path:     filepath.Join(s.customDir, name),
		})
	}

	if skins == nil {
		return "[]"
	}

	b, _ := json.Marshal(skins)
	return string(b)
}

func (s *SkinService) ReadCustomSkin(id string) string {
	path := filepath.Join(s.customDir, id+".html")

	if strings.Contains(id, "/") || strings.Contains(id, "\\") || strings.Contains(id, "..") {
		return ""
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(content)
}
