package print

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	"yoru/utils"
)

// Main reads input from stdin and renders either SIXEL images or a pretty table.
// Use -i to select the input format.
func Main(args []string) {
	fs := flag.NewFlagSet("print", flag.ExitOnError)
	inputFormat := fs.String("i", "png", "Input format: png, jpeg, jpg, gif, bmp, csv, tsv")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: yoru sixel print -i <format> [width] [height]\n")
		fmt.Fprintf(os.Stderr, "\nInput formats:\n")
		fmt.Fprintf(os.Stderr, "  png, jpeg, jpg, gif, bmp   Render stdin as SIXEL\n")
		fmt.Fprintf(os.Stderr, "  csv, tsv                   Render stdin as a pretty table\n")
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  width   Optional target width in pixels (image formats only)\n")
		fmt.Fprintf(os.Stderr, "  height  Optional target height in pixels (image formats only)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  cat image.png | yoru sixel print -i png\n")
		fmt.Fprintf(os.Stderr, "  cat image.jpg | yoru sixel print -i jpeg 100 50\n")
		fmt.Fprintf(os.Stderr, "  cat data.csv | yoru sixel print -i csv\n")
	}

	err := fs.Parse(args)
	utils.Error(err)

	remaining := fs.Args()
	format := normalizeFormat(*inputFormat)
	switch format {
	case "csv", "tsv":
		if len(remaining) > 0 {
			utils.Error(fmt.Errorf("width/height are only supported for image formats"))
		}
		delimiter := ','
		if format == "tsv" {
			delimiter = '\t'
		}
		out, err := renderDelimitedAsTable(os.Stdin, delimiter)
		utils.Error(err)
		_, err = fmt.Fprintln(os.Stdout, out)
		utils.Error(err)
		return
	case "png", "jpeg", "gif", "bmp":
		// continue
	default:
		utils.Error(fmt.Errorf("unsupported input format %q", *inputFormat))
		return
	}

	targetWidth, targetHeight := parseDimensions(remaining)
	img, err := decodeImageByFormat(os.Stdin, format)
	utils.Error(err)
	img = renderToSixelImage(img, targetWidth, targetHeight)
	sixelData := encode(img)
	_, err = fmt.Fprint(os.Stdout, sixelData)
	utils.Error(err)
}

func normalizeFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "jpg" {
		return "jpeg"
	}
	return format
}

func parseDimensions(remaining []string) (int, int) {
	var err error
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
	if len(remaining) > 2 {
		utils.Error(fmt.Errorf("unexpected extra arguments after height"))
	}

	return targetWidth, targetHeight
}

func decodeImageByFormat(r io.Reader, format string) (image.Image, error) {
	switch format {
	case "png":
		img, err := png.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("decoding png image: %w", err)
		}
		return img, nil
	case "jpeg":
		img, err := jpeg.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("decoding jpeg image: %w", err)
		}
		return img, nil
	case "gif":
		img, err := gif.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("decoding gif image: %w", err)
		}
		return img, nil
	case "bmp":
		img, err := bmp.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("decoding bmp image: %w", err)
		}
		return img, nil
	default:
		return nil, fmt.Errorf("unsupported image format %q", format)
	}
}

func renderToSixelImage(img image.Image, targetWidth, targetHeight int) image.Image {
	if targetWidth > 0 || targetHeight > 0 {
		img = scaleImage(img, targetWidth, targetHeight)
	}
	return flattenAlpha(img)
}

func renderDelimitedAsTable(r io.Reader, delimiter rune) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("reading input: %w", err)
	}

	reader := csv.NewReader(strings.NewReader(string(input)))
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("reading delimited input: %w", err)
	}
	if len(records) == 0 {
		return "", nil
	}

	maxCols := 0
	for _, row := range records {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	widths := make([]int, maxCols)
	rows := make([][]string, len(records))
	for i, row := range records {
		normalized := make([]string, maxCols)
		copy(normalized, row)
		rows[i] = normalized
		for colIdx, cell := range normalized {
			if len(cell) > widths[colIdx] {
				widths[colIdx] = len(cell)
			}
		}
	}

	border := func(left, mid, right string) string {
		var b strings.Builder
		b.WriteString(left)
		for i, width := range widths {
			b.WriteString(strings.Repeat("-", width+2))
			if i < len(widths)-1 {
				b.WriteString(mid)
			}
		}
		b.WriteString(right)
		return b.String()
	}

	var out strings.Builder
	out.WriteString(border("+", "+", "+"))
	out.WriteString("\n")

	for rowIdx, row := range rows {
		out.WriteString("|")
		for colIdx, cell := range row {
			padding := widths[colIdx] - len(cell)
			out.WriteString(" ")
			out.WriteString(cell)
			out.WriteString(strings.Repeat(" ", padding+1))
			out.WriteString("|")
		}
		out.WriteString("\n")

		if rowIdx == 0 && len(rows) > 1 {
			out.WriteString(border("+", "+", "+"))
			out.WriteString("\n")
		}
	}

	out.WriteString(border("+", "+", "+"))
	return out.String(), nil
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
//	Indices  0–15  : standard ANSI colours
//	Indices 16–231 : 6×6×6 RGB colour cube
//	Indices 232–255: 24-step greyscale ramp
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
