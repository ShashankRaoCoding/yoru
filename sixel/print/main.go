package print

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"

	_ "golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	"yoru/utils"
)

// Main reads an image from stdin and writes its SIXEL representation to stdout.
// Optional positional arguments set the target width and height in pixels.
func Main(args []string) {
	fs := flag.NewFlagSet("print", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: yoru sixel print [width] [height]\n")
		fmt.Fprintf(os.Stderr, "\nReads an image from stdin and outputs SIXEL escape sequences.\n")
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  width   Optional target width in pixels\n")
		fmt.Fprintf(os.Stderr, "  height  Optional target height in pixels\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  cat image.png | yoru sixel print\n")
		fmt.Fprintf(os.Stderr, "  cat image.jpg | yoru sixel print 100 50\n")
	}

	err := fs.Parse(args)
	utils.Error(err)

	remaining := fs.Args()
	var targetWidth, targetHeight int

	if len(remaining) >= 1 {
		targetWidth, err = strconv.Atoi(remaining[0])
		if err != nil {
			utils.Error(fmt.Errorf("invalid width %q: %w", remaining[0], err))
		}
	}
	if len(remaining) >= 2 {
		targetHeight, err = strconv.Atoi(remaining[1])
		if err != nil {
			utils.Error(fmt.Errorf("invalid height %q: %w", remaining[1], err))
		}
	}

	img, _, err := image.Decode(os.Stdin)
	if err != nil {
		utils.Error(fmt.Errorf("decoding image: %w", err))
	}

	if targetWidth > 0 || targetHeight > 0 {
		img = scaleImage(img, targetWidth, targetHeight)
	}

	// Flatten alpha against a white background before quantisation.
	img = flattenAlpha(img)

	sixelData := encode(img)
	_, err = fmt.Fprint(os.Stdout, sixelData)
	utils.Error(err)
}

// flattenAlpha composites img over a white background so that transparent
// pixels do not confuse the palette quantiser.
func flattenAlpha(img image.Image) image.Image {
	bounds := img.Bounds()
	white := image.NewRGBA(bounds)
	draw.Draw(white, bounds, &image.Uniform{color.White}, bounds.Min, draw.Src)
	draw.Draw(white, bounds, img, bounds.Min, draw.Over)
	return white
}

// scaleImage resizes img to (targetWidth × targetHeight) using bilinear
// interpolation.  If only one dimension is given the other is computed to
// preserve the aspect ratio.
func scaleImage(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	if origW == 0 || origH == 0 {
		return img
	}

	switch {
	case targetWidth <= 0 && targetHeight > 0:
		targetWidth = origW * targetHeight / origH
	case targetHeight <= 0 && targetWidth > 0:
		targetHeight = origH * targetWidth / origW
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	xdraw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, xdraw.Over, nil)
	return dst
}

// xtermPalette builds the standard xterm 256-colour palette.
//
//   Indices  0–15  : standard ANSI colours
//   Indices 16–231 : 6×6×6 RGB colour cube
//   Indices 232–255: 24-step greyscale ramp
func xtermPalette() color.Palette {
	pal := make(color.Palette, 256)

	// Standard ANSI colours (indices 0–15)
	ansi := [16]color.RGBA{
		{0, 0, 0, 255},
		{128, 0, 0, 255},
		{0, 128, 0, 255},
		{128, 128, 0, 255},
		{0, 0, 128, 255},
		{128, 0, 128, 255},
		{0, 128, 128, 255},
		{192, 192, 192, 255},
		{128, 128, 128, 255},
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{255, 255, 0, 255},
		{0, 0, 255, 255},
		{255, 0, 255, 255},
		{0, 255, 255, 255},
		{255, 255, 255, 255},
	}
	for i, c := range ansi {
		pal[i] = c
	}

	// 6×6×6 colour cube (indices 16–231)
	cubeLevel := [6]uint8{0, 95, 135, 175, 215, 255}
	for i := 0; i < 216; i++ {
		r := cubeLevel[i/36%6]
		g := cubeLevel[i/6%6]
		b := cubeLevel[i%6]
		pal[16+i] = color.RGBA{r, g, b, 255}
	}

	// 24-step greyscale ramp (indices 232–255)
	for i := 0; i < 24; i++ {
		v := uint8(8 + 10*i)
		pal[232+i] = color.RGBA{v, v, v, 255}
	}

	return pal
}

// encode converts img to a SIXEL escape sequence string.
func encode(img image.Image) string {
	pal := xtermPalette()

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Quantise the image to the 256-colour palette with Floyd-Steinberg dithering.
	paletted := image.NewPaletted(image.Rect(0, 0, w, h), pal)
	draw.FloydSteinberg.Draw(paletted, paletted.Bounds(), img, bounds.Min)

	var sb strings.Builder

	// DCS header — Pq starts sixel data with default parameters.
	sb.WriteString("\x1bPq")

	// Emit colour definitions for all 256 palette entries.
	for i, c := range pal {
		r16, g16, b16, _ := c.RGBA()
		// Convert 0–65535 to 0–100 percent as required by the SIXEL spec.
		rp := int(r16) * 100 / 65535
		gp := int(g16) * 100 / 65535
		bp := int(b16) * 100 / 65535
		fmt.Fprintf(&sb, "#%d;2;%d;%d;%d", i, rp, gp, bp)
	}

	pix := paletted.Pix
	stride := paletted.Stride

	// Process the image in horizontal bands of 6 rows (one sixel row).
	for y := 0; y < h; y += 6 {
		bandH := 6
		if y+bandH > h {
			bandH = h - y
		}

		// Pre-read the pixel indices for this band to avoid repeated method calls.
		band := make([][]uint8, bandH)
		for dy := 0; dy < bandH; dy++ {
			row := pix[(y+dy)*stride : (y+dy)*stride+w]
			band[dy] = row
		}

		// Determine which palette indices appear in this band.
		used := [256]bool{}
		for dy := 0; dy < bandH; dy++ {
			for x := 0; x < w; x++ {
				used[band[dy][x]] = true
			}
		}

		// For each used colour emit its sixel row, separated by '$' (graphics CR).
		first := true
		for ci := 0; ci < 256; ci++ {
			if !used[ci] {
				continue
			}

			if !first {
				sb.WriteByte('$') // return to left margin within the band
			}
			first = false

			fmt.Fprintf(&sb, "#%d", ci)

			sixelRow := make([]byte, w)
			for x := 0; x < w; x++ {
				var mask byte
				for dy := 0; dy < bandH; dy++ {
					if band[dy][x] == uint8(ci) {
						mask |= 1 << uint(dy)
					}
				}
				sixelRow[x] = 0x3F + mask
			}
			sb.Write(rleSixelRow(sixelRow))
		}

		sb.WriteByte('-') // graphics new-line: advance one sixel band
	}

	// ST (String Terminator) ends the DCS sequence.
	sb.WriteString("\x1b\\")

	return sb.String()
}

// rleSixelRow applies simple run-length encoding to a slice of sixel bytes.
// The SIXEL RLE form is "!<count><char>"; runs of fewer than 4 are emitted verbatim.
func rleSixelRow(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	var out []byte
	i := 0
	for i < len(data) {
		c := data[i]
		n := 1
		for i+n < len(data) && data[i+n] == c {
			n++
		}
		if n >= 4 {
			out = append(out, '!')
			out = append(out, []byte(strconv.Itoa(n))...)
			out = append(out, c)
		} else {
			for j := 0; j < n; j++ {
				out = append(out, c)
			}
		}
		i += n
	}
	return out
}
