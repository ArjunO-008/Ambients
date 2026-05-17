# Ambients Skin API

Skins are plain `.html` files. Ambients renders them inside a sandboxed iframe
and injects the `window.ambientspace` API before your scripts run.

---

## Installation

Drop any `.html` file into the skins folder:

```
%APPDATA%\ambientspace\skins\                        (Windows)
~/.config/ambientspace/skins/                        (Linux)
~/Library/Application Support/ambientspace/skins/   (macOS)
```

Then open **Settings** inside Ambients and select your skin from the list.  
You can also click **Open skins folder** in Settings to open the folder directly.

---

## Layout

The skin iframe fills the **top 80%** of the screen.  
The bottom 20% is reserved for Ambients media controls.  
Use `background: transparent` on `body` to let the user's background show through.

---

## API Reference

### `window.ambientspace.onClock(callback)`

Called once per second with the current time.

```javascript
window.ambientspace.onClock(function(time) {
  time.hours    // string — zero-padded, e.g. "09"
  time.minutes  // string — zero-padded, e.g. "42"
  time.seconds  // string — zero-padded, e.g. "07"
  time.date     // string — e.g. "Wednesday, 22 April 2026"
  time.ampm     // string — "AM" or "PM", empty string in 24h mode
})
```

The 12h/24h format follows the user's setting in Ambients.

---

### `window.ambientspace.onAudio(callback)`

Called approximately 20 times per second with frequency data.

```javascript
window.ambientspace.onAudio(function(bands) {
  // bands — array of 64 numbers, each between 0.0 and 1.0
  // bands[0]  = lowest frequency (bass)
  // bands[63] = highest frequency (treble)
})
```

Values are normalized. `0.0` is silence, `1.0` is peak amplitude.

---

## Sandbox

Skins run with `sandbox="allow-scripts"` only.

| Capability | Available |
|------------|-----------|
| JavaScript | ✅ |
| CSS / Animations | ✅ |
| Inline web fonts | ✅ |
| Canvas / WebGL | ✅ |
| External network | ❌ |
| localStorage | ❌ |
| Parent DOM access | ❌ |
| iframe nesting | ❌ |

All rendering is local. Skins have no access to the internet or the host system.

---

## Minimal example

```html
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  html, body {
    width: 100%; height: 100%;
    background: transparent;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: sans-serif;
    color: white;
  }
</style>
</head>
<body>
<div id="time" style="font-size: 8rem; font-weight: 100;"></div>
<script>
  window.ambientspace.onClock(function(time) {
    document.getElementById('time').textContent =
      time.hours + ':' + time.minutes
  })
</script>
</body>
</html>
```

---

## Browser testing stub

Test your skin in a browser before dropping it into Ambients.  
Paste this stub before your skin's `<script>` tag:

```javascript
window.ambientspace = {
  onClock: function(fn) {
    function tick() {
      var now = new Date()
      var h = now.getHours()
      var ampm = h >= 12 ? 'PM' : 'AM'
      h = h % 12 || 12
      fn({
        hours:   String(h).padStart(2, '0'),
        minutes: String(now.getMinutes()).padStart(2, '0'),
        seconds: String(now.getSeconds()).padStart(2, '0'),
        date:    now.toDateString(),
        ampm:    ampm
      })
    }
    tick()
    setInterval(tick, 1000)
  },
  onAudio: function(fn) {
    setInterval(function() {
      fn(Array.from({ length: 64 }, function() {
        return Math.random() * 0.4
      }))
    }, 50)
  }
}
```

---

## Tips

- Skins fill the top 80% of the screen — design for a wide landscape area
- `background: transparent` lets the user's chosen background show through
- Both `onClock` and `onAudio` can be used independently — you don't need both
- Animations work normally — CSS transitions, Canvas `requestAnimationFrame`, WebGL all work
- The API is injected before your scripts run — no need to wait for a load event

---

## Built-in skins

Two skins ship with Ambients and are seeded into the skins folder on first launch:

| Skin | Description |
|------|-------------|
| `default` | Large clock with mirrored waveform bars below |
| `minimal` | Faded centered time, no waveform |

These can be edited or replaced — they are plain HTML files in the skins folder.