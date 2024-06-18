package imagemanager

import (
	"image"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/maps"
)

type ImageManager struct {
	FS         fs.FS
	imageCache map[string]*ebiten.Image
}

func Create() *ImageManager {
	return CreateWithFS(os.DirFS("."))
}

func CreateWithFS(filesystem fs.FS) *ImageManager {
	return &ImageManager{FS: filesystem, imageCache: make(map[string]*ebiten.Image)}
}

func (im *ImageManager) GetImage(path string) (*ebiten.Image, error) {
	result, exists := im.imageCache[path]
	if exists {
		return result, nil
	}

	file, err := im.FS.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	img2 := ebiten.NewImageFromImage(img)
	im.imageCache[path] = img2
	return img2, nil
}

func (fm *ImageManager) Put(key string, value *ebiten.Image) {
	fm.imageCache[key] = value
}

func (fm *ImageManager) Get(key string) *ebiten.Image {
	return fm.imageCache[key]
}

func (fm *ImageManager) Remove(key string) {
	delete(fm.imageCache, key)
}

func (im *ImageManager) PurgeCache() {
	maps.Clear(im.imageCache)
}
