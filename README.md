# Ambients

A minimal fullscreen ambient overlay for your desktop.  
Clock, waveform visualizer and media controls — all driven by a live HTML skin system.

---

## Features

- Fullscreen ambient overlay with live clock and audio waveform
- Skin system — drop any `.html` file into the skins folder to customize the look entirely
- Built-in skins: Default (waveform + clock) and Minimal (faded clock)
- Background image or video support
- Media playback controls (play/pause, next, previous, stop)
- Auto-launch music player on overlay open
- Prevents system sleep while active
- Exit restores all system state cleanly

---

## Download

Grab the latest release from [Releases](https://github.com/ArjunO-008/Ambients/releases).

| File | Platform |
|------|----------|
| `Ambients-amd64-installer.exe` | Windows — installer (recommended) |
| `Ambients.exe` | Windows — portable executable |
| `Ambients-x86_64.AppImage` | Linux |

> macOS builds coming in a future release.

---

## Windows

Run `Ambients-amd64-installer.exe` — no admin rights required.  
Installs to `%LOCALAPPDATA%\Ambients` and creates desktop and Start Menu shortcuts.

**For waveform visualizer:**  
Right-click speaker → Sounds → Recording tab → right-click empty area → Show Disabled Devices → enable **Stereo Mix**.

---

## Linux

Make the AppImage executable and run it:

```bash
chmod +x Ambients-x86_64.AppImage
./Ambients-x86_64.AppImage
```

**Optional dependencies:**

| Package | Purpose |
|---------|---------|
| `xdotool` | Media keys on X11 |
| `ydotool` | Media keys on Wayland |
| `wmctrl` | Window focus after music player launch |

```bash
sudo apt install xdotool wmctrl
```

**For waveform visualizer:**  
Requires PulseAudio with a monitor source enabled.  
Run `pactl list sources short` and look for a source with "monitor" in the name.

**Known limitations:**
- Video backgrounds require GPU-accelerated WebKit. VMs and minimal installs may show a black background instead.
- Wayland media keys require `ydotool` with `ydotoold` daemon running.

---

## macOS

> No Mac hardware available for testing. Community contributions welcome.

Build from source on a Mac:

```bash
wails build -platform darwin/arm64 -clean   # Apple Silicon
wails build -platform darwin/amd64 -clean   # Intel
```

First launch: right-click `Ambients.app` → Open → Open anyway  
(Required because the app is not notarized)

**For waveform visualizer:**  
Install [BlackHole](https://existential.audio/blackhole/) — a free virtual audio driver for macOS loopback.

---

## Skins

Skins are plain `.html` files. Drop them into the skins folder:

```
%APPDATA%\ambientspace\skins\        (Windows)
~/.config/ambientspace/skins/        (Linux)
~/Library/Application Support/ambientspace/skins/  (macOS)
```

Open the skins folder directly from **Settings → Open skins folder**.

See [SKINS.md](SKINS.md) for the full skin API reference and examples.

---

## Building from source

**Requirements:**
- Go 1.21+
- Node.js 18+
- Wails v2 — `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- PortAudio

```bash
git clone https://github.com/yourusername/ambients
cd ambients
wails build -platform windows/amd64 -clean -nsis
```

---

## License

MIT — see [LICENSE](LICENSE)  
Copyright 2026 Arjun.O
