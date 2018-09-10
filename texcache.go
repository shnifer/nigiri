package nigiri

import (
	"bytes"
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"image"
	"io/ioutil"
	"path"
)

type TexLoaderF func(name string) (*ebiten.Image, error)
type TexCache struct {
	cache  map[string]Tex
	loader TexLoaderF
}

func NewTexCache(loader TexLoaderF) *TexCache {
	return &TexCache{
		cache:  make(map[string]Tex),
		loader: loader,
	}
}

var tCache *TexCache

func init() {
	tCache = NewTexCache(nil)
}

func SetTexLoader(f TexLoaderF) {
	tCache.SetTexLoader(f)
}

func (tc *TexCache) SetTexLoader(f TexLoaderF) {
	tc.loader = f
}

func GetTex(name string) (tex Tex, err error) {
	return tCache.GetTex(name)
}

func (tc *TexCache) GetTex(name string) (tex Tex, err error) {
	tex, ok := tc.cache[name]
	if ok {
		return tex, nil
	}
	if tc.loader == nil {
		return Tex{}, errors.New("texture \"" + name + "\" not found, loader is nil")
	}
	img, err := tc.loader(name)
	if err != nil {
		return Tex{}, errors.Wrap(err, "can't load tex \""+name+"\" with loader")
	}
	tex = NewTex(img)
	tc.cache[name] = tex
	return tex, nil
}

func TexCacheReset() {
	tCache.Reset()
}

func (tc *TexCache) Reset() {
	tc.cache = make(map[string]Tex)
}

func FileTexLoader(pathStr string) TexLoaderF {
	return func(name string) (*ebiten.Image, error) {
		dat, err := ioutil.ReadFile(path.Join(pathStr, name))
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(dat)
		img, _, err := image.Decode(buf)
		if err != nil {
			return nil, err
		}
		ebiImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		if err != nil {
			return nil, err
		}
		return ebiImg, nil
	}
}
