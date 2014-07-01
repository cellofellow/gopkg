// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"fmt"
	"image"
	"image/color"
	"reflect"

	"github.com/chai2010/gopkg/builtin"
	image_ext "github.com/chai2010/gopkg/image"
	color_ext "github.com/chai2010/gopkg/image/color"
)

type Decoder struct {
	Channels int          // 1/3/4
	DataType reflect.Kind // Uint8/Uint16/Float32
	Width    int          // need for Decode
	Height   int          // need for Decode
}

func (p *Decoder) Decode(data []byte) (m image.Image, err error) {
	// Gray/Gray16/Gray32f
	if p.Channels == 1 && p.DataType == reflect.Uint8 {
		return p.decodeGray(data)
	}
	if p.Channels == 1 && p.DataType == reflect.Uint16 {
		return p.decodeGray16(data)
	}
	if p.Channels == 1 && p.DataType == reflect.Float32 {
		return p.decodeGray32f(data)
	}

	// RGB/RGB48/RGB96f
	if p.Channels == 3 && p.DataType == reflect.Uint8 {
		return p.decodeRGB(data)
	}
	if p.Channels == 3 && p.DataType == reflect.Uint16 {
		return p.decodeRGB48(data)
	}
	if p.Channels == 3 && p.DataType == reflect.Float32 {
		return p.decodeRGB96f(data)
	}

	// RGBA/RGBA64/RGBA128f
	if p.Channels == 4 && p.DataType == reflect.Uint8 {
		return p.decodeRGBA(data)
	}
	if p.Channels == 4 && p.DataType == reflect.Uint16 {
		return p.decodeRGBA64(data)
	}
	if p.Channels == 4 && p.DataType == reflect.Float32 {
		return p.decodeRGBA128f(data)
	}

	// Unknown
	err = fmt.Errorf(
		"image/raw: Decode, unknown image format, channels = %v, dataType = %v",
		p.Channels, p.DataType,
	)
	return
}

func (p *Decoder) getPixelSize() int {
	switch p.DataType {
	case reflect.Uint8:
		return p.Channels * 1
	case reflect.Uint16:
		return p.Channels * 2
	case reflect.Float32:
		return p.Channels * 4
	}
	panic("image/raw: getPixelSize, unreachable")
}

func (p *Decoder) getImageDataSize() int {
	return p.getPixelSize() * p.Width * p.Height
}

func (p *Decoder) decodeGray(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeGray, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	gray := image.NewGray(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(gray.Pix[y*gray.Stride:][:p.Width], data[off:])
		off += p.Width
	}
	m = gray
	return
}

func (p *Decoder) decodeGray16(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeGray16, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	gray16 := image.NewGray16(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		u16Pix := builtin.Slice(data[off:], reflect.TypeOf([]uint16(nil))).([]uint16)
		for x := 0; x < p.Width; x++ {
			gray16.SetGray16(x, y, color.Gray16{u16Pix[x]})
		}
		off += p.Width * 2
	}
	m = gray16
	return
}

func (p *Decoder) decodeGray32f(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeGray32f, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	gray32f := image_ext.NewGray32f(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(gray32f.Pix[y*gray32f.Stride:][:p.Width*4], data[off:])
		off += p.Width * 4
	}
	m = gray32f
	return
}

func (p *Decoder) decodeRGB(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGB, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			rgba.SetRGBA(x, y, color.RGBA{
				R: data[off+0],
				G: data[off+1],
				B: data[off+2],
			})
			off += 3
		}
	}
	m = rgba
	return
}

func (p *Decoder) decodeRGB48(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGB48, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba64 := image.NewRGBA64(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		u16Pix := builtin.Slice(data[off:], reflect.TypeOf([]uint16(nil))).([]uint16)
		for x := 0; x < p.Width; x++ {
			rgba64.SetRGBA64(x, y, color.RGBA64{
				R: u16Pix[x*3+0],
				G: u16Pix[x*3+1],
				B: u16Pix[x*3+2],
			})
		}
		off += p.Width * 6
	}
	m = rgba64
	return
}

func (p *Decoder) decodeRGB96f(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGB96f, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba128f := image_ext.NewRGBA128f(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			rgba128f.SetRGBA128f(x, y, color_ext.RGBA128f{
				R: builtin.Float32(data[off+0:]),
				G: builtin.Float32(data[off+4:]),
				B: builtin.Float32(data[off+8:]),
			})
			off += 12
		}
	}
	m = rgba128f
	return
}

func (p *Decoder) decodeRGBA(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGBA, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(rgba.Pix[y*rgba.Stride:][:p.Width*4], data[off:])
		off += p.Width * 4
	}
	m = rgba
	return
}

func (p *Decoder) decodeRGBA64(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGBA64, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba64 := image.NewRGBA64(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		u16Pix := builtin.Slice(data[off:], reflect.TypeOf([]uint16(nil))).([]uint16)
		for x := 0; x < p.Width; x++ {
			rgba64.SetRGBA64(x, y, color.RGBA64{
				R: u16Pix[x*4+0],
				G: u16Pix[x*4+1],
				B: u16Pix[x*4+2],
				A: u16Pix[x*4+3],
			})
		}
		off += p.Width * 8
	}
	m = rgba64
	return
}

func (p *Decoder) decodeRGBA128f(data []byte) (m image.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGBA128f, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba128f := image_ext.NewRGBA128f(image.Rect(0, 0, p.Width, p.Height))
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(rgba128f.Pix[y*rgba128f.Stride:][:p.Width*16], data[off:])
		off += p.Width * 16
	}
	m = rgba128f
	return
}
