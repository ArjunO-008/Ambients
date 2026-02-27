export default function MediaControls() {

  const handlePlayPause = () => window.go.bridge.Bridge.MediaPlayPause()
  const handleNext      = () => window.go.bridge.Bridge.MediaNext()
  const handlePrevious  = () => window.go.bridge.Bridge.MediaPrevious()
  const handleStop      = () => window.go.bridge.Bridge.MediaStop()

  const btnStyle = {
    background: "rgba(255,255,255,0.1)",
    border: "1px solid rgba(255,255,255,0.2)",
    borderRadius: "50%",
    width: "44px",
    height: "44px",
    color: "white",
    fontSize: "1.1rem",
    cursor: "pointer",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    transition: "background 0.2s",
  }

  const hoverStyle = {
    background: "rgba(255,255,255,0.2)",
  }

  return (
    <div style={{
      display: "flex",
      alignItems: "center",
      gap: "16px",
    }}>


      <button
        style={btnStyle}
        onClick={handlePrevious}
        onMouseEnter={e => Object.assign(e.target.style, hoverStyle)}
        onMouseLeave={e => Object.assign(e.target.style, btnStyle)}
        title="Previous"
      >
        ⏮
      </button>

      <button
        style={{ ...btnStyle, width: "56px", height: "56px", fontSize: "1.4rem" }}
        onClick={handlePlayPause}
        onMouseEnter={e => Object.assign(e.target.style, { ...hoverStyle, width: "56px", height: "56px" })}
        onMouseLeave={e => Object.assign(e.target.style, { ...btnStyle, width: "56px", height: "56px" })}
        title="Play / Pause"
      >
        ⏯
      </button>

      <button
        style={btnStyle}
        onClick={handleNext}
        onMouseEnter={e => Object.assign(e.target.style, hoverStyle)}
        onMouseLeave={e => Object.assign(e.target.style, btnStyle)}
        title="Next"
      >
        ⏭
      </button>

      <button
        style={btnStyle}
        onClick={handleStop}
        onMouseEnter={e => Object.assign(e.target.style, hoverStyle)}
        onMouseLeave={e => Object.assign(e.target.style, btnStyle)}
        title="Stop"
      >
        ⏹
      </button>
    </div>
  )
}