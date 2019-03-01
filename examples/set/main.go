// Copyright 2019 The Ebiten Authors
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

// +build example jsgo

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var offscreen = ebiten.NewImage(screenWidth, screenHeight)

func update(screen *ebiten.Image) error {
	w, h := offscreen.Size()
	x := rand.Intn(w)
	y := rand.Intn(h)
	c := color.RGBA{
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(0xff),
	}
	offscreen.Set(x, y, c)

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.DrawImage(offscreen, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Test"); err != nil {
		log.Fatal(err)
	}
}
