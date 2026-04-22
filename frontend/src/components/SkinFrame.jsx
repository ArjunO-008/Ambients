import { useEffect, useRef } from "react"
import { EventsOn } from "../../wailsjs/runtime/runtime"

export default function SkinFrame({ skinId = "default", isCustom = false }) {
  const iframeRef        = useRef(null)
  const audioUnlistenRef = useRef(null)
  const clockUnlistenRef = useRef(null)

  useEffect(() => {
    // Reset hide flags whenever skin changes — default skin has no hide config
    window.dispatchEvent(new CustomEvent("skin-hide", {
      detail: { clock: false, waveform: false }
    }))
    loadSkin(skinId)
  }, [skinId, isCustom])

  useEffect(() => {
    audioUnlistenRef.current = EventsOn("audio:frame", (bands) => {
      postToFrame({ type: "audio", bands })
    })
    clockUnlistenRef.current = EventsOn("clock:tick", (time) => {
      postToFrame({ type: "clock", time })
    })

    // Listen for skin-config messages from iframe
    function handleMessage(e) {
      if (e.data?.type === "skin-config") {
        const hide = e.data.hide || ""
        window.dispatchEvent(new CustomEvent("skin-hide", {
          detail: {
            clock:    hide.includes("clock"),
            waveform: hide.includes("waveform"),
          }
        }))
      }
    }
    window.addEventListener("message", handleMessage)

    return () => {
      audioUnlistenRef.current?.()
      clockUnlistenRef.current?.()
      window.removeEventListener("message", handleMessage)
    }
  }, [])

  function postToFrame(data) {
    iframeRef.current?.contentWindow?.postMessage(data, "*")
  }

  async function loadSkin(id,custom) {
    // All skins now live in the config folder — load via Go
    const html = await window.go.bridge.Bridge.ReadCustomSkin(id)

    if (!html) {
      console.error("Skin not found:", id)
      return
    }

    const apiShim = `
      <script>
        window.ambientspace = {
          _audioCallbacks: [],
          _clockCallbacks: [],
          onAudio: function(fn) { window.ambientspace._audioCallbacks.push(fn) },
          onClock: function(fn) { window.ambientspace._clockCallbacks.push(fn) },
        }
        window.addEventListener('message', function(e) {
          const data = e.data
          if (!data || !data.type) return
          if (data.type === 'audio') window.ambientspace._audioCallbacks.forEach(fn => fn(data.bands))
          if (data.type === 'clock') window.ambientspace._clockCallbacks.forEach(fn => fn(data.time))
        })
      <\/script>
    `

    // inject shim before skin's own scripts
    const injected = html.includes("<script>")
      ? html.replace("<script>", apiShim + "<script>")
      : html.replace("</body>", apiShim + "</body>")

    const iframe = iframeRef.current
    if (!iframe) return
    iframe.srcdoc = injected
  }

  return (
    <iframe
      ref={iframeRef}
      style={{
        position: "absolute",
        inset: 0,
        width: "100%",
        height: "100%",
        border: "none",
        background: "transparent",
        pointerEvents: "none",
        zIndex: 2,
      }}
      sandbox="allow-scripts"
      title="skin"
    />
  )
}