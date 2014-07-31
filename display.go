package main

import (
	"math"
	gl     "github.com/polyfloyd/go-gl"
	glfw   "github.com/go-gl/glfw3"
	input  "polyfloyd/irix/input"
	irix   "polyfloyd/irix"
	mathgl "github.com/go-gl/mathgl/mgl32"
	mesh   "polyfloyd/irix/mesh"
	shader "polyfloyd/irix/shader"
	util   "polyfloyd/irix/util"
)

var (
	showOff = true

	TotalVoxels       int
	DisplayBackBuffer []float32

	frontBuffer []float32
	sphere      *mesh.Mesh
	ledShader   *shader.Program

	camRotX float32
	camRotY float32
	camZoom float32
)

func SwapDisplayBuffer() {
	util.RunGLAsync(func() {
		t := DisplayBackBuffer
		DisplayBackBuffer = frontBuffer
		frontBuffer = t
	})
}

func StartDisplay(title string) {
	TotalVoxels       = CUBE_WIDTH * CUBE_LENGTH * CUBE_HEIGHT
	DisplayBackBuffer = make([]float32, TotalVoxels * 3)
	frontBuffer       = make([]float32, TotalVoxels * 3)

	for i := 0; i < len(frontBuffer); i += 3 {
		frontBuffer[i + 0] = 0.0
		frontBuffer[i + 1] = 0.4
		frontBuffer[i + 2] = 1.0
	}

	irix.UseVSync(true)

	input.OnKeyPress(glfw.KeyS, func(_ glfw.ModifierKey) {
		showOff = !showOff
	})
	input.OnKeyPress(glfw.KeyR, func(_ glfw.ModifierKey) {
		camRotX = 0
		camRotY = 0
		camZoom = -160
	})
	input.OnMouseScroll(func(dx, dy float64) {
		camZoom += float32(dy) * UI_ZOOMACCEL
	})
	input.OnMouseDrag(glfw.MouseButtonLeft, func(x, y float64) {
		camRotX += float32(x) / UI_DRAGDIV
		if camRotX > math.Pi/2 {
			camRotX = math.Pi/2
		} else if camRotX < -math.Pi/2 {
			camRotX = -math.Pi/2
		}
		camRotY += float32(y) / UI_DRAGDIV
		if camRotY > math.Pi/2 {
			camRotY = math.Pi/2
		} else if camRotY < -math.Pi/2 {
			camRotY = -math.Pi/2
		}
	})

	util.RunGLAsync(InitGL)
	util.Check(irix.OpenWindow(1280, 768, false))
	irix.WindowTitle(title)
	irix.Main(UpdateDisplay)
}

func InitGL() {
	gl.ClearColor(0.12, 0.12, 0.12, 1.0)

	m, err := mesh.Build(mesh.GenIcosahedron(2))
	util.Check(err)
	sphere = m[0]
	sphere.Load()

	vert, err := shader.CreateVertexObject(`
		#version 330 core

		{{.vert_position  }}
		{{.vert_normal    }}
		{{.vert_tex2      }}
		{{.vert_color     }}
		{{.mat_modviewproj}}

		out vec3 frag_normal;
		out vec2 frag_tex2;
		out vec3 frag_color;

		void main() {
			frag_normal = vert_normal;
			frag_tex2   = vert_tex2;
			frag_color  = vert_color;
			gl_Position = mat_modviewproj * vec4(vert_position, 1.0);
		}
	`)
	util.Check(err)
	util.Check(vert.Load())
	frag, err := shader.CreateFragmentObject(`
		#version 330 core

		vec3 LIGHT_VEC   = normalize(vec3(1, 1, 1));
		vec3 LIGHT_COLOR = vec3(0.2, 0.2, 0.2);

		in vec3 frag_normal;
		in vec2 frag_tex2;
		in vec3 frag_color;

		uniform vec3 color_led;

		out vec3 color;

		void main() {
			float cosTheta = clamp(dot(frag_normal, LIGHT_VEC), 0, 1);
			color = color_led + LIGHT_COLOR * cosTheta;
		}
	`)
	util.Check(err)
	util.Check(frag.Load())
	ledShader, err = shader.Link(true, vert, frag)
	util.Check(err)

	ledShader.Enable()
	sphere.Enable()
}

func UpdateDisplay(delta float32) {
	projection := mathgl.Perspective(UI_FOVY, irix.WindowAspect(), UI_ZNEAR, UI_ZFAR)

	center := mathgl.Translate3D(
		-(UI_SPACING*float32(CUBE_WIDTH)/2  - UI_SPACING/2),
		-(UI_SPACING*float32(CUBE_HEIGHT)/2 - UI_SPACING/2),
		-(UI_SPACING*float32(CUBE_LENGTH)/2 - UI_SPACING/2),
	)
	view := func() mathgl.Mat4 {
		m := mathgl.Ident4()
		m = m.Mul4(mathgl.Translate3D(0, 0, camZoom))
		m = m.Mul4(mathgl.HomogRotate3DY(camRotX))
		m = m.Mul4(mathgl.HomogRotate3DX(camRotY))
		return m
	}()

	uniformColor := ledShader.Uniform["color_led"]
	uniformMVP   := ledShader.Uniform["mat_modviewproj"]
	for x := 0; x < CUBE_WIDTH; x++ {
		for y := 0; y < CUBE_HEIGHT; y++ {
			for z := 0; z < CUBE_LENGTH; z++ {
				i := x*CUBE_HEIGHT*CUBE_LENGTH + y*CUBE_LENGTH + z

				r := frontBuffer[i*3 + 0]
				g := frontBuffer[i*3 + 1]
				b := frontBuffer[i*3 + 2]
				if !showOff && (r==0 && g==0 && b==0) {
					continue
				}

				model := mathgl.Translate3D(
					float32(x) * UI_SPACING,
					float32(y) * UI_SPACING,
					float32(z) * UI_SPACING,
				).Mul4(center);

				mvp := projection.Mul4(view).Mul4(model)
				uniformMVP.UniformMatrix4f(false, (*[16]float32)(&mvp))
				uniformColor.Uniform3f(r, g, b)
				sphere.Render()
			}
		}
	}
}
