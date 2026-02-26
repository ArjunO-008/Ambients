package audio

import (
	"context"
	"math"
	"math/cmplx"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	sampleRate = 48000
	frameSize  = 1024
	numBands   = 64
	emitEveryN = 3
)

type AudioService struct {
	ctx     context.Context
	startCh chan struct{}
	stopCh  chan struct{}
	running bool
	mu      sync.Mutex
}

func NewAudioService(ctx context.Context) *AudioService {
	a := &AudioService{
		ctx:     ctx,
		startCh: make(chan struct{}, 1),
		stopCh:  make(chan struct{}, 1),
	}
	// audio thread starts waiting immediately —
	// it won't do anything until Start() sends to startCh
	go a.audioThread()
	return a
}

func (a *AudioService) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return nil
	}

	a.running = true
	a.startCh <- struct{}{} // wake up the audio thread
	return nil
}

func (a *AudioService) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		a.stopCh <- struct{}{}
		a.running = false
	}
}

// audioThread runs on its own OS thread.
// PortAudio requires this on some platforms —
// keeping it here prevents crashes in the Wails UI thread.
func (a *AudioService) audioThread() {
	// LockOSThread ensures this goroutine always runs on the same OS thread.
	// PortAudio's WASAPI backend requires this.
	// runtime.LockOSThread() -- only needed if CGO issues arise, skip for now

	for range a.startCh {
		a.runCapture()
	}
}

// runCapture does the actual PortAudio work —
// initialize, find device, open stream, read frames, emit bands.
// Separated from audioThread so it can cleanly return on errors.
func (a *AudioService) runCapture() {
	if err := portaudio.Initialize(); err != nil {
		// emit error event so React can show it
		wailsRuntime.EventsEmit(a.ctx, "audio:error", "PortAudio init failed: "+err.Error())
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
		return
	}
	defer portaudio.Terminate()

	device, err := findLoopbackDevice()
	if err != nil {
		wailsRuntime.EventsEmit(a.ctx, "audio:error", err.Error())
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
		return
	}

	buffer := make([]float32, frameSize)

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   device,
			Channels: 1, // mono — enough for visualization
			Latency:  device.DefaultLowInputLatency,
		},
		SampleRate:      float64(sampleRate),
		FramesPerBuffer: frameSize,
	}, buffer)
	if err != nil {
		wailsRuntime.EventsEmit(a.ctx, "audio:error", "Stream open failed: "+err.Error())
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
		return
	}

	if err := stream.Start(); err != nil {
		stream.Close()
		wailsRuntime.EventsEmit(a.ctx, "audio:error", "Stream start failed: "+err.Error())
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
		return
	}

	defer stream.Close()

	frameCount := 0

	for {
		select {
		case <-a.stopCh:
			// Stop() was called — shut down cleanly
			stream.Stop()
			return

		default:
			if err := stream.Read(); err != nil {
				// small sleep prevents CPU spin on read errors
				time.Sleep(5 * time.Millisecond)
				continue
			}

			frameCount++
			if frameCount%emitEveryN != 0 {
				continue
			}

			bands := computeBands(buffer)
			wailsRuntime.EventsEmit(a.ctx, "audio:frame", bands)
		}
	}
}

func computeBands(samples []float32) []float64 {
	n := len(samples)
	c := make([]complex128, n)
	for i, s := range samples {
		// Hanning window reduces frequency bleed between bands
		window := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(n-1)))
		c[i] = complex(float64(s)*window, 0)
	}
	fft(c)

	half := n / 2
	bandSize := half / numBands
	bands := make([]float64, numBands)

	for i := 0; i < numBands; i++ {
		start := i * bandSize
		end := start + bandSize
		var sum float64
		for j := start; j < end; j++ {
			sum += cmplx.Abs(c[j])
		}
		avg := sum / float64(bandSize)
		// convert to dB scale — closer to how human ears perceive loudness
		db := 20 * math.Log10(avg+1e-10)
		normalized := (db + 80) / 80
		bands[i] = math.Max(0, math.Min(1, normalized))
	}
	return bands
}

func fft(a []complex128) {
	n := len(a)
	if n <= 1 {
		return
	}
	even := make([]complex128, n/2)
	odd := make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		even[i] = a[2*i]
		odd[i] = a[2*i+1]
	}
	fft(even)
	fft(odd)
	for k := 0; k < n/2; k++ {
		angle := -2 * math.Pi * float64(k) / float64(n)
		t := cmplx.Exp(complex(0, angle)) * odd[k]
		a[k] = even[k] + t
		a[k+n/2] = even[k] - t
	}
}
