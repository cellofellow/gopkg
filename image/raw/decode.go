// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package raw

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
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

func (p *Decoder) Decode(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	// Gray/Gray16/Gray32f
	if p.Channels == 1 && p.DataType == reflect.Uint8 {
		return p.decodeGray(data, buf)
	}
	if p.Channels == 1 && p.DataType == reflect.Uint16 {
		return p.decodeGray16(data, buf)
	}
	if p.Channels == 1 && p.DataType == reflect.Float32 {
		return p.decodeGray32f(data, buf)
	}

	// RGB/RGB48/RGB96f
	if p.Channels == 3 && p.DataType == reflect.Uint8 {
		return p.decodeRGB(data, buf)
	}
	if p.Channels == 3 && p.DataType == reflect.Uint16 {
		return p.decodeRGB48(data, buf)
	}
	if p.Channels == 3 && p.DataType == reflect.Float32 {
		return p.decodeRGB96f(data, buf)
	}

	// RGBA/RGBA64/RGBA128f
	if p.Channels == 4 && p.DataType == reflect.Uint8 {
		return p.decodeRGBA(data, buf)
	}
	if p.Channels == 4 && p.DataType == reflect.Uint16 {
		return p.decodeRGBA64(data, buf)
	}
	if p.Channels == 4 && p.DataType == reflect.Float32 {
		return p.decodeRGBA128f(data, buf)
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

func (p *Decoder) decodeGray(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeGray, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	gray := newGray(image.Rect(0, 0, p.Width, p.Height), buf)
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(gray.Pix[y*gray.Stride:][:p.Width], data[off:])
		off += p.Width
	}
	m = gray
	return
}

func (p *Decoder) decodeGray16(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeGray16, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	gray16 := newGray16(image.Rect(0, 0, p.Width, p.Height), buf)
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

func (p *Decoder) decodeGray32f(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeGray32f, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	gray32f := newGray32f(image.Rect(0, 0, p.Width, p.Height), buf)
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(gray32f.Pix[y*gray32f.Stride:][:p.Width*4], data[off:])
		off += p.Width * 4
	}
	m = gray32f
	return
}

func (p *Decoder) decodeRGB(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGB, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgb := newRGB(image.Rect(0, 0, p.Width, p.Height), buf)
	var off = 0
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			rgb.SetRGB(x, y, color_ext.RGB{
				R: data[off+0],
				G: data[off+1],
				B: data[off+2],
			})
			off += 3
		}
	}
	m = rgb
	return
}

func (p *Decoder) decodeRGB48(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGB48, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgb48 := newRGB48(image.Rect(0, 0, p.Width, p.Height), buf)
	var off = 0
	for y := 0; y < p.Height; y++ {
		u16Pix := builtin.Slice(data[off:], reflect.TypeOf([]uint16(nil))).([]uint16)
		for x := 0; x < p.Width; x++ {
			rgb48.SetRGB48(x, y, color_ext.RGB48{
				R: u16Pix[x*3+0],
				G: u16Pix[x*3+1],
				B: u16Pix[x*3+2],
			})
		}
		off += p.Width * 6
	}
	m = rgb48
	return
}

func (p *Decoder) decodeRGB96f(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGB96f, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgb96f := newRGB96f(image.Rect(0, 0, p.Width, p.Height), buf)
	var off = 0
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			rgb96f.SetRGB96f(x, y, color_ext.RGB96f{
				R: builtin.Float32(data[off+0:]),
				G: builtin.Float32(data[off+4:]),
				B: builtin.Float32(data[off+8:]),
			})
			off += 12
		}
	}
	m = rgb96f
	return
}

func (p *Decoder) decodeRGBA(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGBA, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba := newRGBA(image.Rect(0, 0, p.Width, p.Height), buf)
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(rgba.Pix[y*rgba.Stride:][:p.Width*4], data[off:])
		off += p.Width * 4
	}
	m = rgba
	return
}

func (p *Decoder) decodeRGBA64(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGBA64, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba64 := newRGBA64(image.Rect(0, 0, p.Width, p.Height), buf)
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

func (p *Decoder) decodeRGBA128f(data []byte, buf image_ext.ImageBuffer) (m draw.Image, err error) {
	if size := p.getImageDataSize(); len(data) != size {
		err = fmt.Errorf("image/raw: decodeRGBA128f, bad data size, expect = %d, got = %d", size, len(data))
		return
	}
	rgba128f := newRGBA128f(image.Rect(0, 0, p.Width, p.Height), buf)
	var off = 0
	for y := 0; y < p.Height; y++ {
		copy(rgba128f.Pix[y*rgba128f.Stride:][:p.Width*16], data[off:])
		off += p.Width * 16
	}
	m = rgba128f
	return
}
