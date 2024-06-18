package fontmanager

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/exp/maps"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/gofont/goregular"
)

type FontManager struct {
	FS        fs.FS
	fonts     map[string]*truetype.Font
	fontCache map[string]*font.Face
}

const (
	STANDARD_NORMAL      = "standard_normal"
	STANDARD_ITALIC      = "standard_italic"
	STANDARD_BOLD        = "standard_bold"
	STANDARD_BOLD_ITALIC = "standard_bold_italic"

	MONO_NORMAL      = "mono_normal"
	MONO_ITALIC      = "mono_italic"
	MONO_BOLD        = "mono_bold"
	MONO_BOLD_ITALIC = "mono_bold_italic"
)

func Create() *FontManager {
	return CreateWithFS(os.DirFS("."))
}

func CreateWithFS(filesystem fs.FS) *FontManager {
	return &FontManager{FS: filesystem, fonts: make(map[string]*truetype.Font), fontCache: make(map[string]*font.Face)}
}

// This function loads the standard san-serif gofonts into the manager
//
//	It loads the regular, italic, bold, and bolditalic fonts
//
// see https://pkg.go.dev/golang.org/x/image/font/gofont
func (fm *FontManager) LoadStandardFonts() error {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[STANDARD_NORMAL] = ttfFont

	ttfFont, err = truetype.Parse(goitalic.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[STANDARD_ITALIC] = ttfFont

	ttfFont, err = truetype.Parse(gobold.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[STANDARD_BOLD] = ttfFont

	ttfFont, err = truetype.Parse(gobolditalic.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[STANDARD_BOLD_ITALIC] = ttfFont

	return nil
}

// This function loads the mono gofonts into the manager
//
//	It loads the mono, monoitalic, monobold, and monobolditalic fonts
//
// see https://pkg.go.dev/golang.org/x/image/font/gofont
func (fm *FontManager) LoadMonoFonts() error {
	ttfFont, err := truetype.Parse(gomono.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[MONO_NORMAL] = ttfFont

	ttfFont, err = truetype.Parse(gomonoitalic.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[MONO_ITALIC] = ttfFont

	ttfFont, err = truetype.Parse(gomonobold.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[MONO_BOLD] = ttfFont

	ttfFont, err = truetype.Parse(gomonobolditalic.TTF)
	if err != nil {
		return nil
	}
	fm.fonts[MONO_BOLD_ITALIC] = ttfFont

	return nil
}

func (fm *FontManager) LoadFont(name string, path string) error {
	fontFile, _ := fm.FS.Open(path)
	fontData, err := io.ReadAll(fontFile)
	if err != nil {
		return err
	}
	return fm.LoadFontData(name, fontData)
}

func (fm *FontManager) LoadFontData(name string, fontData []byte) error {
	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return err
	}
	fm.fonts[name] = ttfFont

	return nil
}

func (fm *FontManager) GetFace(name string, size float64) (font.Face, error) {

	cachedFace, exists := fm.fontCache[name+strconv.FormatFloat(size, 'f', -1, 64)]
	if exists {
		return *cachedFace, nil
	}

	ttfFont, exists := fm.fonts[name]
	if !exists {
		return nil, errors.New("Font not found: " + name)
	}
	face := truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	fm.fontCache[name+strconv.FormatFloat(size, 'f', -1, 64)] = &face
	return face, nil
}

func (fm *FontManager) GetFaceWithSpacing(name string, size float64, lineHeight float64) (font.Face, error) {
	face, err := fm.GetFace(name, size)
	if err != nil {
		return nil, err
	}
	if lineHeight > 1 {
		h := float64(face.Metrics().Height.Round()) * lineHeight
		face = text.FaceWithLineHeight(face, h)
	}
	return face, nil
}

func (fm *FontManager) PurgeCache() {
	maps.Clear(fm.fontCache)
}
