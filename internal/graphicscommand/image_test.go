// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphicscommand_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/affine"
	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicscommand"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
	etesting "github.com/hajimehoshi/ebiten/v2/internal/testing"
)

func TestMain(m *testing.M) {
	etesting.MainWithRunLoop(m)
}

func quadVertices(w, h float32) []float32 {
	return []float32{
		0, 0, 0, 0, 1, 1, 1, 1,
		w, 0, w, 0, 1, 1, 1, 1,
		0, w, 0, h, 1, 1, 1, 1,
		w, h, w, h, 1, 1, 1, 1,
	}
}

func TestClear(t *testing.T) {
	const w, h = 1024, 1024
	src := graphicscommand.NewImage(w/2, h/2)
	dst := graphicscommand.NewImage(w, h)

	vs := quadVertices(w/2, h/2)
	is := graphics.QuadIndices()
	dr := graphicsdriver.Region{
		X:      0,
		Y:      0,
		Width:  w,
		Height: h,
	}
	dst.DrawTriangles([graphics.ShaderImageNum]*graphicscommand.Image{src}, [graphics.ShaderImageNum - 1][2]float32{}, vs, is, affine.ColorMIdentity{}, graphicsdriver.CompositeModeClear, graphicsdriver.FilterNearest, graphicsdriver.AddressUnsafe, dr, graphicsdriver.Region{}, nil, nil, false)

	pix := make([]byte, 4*w*h)
	if err := dst.ReadPixels(pix); err != nil {
		t.Fatal(err)
	}
	for j := 0; j < h/2; j++ {
		for i := 0; i < w/2; i++ {
			idx := 4 * (i + w*j)
			got := color.RGBA{pix[idx], pix[idx+1], pix[idx+2], pix[idx+3]}
			want := color.RGBA{}
			if got != want {
				t.Errorf("dst.At(%d, %d) after DrawTriangles: got %v, want: %v", i, j, got, want)
			}
		}
	}
}

func TestReplacePixelsPartAfterDrawTriangles(t *testing.T) {
	const w, h = 32, 32
	clr := graphicscommand.NewImage(w, h)
	src := graphicscommand.NewImage(w/2, h/2)
	dst := graphicscommand.NewImage(w, h)
	vs := quadVertices(w/2, h/2)
	is := graphics.QuadIndices()
	dr := graphicsdriver.Region{
		X:      0,
		Y:      0,
		Width:  w,
		Height: h,
	}
	dst.DrawTriangles([graphics.ShaderImageNum]*graphicscommand.Image{clr}, [graphics.ShaderImageNum - 1][2]float32{}, vs, is, affine.ColorMIdentity{}, graphicsdriver.CompositeModeClear, graphicsdriver.FilterNearest, graphicsdriver.AddressUnsafe, dr, graphicsdriver.Region{}, nil, nil, false)
	dst.DrawTriangles([graphics.ShaderImageNum]*graphicscommand.Image{src}, [graphics.ShaderImageNum - 1][2]float32{}, vs, is, affine.ColorMIdentity{}, graphicsdriver.CompositeModeSourceOver, graphicsdriver.FilterNearest, graphicsdriver.AddressUnsafe, dr, graphicsdriver.Region{}, nil, nil, false)
	dst.ReplacePixels(make([]byte, 4), 0, 0, 1, 1)

	// TODO: Check the result.
}

func TestShader(t *testing.T) {
	const w, h = 16, 16
	clr := graphicscommand.NewImage(w, h)
	dst := graphicscommand.NewImage(w, h)
	vs := quadVertices(w, h)
	is := graphics.QuadIndices()
	dr := graphicsdriver.Region{
		X:      0,
		Y:      0,
		Width:  w,
		Height: h,
	}
	dst.DrawTriangles([graphics.ShaderImageNum]*graphicscommand.Image{clr}, [graphics.ShaderImageNum - 1][2]float32{}, vs, is, affine.ColorMIdentity{}, graphicsdriver.CompositeModeClear, graphicsdriver.FilterNearest, graphicsdriver.AddressUnsafe, dr, graphicsdriver.Region{}, nil, nil, false)

	ir := etesting.ShaderProgramFill(graphicscommand.NeedsInvertY(), 0xff, 0, 0, 0xff)
	s := graphicscommand.NewShader(&ir)
	dst.DrawTriangles([graphics.ShaderImageNum]*graphicscommand.Image{}, [graphics.ShaderImageNum - 1][2]float32{}, vs, is, affine.ColorMIdentity{}, graphicsdriver.CompositeModeSourceOver, graphicsdriver.FilterNearest, graphicsdriver.AddressUnsafe, dr, graphicsdriver.Region{}, s, nil, false)

	pix := make([]byte, 4*w*h)
	if err := dst.ReadPixels(pix); err != nil {
		t.Fatal(err)
	}
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			idx := 4 * (i + w*j)
			got := color.RGBA{pix[idx], pix[idx+1], pix[idx+2], pix[idx+3]}
			want := color.RGBA{0xff, 0, 0, 0xff}
			if got != want {
				t.Errorf("dst.At(%d, %d) after DrawTriangles: got %v, want: %v", i, j, got, want)
			}
		}
	}
}
