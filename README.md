# godesktopdup

**Version:** v0.1.0

High-performance library for screen capture via Desktop Duplication API (DDA) on Windows. Returns frames in BGRA format with GPU acceleration and optimized CPU usage.

## Features

- **GPU-accelerated capture** via Desktop Duplication API
- **Cursor capture** with automatic rendering on frames
- **Multi-monitor support** with proper coordinate handling
- **Optimized performance** with dirty rects and efficient memory operations
- **Rotated monitor support** (90°, 180°, 270°)
- **Low CPU usage** with hardware-accelerated operations

## Requirements

- Windows 8 or higher
- DirectX 11.1+ compatible GPU
- Go 1.21 or higher

## Installation

```bash
go get github.com/shinkar94/godesktopdup@v0.1.0
```

Or use a local version:

```go
// In your go.mod
replace github.com/shinkar94/godesktopdup => ../godesktopdup

require github.com/shinkar94/godesktopdup v0.1.0
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/shinkar94/godesktopdup"
)

func main() {
    // Create capture for monitor 0
    dd, err := dda.New(0)
    if err != nil {
        panic(err)
    }
    defer dd.Release()

    // Get screen size
    width, height, err := dd.GetSize()
    if err != nil {
        panic(err)
    }

    // Create buffer for frame (BGRA = 4 bytes per pixel)
    buffer := make([]byte, width*height*4)

    // Capture frame (timeout in milliseconds)
    err = dd.GetFrameBGRA(buffer, 1000)
    if err != nil {
        panic(err)
    }

    // buffer now contains data in BGRA format
    fmt.Printf("Captured %dx%d frame\n", width, height)
}
```

## Working with FPS

The library doesn't control FPS directly - it captures frames as fast as Windows provides them. To achieve a target FPS, calculate the timeout based on your desired frame rate:

```go
package main

import (
    "time"
    "github.com/shinkar94/godesktopdup"
)

func captureAtFPS(dd *dda.DesktopDuplication, targetFPS int) error {
    width, height, err := dd.GetSize()
    if err != nil {
        return err
    }

    buffer := make([]byte, width*height*4)
    
    // Calculate timeout: half of frame interval for quick response
    // Example: 30 FPS → 33ms interval → 16ms timeout
    //          60 FPS → 16ms interval → 8ms timeout
    frameInterval := time.Duration(1000/targetFPS) * time.Millisecond
    timeoutMs := uint(frameInterval.Milliseconds() / 2)
    if timeoutMs < 5 {
        timeoutMs = 5 // Minimum 5ms
    }

    // Capture loop
    ticker := time.NewTicker(frameInterval)
    defer ticker.Stop()

    for range ticker.C {
        err := dd.GetFrameBGRA(buffer, timeoutMs)
        if err != nil {
            if err.Error() == "no image yet" {
                // Screen hasn't changed, skip this frame
                continue
            }
            return err
        }
        
        // Process frame here
        processFrame(buffer, width, height)
    }
    
    return nil
}

func processFrame(buffer []byte, width, height int) {
    // Your frame processing logic
}
```

### Recommended FPS Settings

- **30 FPS**: `timeoutMs = 16` (good balance between quality and performance)
- **60 FPS**: `timeoutMs = 8` (smooth video, higher CPU usage)
- **15 FPS**: `timeoutMs = 33` (lower CPU usage, acceptable for static content)

## Cursor Capture

Enable cursor capture to automatically render the mouse cursor on frames:

```go
package main

import "github.com/shinkar94/godesktopdup"

func main() {
    dd, err := dda.New(0)
    if err != nil {
        panic(err)
    }
    defer dd.Release()

    // Enable cursor capture
    dd.SetCaptureCursor(true)

    width, height, _ := dd.GetSize()
    buffer := make([]byte, width*height*4)

    // Cursor will be automatically drawn on captured frames
    err = dd.GetFrameBGRA(buffer, 1000)
    if err != nil {
        panic(err)
    }
    
    // buffer contains frame with cursor rendered
}
```

**Note**: Cursor capture adds a small CPU overhead. Disable it if you don't need cursor visibility:

```go
dd.SetCaptureCursor(false) // Disable cursor capture
```

## Multi-Monitor Support

Capture from multiple monitors:

```go
package main

import (
    "fmt"
    "github.com/shinkar94/godesktopdup"
)

func main() {
    // Capture from monitor 0 (primary)
    dd0, err := dda.New(0)
    if err != nil {
        panic(err)
    }
    defer dd0.Release()

    // Capture from monitor 1 (secondary)
    dd1, err := dda.New(1)
    if err != nil {
        panic(err)
    }
    defer dd1.Release()

    width0, height0, _ := dd0.GetSize()
    width1, height1, _ := dd1.GetSize()

    buffer0 := make([]byte, width0*height0*4)
    buffer1 := make([]byte, width1*height1*4)

    // Capture from both monitors
    err = dd0.GetFrameBGRA(buffer0, 1000)
    if err != nil {
        panic(err)
    }

    err = dd1.GetFrameBGRA(buffer1, 1000)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Monitor 0: %dx%d\n", width0, height0)
    fmt.Printf("Monitor 1: %dx%d\n", width1, height1)
}
```

## Monitor Bounds

Set monitor bounds for proper cursor positioning when using multiple monitors:

```go
package main

import "github.com/shinkar94/godesktopdup"

func main() {
    dd, err := dda.New(0)
    if err != nil {
        panic(err)
    }
    defer dd.Release()

    // Get monitor bounds
    left, top, right, bottom, err := dd.GetBounds()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Monitor bounds: (%d, %d) to (%d, %d)\n", left, top, right, bottom)

    // Set monitor bounds manually (useful for multi-monitor setups)
    // This helps with cursor positioning across monitors
    dd.SetMonitorBounds(1920, 0, 3840, 1080) // Example: second monitor at 1920x0

    // Enable cursor capture - cursor will be positioned correctly
    dd.SetCaptureCursor(true)
}
```

## Error Handling

The library returns specific errors that you should handle:

```go
package main

import (
    "errors"
    "fmt"
    "github.com/shinkar94/godesktopdup"
)

func captureFrame(dd *dda.DesktopDuplication) error {
    width, height, err := dd.GetSize()
    if err != nil {
        return fmt.Errorf("failed to get size: %w", err)
    }

    buffer := make([]byte, width*height*4)
    err = dd.GetFrameBGRA(buffer, 1000)

    if err != nil {
        if errors.Is(err, errors.New("no image yet")) {
            // Screen hasn't changed since last capture
            // This is normal and not an error
            return nil
        }
        return fmt.Errorf("capture failed: %w", err)
    }

    return nil
}
```

### Common Errors

- **"no image yet"**: Screen hasn't changed since last capture. This is normal and not an error.
- **"timeout waiting for frame"**: Timeout expired. Increase timeout or check if screen is updating.
- **"outputDuplication is nil"**: Capture object was released or not initialized properly.

## Complete Example: Screen Capture Service

```go
package main

import (
    "fmt"
    "time"
    "github.com/shinkar94/godesktopdup"
)

type CaptureService struct {
    dd        *dda.DesktopDuplication
    width     int
    height    int
    buffer    []byte
    targetFPS int
}

func NewCaptureService(monitorIndex uint, targetFPS int, captureCursor bool) (*CaptureService, error) {
    dd, err := dda.New(monitorIndex)
    if err != nil {
        return nil, fmt.Errorf("failed to create capture: %w", err)
    }

    width, height, err := dd.GetSize()
    if err != nil {
        dd.Release()
        return nil, fmt.Errorf("failed to get size: %w", err)
    }

    dd.SetCaptureCursor(captureCursor)

    return &CaptureService{
        dd:        dd,
        width:     width,
        height:    height,
        buffer:    make([]byte, width*height*4),
        targetFPS: targetFPS,
    }, nil
}

func (cs *CaptureService) Start() error {
    frameInterval := time.Duration(1000/cs.targetFPS) * time.Millisecond
    timeoutMs := uint(frameInterval.Milliseconds() / 2)
    if timeoutMs < 5 {
        timeoutMs = 5
    }

    ticker := time.NewTicker(frameInterval)
    defer ticker.Stop()

    for range ticker.C {
        err := cs.dd.GetFrameBGRA(cs.buffer, timeoutMs)
        if err != nil {
            if err.Error() == "no image yet" {
                continue
            }
            return fmt.Errorf("capture failed: %w", err)
        }

        // Process frame
        cs.onFrame(cs.buffer)
    }

    return nil
}

func (cs *CaptureService) onFrame(buffer []byte) {
    // Your frame processing logic here
    // For example: encode to video, send over network, etc.
    fmt.Printf("Frame captured: %dx%d\n", cs.width, cs.height)
}

func (cs *CaptureService) Release() {
    if cs.dd != nil {
        cs.dd.Release()
    }
}

func main() {
    service, err := NewCaptureService(0, 30, true)
    if err != nil {
        panic(err)
    }
    defer service.Release()

    // Capture for 10 seconds
    go func() {
        time.Sleep(10 * time.Second)
    }()

    if err := service.Start(); err != nil {
        panic(err)
    }
}
```

## API Reference

### New(outputIndex uint) (*DesktopDuplication, error)

Creates a new capture instance for the specified monitor.

- `outputIndex`: Monitor index (0 for first monitor, 1 for second, etc.)
- Returns: `*DesktopDuplication` instance or error

### GetFrameBGRA(buffer []byte, timeoutMs uint) error

Captures a screen frame and writes it to the buffer in BGRA format.

- `buffer`: Pre-allocated buffer (must be at least `width * height * 4` bytes)
- `timeoutMs`: Timeout in milliseconds (recommended: 16ms for 30 FPS, 8ms for 60 FPS)
- Returns: Error if capture fails

**Note**: Returns `"no image yet"` error if screen hasn't changed. This is normal and not a failure.

### GetSize() (width, height int, error)

Returns the size of the captured screen in pixels.

### GetBounds() (left, top, right, bottom int, error)

Returns the monitor bounds in desktop coordinates.

### SetMonitorBounds(left, top, right, bottom int)

Sets monitor bounds for proper cursor positioning in multi-monitor setups.

### SetCaptureCursor(enabled bool)

Enables or disables cursor capture. When enabled, the mouse cursor is automatically rendered on captured frames.

### Release()

Releases all resources associated with the capture. Always call this when done.

## Data Format

Data is returned in **BGRA format** (Blue-Green-Red-Alpha), where each pixel is represented by 4 bytes:

- **Byte 0**: Blue (0-255)
- **Byte 1**: Green (0-255)
- **Byte 2**: Red (0-255)
- **Byte 3**: Alpha (0-255, usually 255 for opaque pixels)

Data is arranged row by row, starting from the top-left corner.

### Converting BGRA to RGBA

If you need RGBA format:

```go
func convertBGRAtoRGBA(bgra, rgba []byte) {
    for i := 0; i < len(bgra); i += 4 {
        rgba[i] = bgra[i+2]     // R
        rgba[i+1] = bgra[i+1]   // G
        rgba[i+2] = bgra[i]     // B
        rgba[i+3] = bgra[i+3]   // A
    }
}
```

## Performance Tips

1. **Reuse buffers**: Allocate buffer once and reuse it for all captures
2. **Optimize timeout**: Use shorter timeouts for higher FPS, but not too short (minimum 5ms)
3. **Disable cursor**: If you don't need cursor, disable it to save CPU
4. **Handle "no image yet"**: This is normal - screen hasn't changed, skip processing
5. **Keep instance alive**: Don't create new instances for each frame - reuse the same instance

## Troubleshooting

### "Failed to create device"

- Ensure DirectX 11.1+ is available
- Check GPU drivers are up to date
- Verify Windows 8+ is installed

### "no image yet" errors

This is normal - the screen hasn't changed. Handle it gracefully:

```go
err := dd.GetFrameBGRA(buffer, timeoutMs)
if err != nil && err.Error() != "no image yet" {
    return err
}
```

### Low FPS

- Increase timeout value
- Check if screen is actually updating
- Verify GPU is not overloaded
- Consider disabling cursor capture

### Cursor not visible

- Ensure `SetCaptureCursor(true)` is called
- Check if cursor is within monitor bounds
- Verify monitor bounds are set correctly for multi-monitor setups

## License

See [LICENSE](LICENSE) file for details.

