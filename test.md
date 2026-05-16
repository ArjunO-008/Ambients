Run these one by one in your Linux terminal:

### 1. Build essentials
```bash
sudo apt update
sudo apt install -y build-essential git curl
```

### 2. PortAudio (required for waveform)
```bash
sudo apt install -y libportaudio2 portaudio19-dev
```

### 3. Wails system dependencies
```bash
sudo apt install -y libgtk-3-dev libwebkit2gtk-4.0-dev
```

### 4. Media keys
```bash
sudo apt install -y xdotool
```

### 5. Go
```bash
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### 6. Node.js
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node --version
npm --version
```

### 7. Wails
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
wails version
```

### 8. Verify everything
```bash
wails doctor
```

`wails doctor` will tell you exactly if anything is missing. Paste the output here if anything shows as failed and we'll fix it.
