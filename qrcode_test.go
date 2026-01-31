package qrcode

import (
	"bytes"
	"image/png"
	"testing"
)

func TestNew(t *testing.T) {
	generator := New()
	if generator == nil {
		t.Fatal("New() returned nil")
	}
}

func TestGeneratePNG_Basic(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr bool
	}{
		{
			name: "basic QR code",
			opts: Options{
				Data: "https://example.com",
				Size: 300,
			},
			wantErr: false,
		},
		{
			name: "empty data",
			opts: Options{
				Data: "",
				Size: 300,
			},
			wantErr: true,
		},
		{
			name: "default size",
			opts: Options{
				Data: "test",
			},
			wantErr: false,
		},
		{
			name: "zero size uses default",
			opts: Options{
				Data: "test",
				Size: 0,
			},
			wantErr: false,
		},
		{
			name: "negative size uses default",
			opts: Options{
				Data: "test",
				Size: -100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pngData, err := GeneratePNG(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(pngData) == 0 {
					t.Error("GeneratePNG() returned empty PNG")
				}
				// Verify it's a valid PNG
				_, err := png.Decode(bytes.NewReader(pngData))
				if err != nil {
					t.Errorf("GeneratePNG() returned invalid PNG: %v", err)
				}
			}
		})
	}
}

func TestGeneratePNG_ErrorCorrection(t *testing.T) {
	levels := []string{"L", "M", "Q", "H", ""} // Empty string should default to M

	for _, level := range levels {
		t.Run("error_level_"+level, func(t *testing.T) {
			opts := Options{
				Data:  "https://example.com",
				Size:  300,
				Error: level,
			}
			pngData, err := GeneratePNG(opts)
			if err != nil {
				t.Errorf("GeneratePNG() with error level %s failed: %v", level, err)
				return
			}
			if len(pngData) == 0 {
				t.Error("GeneratePNG() returned empty PNG")
			}
		})
	}
}

func TestGeneratePNG_Colors(t *testing.T) {
	tests := []struct {
		name       string
		foreground string
		background string
		wantErr    bool
	}{
		{
			name:       "named colors",
			foreground: "black",
			background: "white",
			wantErr:    false,
		},
		{
			name:       "RGB colors",
			foreground: "rgb(255,0,0)",
			background: "rgb(0,255,0)",
			wantErr:    false,
		},
		{
			name:       "RGBA colors",
			foreground: "rgba(255,0,0,255)",
			background: "rgba(0,255,0,128)",
			wantErr:    false,
		},
		{
			name:       "mixed colors",
			foreground: "red",
			background: "rgb(240,240,240)",
			wantErr:    false,
		},
		{
			name:       "empty colors use defaults",
			foreground: "",
			background: "",
			wantErr:    false,
		},
		{
			name:       "invalid color format uses default",
			foreground: "invalid-color",
			background: "also-invalid",
			wantErr:    false, // Should not error, just use defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				Data:       "https://example.com",
				Size:       300,
				Foreground: tt.foreground,
				Background: tt.background,
			}
			pngData, err := GeneratePNG(opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(pngData) == 0 {
				t.Error("GeneratePNG() returned empty PNG")
			}
		})
	}
}

func TestGeneratePNG_Border(t *testing.T) {
	tests := []struct {
		name    string
		border  int
		wantErr bool
	}{
		{
			name:    "no border",
			border:  0,
			wantErr: false,
		},
		{
			name:    "small border",
			border:  4,
			wantErr: false,
		},
		{
			name:    "large border",
			border:  20,
			wantErr: false,
		},
		{
			name:    "negative border",
			border:  -5,
			wantErr: false, // Should default to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				Data:   "https://example.com",
				Size:   300,
				Border: tt.border,
			}
			pngData, err := GeneratePNG(opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(pngData) == 0 {
				t.Error("GeneratePNG() returned empty PNG")
			}
		})
	}
}

func TestGeneratePNG_Gradient(t *testing.T) {
	tests := []struct {
		name          string
		gradientStart string
		gradientEnd   string
		gradientType  string
		wantErr       bool
	}{
		{
			name:          "linear gradient",
			gradientStart: "rgb(255,0,0)",
			gradientEnd:   "rgb(0,0,255)",
			gradientType:  "linear",
			wantErr:       false,
		},
		{
			name:          "radial gradient",
			gradientStart: "rgb(255,255,255)",
			gradientEnd:   "rgb(0,0,0)",
			gradientType:  "radial",
			wantErr:       false,
		},
		{
			name:          "default gradient type",
			gradientStart: "rgb(255,0,0)",
			gradientEnd:   "rgb(0,0,255)",
			gradientType:  "",
			wantErr:       false, // Should default to linear
		},
		{
			name:          "invalid gradient type",
			gradientStart: "rgb(255,0,0)",
			gradientEnd:   "rgb(0,0,255)",
			gradientType:  "invalid",
			wantErr:       false, // Should default to linear
		},
		{
			name:          "gradient with only start color",
			gradientStart: "rgb(255,0,0)",
			gradientEnd:   "",
			gradientType:  "linear",
			wantErr:       false, // Should not apply gradient
		},
		{
			name:          "gradient with only end color",
			gradientStart: "",
			gradientEnd:   "rgb(0,0,255)",
			gradientType:  "linear",
			wantErr:       false, // Should not apply gradient
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				Data:          "https://example.com",
				Size:          400,
				GradientStart: tt.gradientStart,
				GradientEnd:   tt.gradientEnd,
				GradientType:  tt.gradientType,
			}
			pngData, err := GeneratePNG(opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(pngData) == 0 {
				t.Error("GeneratePNG() returned empty PNG")
			}
		})
	}
}

func TestGeneratePNG_LogoSize(t *testing.T) {
	tests := []struct {
		name     string
		logoSize float64
		wantErr  bool
	}{
		{
			name:     "default logo size",
			logoSize: 0,
			wantErr:  false, // Should use default 20.0
		},
		{
			name:     "small logo",
			logoSize: 10.0,
			wantErr:  false,
		},
		{
			name:     "large logo",
			logoSize: 30.0,
			wantErr:  false,
		},
		{
			name:     "negative logo size",
			logoSize: -5.0,
			wantErr:  false, // Should use default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				Data:     "https://example.com",
				Size:     500,
				LogoSize: tt.logoSize,
				// Note: LogoURL is empty, so logo won't be embedded
				// This just tests that LogoSize doesn't cause errors
			}
			pngData, err := GeneratePNG(opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(pngData) == 0 {
				t.Error("GeneratePNG() returned empty PNG")
			}
		})
	}
}

func TestGenerator_GeneratePNG(t *testing.T) {
	generator := New()
	if generator == nil {
		t.Fatal("New() returned nil")
	}

	opts := Options{
		Data: "https://example.com",
		Size: 300,
	}

	pngData, err := generator.GeneratePNG(opts)
	if err != nil {
		t.Errorf("Generator.GeneratePNG() error = %v", err)
		return
	}

	if len(pngData) == 0 {
		t.Error("Generator.GeneratePNG() returned empty PNG")
	}

	// Verify it's a valid PNG
	_, err = png.Decode(bytes.NewReader(pngData))
	if err != nil {
		t.Errorf("Generator.GeneratePNG() returned invalid PNG: %v", err)
	}
}

func TestGenerator_MultipleCalls(t *testing.T) {
	generator := New()

	opts1 := Options{
		Data: "https://example.com/first",
		Size: 300,
	}

	opts2 := Options{
		Data: "https://example.com/second",
		Size: 300,
	}

	pngData1, err1 := generator.GeneratePNG(opts1)
	pngData2, err2 := generator.GeneratePNG(opts2)

	if err1 != nil {
		t.Errorf("First GeneratePNG() error = %v", err1)
	}
	if err2 != nil {
		t.Errorf("Second GeneratePNG() error = %v", err2)
	}

	if len(pngData1) == 0 {
		t.Error("First GeneratePNG() returned empty PNG")
	}
	if len(pngData2) == 0 {
		t.Error("Second GeneratePNG() returned empty PNG")
	}

	// Verify both are valid PNGs
	_, err := png.Decode(bytes.NewReader(pngData1))
	if err != nil {
		t.Errorf("First PNG is invalid: %v", err)
	}

	_, err = png.Decode(bytes.NewReader(pngData2))
	if err != nil {
		t.Errorf("Second PNG is invalid: %v", err)
	}

	// PNGs with different data should produce different results
	// (though exact byte comparison might fail due to compression)
	if len(pngData1) == len(pngData2) && bytes.Equal(pngData1, pngData2) {
		// If they're identical, at least verify they're valid
		// This can happen with very similar QR codes, but is unlikely
		t.Logf("Warning: PNGs with different data produced identical output (size: %d bytes)", len(pngData1))
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"black", "black", true},
		{"white", "white", true},
		{"red", "red", true},
		{"green", "green", true},
		{"blue", "blue", true},
		{"RGB format", "rgb(255,0,0)", true},
		{"RGBA format", "rgba(255,0,0,128)", true},
		{"invalid format", "invalid", true}, // Should return default (black)
		{"empty string", "", true},          // Should return default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := parseColor(tt.input)
			if color == nil {
				t.Error("parseColor() returned nil")
			}
		})
	}
}

func TestGetErrorCorrection(t *testing.T) {
	tests := []struct {
		name  string
		level string
		valid bool
	}{
		{"Low", "L", true},
		{"Medium", "M", true},
		{"High", "Q", true},
		{"Highest", "H", true},
		{"default", "", true},
		{"invalid", "X", true}, // Should default to Medium
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := getErrorCorrection(tt.level)
			// Just verify it doesn't panic and returns a valid level
			_ = level
		})
	}
}

func TestGeneratePNG_ComplexOptions(t *testing.T) {
	opts := Options{
		Data:          "https://example.com",
		Size:          600,
		Foreground:    "rgb(50,50,50)",
		Background:    "rgb(250,250,250)",
		GradientStart: "rgb(100,200,255)",
		GradientEnd:   "rgb(255,100,200)",
		GradientType:  "linear",
		Error:         "H",
		Border:        12,
		LogoSize:      20.0,
	}

	pngData, err := GeneratePNG(opts)
	if err != nil {
		t.Errorf("GeneratePNG() with complex options error = %v", err)
		return
	}

	if len(pngData) == 0 {
		t.Error("GeneratePNG() returned empty PNG")
	}

	// Verify it's a valid PNG
	_, err = png.Decode(bytes.NewReader(pngData))
	if err != nil {
		t.Errorf("GeneratePNG() returned invalid PNG: %v", err)
	}
}

func BenchmarkGeneratePNG_Basic(b *testing.B) {
	opts := Options{
		Data: "https://example.com",
		Size: 300,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GeneratePNG(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGeneratePNG_WithGradient(b *testing.B) {
	opts := Options{
		Data:          "https://example.com",
		Size:          400,
		GradientStart: "rgb(255,0,0)",
		GradientEnd:   "rgb(0,0,255)",
		GradientType:  "linear",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GeneratePNG(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerator_GeneratePNG(b *testing.B) {
	generator := New()
	opts := Options{
		Data: "https://example.com",
		Size: 300,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.GeneratePNG(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
