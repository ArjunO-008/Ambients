import { useEffect, useState, useRef } from "react"
import { WindowFullscreen } from "../../wailsjs/runtime/runtime"
import SkinFrame from "../components/SkinFrame"
import Settings from "./Settings"

export default function Overlay() {
  const [showSettings, setShowSettings] = useState(false)
  const [appSettings, setAppSettings]   = useState(null)
  const settingsLoadedRef               = useRef(false)

  // Load settings ONCE
  useEffect(() => {
    if (settingsLoadedRef.current) return
    settingsLoadedRef.current = true
    window.go.bridge.Bridge.LoadSettings().then((raw) => {
      setAppSettings(JSON.parse(raw))
    })
  }, [])

  // Fullscreen
  useEffect(() => {
    WindowFullscreen()
  }, [])

  // Prevent sleep
  useEffect(() => {
    window.go.bridge.Bridge.PreventSleep()
    return () => window.go.bridge.Bridge.RestoreSleep()
  }, [])

  // Auto-start music player
  useEffect(() => {
    window.go.bridge.Bridge.LaunchMusicPlayer()
  }, [])

  // Start clock + audio — feeds data into the skin via SkinFrame's EventsOn listeners
  // Depends on appSettings so we pass the correct use24Hour value
  // Re-runs if use24Hour changes (toggled in settings)
  useEffect(() => {
    if (!appSettings) return
    window.go.bridge.Bridge.StartClock(appSettings.use24Hour)
    window.go.bridge.Bridge.StartAudio()
    return () => {
      window.go.bridge.Bridge.StopClock()
      window.go.bridge.Bridge.StopAudio()
    }
  }, [appSettings?.use24Hour])

  function handleSettingsClose() {
    setShowSettings(false)
    window.go.bridge.Bridge.LoadSettings().then((raw) => {
      setAppSettings(JSON.parse(raw))
    })
  }

  if (!appSettings) {
    return <div style={{ position: "fixed", inset: 0, background: "#080808" }} />
  }

  return (
    <div style={styles.root}>

      {/* Background image or video */}
      <Background settings={appSettings} />

      {/* Grain overlay */}
      <div style={styles.grain} />

      {/* Top 80% — skin iframe renders clock + waveform */}
      <div style={styles.skinArea}>
        <SkinFrame
          skinId={appSettings.activeSkinID || "default"}
          isCustom={appSettings.activeSkinCustom || false}
        />
      </div>

      {/* Bottom 20% — media controls */}
      <div style={styles.mediaArea}>
        <MediaBar />
      </div>

      {/* Settings gear — fixed top right, above everything */}
      <button
        style={styles.settingsBtn}
        onClick={() => setShowSettings(true)}
        title="Settings"
      >
        ⚙
      </button>

      {showSettings && <Settings onClose={handleSettingsClose} />}

    </div>
  )
}

function Background({ settings }) {
  if (!settings?.backgroundPath || !settings?.backgroundType) return null
  const src = `/media?path=${encodeURIComponent(settings.backgroundPath)}`
  const base = {
    position: "fixed", inset: 0,
    width: "100%", height: "100%",
    objectFit: "cover", zIndex: 0, pointerEvents: "none",
  }
  if (settings.backgroundType === "video") {
    return <video key={src} style={base} src={src} autoPlay loop muted playsInline />
  }
  return <img key={src} style={base} src={src} alt="" />
}

function MediaBar() {
  const btn = (label, fn, large = false) => (
    <button
      style={{ ...styles.mediaBtn, ...(large ? styles.mediaBtnLarge : {}) }}
      onClick={fn}
      onMouseEnter={e => e.currentTarget.style.background = "rgba(255,255,255,0.12)"}
      onMouseLeave={e => e.currentTarget.style.background = "rgba(255,255,255,0.05)"}
    >
      {label}
    </button>
  )
  return (
    <div style={styles.mediaBar}>
      {btn("⏮", () => window.go.bridge.Bridge.MediaPrevious())}
      {btn("⏯", () => window.go.bridge.Bridge.MediaPlayPause(), true)}
      {btn("⏭", () => window.go.bridge.Bridge.MediaNext())}
      {btn("⏹", () => window.go.bridge.Bridge.MediaStop())}
    </div>
  )
}

const styles = {
  root: {
    position: "fixed", inset: 0,
    background: "#080808",
    display: "flex",
    flexDirection: "column",
    overflow: "hidden",
    userSelect: "none",
    fontFamily: "'Segoe UI', Arial, sans-serif",
  },
  grain: {
    position: "fixed", inset: 0,
    backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)' opacity='0.03'/%3E%3C/svg%3E")`,
    backgroundRepeat: "repeat",
    backgroundSize: "128px 128px",
    pointerEvents: "none",
    opacity: 0.4,
    zIndex: 1,
  },
  skinArea: {
    flex: "0 0 80%",
    position: "relative",
    zIndex: 2,
    overflow: "hidden",
  },
  mediaArea: {
    flex: "0 0 20%",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    zIndex: 2,
    borderTop: "1px solid rgba(255,255,255,0.04)",
  },
  mediaBar: {
    display: "flex",
    alignItems: "center",
    gap: "16px",
  },
  mediaBtn: {
    background: "rgba(255,255,255,0.05)",
    border: "1px solid rgba(255,255,255,0.08)",
    borderRadius: "50%",
    width: "44px", height: "44px",
    color: "rgba(255,255,255,0.6)",
    fontSize: "1rem",
    cursor: "pointer",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    transition: "background 0.2s",
    outline: "none",
  },
  mediaBtnLarge: {
    width: "56px", height: "56px",
    fontSize: "1.3rem",
    color: "rgba(255,255,255,0.9)",
    border: "1px solid rgba(255,255,255,0.15)",
  },
  settingsBtn: {
    position: "fixed",
    top: "16px", right: "16px",
    background: "rgba(255,255,255,0.05)",
    border: "1px solid rgba(255,255,255,0.08)",
    borderRadius: "50%",
    width: "36px", height: "36px",
    color: "rgba(255,255,255,0.3)",
    fontSize: "1rem",
    cursor: "pointer",
    zIndex: 50,
    outline: "none",
  },
}