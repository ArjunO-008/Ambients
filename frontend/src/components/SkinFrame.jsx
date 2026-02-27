import { useEffect, useRef } from "react"
import { EventsOn } from "../../wailsjs/runtime/runtime"

// Built-in skins imported directly as raw HTML strings
import defaultSkin from "../skins/default.html?raw"
import minimalSkin from "../skins/minimal.html?raw"

const BUILT_IN_SKINS = {
  default: defaultSkin,
  minimal: minimalSkin,
}

// SkinFrame renders the active skin in a sandboxed iframe.
// It injects window.ambientspace API so skins can receive live data
// without knowing anything about the app internals.
//
// Props:
//   skinId     — "default" | "minimal" | custom skin ID
//   isCustom   — if true, loads HTML from Go via ReadCustomSkin
export default function SkinFrame({ skinId = "default", isCustom = false }) {
  const iframeRef = useRef(null)
  const audioUnlistenRef = useRef(null)
  const clockUnlistenRef = useRef(null)

  useEffect(() => {
    loadSkin(skinId, isCustom)
  }, [skinId, isCustom])

  useEffect(() => {
    // Forward audio events into the iframe via postMessage
    audioUnlistenRef.current = EventsOn("audio:frame", (bands) => {
      postToFrame({ type: "audio", bands })
    })

    // Forward clock events into the iframe
    clockUnlistenRef.current = EventsOn("clock:tick", (time) => {
      postToFrame({ type: "clock", time })
    })

    return () => {
      audioUnlistenRef.current?.()
      clockUnlistenRef.current?.()
    }
  }, [])

  // postToFrame sends data to the iframe via window.postMessage
  function postToFrame(data) {
    iframeRef.current?.contentWindow?.postMessage(data, "*")
  }

  // loadSkin fetches HTML and injects it into the iframe with the API wrapper
  async function loadSkin(id, custom) {
    let html = ""

    if (custom) {
      // load from Go (reads from user's skins folder)
      html = await window.go.bridge.Bridge.ReadCustomSkin(id)
    } else {
      // load from built-in map
      html = BUILT_IN_SKINS[id] || BUILT_IN_SKINS["default"]
    }

    // Inject the window.ambientspace API shim into the skin HTML.
    // This runs inside the iframe — skins call these functions,
    // and under the hood they listen to postMessage from the parent.
    const apiShim = `
      <script>
        window.ambientspace = {
          _audioCallbacks: [],
          _clockCallbacks: [],

          onAudio: function(fn) {
            window.ambientspace._audioCallbacks.push(fn)
          },
          onClock: function(fn) {
            window.ambientspace._clockCallbacks.push(fn)
          },
        }

        window.addEventListener('message', function(e) {
          const data = e.data
          if (!data || !data.type) return

          if (data.type === 'audio') {
            window.ambientspace._audioCallbacks.forEach(fn => fn(data.bands))
          }
          if (data.type === 'clock') {
            window.ambientspace._clockCallbacks.forEach(fn => fn(data.time))
          }
        })
      <\/script>
    `

    // inject the shim right after <head> opens so it's available immediately
    const injected = html.replace("<head>", "<head>" + apiShim)

    // write the full HTML into the iframe
    const iframe = iframeRef.current
    if (!iframe) return
    iframe.srcdoc = injected
  }

  return (
    <iframe
      ref={iframeRef}
      style={{
        position: "fixed",
        inset: 0,
        width: "100%",
        height: "100%",
        border: "none",
        background: "transparent",
        pointerEvents: "none", // skin is visual only — clicks pass through to overlay
        zIndex: 0,             // sits behind the main UI
      }}
      sandbox="allow-scripts" // scripts yes, network/storage/popups no
      title="skin"
    />
  )
}