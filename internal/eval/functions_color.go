// internal/eval/functions_color.go

package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// COLOR CREATION
// ════════════════════════════════════════════════════════════════

// FnRGB creates a color from RGB components (0-255).
// Returns hex string representation.
func FnRGB(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("rgb requires 3 or 4 arguments: r, g, b, [a]")
	}

	r := clampInt(int(args[0].AsFloat()), 0, 255)
	g := clampInt(int(args[1].AsFloat()), 0, 255)
	b := clampInt(int(args[2].AsFloat()), 0, 255)

	if len(args) == 4 {
		a := clampFloat(args[3].AsFloat(), 0, 1)
		return types.StringValue(fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, a))
	}

	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

// FnRGBA creates a color from RGBA components.
func FnRGBA(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("rgba requires 4 arguments: r, g, b, a")
	}

	r := clampInt(int(args[0].AsFloat()), 0, 255)
	g := clampInt(int(args[1].AsFloat()), 0, 255)
	b := clampInt(int(args[2].AsFloat()), 0, 255)
	a := clampFloat(args[3].AsFloat(), 0, 1)

	return types.StringValue(fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, a))
}

// FnHSL creates a color from HSL components.
// h: 0-360, s: 0-100, l: 0-100
func FnHSL(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("hsl requires 3 or 4 arguments: h, s, l, [a]")
	}

	h := math.Mod(args[0].AsFloat(), 360)
	if h < 0 {
		h += 360
	}
	s := clampFloat(args[1].AsFloat(), 0, 100)
	l := clampFloat(args[2].AsFloat(), 0, 100)

	if len(args) == 4 {
		a := clampFloat(args[3].AsFloat(), 0, 1)
		return types.StringValue(fmt.Sprintf("hsla(%.0f, %.0f%%, %.0f%%, %.2f)", h, s, l, a))
	}

	// Convert to hex for consistency
	r, g, b := hslToRGB(h, s/100, l/100)
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

// FnHSLA creates a color from HSLA components.
func FnHSLA(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("hsla requires 4 arguments: h, s, l, a")
	}

	h := math.Mod(args[0].AsFloat(), 360)
	if h < 0 {
		h += 360
	}
	s := clampFloat(args[1].AsFloat(), 0, 100)
	l := clampFloat(args[2].AsFloat(), 0, 100)
	a := clampFloat(args[3].AsFloat(), 0, 1)

	return types.StringValue(fmt.Sprintf("hsla(%.0f, %.0f%%, %.0f%%, %.2f)", h, s, l, a))
}

// FnHSV creates a color from HSV/HSB components.
// h: 0-360, s: 0-100, v: 0-100
func FnHSV(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("hsv requires 3 arguments: h, s, v")
	}

	h := math.Mod(args[0].AsFloat(), 360)
	if h < 0 {
		h += 360
	}
	s := clampFloat(args[1].AsFloat(), 0, 100) / 100
	v := clampFloat(args[2].AsFloat(), 0, 100) / 100

	r, g, b := hsvToRGB(h, s, v)
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

// FnHex parses a hex color string and returns it normalized.
// Accepts: #RGB, #RRGGBB, #RRGGBBAA, RGB, RRGGBB
func FnColorHex(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("color requires exactly 1 argument")
	}

	hex := args[0].AsString()
	r, g, b, a, err := parseHexColor(hex)
	if err != nil {
		return types.Error(err.Error())
	}

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", r, g, b, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

// ════════════════════════════════════════════════════════════════
// COLOR CONVERSION
// ════════════════════════════════════════════════════════════════

// FnRGB2HSL converts RGB to HSL.
func FnRGB2HSL(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("rgb2hsl requires 3 arguments: r, g, b")
	}

	r := clampFloat(args[0].AsFloat(), 0, 255) / 255
	g := clampFloat(args[1].AsFloat(), 0, 255) / 255
	b := clampFloat(args[2].AsFloat(), 0, 255) / 255

	h, s, l := rgbToHSL(r, g, b)

	return types.StringValue(fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", h, s*100, l*100))
}

// FnHSL2RGB converts HSL to RGB.
func FnHSL2RGB(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("hsl2rgb requires 3 arguments: h, s, l")
	}

	h := math.Mod(args[0].AsFloat(), 360)
	if h < 0 {
		h += 360
	}
	s := clampFloat(args[1].AsFloat(), 0, 100) / 100
	l := clampFloat(args[2].AsFloat(), 0, 100) / 100

	r, g, b := hslToRGB(h, s, l)

	return types.StringValue(fmt.Sprintf("rgb(%d, %d, %d)", r, g, b))
}

// FnRGB2HSV converts RGB to HSV.
func FnRGB2HSV(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("rgb2hsv requires 3 arguments: r, g, b")
	}

	r := clampFloat(args[0].AsFloat(), 0, 255) / 255
	g := clampFloat(args[1].AsFloat(), 0, 255) / 255
	b := clampFloat(args[2].AsFloat(), 0, 255) / 255

	h, s, v := rgbToHSV(r, g, b)

	return types.StringValue(fmt.Sprintf("hsv(%.0f, %.0f%%, %.0f%%)", h, s*100, v*100))
}

// FnHSV2RGB converts HSV to RGB.
func FnHSV2RGB(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("hsv2rgb requires 3 arguments: h, s, v")
	}

	h := math.Mod(args[0].AsFloat(), 360)
	if h < 0 {
		h += 360
	}
	s := clampFloat(args[1].AsFloat(), 0, 100) / 100
	v := clampFloat(args[2].AsFloat(), 0, 100) / 100

	r, g, b := hsvToRGB(h, s, v)

	return types.StringValue(fmt.Sprintf("rgb(%d, %d, %d)", r, g, b))
}

// FnHex2RGB converts hex to RGB components.
func FnHex2RGB(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hex2rgb requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	return types.StringValue(fmt.Sprintf("rgb(%d, %d, %d)", r, g, b))
}

// FnRGB2Hex converts RGB to hex.
func FnRGB2Hex(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("rgb2hex requires 3 arguments: r, g, b")
	}

	r := clampInt(int(args[0].AsFloat()), 0, 255)
	g := clampInt(int(args[1].AsFloat()), 0, 255)
	b := clampInt(int(args[2].AsFloat()), 0, 255)

	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

// FnHex2HSL converts hex to HSL.
func FnHex2HSL(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hex2hsl requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	return types.StringValue(fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", h, s*100, l*100))
}

// ════════════════════════════════════════════════════════════════
// COLOR MANIPULATION
// ════════════════════════════════════════════════════════════════

// FnDarken darkens a color by a percentage.
func FnDarken(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("darken requires 2 arguments: color, amount")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	amount := args[1].AsFloat()
	if amount > 1 {
		amount = amount / 100 // Treat as percentage
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)
	l = clampFloat(l-amount, 0, 1)

	nr, ng, nb := hslToRGB(h, s, l)

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnLighten lightens a color by a percentage.
func FnLighten(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("lighten requires 2 arguments: color, amount")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	amount := args[1].AsFloat()
	if amount > 1 {
		amount = amount / 100
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)
	l = clampFloat(l+amount, 0, 1)

	nr, ng, nb := hslToRGB(h, s, l)

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnSaturate increases saturation by a percentage.
func FnSaturate(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("saturate requires 2 arguments: color, amount")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	amount := args[1].AsFloat()
	if amount > 1 {
		amount = amount / 100
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)
	s = clampFloat(s+amount, 0, 1)

	nr, ng, nb := hslToRGB(h, s, l)

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnDesaturate decreases saturation by a percentage.
func FnDesaturate(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("desaturate requires 2 arguments: color, amount")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	amount := args[1].AsFloat()
	if amount > 1 {
		amount = amount / 100
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)
	s = clampFloat(s-amount, 0, 1)

	nr, ng, nb := hslToRGB(h, s, l)

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnRotateHue rotates the hue by degrees.
func FnRotateHue(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rotatehue requires 2 arguments: color, degrees")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	degrees := args[1].AsFloat()

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)
	h = math.Mod(h+degrees, 360)
	if h < 0 {
		h += 360
	}

	nr, ng, nb := hslToRGB(h, s, l)

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnSetAlpha sets the alpha channel of a color.
func FnSetAlpha(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("setalpha requires 2 arguments: color, alpha")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	alpha := args[1].AsFloat()
	if alpha > 1 {
		alpha = alpha / 100
	}
	a := clampInt(int(alpha*255), 0, 255)

	return types.StringValue(fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, float64(a)/255))
}

// ════════════════════════════════════════════════════════════════
// COLOR UTILITIES
// ════════════════════════════════════════════════════════════════

// FnComplement returns the complementary color (180° hue rotation).
func FnComplement(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("complement requires exactly 1 argument")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)
	h = math.Mod(h+180, 360)

	nr, ng, nb := hslToRGB(h, s, l)

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnInvert inverts a color.
func FnInvert(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("invert requires exactly 1 argument")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	nr := 255 - r
	ng := 255 - g
	nb := 255 - b

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnGrayscale converts a color to grayscale.
func FnGrayscale(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("grayscale requires exactly 1 argument")
	}

	r, g, b, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	// Use luminosity method for perceptual grayscale
	gray := int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

	if a < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", gray, gray, gray, a))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", gray, gray, gray))
}

// FnBlend blends two colors together.
// Args: color1, color2, [weight] (0-1, default 0.5)
func FnBlend(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("blend requires 2 or 3 arguments: color1, color2, [weight]")
	}

	r1, g1, b1, a1, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	r2, g2, b2, a2, err := parseHexColor(args[1].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	weight := 0.5
	if len(args) == 3 {
		weight = args[2].AsFloat()
		if weight > 1 {
			weight = weight / 100
		}
		weight = clampFloat(weight, 0, 1)
	}

	nr := int(float64(r1)*(1-weight) + float64(r2)*weight)
	ng := int(float64(g1)*(1-weight) + float64(g2)*weight)
	nb := int(float64(b1)*(1-weight) + float64(b2)*weight)
	na := int(float64(a1)*(1-weight) + float64(a2)*weight)

	if na < 255 {
		return types.StringValue(fmt.Sprintf("#%02X%02X%02X%02X", nr, ng, nb, na))
	}
	return types.StringValue(fmt.Sprintf("#%02X%02X%02X", nr, ng, nb))
}

// FnContrast returns black or white depending on which has better contrast.
func FnContrast(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("contrast requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	// Calculate relative luminance
	luminance := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)

	if luminance > 128 {
		return types.StringValue("#000000")
	}
	return types.StringValue("#FFFFFF")
}

// FnContrastRatio calculates the contrast ratio between two colors.
func FnContrastRatio(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("contrastratio requires 2 arguments: color1, color2")
	}

	r1, g1, b1, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	r2, g2, b2, _, err := parseHexColor(args[1].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	l1 := relativeLuminance(r1, g1, b1)
	l2 := relativeLuminance(r2, g2, b2)

	if l1 > l2 {
		l1, l2 = l2, l1
	}

	ratio := (l2 + 0.05) / (l1 + 0.05)

	return types.Number(math.Round(ratio*100) / 100)
}

// ════════════════════════════════════════════════════════════════
// COLOR COMPONENT EXTRACTION
// ════════════════════════════════════════════════════════════════

// FnRed extracts the red component (0-255).
func FnRed(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("red requires exactly 1 argument")
	}

	r, _, _, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	return types.Number(float64(r))
}

// FnGreen extracts the green component (0-255).
func FnGreen(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("green requires exactly 1 argument")
	}

	_, g, _, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	return types.Number(float64(g))
}

// FnBlue extracts the blue component (0-255).
func FnBlue(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("blue requires exactly 1 argument")
	}

	_, _, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	return types.Number(float64(b))
}

// FnAlpha extracts the alpha component (0-1).
func FnAlpha(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("alpha requires exactly 1 argument")
	}

	_, _, _, a, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	return types.Number(float64(a) / 255)
}

// FnHue extracts the hue component (0-360).
func FnHue(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hue requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	h, _, _ := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	return types.Number(math.Round(h))
}

// FnSaturation extracts the saturation component (0-100).
func FnSaturation(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("saturation requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	_, s, _ := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	return types.Number(math.Round(s * 100))
}

// FnLightness extracts the lightness component (0-100).
func FnLightness(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lightness requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	_, _, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	return types.Number(math.Round(l * 100))
}

// FnLuminance calculates relative luminance (0-1).
func FnLuminance(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("luminance requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	lum := relativeLuminance(r, g, b)

	return types.Number(math.Round(lum*1000) / 1000)
}

// ════════════════════════════════════════════════════════════════
// COLOR PALETTE GENERATION
// ════════════════════════════════════════════════════════════════

// FnTriadic returns the triadic colors (120° apart).
func FnTriadic(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("triadic requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	h2 := math.Mod(h+120, 360)
	h3 := math.Mod(h+240, 360)

	r2, g2, b2 := hslToRGB(h2, s, l)
	r3, g3, b3 := hslToRGB(h3, s, l)

	return types.StringValue(fmt.Sprintf("#%02X%02X%02X, #%02X%02X%02X, #%02X%02X%02X",
		r, g, b, r2, g2, b2, r3, g3, b3))
}

// FnTetradic returns the tetradic colors (90° apart).
func FnTetradic(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tetradic requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	h2 := math.Mod(h+90, 360)
	h3 := math.Mod(h+180, 360)
	h4 := math.Mod(h+270, 360)

	r2, g2, b2 := hslToRGB(h2, s, l)
	r3, g3, b3 := hslToRGB(h3, s, l)
	r4, g4, b4 := hslToRGB(h4, s, l)

	return types.StringValue(fmt.Sprintf("#%02X%02X%02X, #%02X%02X%02X, #%02X%02X%02X, #%02X%02X%02X",
		r, g, b, r2, g2, b2, r3, g3, b3, r4, g4, b4))
}

// FnSplitComplement returns split-complementary colors (150° and 210°).
func FnSplitComplement(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("splitcomplement requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	h2 := math.Mod(h+150, 360)
	h3 := math.Mod(h+210, 360)

	r2, g2, b2 := hslToRGB(h2, s, l)
	r3, g3, b3 := hslToRGB(h3, s, l)

	return types.StringValue(fmt.Sprintf("#%02X%02X%02X, #%02X%02X%02X, #%02X%02X%02X",
		r, g, b, r2, g2, b2, r3, g3, b3))
}

// FnAnalogous returns analogous colors (30° apart).
func FnAnalogous(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("analogous requires exactly 1 argument")
	}

	r, g, b, _, err := parseHexColor(args[0].AsString())
	if err != nil {
		return types.Error(err.Error())
	}

	h, s, l := rgbToHSL(float64(r)/255, float64(g)/255, float64(b)/255)

	h2 := math.Mod(h+30, 360)
	h3 := math.Mod(h-30+360, 360)

	r2, g2, b2 := hslToRGB(h2, s, l)
	r3, g3, b3 := hslToRGB(h3, s, l)

	return types.StringValue(fmt.Sprintf("#%02X%02X%02X, #%02X%02X%02X, #%02X%02X%02X",
		r3, g3, b3, r, g, b, r2, g2, b2))
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// parseHexColor parses various hex color formats.
func parseHexColor(s string) (r, g, b, a int, err error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")
	s = strings.ToLower(s)

	a = 255 // Default alpha

	var r64, g64, b64, a64 int64

	switch len(s) {
	case 3: // RGB
		r64, _ = strconv.ParseInt(string(s[0])+string(s[0]), 16, 64)
		g64, _ = strconv.ParseInt(string(s[1])+string(s[1]), 16, 64)
		b64, _ = strconv.ParseInt(string(s[2])+string(s[2]), 16, 64)
		return int(r64), int(g64), int(b64), a, nil

	case 4: // RGBA
		r64, _ = strconv.ParseInt(string(s[0])+string(s[0]), 16, 64)
		g64, _ = strconv.ParseInt(string(s[1])+string(s[1]), 16, 64)
		b64, _ = strconv.ParseInt(string(s[2])+string(s[2]), 16, 64)
		a64, _ = strconv.ParseInt(string(s[3])+string(s[3]), 16, 64)
		return int(r64), int(g64), int(b64), int(a64), nil

	case 6: // RRGGBB
		r64, _ = strconv.ParseInt(s[0:2], 16, 64)
		g64, _ = strconv.ParseInt(s[2:4], 16, 64)
		b64, _ = strconv.ParseInt(s[4:6], 16, 64)
		return int(r64), int(g64), int(b64), a, nil

	case 8: // RRGGBBAA
		r64, _ = strconv.ParseInt(s[0:2], 16, 64)
		g64, _ = strconv.ParseInt(s[2:4], 16, 64)
		b64, _ = strconv.ParseInt(s[4:6], 16, 64)
		a64, _ = strconv.ParseInt(s[6:8], 16, 64)
		return int(r64), int(g64), int(b64), int(a64), nil

	default:
		return 0, 0, 0, 0, fmt.Errorf("invalid color format: %s", s)
	}
}

// rgbToHSL converts RGB (0-1) to HSL.
func rgbToHSL(r, g, b float64) (h, s, l float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	l = (max + min) / 2

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min

		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}

		h *= 60
	}

	return h, s, l
}

// hslToRGB converts HSL to RGB (0-255).
func hslToRGB(h, s, l float64) (r, g, b int) {
	if s == 0 {
		v := int(l * 255)
		return v, v, v
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	hk := h / 360

	tr := hk + 1.0/3.0
	tg := hk
	tb := hk - 1.0/3.0

	r = int(hueToRGB(p, q, tr) * 255)
	g = int(hueToRGB(p, q, tg) * 255)
	b = int(hueToRGB(p, q, tb) * 255)

	return r, g, b
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// rgbToHSV converts RGB (0-1) to HSV.
func rgbToHSV(r, g, b float64) (h, s, v float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	v = max
	d := max - min

	if max == 0 {
		s = 0
	} else {
		s = d / max
	}

	if max == min {
		h = 0
	} else {
		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h *= 60
	}

	return h, s, v
}

// hsvToRGB converts HSV to RGB (0-255).
func hsvToRGB(h, s, v float64) (r, g, b int) {
	if s == 0 {
		val := int(v * 255)
		return val, val, val
	}

	h = h / 60
	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	var rf, gf, bf float64
	switch int(i) % 6 {
	case 0:
		rf, gf, bf = v, t, p
	case 1:
		rf, gf, bf = q, v, p
	case 2:
		rf, gf, bf = p, v, t
	case 3:
		rf, gf, bf = p, q, v
	case 4:
		rf, gf, bf = t, p, v
	case 5:
		rf, gf, bf = v, p, q
	}

	return int(rf * 255), int(gf * 255), int(bf * 255)
}

// relativeLuminance calculates WCAG relative luminance.
func relativeLuminance(r, g, b int) float64 {
	rs := float64(r) / 255
	gs := float64(g) / 255
	bs := float64(b) / 255

	if rs <= 0.03928 {
		rs = rs / 12.92
	} else {
		rs = math.Pow((rs+0.055)/1.055, 2.4)
	}

	if gs <= 0.03928 {
		gs = gs / 12.92
	} else {
		gs = math.Pow((gs+0.055)/1.055, 2.4)
	}

	if bs <= 0.03928 {
		bs = bs / 12.92
	} else {
		bs = math.Pow((bs+0.055)/1.055, 2.4)
	}

	return 0.2126*rs + 0.7152*gs + 0.0722*bs
}

// clampInt clamps an integer to a range.
func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// clampFloat clamps a float to a range.
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
