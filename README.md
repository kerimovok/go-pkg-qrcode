# go-pkg-qrcode

A powerful Go package for generating QR codes with advanced customization options including gradient effects, logo embedding, and flexible styling.

## 🚀 Features

- **Advanced Customization**: Colors, gradients, borders, and sizing options
- **Logo Embedding**: Automatic logo download and embedding with size control
- **Gradient Effects**: Linear and radial gradient support for QR codes
- **Error Correction**: Configurable error correction levels (L, M, Q, H)
- **Flexible Sizing**: Customizable QR code dimensions
- **Color Support**: RGB, RGBA, and named color formats
- **Border Control**: Configurable border width and disable options
- **High Performance**: Fast QR code generation with optimized image processing

## 📦 Installation

```bash
go get github.com/kerimovok/go-pkg-qrcode
```

## 📖 Quick Start

```go
package main

import (
    "os"
    "github.com/kerimovok/go-pkg-qrcode"
)

func main() {
    // Generate a basic QR code
    png, err := qrcode.GeneratePNG(qrcode.Options{
        Data: "https://example.com",
        Size: 300,
    })
    if err != nil {
        panic(err)
    }
    
    os.WriteFile("qrcode.png", png, 0644)
}
```

## 🎯 Usage

### Basic Usage

```go
import "github.com/kerimovok/go-pkg-qrcode"

// Simple function call
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data: "https://example.com",
    Size: 300,
})
```

### Using Generator Instance

```go
generator := qrcode.New()

png, err := generator.GeneratePNG(qrcode.Options{
    Data: "https://example.com",
    Size: 300,
})
```

### Customized QR Code

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data:       "https://example.com",
    Size:       400,
    Foreground: "rgb(0,100,200)",
    Background: "rgb(240,240,240)",
    Error:      "H",
    Border:     8,
})
```

### QR Code with Gradient

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data:         "https://example.com",
    Size:         400,
    GradientStart: "rgb(255,0,0)",
    GradientEnd:   "rgb(0,0,255)",
    GradientType:  "linear", // or "radial"
    Error:        "H",
})
```

### QR Code with Logo

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data:     "https://example.com",
    Size:     500,
    LogoURL:  "https://example.com/logo.png",
    LogoSize: 25.0, // percentage of QR code size
    Error:    "H",
})
```

## ⚙️ Options

### Options Struct

```go
type Options struct {
    // Data is the content to encode in the QR code (required)
    Data string

    // Size is the QR code dimensions in pixels (default: 300)
    Size int

    // Foreground is the foreground color (QR code pattern)
    // Supports: rgb(r,g,b), rgba(r,g,b,a), or named colors
    // Default: black
    Foreground string

    // Background is the background color
    // Supports: rgb(r,g,b), rgba(r,g,b,a), or named colors
    // Default: white
    Background string

    // Error is the error correction level: L, M, Q, H
    // Default: M
    Error string

    // Border is the border width in pixels (0 = no border)
    // Default: 0
    Border int

    // LogoURL is the URL to a logo image to embed
    LogoURL string

    // LogoSize is the logo size as a percentage (default: 20.0)
    LogoSize float64

    // GradientStart is the start color for gradient effect
    GradientStart string

    // GradientEnd is the end color for gradient effect
    GradientEnd string

    // GradientType is the type of gradient: "linear" or "radial"
    GradientType string
}
```

### Color Formats

The package supports multiple color formats:

- **RGB**: `rgb(255,0,0)`
- **RGBA**: `rgba(255,0,0,128)`
- **Named Colors**: `black`, `white`, `red`, `green`, `blue`

### Error Correction Levels

| Level | Description            | Data Recovery |
|-------|------------------------|---------------|
| L     | Low                    | ~7%           |
| M     | Medium (default)       | ~15%          |
| Q     | High                   | ~25%          |
| H     | Highest | ~30%          |

### Gradient Types

- **linear**: Horizontal gradient from start to end color
- **radial**: Circular gradient from center outward

## 📚 Examples

### Example 1: Basic QR Code

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data: "https://example.com",
    Size: 300,
})
```

### Example 2: Customized QR Code

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data:       "https://example.com",
    Size:       400,
    Foreground: "rgb(0,100,200)",
    Background: "rgb(240,240,240)",
    Error:      "H",
    Border:     8,
})
```

### Example 3: QR Code with Gradient

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data:         "https://example.com",
    Size:         400,
    GradientStart: "rgb(255,0,0)",
    GradientEnd:   "rgb(0,0,255)",
    GradientType:  "radial",
    Error:        "Q",
})
```

### Example 4: QR Code with Logo

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data:     "https://example.com",
    Size:     500,
    LogoURL:  "https://example.com/logo.png",
    LogoSize: 30.0,
    Error:    "H",
})
```

### Example 5: Complex QR Code

```go
png, err := qrcode.GeneratePNG(qrcode.Options{
    Data:         "https://example.com",
    Size:         600,
    Foreground:    "rgb(50,50,50)",
    Background:    "rgb(250,250,250)",
    GradientStart: "rgb(100,200,255)",
    GradientEnd:   "rgb(255,100,200)",
    GradientType:  "linear",
    LogoURL:       "https://example.com/logo.png",
    LogoSize:      20.0,
    Error:         "H",
    Border:        12,
})
```

## 🔧 API Reference

### Functions

#### `GeneratePNG(opts Options) ([]byte, error)`

Convenience function that creates a generator and generates a QR code.

**Returns**: PNG image byte array and error

#### `New() *Generator`

Creates a new QR code generator instance.

### Methods

#### `(*Generator) GeneratePNG(opts Options) ([]byte, error)`

Generates a QR code as a PNG image byte array.

**Returns**: PNG image byte array and error

## 🛠️ Dependencies

- `github.com/skip2/go-qrcode` - QR code generation
- `github.com/disintegration/imaging` - Image processing

## 📝 License

[Your License Here]

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- Built on top of [go-qrcode](https://github.com/skip2/go-qrcode) for QR code generation
- Uses [imaging](https://github.com/disintegration/imaging) for high-quality image processing
- Designed for production reliability and developer experience

---

**Note**: This package requires Go 1.25 or later.
