package nigiriutil

import (
	"github.com/gobuffalo/packr"
	"github.com/shnifer/nigiri"
	"github.com/hajimehoshi/ebiten"
	"bytes"
	"image"
	"path"
	"golang.org/x/image/font"
	"github.com/golang/freetype/truetype"
)

func PackrTexLoader(box packr.Box, pathStr string) nigiri.TexLoaderF{
	return func(name string) (*ebiten.Image, error) {
		dat, err := box.MustBytes(path.Join(pathStr,name))
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
func PackrFaceLoader(box packr.Box, pathStr string) nigiri.FaceLoaderF{
	return func(name string, size float64) (font.Face, error) {
		dat, err := box.MustBytes(path.Join(pathStr,name))
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