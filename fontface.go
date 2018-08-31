package nigiri

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io/ioutil"
	"path"
)

type faceSign struct {
	filename string
	size     float64
}

type FaceLoaderF func(name string, size float64) (font.Face, error)

type faceCache struct {
	cache  map[faceSign]font.Face
	loader FaceLoaderF
}

var fCache *faceCache

func init() {
	fCache = newFaceCache(nil)
}

func GetFace(name string, size float64) (font.Face, error) {
	return fCache.GetFace(name, size)
}
func SetFaceLoader(loader FaceLoaderF) {
	fCache.SetLoader(loader)
}
func ResetFaceCache() {
	fCache.Reset()
}

func newFaceCache(loader FaceLoaderF) *faceCache {
	return &faceCache{
		cache:  make(map[faceSign]font.Face),
		loader: loader,
	}
}

func (fc *faceCache) SetLoader(loader FaceLoaderF) {
	fc.loader = loader
}

func newFace(b []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	tto := &truetype.Options{
		Size: size,
	}
	face := truetype.NewFace(f, tto)
	return face, nil
}

func (fc *faceCache) GetFace(name string, size float64) (font.Face, error) {
	sign := faceSign{
		filename: name,
		size:     size,
	}
	if face, ok := fc.cache[sign]; ok {
		return face, nil
	}
	face, err := fc.loader(name, size)
	if err != nil {
		return nil, err
	}
	fc.cache[sign] = face
	return face, nil
}

func (fc *faceCache) Reset() {
	fc.cache = make(map[faceSign]font.Face)
}

func FileFaceLoader(pathStr string) FaceLoaderF {
	return func(name string, size float64) (font.Face, error) {
		dat, err := ioutil.ReadFile(path.Join(pathStr, name))
		if err != nil {
			return nil, err
		}
		face, err := newFace(dat, size)
		if err != nil {
			return nil, err
		}

		return face, nil
	}
}
