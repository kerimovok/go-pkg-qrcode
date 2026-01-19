package qrcode

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/skip2/go-qrcode"
	_ "golang.org/x/image/webp"
)

// Options represents the configuration options for QR code generation
type Options struct {
	// Data is the content to encode in the QR code (required)
	Data string

	// Size is the QR code dimensions in pixels (default: 300)
	Size int

	// Foreground is the foreground color (QR code pattern)
	// Supports: rgb(r,g,b), rgba(r,g,b,a), or named colors (black, white, red, green, blue)
	// Default: black
	Foreground string

	// Background is the background color
	// Supports: rgb(r,g,b), rgba(r,g,b,a), or named colors (black, white, red, green, blue)
	// Default: white
	Background string

	// Error is the error correction level: L (Low ~7%), M (Medium ~15%), Q (High ~25%), H (Highest ~30%)
	// Default: M
	Error string

	// Border is the border width in pixels (0 = no border)
	// Default: 0
	Border int

	// LogoURL is the URL to a logo image to embed in the center of the QR code
	LogoURL string

	// LogoSize is the logo size as a percentage of the QR code (default: 20.0)
	LogoSize float64

	// GradientStart is the start color for gradient effect
	// Requires GradientEnd to be set
	GradientStart string

	// GradientEnd is the end color for gradient effect
	// Requires GradientStart to be set
	GradientEnd string

	// GradientType is the type of gradient: "linear" or "radial" (default: "linear")
	GradientType string
}

// Generator provides QR code generation functionality
type Generator struct{}

// New creates a new QR code generator
func New() *Generator {
	return &Generator{}
}

// GeneratePNG generates a QR code as a PNG image byte array
func (g *Generator) GeneratePNG(opts Options) ([]byte, error) {
	if opts.Data == "" {
		return nil, fmt.Errorf("data is required")
	}

	if opts.Size <= 0 {
		opts.Size = 300
	}
	if opts.Error == "" {
		opts.Error = "M"
	}
	if opts.Border < 0 {
		opts.Border = 0
	}
	if opts.LogoSize <= 0 {
		opts.LogoSize = 20.0
	}

	qr, err := qrcode.New(opts.Data, getErrorCorrection(opts.Error))
	if err != nil {
		return nil, fmt.Errorf("failed to init qrcode: %w", err)
	}

	qr.ForegroundColor = parseColor(opts.Foreground)
	qr.BackgroundColor = parseColor(opts.Background)

	if opts.Border == 0 {
		qr.DisableBorder = true
	} else {
		qr.DisableBorder = false
		extra := opts.Border - 4
		if extra > 0 {
			opts.Size += extra * 2
		}
	}

	var buf bytes.Buffer
	if err := qr.Write(opts.Size, &buf); err != nil {
		return nil, fmt.Errorf("failed to render qrcode: %w", err)
	}

	img, err := png.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to decode qrcode: %w", err)
	}

	if opts.GradientStart != "" && opts.GradientEnd != "" {
		start := parseColor(opts.GradientStart)
		end := parseColor(opts.GradientEnd)
		gradient := createGradient(img.Bounds().Dx(), img.Bounds().Dy(), start, end, opts.GradientType)
		finalImg := image.NewRGBA(img.Bounds())
		draw.Draw(finalImg, finalImg.Bounds(), gradient, image.Point{}, draw.Src)
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				fr, fg, fb, _ := qr.ForegroundColor.RGBA()
				if r == fr && g == fg && b == fb {
					finalImg.Set(x, y, gradient.At(x, y))
				} else {
					finalImg.Set(x, y, qr.BackgroundColor)
				}
			}
		}
		img = finalImg
	}

	if opts.LogoURL != "" {
		withLogo, err := embedLogo(img, opts.LogoURL, opts.LogoSize)
		if err != nil {
			return nil, fmt.Errorf("failed to embed logo: %w", err)
		}
		img = withLogo
	}

	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}
	return out.Bytes(), nil
}

// GeneratePNG is a convenience function that creates a generator and generates a QR code
func GeneratePNG(opts Options) ([]byte, error) {
	g := New()
	return g.GeneratePNG(opts)
}

func parseColor(colorStr string) color.Color {
	var r, g, b, a uint8 = 0, 0, 0, 255
	if n, err := fmt.Sscanf(colorStr, "rgb(%d,%d,%d)", &r, &g, &b); err == nil && n == 3 {
		return color.RGBA{R: r, G: g, B: b, A: a}
	}
	if n, err := fmt.Sscanf(colorStr, "rgba(%d,%d,%d,%d)", &r, &g, &b, &a); err == nil && n == 4 {
		return color.RGBA{R: r, G: g, B: b, A: a}
	}
	switch strings.ToLower(colorStr) {
	case "black":
		return color.Black
	case "white":
		return color.White
	case "red":
		return color.RGBA{R: 255, A: 255}
	case "green":
		return color.RGBA{G: 255, A: 255}
	case "blue":
		return color.RGBA{B: 255, A: 255}
	default:
		return color.Black
	}
}

func getErrorCorrection(level string) qrcode.RecoveryLevel {
	switch level {
	case "L":
		return qrcode.Low
	case "M":
		return qrcode.Medium
	case "Q":
		return qrcode.High
	case "H":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

func embedLogo(qrImage image.Image, logoURL string, sizePercent float64) (image.Image, error) {
	resp, err := http.Get(logoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logo: %w", err)
	}
	defer resp.Body.Close()

	// Use imaging.Decode which supports multiple formats (JPEG, PNG, GIF, WebP, etc.)
	// It will automatically detect the format regardless of Content-Type header
	logoImg, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode logo image: %w", err)
	}

	qrSize := qrImage.Bounds().Size()
	logoWidth := int(float64(qrSize.X) * sizePercent / 100)
	logoHeight := int(float64(qrSize.Y) * sizePercent / 100)
	logoImg = imaging.Fit(logoImg, logoWidth, logoHeight, imaging.Lanczos)
	finalImg := image.NewRGBA(qrImage.Bounds())
	draw.Draw(finalImg, finalImg.Bounds(), qrImage, image.Point{}, draw.Over)
	x := (qrSize.X - logoWidth) / 2
	y := (qrSize.Y - logoHeight) / 2
	logoPos := image.Rect(x, y, x+logoWidth, y+logoHeight)
	draw.Draw(finalImg, logoPos, logoImg, image.Point{}, draw.Over)
	return finalImg, nil
}

func createGradient(width, height int, startColor, endColor color.Color, gradientType string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	startR, startG, startB, _ := startColor.RGBA()
	endR, endG, endB, _ := endColor.RGBA()
	startR, startG, startB = startR>>8, startG>>8, startB>>8
	endR, endG, endB = endR>>8, endG>>8, endB>>8
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var ratio float64
			switch gradientType {
			case "linear":
				ratio = float64(x) / float64(width-1)
			case "radial":
				centerX, centerY := float64(width)/2, float64(height)/2
				distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
				maxDistance := math.Sqrt(math.Pow(centerX, 2) + math.Pow(centerY, 2))
				ratio = math.Min(distance/maxDistance, 1.0)
			default:
				ratio = float64(x) / float64(width-1)
			}
			r := uint8(float64(startR) + ratio*float64(int(endR)-int(startR)))
			g := uint8(float64(startG) + ratio*float64(int(endG)-int(startG)))
			b := uint8(float64(startB) + ratio*float64(int(endB)-int(startB)))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return img
}
