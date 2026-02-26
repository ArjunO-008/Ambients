import { useEffect, useRef, useState } from "react"
import { EventsOn } from "../../wailsjs/runtime/runtime"

const NUM_BANDS = 64

export default function Waveform() {
  const [bands, setBands] = useState(() => Array(NUM_BANDS).fill(0))
  const [error, setError] = useState("")
  const startedRef = useRef(false)

  useEffect(() => {
    // Listen for frequency data from Go
    const unlisten = EventsOn("audio:frame", (data) => {
      if (Array.isArray(data) && data.length === NUM_BANDS) {
        setBands(data)
      }
    })

    // Listen for errors emitted from Go
    const unlistenError = EventsOn("audio:error", (msg) => {
      console.error("Audio error:", msg)
      setError(msg)
    })

    // Start audio capture
    if (!startedRef.current) {
      startedRef.current = true
      window.go.bridge.Bridge.StartAudio().then((err) => {
        if (err) {
          console.error("StartAudio error:", err)
          setError(err)
        }
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
    return <div style={{ color: "red", fontSize: "0.8rem", padding: "8px" }}>{error}</div>
  }

  return (
    <div style={{
      display: "flex",
      alignItems: "flex-end",
      gap: "3px",
      height: "160px",
      width: "500px",
      padding: "0 8px",
    }}>
      {bands.map((amplitude, i) => (
        <div
          key={i}
          style={{
            flex: 1,
            height: `${Math.max(8, amplitude * 100)}%`,
            background: `rgba(99, 179, 237, ${0.5 + amplitude * 0.5})`,
            borderRadius: "2px 2px 0 0",
            transition: "height 0.05s ease-out",
          }}
        />
      ))}
    </div>
  )
}