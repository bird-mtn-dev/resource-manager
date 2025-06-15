package audiomanager

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"golang.org/x/exp/maps"
)

type AudioManager struct {
	FS         fs.FS
	audioFiles map[string]string
}

type AudioOptions struct {
	looping     bool
	loopLength  int64
	introLength int64
	volume      float64
}

type AudioOpt func(l *AudioOptions)

type ReadSeekLength interface {
	io.ReadSeeker
	Length() int64
}

func Looping() AudioOpt {
	return func(l *AudioOptions) {
		l.looping = true
	}
}
func LoopLength(length int64) AudioOpt {
	return func(l *AudioOptions) {
		l.loopLength = length
		l.looping = true
	}
}
func IntroLength(length int64) AudioOpt {
	return func(l *AudioOptions) {
		l.introLength = length
		l.looping = true
	}
}

// Volume must be between 0 and 1
func Volume(volume float64) AudioOpt {
	return func(l *AudioOptions) {
		l.volume = volume
	}
}

func Create() *AudioManager {
	return CreateWithFS(os.DirFS("."))
}

func CreateWithFS(filesystem fs.FS) *AudioManager {
	return &AudioManager{FS: filesystem, audioFiles: make(map[string]string)}
}

func (fm *AudioManager) LoadPlayer(name string, path string) {
	fm.audioFiles[name] = path
}

func (fm *AudioManager) GetPlayer(name string, options ...AudioOpt) (*audio.Player, error) {
	opts := AudioOptions{loopLength: -1, introLength: -1, volume: 1}
	for _, o := range options {
		o(&opts)
	}

	path, exists := fm.audioFiles[name]
	if !exists {
		panic("Audio file: " + name + " not loaded")
	}

	file, err := fm.FS.Open(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var stream ReadSeekLength
	switch filepath.Ext(path) {
	case ".mp3":
		stream, _ = mp3.DecodeWithoutResampling(file)
	case ".ogg":
		stream, _ = vorbis.DecodeWithoutResampling(file)
	case ".wav":
		stream, _ = wav.DecodeWithoutResampling(file)
	default:
		panic("Audio File Format Unknown: " + path)
	}
	var finalStream io.Reader
	if opts.looping {
		var loopLength int64
		if opts.loopLength != -1 {
			loopLength = opts.loopLength
		} else {
			loopLength = stream.Length()
		}
		if opts.introLength != -1 {

			finalStream = audio.NewInfiniteLoopWithIntro(stream, opts.introLength, loopLength)
		} else {
			finalStream = audio.NewInfiniteLoop(stream, loopLength)
		}
	} else {
		finalStream = stream
	}

	player, err := audio.CurrentContext().NewPlayer(finalStream)
	if err == nil {
		player.SetVolume(opts.volume)
	}
	return player, err
}

// This function will remove all loaded audio file paths in this Manager.
func (fm *AudioManager) Clear() {
	maps.Clear(fm.audioFiles)
}
