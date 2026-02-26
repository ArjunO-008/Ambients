import Clock from "./components/Clock"
import Waveform from "./components/Waveform"

function App() {
  return (
    <div style={{
      background: "#0a0a0a",
      minHeight: "100vh",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      gap: "40px",
    }}>
      <Clock use24Hour={false} />
      <Waveform />
    </div>
  )
}

export default App