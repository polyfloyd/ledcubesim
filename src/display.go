/*
 * Copyright (c) 2014 PolyFloyd
 */

package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"
	"unsafe"

	gl "github.com/go-gl/gl/v3.3-core/gl"
	glfw "github.com/go-gl/glfw/v3.2/glfw"
	mathgl "github.com/go-gl/mathgl/mgl32"
)

type Display struct {
	Buffer  []float32
	HideOff bool

	CubeHeight int
	CubeLength int
	CubeWidth  int

	camRot  mathgl.Quat
	camZoom float32

	detail            int
	frontBuffer       []float32
	shader            uint32
	shouldSwapBuffers bool
	voxelBuffer       uint32
	voxelLen          int
	win               *glfw.Window
}

func NewDisplay(w, h, l, detail int) *Display {
	disp := &Display{
		CubeHeight:  h,
		CubeLength:  l,
		CubeWidth:   w,
		Buffer:      make([]float32, w*h*l*3),
		frontBuffer: make([]float32, w*h*l*3),
		detail:      detail,
	}
	disp.ResetView()
	return disp
}

func (disp *Display) Run() {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	for i := 0; i < len(disp.frontBuffer); i += 3 {
		disp.frontBuffer[i+0] = 0.0
		disp.frontBuffer[i+1] = 0.4
		disp.frontBuffer[i+2] = 1.0
	}

	if err := disp.init(); err != nil {
		panic(err)
	}

	for !disp.win.ShouldClose() {
		if disp.shouldSwapBuffers {
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
	uniformColor := gl.GetUniformLocation(disp.shader, gl.Str("color_led\x00"))
	uniformMVP := gl.GetUniformLocation(disp.shader, gl.Str("mat_mvp\x00"))

	projection := mathgl.Perspective(
		UI_FOVY,
		func(w, h int) float32 {
			return float32(w) / float32(h)
		}(disp.win.GetSize()),
		UI_ZNEAR,
		UI_ZFAR,
	)
	center := mathgl.Translate3D(
		-(UI_SPACING*float32(disp.CubeWidth)/2 - UI_SPACING/2),
		-(UI_SPACING*float32(disp.CubeHeight)/2 - UI_SPACING/2),
		-(UI_SPACING*float32(disp.CubeLength)/2 - UI_SPACING/2),
	)
	view := func() mathgl.Mat4 {
		m := mathgl.Ident4()
		m = m.Mul4(mathgl.Translate3D(0, 0, disp.camZoom))
		m = m.Mul4(disp.camRot.Mat4())
		return m
	}()

	for x := 0; x < disp.CubeWidth; x++ {
		for y := 0; y < disp.CubeLength; y++ {
			for z := 0; z < disp.CubeHeight; z++ {
				i := x*disp.CubeLength*disp.CubeHeight + y*disp.CubeHeight + z
				r := disp.frontBuffer[i*3+0]
				g := disp.frontBuffer[i*3+1]
				b := disp.frontBuffer[i*3+2]
				if disp.HideOff && (r == 0 && g == 0 && b == 0) {
					continue
				}

				model := mathgl.Translate3D(
					float32(x)*UI_SPACING,
					float32(z)*UI_SPACING,
					float32(y)*UI_SPACING,
				).Mul4(center)

				mvp := projection.Mul4(view).Mul4(model)
				gl.UniformMatrix4fv(uniformMVP, 1, false, (*float32)(&mvp[0]))
				gl.Uniform3f(uniformColor, r, g, b)
				gl.DrawArrays(gl.TRIANGLES, 0, int32(disp.voxelLen))
			}
		}
	}
}

func (disp *Display) init() error {
	// Initialize GLFW and create a window
	if err := glfw.Init(); err != nil {
		return err
	}
	var err error
	disp.win, err = glfw.CreateWindow(UI_WIN_W, UI_WIN_H, INFO, nil, nil)
	if err != nil {
		return err
	}
	disp.win.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		return err
	}

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)
	gl.DebugMessageCallback(func(source uint32, typ uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
		var sevStr string
		switch severity {
		case gl.DEBUG_SEVERITY_HIGH:
			sevStr = "high"
		case gl.DEBUG_SEVERITY_MEDIUM:
			sevStr = "medium"
		case gl.DEBUG_SEVERITY_LOW:
			sevStr = "low"
		case gl.DEBUG_SEVERITY_NOTIFICATION:
			sevStr = "note"
		}
		if severity == gl.DEBUG_SEVERITY_HIGH {
			log.Fatalf("OpenGL [%s] %s", sevStr, message)
		} else {
			log.Printf("OpenGL [%s] %s", sevStr, message)
		}
	}, nil)

	resize := func(w, h int) { gl.Viewport(0, 0, int32(w), int32(h)) }
	disp.win.SetSizeCallback(func(win *glfw.Window, w, h int) {
		resize(w, h)
	})
	resize(disp.win.GetSize())
	glfw.SwapInterval(1)

	// Initialize user input
	var dragButtonDown bool
	var mousePosLastX float64
	var mousePosLastY float64
	disp.win.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		deltaX := x - mousePosLastX
		deltaY := y - mousePosLastY
		mousePosLastX = x
		mousePosLastY = y
		if dragButtonDown {
			disp.camRot = mathgl.QuatRotate(float32(deltaX)/UI_DRAGDIV, mathgl.Vec3{0, 1, 0}).Mul(disp.camRot)
			disp.camRot = mathgl.QuatRotate(float32(deltaY)/UI_DRAGDIV, mathgl.Vec3{1, 0, 0}).Mul(disp.camRot)
		}
	})
	disp.win.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton,
		action glfw.Action, mods glfw.ModifierKey) {
		dragButtonDown = action == glfw.Press && button == glfw.MouseButtonLeft
	})
	disp.win.SetScrollCallback(func(_ *glfw.Window, dx, dy float64) {
		disp.camZoom += float32(dy) * UI_ZOOMACCEL
	})
	disp.win.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int,
		action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Release {
			switch key {
			case glfw.KeyS:
				disp.HideOff = !disp.HideOff
			case glfw.KeyR:
				disp.ResetView()
			}
		}
	})

	// Initialize OpenGL
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.12, 0.12, 0.12, 1.0)

	// Compile the voxel shader
	compileShader := func(typ uint32, src string) (uint32, error) {
		shader := gl.CreateShader(typ)
		csources, free := gl.Strs(src + "\x00")
		gl.ShaderSource(shader, 1, csources, nil)
		free()
		gl.CompileShader(shader)

		var status int32
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)
			log := strings.Repeat("\x00", int(logLen+1))
			gl.GetShaderInfoLog(shader, logLen, nil, gl.Str(log))
			gl.DeleteShader(shader)
			return 0, fmt.Errorf("%s", log)
		}
		return shader, nil
	}
	vx, err := compileShader(gl.VERTEX_SHADER, vertexShaderSource)
	if err != nil {
		return err
	}
	fg, err := compileShader(gl.FRAGMENT_SHADER, fragmentShaderSource)
	if err != nil {
		return err
	}

	disp.shader = gl.CreateProgram()
	gl.AttachShader(disp.shader, vx)
	gl.AttachShader(disp.shader, fg)
	gl.LinkProgram(disp.shader)
	gl.DetachShader(disp.shader, vx)
	gl.DetachShader(disp.shader, fg)
	gl.DeleteShader(vx)
	gl.DeleteShader(fg)
	var status int32
	gl.GetProgramiv(disp.shader, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(disp.shader, gl.INFO_LOG_LENGTH, &logLen)
		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(disp.shader, logLen, nil, gl.Str(log))
		return fmt.Errorf("%s", log)
	}

	gl.UseProgram(disp.shader)

	// Generate and initialize the voxel model
	bufferData := getVoxelBuffer(disp.detail)
	gl.GenBuffers(1, &disp.voxelBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, disp.voxelBuffer)
	disp.voxelLen = len(bufferData) * 3
	gl.BufferData(gl.ARRAY_BUFFER, disp.voxelLen*4, gl.Ptr(&bufferData[0][0]), gl.STATIC_DRAW)
	attribVert := uint32(gl.GetAttribLocation(disp.shader, gl.Str("voxel\x00")))
	gl.EnableVertexAttribArray(attribVert)
	gl.VertexAttribPointer(attribVert, 3, gl.FLOAT, false, 3*4, nil)

	// Initialize the shader by setting some constant uniforms
	gl.Uniform3f(gl.GetUniformLocation(disp.shader, gl.Str("light_color\x00")), 0.15, 0.15, 0.15)
	lx, ly, lz := mathgl.Vec3{1, 1, 1}.Normalize().Elem()
	gl.Uniform3f(gl.GetUniformLocation(disp.shader, gl.Str("light_vec\x00")), lx, ly, lz)
	gl.Uniform1f(gl.GetUniformLocation(disp.shader, gl.Str("voxel_radius\x00")), 1)

	return nil
}

func (disp *Display) NumVoxels() int {
	return disp.CubeHeight * disp.CubeLength * disp.CubeWidth
}

func (disp *Display) SwapBuffers() {
	disp.shouldSwapBuffers = true
}

func (disp *Display) ResetView() {
	disp.camRot = mathgl.QuatIdent()
	disp.camZoom = -160
}

func getVoxelBuffer(detail int) []mathgl.Vec3 {
	ico := float32(1+math.Sqrt(5)) / 2
	verts := []mathgl.Vec3{
		{-1.0, ico, 0.0},
		{1.0, ico, 0.0},
		{-1.0, -ico, 0.0},
		{1.0, -ico, 0.0},
		{0.0, -1.0, ico},
		{0.0, 1.0, ico},
		{0.0, -1.0, -ico},
		{0.0, 1.0, -ico},
		{ico, 0.0, -1.0},
		{ico, 0.0, 1.0},
		{-ico, 0.0, -1.0},
		{-ico, 0.0, 1.0},
	}
	polys := []int{
		0, 11, 5,
		0, 5, 1,
		0, 1, 7,
		0, 7, 10,
		0, 10, 11,
		1, 5, 9,
		5, 11, 4,
		11, 10, 2,
		10, 7, 6,
		7, 1, 8,
		3, 9, 4,
		3, 4, 2,
		3, 2, 6,
		3, 6, 8,
		3, 8, 9,
		4, 9, 5,
		2, 4, 11,
		6, 2, 10,
		8, 6, 7,
		9, 8, 1,
	}
	bufferData := make([]mathgl.Vec3, len(polys))[:0]
	for _, p := range polys {
		v := verts[p].Normalize()
		bufferData = append(bufferData, v)
	}
	var tessellate func(data []mathgl.Vec3, level int) []mathgl.Vec3
	tessellate = func(data []mathgl.Vec3, level int) []mathgl.Vec3 {
		if level == 0 {
			return data
		}
		newData := make([]mathgl.Vec3, len(data)*3)[:0]
		for i := 0; i < len(data); i += 3 {
			old := data[i : i+3]
			new := [3]mathgl.Vec3{}
			for j := range old {
				a := data[i+j]
				b := data[i+(j+1)%3]
				new[j] = mathgl.Vec3{
					a.X() - (a.X()-b.X())/2,
					a.Y() - (a.Y()-b.Y())/2,
					a.Z() - (a.Z()-b.Z())/2,
				}.Normalize()
			}
			newData = append(newData, old[0], new[0], new[2])
			newData = append(newData, old[1], new[1], new[0])
			newData = append(newData, old[2], new[2], new[1])
			newData = append(newData, new[0], new[1], new[2])
		}
		return tessellate(newData, level-1)
	}
	return tessellate(bufferData, detail)
}

const vertexShaderSource = `
	#version 120

	uniform float voxel_radius;
	uniform mat4 mat_mvp;

	attribute vec3 voxel;
	varying vec3 frag_normal;

	void main() {
		frag_normal = voxel;
		gl_Position = mat_mvp * vec4(voxel*voxel_radius, 1.0);
	}
`

const fragmentShaderSource = `
	#version 120

	uniform vec3 light_color;
	uniform vec3 light_vec;
	uniform vec3 color_led;

	varying vec3 frag_normal;

	void main() {
		float cosTheta = clamp(dot(frag_normal, light_vec), 0, 1);
		gl_FragColor = vec4(color_led + light_color * cosTheta, 1.0);
	}
`
