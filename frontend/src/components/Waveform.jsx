import { useEffect, useRef, useState } from "react"
import { EventsOn } from "../../wailsjs/runtime/runtime"

const NUM_BANDS = 64

export default function Waveform() {
  const [bands, setBands] = useState(() => Array(NUM_BANDS).fill(0))
  const [error, setError] = useState("")
  const startedRef = useRef(false)

  useEffect(() => {
    if (!startedRef.current) {
      startedRef.current = true
      window.go.bridge.Bridge.StartAudio().then((err) => {
        if (err) setError(err)
      })
    }

    const unlisten = EventsOn("audio:frame", (data) => {
      if (Array.isArray(data) && data.length === NUM_BANDS) {
        setBands(data)
      }
    })

    return () => {
      unlisten?.()
      // Only stop audio when component truly unmounts
      window.go.bridge.Bridge.StopAudio()
      startedRef.current = false
    }
  }, [])

  if (error) {
    return <div style={{ color: "red", fontSize: "0.8rem" }}>{error}</div>
  }

  return (
    <div
      style={{
        display: "flex",
        alignItems: "flex-end",
        gap: "3px",
        height: "120px",
        padding: "0 8px",
      }}
    >
      {bands.map((amplitude, i) => (
        <div
          key={i}
          style={{
            flex: 1,
            height: `${Math.max(4, amplitude * 100)}%`,
            background: `rgba(255,255,255,${0.4 + amplitude * 0.6})`,
            borderRadius: "2px 2px 0 0",
            transition: "height 0.05s ease-out",
          }}
        />
      ))}
    </div>
  )
}