package audio

import (
	"context"
	"math"
	"math/cmplx"
	"runtime"
	"sync"

	"github.com/gordonklaus/portaudio"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	sampleRate = 44100
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

	go a.audioThread()
	return a
}

// SAFE: called from Wails
func (a *AudioService) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return nil
	}

	a.startCh <- struct{}{}
	a.running = true
	return nil
}

// SAFE: called from Wails
func (a *AudioService) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		a.stopCh <- struct{}{}
		a.running = false
	}
}

// Native audio lives HERE — not in Wails thread
func (a *AudioService) audioThread() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for range a.startCh {

		if err := portaudio.Initialize(); err != nil {
			a.running = false
			continue
		}

		device, err := findLoopbackDevice()
		if err != nil {
			portaudio.Terminate()
			a.running = false
			continue
		}

		buffer := make([]float32, frameSize)

		stream, err := portaudio.OpenStream(portaudio.StreamParameters{
			Input: portaudio.StreamDeviceParameters{
				Device:   device,
				Channels: 1,
				Latency:  device.DefaultLowInputLatency,
			},
			SampleRate:      float64(sampleRate),
			FramesPerBuffer: frameSize,
		}, buffer)

		if err != nil {
			portaudio.Terminate()
			a.running = false
			continue
		}

		if err := stream.Start(); err != nil {
			stream.Close()
			portaudio.Terminate()
			a.running = false
			continue
		}

		frameCount := 0

	audioLoop:
		for {
			select {
			case <-a.stopCh:
				stream.Stop()
				break audioLoop
			default:
				if err := stream.Read(); err != nil {
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

		stream.Close()
		portaudio.Terminate()
	}
}

func computeBands(samples []float32) []float64 {
	n := len(samples)
	c := make([]complex128, n)
	for i, s := range samples {
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
