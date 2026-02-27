package bridge

import (
	"Ambients/core/audio"
	"Ambients/core/clock"
	"Ambients/core/config"
	"Ambients/core/media"
	"Ambients/core/skin"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type Bridge struct {
	ctx context.Context

	clockService  *clock.ClockService
	audioService  *audio.AudioService
	mediaService  *media.MediaService
	SkinService   *skin.SkinService
	configService *config.ConfigService

	once sync.Once
}

func NewBridge() *Bridge {
	return &Bridge{}
}

// Called by Wails once context is ready
func (b *Bridge) SetContext(ctx context.Context) {
	b.once.Do(func() {
		b.ctx = ctx
		b.clockService = clock.NewClockService(ctx)
		b.audioService = audio.NewAudioService(ctx)
		b.mediaService = media.NewMediaService()
		b.SkinService = skin.NewSkinService()
		b.configService = config.NewConfigService()

		b.SkinService.EnsureCustomDir()
		b.configService.EnsureDir()
	})
}

// CLOCK
func (b *Bridge) StartClock(user24Hour bool) {
	if b.clockService != nil {
		b.clockService.Start(user24Hour)
	}
}

func (b *Bridge) StopClock() {
	if b.clockService != nil {
		b.clockService.Stop()
	}
}

// AUDIO
func (b *Bridge) StartAudio() string {
	if b.audioService == nil {
		return "audio service not initialized"
	}

	if err := b.audioService.Start(); err != nil {
		return fmt.Sprintf("audio error: %s", err.Error())
	}
	return ""
}

func (b *Bridge) StopAudio() {
	if b.audioService != nil {
		b.audioService.Stop()
	}
}

// MEDIA
func (b *Bridge) MediaPlayPause() {
	if b.mediaService != nil {
		b.mediaService.PlayPause()
	}
}
func (b *Bridge) MediaNext() {
	if b.mediaService != nil {
		b.mediaService.Next()
	}
}

func (b *Bridge) MediaPrevious() {
	if b.mediaService != nil {
		b.mediaService.Previous()
	}
}

func (b *Bridge) MediaStop() {
	if b.mediaService != nil {
		b.mediaService.Stop()
	}
}

// SKINS
func (b *Bridge) ListCustomSkins() string {
	if b.SkinService == nil {
		return "[]"
	}
	return b.SkinService.ListCustomSkin()
}

func (b *Bridge) ReadCustomSkin(id string) string {
	if b.SkinService == nil {
		return ""
	}
	return b.SkinService.ReadCustomSkin(id)
}

func (b *Bridge) GetSkinsDir() string {
	if b.SkinService == nil {
		return ""
	}
	return b.SkinService.GetCustomDir()
}

// CONFIG
func (b *Bridge) LoadSettings() string {
	s := b.configService.Load()
	data, _ := json.Marshal(s)
	return string(data)
}

func (b *Bridge) SaveSettings(jsonStr string) string {
	var s config.Settings
	if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
		return "invalid settings: " + err.Error()
	}
	if err := b.configService.Save(s); err != nil {
		return "save failed: " + err.Error()
	}
	return ""
}

func (b *Bridge) PickBackground() string {
	path, mediaType, errMsg := config.PickBackgroundFile()
	result := map[string]string{
		"path":      path,
		"mediaType": mediaType,
		"error":     errMsg,
	}
	data, _ := json.Marshal(result)
	return string(data)
}

func (b *Bridge) OpenFolder(path string) {
	openFolder(path)
}
