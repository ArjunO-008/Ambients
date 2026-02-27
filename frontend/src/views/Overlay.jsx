import { useState, useEffect, useRef } from "react"
import { EventsOn } from "../../wailsjs/runtime/runtime"

const NUM_BANDS = 64

//____OVERLAY_______________________________________

export default function Overlay() {
    return (
        <div style={styles.root}>
            {/* Subtle grain texture overlay — adds depth without color */}
            <div style={styles.grain} />

            {/* Main content — vertically centered */}
            <div style={styles.center}>
                <ClockDisplay />
                <WaveformDisplay />
            </div>

            {/* Media controls pinned to bottom */}
            <div style={styles.bottom}>
                <MediaBar />
            </div>
        </div>
    )
}

//____CLOCK_______________________________________

function ClockDisplay() {
    const [time, setTime] = useState({
        hours: "00",
        minutes: "00",
        seconds: "00",
        date: "",
        ampm: "",
    })

    useEffect(() => {
        window.go.bridge.Bridge.StartClock(false)
        const unlisten = EventsOn("clock:tick", setTime)
        return () => {
            unlisten?.()
            window.go.bridge.Bridge.StopClock()
        }
    }, [])

    return (
        <div style={styles.clockWrap}>
            {/* Main time display */}
            <div style={styles.clockRow}>
                <span style={styles.clockDigits}>
                    {time.hours}
                </span>
                <span style={styles.clockColon}>:</span>
                <span style={styles.clockDigits}>
                    {time.minutes}
                </span>
                {time.ampm && (
                    <span style={styles.clockAmPm}>{time.ampm}</span>
                )}
            </div>

            {/* Seconds — smaller, dimmer */}
            <div style={styles.clockSeconds}>
                :{time.seconds}
            </div>

            {/* Date */}
            <div style={styles.clockDate}>
                {time.date}
            </div>
        </div>
    )

}

// ─── Waveform ─────────────────────────────────────────────────────────────────

function WaveformDisplay() {
    const [bands, setBands] = useState(() => Array(NUM_BANDS).fill(0))
    const [error, setError] = useState("")
    const startedRef = useRef(false)

    useEffect(() => {
        const unlisten = EventsOn("audio:frame", (data) => {
            if (Array.isArray(data) && data.length === NUM_BANDS) {
                setBands(data)
            }
        })

        const unlistenError = EventsOn("audio:error", (msg) => {
            setError(msg)
        })

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

    if (error) {
        return <div style={styles.error}>{error}</div>
    }

    return (
        <div style={styles.waveformWrap}>
            {bands.map((amplitude, i) => {
                // mirror the waveform — bars grow from center outward
                // gives a more organic, symmetric feel
                const mirrored = i < NUM_BANDS / 2
                    ? bands[NUM_BANDS / 2 - 1 - i]
                    : bands[i - NUM_BANDS / 2]

                const height = Math.max(3, mirrored * 100)
                const opacity = 0.25 + mirrored * 0.75

                return (
                    <div
                        key={i}
                        style={{
                            ...styles.bar,
                            height: `${height}%`,
                            opacity,
                        }}
                    />
                )
            })}
        </div>
    )
}


// ─── Media Controls ───────────────────────────────────────────────────────────

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
        fontFamily: "'Segoe UI', 'Arial', sans-serif",
        overflow: "hidden",
        userSelect: "none",
    },

    // SVG noise grain — pure CSS, no image needed
    grain: {
        position: "fixed",
        inset: 0,
        backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)' opacity='0.03'/%3E%3C/svg%3E")`,
        backgroundRepeat: "repeat",
        backgroundSize: "128px 128px",
        pointerEvents: "none",
        opacity: 0.4,
    },

    center: {
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "48px",
        zIndex: 1,
    },

    // Clock
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
        fontWeight: 500,       // ultra-thin — elegant on dark background
        color: "#ffffff",
        letterSpacing: "-0.02em",
        fontVariantNumeric: "tabular-nums", // prevents layout shift as digits change
    },

    clockColon: {
        fontSize: "clamp(4rem, 11vw, 8rem)",
        fontWeight: 200,
        color: "rgba(255,255,255,0.3)",
        margin: "0 4px",
        // blink the colon every second for a subtle pulse
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

    // Waveform
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
        transformOrigin: "bottom",
    },

    // Media controls
    bottom: {
        position: "fixed",
        bottom: "60px",
        zIndex: 1,
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
        transition: "background 0.2s, color 0.2s",
        outline: "none",
    },

    mediaBtnLarge: {
        width: "56px",
        height: "56px",
        fontSize: "1.3rem",
        color: "rgba(255,255,255,0.9)",
        border: "1px solid rgba(255,255,255,0.15)",
    },

    error: {
        color: "rgba(255,100,100,0.8)",
        fontSize: "0.8rem",
        maxWidth: "400px",
        textAlign: "center",
    },
}