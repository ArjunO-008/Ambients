import { useEffect, useState, useRef } from "react"
import { EventsOn } from "../../wailsjs/runtime/runtime"
import SkinFrame from "../components/SkinFrame"
import Settings from "./Settings"

const NUM_BANDS = 64

export default function Overlay() {
  const [showSettings, setShowSettings] = useState(false)
  const [appSettings, setAppSettings]   = useState(null)

  // Load settings on mount
  useEffect(() => {
    window.go.bridge.Bridge.LoadSettings().then((raw) => {
      setAppSettings(JSON.parse(raw))
    })
  }, [])

  // Reload settings when settings panel closes
  function handleSettingsClose() {
    setShowSettings(false)
    window.go.bridge.Bridge.LoadSettings().then((raw) => {
      setAppSettings(JSON.parse(raw))
    })
  }

  if (!appSettings) return null // wait for settings to load

  return (
    <div style={styles.root}>

      {/* Background — image or video */}
      <Background settings={appSettings} />

      {/* Active skin */}
      <SkinFrame
        skinId={appSettings.activeSkinID || "default"}
        isCustom={appSettings.activeSkinCustom || false}
      />

      {/* Grain texture */}
      <div style={styles.grain} />

      {/* Main content */}
      <div style={styles.center}>
        <ClockDisplay use24Hour={appSettings.use24Hour} />
        <WaveformDisplay />
      </div>

      {/* Media controls */}
      <div style={styles.bottom}>
        <MediaBar />
      </div>

      {/* Temporary settings toggle — will move to tray in Step 9 */}
      <button
        style={styles.settingsBtn}
        onClick={() => setShowSettings(true)}
        title="Settings"
      >
        ⚙
      </button>

      {/* Settings panel */}
      {showSettings && <Settings onClose={handleSettingsClose} />}

    </div>
  )
}

// ─── Background ───────────────────────────────────────────────────────────────

function Background({ settings }) {
  if (!settings.backgroundPath || !settings.backgroundType) return null

  // encode the path so backslashes and spaces are safe in a URL
  const src = `/media?path=${encodeURIComponent(settings.backgroundPath)}`

  const baseStyle = {
    position: "fixed",
    inset: 0,
    width: "100%",
    height: "100%",
    objectFit: "cover",
    zIndex: 0,
    pointerEvents: "none",
  }

  if (settings.backgroundType === "video") {
    return (
      <video
        key={src}
        style={baseStyle}
        src={src}
        autoPlay
        loop
        muted
        playsInline
      />
    )
  }

  return (
    <img
      key={src}
      style={baseStyle}
      src={src}
      alt=""
    />
  )
}
// ─── Clock, Waveform, MediaBar — same as before, just pass use24Hour ─────────

function ClockDisplay({ use24Hour }) {
  const [time, setTime] = useState({
    hours: "00", minutes: "00", seconds: "00", date: "", ampm: "",
  })

  useEffect(() => {
    window.go.bridge.Bridge.StartClock(use24Hour)
    const unlisten = EventsOn("clock:tick", setTime)
    return () => {
      unlisten?.()
      window.go.bridge.Bridge.StopClock()
    }
  }, [use24Hour])

  return (
    <div style={styles.clockWrap}>
      <div style={styles.clockRow}>
        <span style={styles.clockDigits}>{time.hours}</span>
        <span style={styles.clockColon}>:</span>
        <span style={styles.clockDigits}>{time.minutes}</span>
        {time.ampm && <span style={styles.clockAmPm}>{time.ampm}</span>}
      </div>
      <div style={styles.clockSeconds}>:{time.seconds}</div>
      <div style={styles.clockDate}>{time.date}</div>
    </div>
  )
}

function WaveformDisplay() {
  const [bands, setBands]   = useState(() => Array(NUM_BANDS).fill(0))
  const [error, setError]   = useState("")
  const startedRef          = useRef(false)

  useEffect(() => {
    const unlisten      = EventsOn("audio:frame", (data) => {
      if (Array.isArray(data) && data.length === NUM_BANDS) setBands(data)
    })
    const unlistenError = EventsOn("audio:error", setError)

    if (!startedRef.current) {
      startedRef.current = true
      window.go.bridge.Bridge.StartAudio().then((err) => {
        if (err) setError(err)
      })
    }
    return () => {
      unlisten?.()
      unlistenError?.()
      window.go.bridge.Bridge.StopAudio()
      startedRef.current = false
    }
  }, [])

  if (error) return <div style={styles.error}>{error}</div>

  return (
    <div style={styles.waveformWrap}>
      {bands.map((amplitude, i) => {
        const mirrored = i < NUM_BANDS / 2
          ? bands[NUM_BANDS / 2 - 1 - i]
          : bands[i - NUM_BANDS / 2]
        return (
          <div key={i} style={{
            ...styles.bar,
            height: `${Math.max(3, mirrored * 100)}%`,
            opacity: 0.25 + mirrored * 0.75,
          }} />
        )
      })}
    </div>
  )
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

// ─── Styles ───────────────────────────────────────────────────────────────────

const styles = {
  root: {
    position: "fixed",
    inset: 0,
    background: "#080808",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    fontFamily: "'Segoe UI', Arial, sans-serif",
    overflow: "hidden",
    userSelect: "none",
  },
  grain: {
    position: "fixed",
    inset: 0,
    backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)' opacity='0.03'/%3E%3C/svg%3E")`,
    backgroundRepeat: "repeat",
    backgroundSize: "128px 128px",
    pointerEvents: "none",
    opacity: 0.4,
    zIndex: 1,
  },
  center: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    gap: "48px",
    zIndex: 2,
  },
  clockWrap: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    gap: "4px",
  },
  clockRow: {
    display: "flex",
    alignItems: "baseline",
    gap: "2px",
    lineHeight: 1,
  },
  clockDigits: {
    fontSize: "clamp(5rem, 14vw, 10rem)",
    fontWeight: 200,
    color: "#ffffff",
    letterSpacing: "-0.02em",
    fontVariantNumeric: "tabular-nums",
  },
  clockColon: {
    fontSize: "clamp(4rem, 11vw, 8rem)",
    fontWeight: 200,
    color: "rgba(255,255,255,0.3)",
    margin: "0 4px",
    animation: "blink 1s step-end infinite",
  },
  clockAmPm: {
    fontSize: "clamp(1.2rem, 3vw, 2rem)",
    fontWeight: 300,
    color: "rgba(255,255,255,0.4)",
    marginLeft: "12px",
    alignSelf: "flex-end",
    marginBottom: "12px",
    letterSpacing: "0.1em",
  },
  clockSeconds: {
    fontSize: "clamp(1rem, 2.5vw, 1.6rem)",
    fontWeight: 300,
    color: "rgba(255,255,255,0.2)",
    letterSpacing: "0.05em",
    marginTop: "4px",
  },
  clockDate: {
    fontSize: "clamp(0.75rem, 1.5vw, 1rem)",
    fontWeight: 300,
    color: "rgba(255,255,255,0.25)",
    letterSpacing: "0.2em",
    textTransform: "uppercase",
    marginTop: "8px",
  },
  waveformWrap: {
    display: "flex",
    alignItems: "flex-end",
    gap: "2px",
    height: "120px",
    width: "min(700px, 75vw)",
  },
  bar: {
    flex: 1,
    background: "#ffffff",
    borderRadius: "1px 1px 0 0",
    transition: "height 0.06s ease-out, opacity 0.06s ease-out",
  },
  bottom: {
    position: "fixed",
    bottom: "60px",
    zIndex: 2,
  },
  mediaBar: {
    display: "flex",
    alignItems: "center",
    gap: "12px",
  },
  mediaBtn: {
    background: "rgba(255,255,255,0.05)",
    border: "1px solid rgba(255,255,255,0.08)",
    borderRadius: "50%",
    width: "44px",
    height: "44px",
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
    width: "56px",
    height: "56px",
    fontSize: "1.3rem",
    color: "rgba(255,255,255,0.9)",
    border: "1px solid rgba(255,255,255,0.15)",
  },
  settingsBtn: {
    position: "fixed",
    top: "20px",
    right: "20px",
    background: "rgba(255,255,255,0.05)",
    border: "1px solid rgba(255,255,255,0.08)",
    borderRadius: "50%",
    width: "36px",
    height: "36px",
    color: "rgba(255,255,255,0.3)",
    fontSize: "1rem",
    cursor: "pointer",
    zIndex: 10,
    outline: "none",
  },
  error: {
    color: "rgba(255,100,100,0.8)",
    fontSize: "0.8rem",
    maxWidth: "400px",
    textAlign: "center",
  },
}