package nigiri

import (
	"bytes"
	"github.com/hajimehoshi/ebiten"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"path"
)

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
