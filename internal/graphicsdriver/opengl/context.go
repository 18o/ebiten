// Copyright 2016 Hajime Hoshi
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

package opengl

import (
	"fmt"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

type operation int

func convertOperation(op graphicsdriver.Operation) operation {
	switch op {
	case graphicsdriver.Zero:
		return zero
	case graphicsdriver.One:
		return one
	case graphicsdriver.SrcAlpha:
		return srcAlpha
	case graphicsdriver.DstAlpha:
		return dstAlpha
	case graphicsdriver.OneMinusSrcAlpha:
		return oneMinusSrcAlpha
	case graphicsdriver.OneMinusDstAlpha:
		return oneMinusDstAlpha
	case graphicsdriver.DstColor:
		return dstColor
	default:
		panic(fmt.Sprintf("opengl: invalid operation %d at convertOperation", op))
	}
}

type context struct {
	locationCache      *locationCache
	screenFramebuffer  framebufferNative // This might not be the default frame buffer '0' (e.g. iOS).
	lastFramebuffer    framebufferNative
	lastTexture        textureNative
	lastRenderbuffer   renderbufferNative
	lastViewportWidth  int
	lastViewportHeight int
	lastCompositeMode  graphicsdriver.CompositeMode
	maxTextureSize     int
	maxTextureSizeOnce sync.Once
	highp              bool
	highpOnce          sync.Once

	contextImpl
}

func (c *context) bindTexture(t textureNative) {
	if c.lastTexture.equal(t) {
		return
	}
	c.bindTextureImpl(t)
	c.lastTexture = t
}

func (c *context) bindRenderbuffer(r renderbufferNative) {
	if c.lastRenderbuffer.equal(r) {
		return
	}
	c.bindRenderbufferImpl(r)
	c.lastRenderbuffer = r
}

func (c *context) bindFramebuffer(f framebufferNative) {
	if c.lastFramebuffer.equal(f) {
		return
	}
	c.bindFramebufferImpl(f)
	c.lastFramebuffer = f
}

func (c *context) setViewport(f *framebuffer) {
	c.bindFramebuffer(f.native)
	if c.lastViewportWidth != f.width || c.lastViewportHeight != f.height {
		// On some environments, viewport size must be within the framebuffer size.
		// e.g. Edge (#71), Chrome on GPD Pocket (#420), macOS Mojave (#691).
		// Use the same size of the framebuffer here.
		c.setViewportImpl(f.width, f.height)

		// glViewport must be called at least at every frame on iOS.
		// As the screen framebuffer is the last render target, next SetViewport should be
		// the first call at a frame.
		if f.native.equal(c.screenFramebuffer) {
			c.lastViewportWidth = 0
			c.lastViewportHeight = 0
		} else {
			c.lastViewportWidth = f.width
			c.lastViewportHeight = f.height
		}
	}
}

func (c *context) getScreenFramebuffer() framebufferNative {
	return c.screenFramebuffer
}

func (c *context) getMaxTextureSize() int {
	c.maxTextureSizeOnce.Do(func() {
		c.maxTextureSize = c.maxTextureSizeImpl()
	})
	return c.maxTextureSize
}

// highpPrecision represents an enough mantissa of float values in a shader.
const highpPrecision = 23

func (c *context) hasHighPrecisionFloat() bool {
	c.highpOnce.Do(func() {
		c.highp = c.getShaderPrecisionFormatPrecision() >= highpPrecision
	})
	return c.highp
}
