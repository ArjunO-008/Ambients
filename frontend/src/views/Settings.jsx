import { useEffect, useState } from "react"

export default function Settings({ onClose }) {
  const [settings, setSettings] = useState(null)
  const [customSkins, setCustomSkins] = useState([])
  const [status, setStatus] = useState("")
  const [loading, setLoading] = useState(true)

  // Load settings and custom skins on mount
  useEffect(() => {
    Promise.all([
      window.go.bridge.Bridge.LoadSettings(),
      window.go.bridge.Bridge.ListCustomSkins(),
    ]).then(([settingsJson, skinsJson]) => {
      setSettings(JSON.parse(settingsJson))
      setCustomSkins(JSON.parse(skinsJson))
      setLoading(false)
    })
  }, [])

  // Save settings to disk whenever they change
  async function save(updated) {
    setSettings(updated)
    const err = await window.go.bridge.Bridge.SaveSettings(JSON.stringify(updated))
    if (err) {
      setStatus("Save failed: " + err)
    } else {
      setStatus("Saved")
      setTimeout(() => setStatus(""), 1500)
    }
  }

  async function handlePickBackground() {
    setStatus("Opening file picker...")
    const raw = await window.go.bridge.Bridge.PickBackground()
    const result = JSON.parse(raw)

    if (result.error) {
      setStatus("Error: " + result.error)
      return
    }
    if (!result.path) {
      setStatus("") // cancelled
      return
    }

    save({
      ...settings,
      backgroundPath: result.path,
      backgroundType: result.mediaType,
    })
  }

  function handleClearBackground() {
    save({ ...settings, backgroundPath: "", backgroundType: "" })
  }

  async function handleOpenSkinsFolder() {
    const dir = await window.go.bridge.Bridge.GetSkinsDir()
    // open folder in file explorer
    window.go.bridge.Bridge.OpenFolder(dir)
  }

  if (loading) {
    return <div style={styles.root}><div style={styles.loading}>Loading...</div></div>
  }

  return (
    <div style={styles.root}>
      <div style={styles.panel}>

        {/* Header */}
        <div style={styles.header}>
          <span style={styles.title}>Settings</span>
          <button style={styles.closeBtn} onClick={onClose}>✕</button>
        </div>

        {/* Music Player Auto-start */}
        <Section label="Music Player">
          <ToggleRow
            label="Auto-start on overlay open"
            value={settings.musicPlayerEnabled || false}
            onChange={v => save({ ...settings, musicPlayerEnabled: v })}
          />
          {settings.musicPlayerEnabled && (
            <>
              <div style={styles.bgPreview}>
                <span style={styles.bgPath}>
                  {settings.musicPlayerPath
                    ? shortPath(settings.musicPlayerPath)
                    : "No player selected"}
                </span>
              </div>
              <button
                style={styles.primaryBtn}
                onClick={async () => {
                  const path = await window.go.bridge.Bridge.PickMusicPlayer()
                  if (path) save({ ...settings, musicPlayerPath: path })
                }}
              >
                {settings.musicPlayerPath ? "Change player" : "Select music player"}
              </button>
              <span style={styles.hint}>
                App will launch and start playing when overlay activates
              </span>
            </>
          )}
        </Section>

        {/* Clock format */}
        <Section label="Clock">
          <ToggleRow
            label="24-hour format"
            value={settings.use24Hour}
            onChange={v => save({ ...settings, use24Hour: v })}
          />
        </Section>

        {/* Skin selector */}
        <Section label="Skin">
          <div style={styles.skinGrid}>
            {customSkins.map(skin => (
              <SkinCard
                key={skin.id}
                skin={skin}
                active={settings.activeSkinID === skin.id}
                onSelect={() => save({
                  ...settings,
                  activeSkinID: skin.id,
                  activeSkinCustom: true,
                })}
              />
            ))}
          </div>
          <button style={styles.ghostBtn} onClick={handleOpenSkinsFolder}>
            Open skins folder
          </button>
        </Section>

        {/* Background */}
        <Section label="Background">
          {settings.backgroundPath ? (
            <div style={styles.bgPreview}>
              <span style={styles.bgPath}>
                {settings.backgroundType === "video" ? "🎬" : "🖼"} {shortPath(settings.backgroundPath)}
              </span>
              <button style={styles.ghostBtn} onClick={handleClearBackground}>
                Remove
              </button>
            </div>
          ) : (
            <span style={styles.dimText}>No background set</span>
          )}
          <button style={styles.primaryBtn} onClick={handlePickBackground}>
            {settings.backgroundPath ? "Change background" : "Choose image or video"}
          </button>
          <span style={styles.hint}>
            Images: jpg, png, webp, gif — Videos: mp4, webm, mov (max 60s, loops automatically)
          </span>
        </Section>

        {/* Status message */}
        {status && <div style={styles.status}>{status}</div>}  
        {/* Exit button */}
        <button
          style={styles.exitBtn}
          onClick={async () => {
            await window.go.bridge.Bridge.RestoreSleep()
            window.go.bridge.Bridge.QuitApp()
          }}
        >
          Exit AmbientSpace
        </button>     
      </div>      
    </div>
  )
}

// ─── Sub-components ───────────────────────────────────────────────────────────

function Section({ label, children }) {
  return (
    <div style={styles.section}>
      <div style={styles.sectionLabel}>{label}</div>
      <div style={styles.sectionContent}>{children}</div>
    </div>
  )
}

function ToggleRow({ label, value, onChange }) {
  return (
    <div style={styles.toggleRow}>
      <span style={styles.toggleLabel}>{label}</span>
      <div
        style={{
          ...styles.toggle,
          background: value ? "rgba(255,255,255,0.9)" : "rgba(255,255,255,0.1)",
        }}
        onClick={() => onChange(!value)}
      >
        <div style={{
          ...styles.toggleThumb,
          transform: value ? "translateX(20px)" : "translateX(2px)",
          background: value ? "#080808" : "rgba(255,255,255,0.4)",
        }} />
      </div>
    </div>
  )
}

function SkinCard({ skin, active, onSelect }) {
  return (
    <div
      style={{
        ...styles.skinCard,
        border: active
          ? "1px solid rgba(255,255,255,0.6)"
          : "1px solid rgba(255,255,255,0.1)",
      }}
      onClick={onSelect}
    >
      <span style={styles.skinName}>{skin.name}</span>
      {skin.isCustom && <span style={styles.customBadge}>custom</span>}
      {active && <span style={styles.activeDot}>●</span>}
    </div>
  )
}


// ─── Helpers ──────────────────────────────────────────────────────────────────

// shortPath trims long file paths for display
function shortPath(p) {
  const parts = p.replace(/\\/g, "/").split("/")
  return parts.length > 3
    ? "…/" + parts.slice(-2).join("/")
    : p
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const styles = {
  root: {
    position: "fixed",
    inset: 0,
    background: "rgba(0,0,0,0.85)",
    backdropFilter: "blur(12px)",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    zIndex: 100,
    fontFamily: "'Segoe UI', Arial, sans-serif",
  },

  panel: {
    background: "#111",
    border: "1px solid rgba(255,255,255,0.08)",
    borderRadius: "12px",
    width: "440px",
    maxHeight: "80vh",
    overflowY: "auto",
    padding: "24px",
    display: "flex",
    flexDirection: "column",
    gap: "24px",
  },

  header: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },

  title: {
    color: "#fff",
    fontSize: "1rem",
    fontWeight: 500,
    letterSpacing: "0.15em",
    textTransform: "uppercase",
  },

  closeBtn: {
    background: "none",
    border: "none",
    color: "rgba(255,255,255,0.4)",
    fontSize: "1rem",
    cursor: "pointer",
    padding: "4px 8px",
  },

  section: {
    display: "flex",
    flexDirection: "column",
    gap: "12px",
  },

  sectionLabel: {
    color: "rgba(255,255,255,0.3)",
    fontSize: "0.7rem",
    letterSpacing: "0.2em",
    textTransform: "uppercase",
  },

  sectionContent: {
    display: "flex",
    flexDirection: "column",
    gap: "10px",
  },

  // Toggle
  toggleRow: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
  },

  toggleLabel: {
    color: "rgba(255,255,255,0.7)",
    fontSize: "0.9rem",
  },

  toggle: {
    width: "44px",
    height: "24px",
    borderRadius: "12px",
    cursor: "pointer",
    position: "relative",
    transition: "background 0.2s",
  },

  toggleThumb: {
    position: "absolute",
    top: "3px",
    width: "18px",
    height: "18px",
    borderRadius: "50%",
    transition: "transform 0.2s, background 0.2s",
  },

  // Skin grid
  skinGrid: {
    display: "grid",
    gridTemplateColumns: "1fr 1fr",
    gap: "8px",
  },

  skinCard: {
    padding: "12px",
    borderRadius: "8px",
    cursor: "pointer",
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    gap: "6px",
    transition: "border 0.15s",
  },

  skinName: {
    color: "rgba(255,255,255,0.8)",
    fontSize: "0.85rem",
  },

  customBadge: {
    fontSize: "0.6rem",
    color: "rgba(255,255,255,0.3)",
    border: "1px solid rgba(255,255,255,0.15)",
    borderRadius: "4px",
    padding: "1px 5px",
    letterSpacing: "0.05em",
  },

  activeDot: {
    color: "rgba(255,255,255,0.9)",
    fontSize: "0.5rem",
  },

  // Background
  bgPreview: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    gap: "8px",
  },

  bgPath: {
    color: "rgba(255,255,255,0.5)",
    fontSize: "0.8rem",
    overflow: "hidden",
    textOverflow: "ellipsis",
    whiteSpace: "nowrap",
  },

  dimText: {
    color: "rgba(255,255,255,0.25)",
    fontSize: "0.85rem",
  },

  hint: {
    color: "rgba(255,255,255,0.2)",
    fontSize: "0.72rem",
    lineHeight: 1.5,
  },

  // Buttons
  primaryBtn: {
    background: "rgba(255,255,255,0.08)",
    border: "1px solid rgba(255,255,255,0.15)",
    borderRadius: "6px",
    color: "rgba(255,255,255,0.8)",
    padding: "8px 14px",
    fontSize: "0.85rem",
    cursor: "pointer",
    textAlign: "left",
  },

  ghostBtn: {
    background: "none",
    border: "none",
    color: "rgba(255,255,255,0.3)",
    fontSize: "0.8rem",
    cursor: "pointer",
    padding: "0",
    textDecoration: "underline",
    textUnderlineOffset: "3px",
  },

  status: {
    color: "rgba(255,255,255,0.4)",
    fontSize: "0.8rem",
    textAlign: "center",
  },

  loading: {
    color: "rgba(255,255,255,0.3)",
    fontSize: "0.9rem",
  },
  input: {
    background: "rgba(255,255,255,0.05)",
    border: "1px solid rgba(255,255,255,0.12)",
    borderRadius: "6px",
    color: "rgba(255,255,255,0.8)",
    padding: "8px 12px",
    fontSize: "0.85rem",
    outline: "none",
    width: "100%",
  },
  exitBtn: {
    background: "rgba(220,50,50,0.15)",
    border: "1px solid rgba(220,50,50,0.3)",
    borderRadius: "6px",
    color: "rgba(220,100,100,0.9)",
    padding: "10px 14px",
    fontSize: "0.85rem",
    cursor: "pointer",
    width: "100%",
    textAlign: "center",
    marginTop: "8px",
  },
}