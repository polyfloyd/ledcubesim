package main

import (
	"runtime"
	gl     "github.com/polyfloyd/go-gl"
	glfw   "github.com/go-gl/glfw3"
	input  "polyfloyd/irix/input"
	mathgl "github.com/go-gl/mathgl/mgl32"
	mesh   "polyfloyd/irix/mesh"
	shader "polyfloyd/irix/shader"
)

type Display struct {
	Buffer  []float32
	HideOff bool

	cubeHeight  int
	cubeLength  int
	cubeWidth   int

	camRot  mathgl.Quat
	camZoom float32

	frontBuffer       []float32
	ledModel          *mesh.Mesh
	shader            *shader.Program
	shouldSwapBuffers bool
	win  *glfw.Window
}

func NewDisplay(w, h, l int) *Display {
	disp := &Display{
		cubeHeight:  h,
		cubeLength:  l,
		cubeWidth:   w,
		Buffer:      make([]float32, w*h*l * 3),
		frontBuffer: make([]float32, w*h*l * 3),
	}
	disp.ResetView()
	return disp
}

func (disp *Display) Start() {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	for i := 0; i < len(disp.frontBuffer); i += 3 {
		disp.frontBuffer[i + 0] = 0.0
		disp.frontBuffer[i + 1] = 0.4
		disp.frontBuffer[i + 2] = 1.0
	}

	input.OnKeyPress(glfw.KeyS, func(_ glfw.ModifierKey) {
		disp.HideOff = !disp.HideOff
	})
	input.OnKeyPress(glfw.KeyR, func(_ glfw.ModifierKey) {
		disp.ResetView()
	})
	input.OnMouseScroll(func(dx, dy float64) {
		disp.camZoom += float32(dy) * UI_ZOOMACCEL
	})
	input.OnMouseDrag(glfw.MouseButtonLeft, func(x, y float64) {
		disp.camRot = mathgl.QuatRotate(float32(x) / UI_DRAGDIV, mathgl.Vec3{0, 1, 0}).Mul(disp.camRot)
		disp.camRot = mathgl.QuatRotate(float32(y) / UI_DRAGDIV, mathgl.Vec3{1, 0, 0}).Mul(disp.camRot)
	})

	if !glfw.Init() {
		panic("Can't init GLFW!")
	}
	glfw.SwapInterval(1)

	var err error
	disp.win, err = glfw.CreateWindow(UI_WIN_W, UI_WIN_H, INFO, nil, nil)
	if err != nil {
		panic(err)
	}
	disp.win.MakeContextCurrent()

	resize := func(w, h int) {
		gl.Viewport(0, 0, w, h)
	}
	disp.win.SetSizeCallback(func(win *glfw.Window, w, h int) {
		resize(w, h)
	})
	input.SetInputWindow(disp.win)
	resize(disp.win.GetSize())

	if err := disp.initGL(); err != nil {
		panic(err)
	}

	for !disp.win.ShouldClose() {
		if (disp.shouldSwapBuffers) {
			disp.frontBuffer, disp.Buffer = disp.Buffer, disp.frontBuffer
			disp.shouldSwapBuffers = false
		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		disp.render()

		disp.win.SwapBuffers()
		glfw.PollEvents()
	}
	disp.win.Destroy()
	runtime.UnlockOSThread()
}

func (disp *Display) render() {
	uniformColor := disp.shader.Uniform["color_led"]
	uniformMVP   := disp.shader.Uniform["mat_modviewproj"]

	projection := mathgl.Perspective(
		UI_FOVY,
		func(w, h int) float32 {
			return float32(w) / float32(h)
		}(disp.win.GetSize()),
		UI_ZNEAR,
		UI_ZFAR,
	)
	center := mathgl.Translate3D(
		-(UI_SPACING*float32(disp.cubeWidth)/2  - UI_SPACING/2),
		-(UI_SPACING*float32(disp.cubeHeight)/2 - UI_SPACING/2),
		-(UI_SPACING*float32(disp.cubeLength)/2 - UI_SPACING/2),
	)
	view := func() mathgl.Mat4 {
		m := mathgl.Ident4()
		m = m.Mul4(mathgl.Translate3D(0, 0, disp.camZoom))
		m = m.Mul4(disp.camRot.Mat4())
		return m
	}()

	for x := 0; x < disp.cubeWidth; x++ {
		for y := 0; y < disp.cubeHeight; y++ {
			for z := 0; z < disp.cubeLength; z++ {
				i := x*disp.cubeHeight*disp.cubeLength + y*disp.cubeLength + z

				r := disp.frontBuffer[i*3 + 0]
				g := disp.frontBuffer[i*3 + 1]
				b := disp.frontBuffer[i*3 + 2]
				if disp.HideOff && (r==0 && g==0 && b==0) {
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
				disp.ledModel.Render()
			}
		}
	}
}

func (disp *Display) initGL() error {
	if err := gl.Init(); err != nil { return err }

	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.12, 0.12, 0.12, 1.0)

	m, err := mesh.Build(mesh.GenIcosahedron(2))
	if (err != nil) { return err }
	disp.ledModel = m[0]
	disp.ledModel.Load()

	vert, err := shader.CreateVertexObject(SHADER_SRC_VX)
	if (err != nil) { return err }
	if err := vert.Load(); err != nil { return err }
	frag, err := shader.CreateFragmentObject(SHADER_SRC_FG)
	if (err != nil) { return err }
	if err := frag.Load(); err != nil { return err }
	disp.shader, err = shader.Link(true, vert, frag)
	if (err != nil) { return err }

	disp.shader.Enable()
	disp.ledModel.Enable()
	return nil
}

func (disp *Display) SwapBuffers() {
	disp.shouldSwapBuffers = true
}

func (disp *Display) ResetView() {
	disp.camRot  = mathgl.QuatIdent()
	disp.camZoom = -160
}

const SHADER_SRC_VX = `
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
`

const SHADER_SRC_FG = `
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
`
