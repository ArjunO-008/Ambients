import { useEffect, useState } from "react";
import { EventsOn } from "../../wailsjs/runtime"

export default function Clock({ user24Hours = false }) {
    const [time, setTime] = useState({
        hours: "00",
        minutes: "00",
        seconds: "00",
        date: "",
        ampm: "",
    })

    useEffect(() => {
        window.go.bridge.Bridge.StartClock(user24Hours)
        const unlisten = EventsOn("clock:tick", (payload) => {
            setTime(payload)
        })

        return () => {
             window.go.bridge.Bridge.StopClock()
             unlisten()
        }
    }, [user24Hours])

     return (
    <div style={{ textAlign: "center", color: "white", fontFamily: "monospace" }}>
      <div style={{ fontSize: "6rem", letterSpacing: "0.1em" }}>
        {time.hours}:{time.minutes}
        {time.ampm && (
          <span style={{ fontSize: "2rem", marginLeft: "0.5rem" }}>
            {time.ampm}
          </span>
        )}
      </div>
      <div style={{ fontSize: "1.2rem", opacity: 0.6 }}>
        {time.date}
      </div>
    </div>
  )
}