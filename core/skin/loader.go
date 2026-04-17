package skin

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

//go:embed builtin/default.html
var builtinDefault []byte

//go:embed builtin/minimal.html
var builtinMinimal []byte

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

// SeedBuiltinSkins copies default.html and minimal.html into the skins folder
// only if the folder is empty — first launch only.
func (s *SkinService) SeedBuiltinSkins() error {
	entries, err := os.ReadDir(s.customDir)
	if err != nil {
		return err
	}

	// only seed if folder has no .html files
	for _, e := range entries {
		if strings.HasSuffix(strings.ToLower(e.Name()), ".html") {
			return nil // already has skins, don't overwrite
		}
	}

	skins := map[string][]byte{
		"default.html": builtinDefault,
		"minimal.html": builtinMinimal,
	}

	for name, content := range skins {
		path := filepath.Join(s.customDir, name)
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}

	return nil
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
	if strings.Contains(id, "/") || strings.Contains(id, "\\") || strings.Contains(id, "..") {
		return ""
	}
	path := filepath.Join(s.customDir, id+".html")
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}
