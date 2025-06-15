package fontmanager

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/exp/maps"
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
	fonts     map[string]*text.GoTextFaceSource
	fontCache map[string]*text.GoTextFace
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
	return &FontManager{FS: filesystem, fonts: make(map[string]*text.GoTextFaceSource), fontCache: make(map[string]*text.GoTextFace)}
}

// This function loads the standard san-serif gofonts into the manager
//
//	It loads the regular, italic, bold, and bolditalic fonts
//
// see https://pkg.go.dev/golang.org/x/image/font/gofont
func (fm *FontManager) LoadStandardFonts() error {
	ttfFont, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil
	}
	fm.fonts[STANDARD_NORMAL] = ttfFont

	ttfFont, err = text.NewGoTextFaceSource(bytes.NewReader(goitalic.TTF))
	if err != nil {
		return nil
	}
	fm.fonts[STANDARD_ITALIC] = ttfFont

	ttfFont, err = text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		return nil
	}
	fm.fonts[STANDARD_BOLD] = ttfFont

	ttfFont, err = text.NewGoTextFaceSource(bytes.NewReader(gobolditalic.TTF))
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
	ttfFont, err := text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
	if err != nil {
		return nil
	}
	fm.fonts[MONO_NORMAL] = ttfFont

	ttfFont, err = text.NewGoTextFaceSource(bytes.NewReader(gomonoitalic.TTF))
	if err != nil {
		return nil
	}
	fm.fonts[MONO_ITALIC] = ttfFont

	ttfFont, err = text.NewGoTextFaceSource(bytes.NewReader(gomonobold.TTF))
	if err != nil {
		return nil
	}
	fm.fonts[MONO_BOLD] = ttfFont

	ttfFont, err = text.NewGoTextFaceSource(bytes.NewReader(gomonobolditalic.TTF))
	if err != nil {
		return nil
	}
	fm.fonts[MONO_BOLD_ITALIC] = ttfFont

	return nil
}

// This function loads a font at the provided filesystem path.
func (fm *FontManager) LoadFont(name string, path string) error {
	fontFile, _ := fm.FS.Open(path)
	fontData, err := io.ReadAll(fontFile)
	if err != nil {
		return err
	}
	return fm.LoadFontData(name, fontData)
}

// This function loads a font from the provided byte array
func (fm *FontManager) LoadFontData(name string, fontData []byte) error {
	ttfFont, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		return err
	}
	fm.fonts[name] = ttfFont

	return nil
}

// This function returns a font face for the loaded font with 'name'. It will cache this face for future use.
func (fm *FontManager) GetFace(name string, size float64) (text.Face, error) {

	cachedFace, exists := fm.fontCache[name+strconv.FormatFloat(size, 'f', -1, 64)]
	if exists {
		return cachedFace, nil
	}

	ttfFont, exists := fm.fonts[name]
	if !exists {
		return nil, errors.New("Font not found: " + name)
	}

	face := &text.GoTextFace{
		Source: ttfFont,
		Size:   size,
	}
	fm.fontCache[name+strconv.FormatFloat(size, 'f', -1, 64)] = face
	return face, nil
}

// This function will clear the font cache
func (fm *FontManager) PurgeCache() {
	maps.Clear(fm.fontCache)
}

// This function will remove the specified font data
func (fm *FontManager) Remove(key string) {
	delete(fm.fonts, key)
}

// This function will remove all font data and the font cache in this Manager.
func (fm *FontManager) Clear() {
	maps.Clear(fm.fontCache)
	maps.Clear(fm.fonts)
}
