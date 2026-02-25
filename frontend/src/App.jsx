import Clock from "./components/Clock"

function App() {
    return (
        <div style={{ background: "#0a0a0a", minHeight: "100vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
            <Clock user24Hours={false} />
        </div>
    )
}

export default App
