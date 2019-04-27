/*
 * Copyright (c) 2014 PolyFloyd
 */

package main

import (
	"context"
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
	CubeHeight int
	CubeLength int
	CubeWidth  int

	camRot  mathgl.Quat
	camZoom float32

	swap      chan []float32
	showBlack bool

	voxelLen int
	shader   uint32
	win      *glfw.Window

	vertVAO, vertVBO uint32
	colorVBO         uint32
	translationVBO   uint32
}

func NewDisplay(w, h, l int) *Display {
	disp := &Display{
		CubeHeight: h,
		CubeLength: l,
		CubeWidth:  w,
		swap:       make(chan []float32, 1),
		showBlack:  true,
	}
	disp.ResetView()
	return disp
}

func (disp *Display) Run(ctx context.Context) {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := disp.init(); err != nil {
		panic(err)
	}

	for !disp.win.ShouldClose() {
		select {
		case f := <-disp.swap:
			gl.BindBuffer(gl.ARRAY_BUFFER, disp.colorVBO)
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(f)*4, gl.Ptr(&f[0]))
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		case <-ctx.Done():
			break
		default:
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
	uniformView := gl.GetUniformLocation(disp.shader, gl.Str("view\x00"))
	uniformProjection := gl.GetUniformLocation(disp.shader, gl.Str("projection\x00"))
	uniformShowBlack := gl.GetUniformLocation(disp.shader, gl.Str("show_black\x00"))

	projection := mathgl.Perspective(
		45.0,
		func(w, h int) float32 {
			return float32(w) / float32(h)
		}(disp.win.GetSize()),
		1.0,
		2048.0,
	)
	view := func() mathgl.Mat4 {
		m := mathgl.Ident4()
		m = m.Mul4(mathgl.Translate3D(0, 0, disp.camZoom))
		m = m.Mul4(disp.camRot.Mat4())
		return m
	}()
	gl.UniformMatrix4fv(uniformView, 1, false, (*float32)(&view[0]))
	gl.UniformMatrix4fv(uniformProjection, 1, false, (*float32)(&projection[0]))

	if disp.showBlack {
		gl.Uniform1f(uniformShowBlack, 1)
	} else {
		gl.Uniform1f(uniformShowBlack, 0)
	}

	gl.UseProgram(disp.shader)
	gl.BindVertexArray(disp.vertVAO)
	count := disp.CubeWidth * disp.CubeLength * disp.CubeHeight
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(disp.voxelLen), int32(count))
}

func (disp *Display) init() error {
	// Initialize GLFW and create a window
	if err := glfw.Init(); err != nil {
		return err
	}
	var err error
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	disp.win, err = glfw.CreateWindow(1280, 768, "ledcubesim", nil, nil)
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
			panic(fmt.Errorf("OpenGL [%s] %s", sevStr, message))
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
			disp.camRot = mathgl.QuatRotate(float32(deltaX)/240, mathgl.Vec3{0, 1, 0}).Mul(disp.camRot)
			disp.camRot = mathgl.QuatRotate(float32(deltaY)/240, mathgl.Vec3{1, 0, 0}).Mul(disp.camRot)
		}
	})
	disp.win.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton,
		action glfw.Action, mods glfw.ModifierKey) {
		dragButtonDown = action == glfw.Press && button == glfw.MouseButtonLeft
	})
	disp.win.SetScrollCallback(func(_ *glfw.Window, dx, dy float64) {
		disp.camZoom += float32(dy) * 12.0
	})
	disp.win.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int,
		action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Release {
			switch key {
			case glfw.KeyR:
				disp.ResetView()
			case glfw.KeyS:
				disp.ToggleShowBlack()
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
	vertices := getVoxelBuffer(0)
	disp.voxelLen = len(vertices)

	gl.CreateVertexArrays(1, &disp.vertVAO)
	gl.BindVertexArray(disp.vertVAO)
	gl.CreateBuffers(1, &disp.vertVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, disp.vertVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4*3, gl.Ptr(&vertices[0][0]), gl.STATIC_DRAW)
	vertAttrib := uint32(gl.GetAttribLocation(disp.shader, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	translations := getTranslationsBuffer(disp.CubeWidth, disp.CubeHeight, disp.CubeLength)
	gl.CreateBuffers(1, &disp.translationVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, disp.translationVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(translations)*4*3, gl.Ptr(&translations[0][0]), gl.STATIC_DRAW)
	translationAttrib := uint32(gl.GetAttribLocation(disp.shader, gl.Str("translation\x00")))
	gl.EnableVertexAttribArray(translationAttrib)
	gl.VertexAttribPointer(translationAttrib, 3, gl.FLOAT, false, 3*4, nil)
	gl.VertexAttribDivisor(translationAttrib, 1)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	colors := getInitialColorBuffer(disp.CubeWidth, disp.CubeHeight, disp.CubeLength)
	gl.CreateBuffers(1, &disp.colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, disp.colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4*3, gl.Ptr(&colors[0][0]), gl.DYNAMIC_DRAW)
	colorAttrib := uint32(gl.GetAttribLocation(disp.shader, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 3, gl.FLOAT, false, 3*4, nil)
	gl.VertexAttribDivisor(colorAttrib, 1)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)

	// Initialize the shader by setting some constant uniforms
	gl.Uniform3f(gl.GetUniformLocation(disp.shader, gl.Str("light_color\x00")), 0.15, 0.15, 0.15)
	lx, ly, lz := mathgl.Vec3{1, 1, 1}.Normalize().Elem()
	gl.Uniform3f(gl.GetUniformLocation(disp.shader, gl.Str("light_vec\x00")), lx, ly, lz)
	gl.Uniform1f(gl.GetUniformLocation(disp.shader, gl.Str("radius\x00")), 1)

	return nil
}

func (disp *Display) NumVoxels() int {
	return disp.CubeHeight * disp.CubeLength * disp.CubeWidth
}

func (disp *Display) Show(frame []float32) {
	select {
	case disp.swap <- frame:
	default:
	}
}

func (disp *Display) ResetView() {
	disp.camRot = mathgl.QuatIdent()
	disp.camRot = mathgl.QuatRotate(.5, mathgl.Vec3{0, -1, 0}).Mul(disp.camRot)
	disp.camRot = mathgl.QuatRotate(.5, mathgl.Vec3{1, 0, 0}).Mul(disp.camRot)
	maxDim := 0
	for _, dim := range []int{disp.CubeHeight, disp.CubeLength, disp.CubeWidth} {
		if dim > maxDim {
			maxDim = dim
		}
	}
	disp.camZoom = float32(maxDim) * -12
}

func (disp *Display) ToggleShowBlack() {
	disp.showBlack = !disp.showBlack
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

func getTranslationsBuffer(sx, sy, sz int) []mathgl.Vec3 {
	const spacing = 8.0
	buf := make([]mathgl.Vec3, 0, sx*sy*sz)
	for x := 0; x < sx; x++ {
		for y := 0; y < sy; y++ {
			for z := 0; z < sz; z++ {
				buf = append(buf, mathgl.Vec3{
					spacing*float32(x) - (spacing*float32(sx)*.5 - spacing*.5),
					spacing*float32(y) - (spacing*float32(sy)*.5 - spacing*.5),
					spacing*float32(z) - (spacing*float32(sx)*.5 - spacing*.5),
				})
			}
		}
	}
	return buf
}

func getInitialColorBuffer(sx, sy, sz int) []mathgl.Vec3 {
	buf := make([]mathgl.Vec3, 0, sx*sy*sz)
	for x := 0; x < sx; x++ {
		for y := 0; y < sy; y++ {
			for z := 0; z < sz; z++ {
				buf = append(buf, mathgl.Vec3{0.0, 0.4, 1.0})
			}
		}
	}
	return buf
}

const vertexShaderSource = `
	#version 330 core

	uniform float radius;
	uniform mat4 view;
	uniform mat4 projection;

	in vec3 vert;
	in vec3 translation;
	in vec3 color;

	out vec3 fNormal;
	out vec3 fColor;

	void main() {
		fNormal = vert;
		fColor = color;
		mat4 mvp = projection * view;

		gl_Position = mvp * vec4(vert*radius + translation, 1.0);
	}
`

const fragmentShaderSource = `
	#version 330 core

	uniform vec3 light_color;
	uniform vec3 light_vec;
	uniform float show_black;

	in vec3 fNormal;
	in vec3 fColor;

	void main() {
		if (fColor.r + fColor.g + fColor.b + show_black < 1e-18) {
			discard;
		}
		float cosTheta = clamp(dot(fNormal, light_vec), 0, 1);
		gl_FragColor = vec4(fColor + light_color * cosTheta, 1.0);
	}
`
