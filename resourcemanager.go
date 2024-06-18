package resourcemanager

import (
	"io/fs"
	"os"

	audiomanager "github.com/bird-mtn-dev/resource-manager/audio-manager"
	fontmanager "github.com/bird-mtn-dev/resource-manager/font-manager"
	imagemanager "github.com/bird-mtn-dev/resource-manager/image-manager"
)

type CustomID int

type JSONID int

type ResourceManager struct {
	fs     fs.FS
	Audio  *audiomanager.AudioManager
	Font   *fontmanager.FontManager
	Image  *imagemanager.ImageManager
	custom map[CustomID]any
	json   map[JSONID]any
}

func Create() ResourceManager {
	return CreateWithFS(os.DirFS("."))
}

func CreateWithFS(fs fs.FS) ResourceManager {
	return ResourceManager{
		fs:     fs,
		Audio:  audiomanager.CreateWithFS(fs),
		Font:   fontmanager.CreateWithFS(fs),
		Image:  imagemanager.CreateWithFS(fs),
		custom: make(map[CustomID]any),
		json:   make(map[JSONID]any),
	}
}

func (rm *ResourceManager) AddCustomManager(id CustomID, custom any) {
	rm.custom[id] = custom
}

func (rm *ResourceManager) GetCustomManager(id CustomID) any {
	return rm.custom[id]
}

func (rm *ResourceManager) RemoveCustomManager(id CustomID) {
	delete(rm.custom, id)
}

func (rm *ResourceManager) AddJSONManager(id JSONID, json any) {
	rm.json[id] = json
}

func (rm *ResourceManager) GetJSONManager(id JSONID) any {
	return rm.json[id]
}

func (rm *ResourceManager) RemoveJSONManager(id JSONID) {
	delete(rm.json, id)
}
