package main

import (
	glfw    "github.com/go-gl/glfw3"
	gl      "github.com/polyfloyd/gl"
	irix    "polyfloyd/irix"
	input   "polyfloyd/irix/input"
	xmath   "polyfloyd/irix/math"
	matreex "polyfloyd/irix/math/matreex"
	mesh    "polyfloyd/irix/mesh"
	shader  "polyfloyd/irix/shader"
	util    "polyfloyd/irix/util"
)

var (
	showOff = true

	TotalVoxels       int
	DisplayBackBuffer []float32

	frontBuffer []float32

	cam           = NewCubeCamera()
	projection    = matreex.NewElement()
	sphere        *mesh.Mesh
	ledShader     *shader.Program
	ledTransforms []matreex.Element
)

func SwapDisplayBuffer() {
	util.RunGLAsync(func() {
		t := DisplayBackBuffer
		DisplayBackBuffer = frontBuffer
		frontBuffer = t
	})
}

func StartDisplay(title string) {
	width  := Config.Int("cube.width")
	length := Config.Int("cube.length")
	height := Config.Int("cube.height")
	TotalVoxels       = width * length * height
	DisplayBackBuffer = make([]float32, TotalVoxels * 3)
	frontBuffer       = make([]float32, TotalVoxels * 3)
	ledTransforms     = make([]matreex.Element, TotalVoxels)

	irix.UseVSync(true)

	input.OnKeyPress(glfw.KeyS, func(_ glfw.ModifierKey) {
		showOff = !showOff
	})
	input.OnKeyPress(glfw.KeyR, func(_ glfw.ModifierKey) {
		cam.RotX = 0
		cam.RotY = 0
		cam.Zoom = -160
	})
	input.OnMouseScroll(func(dx, dy float64) {
		cam.Zoom += float32(dy) * Config.Float32("ui.zoomAccel")
	})
	input.OnMouseDrag(glfw.MouseButtonLeft, func(x, y float64) {
		cam.RotX += float32(x / 10)
		if cam.RotX > 90 {
			cam.RotX = 90
		} else if cam.RotX < -90 {
			cam.RotX = -90
		}
		cam.RotY += float32(y / 10)
		if cam.RotY > 90 {
			cam.RotY = 90
		} else if cam.RotY < -90 {
			cam.RotY = -90
		}
	})

	for i := 0; i < len(frontBuffer); i += 3 {
		frontBuffer[i]     = 0.0
		frontBuffer[i + 1] = 0.4
		frontBuffer[i + 2] = 1.0
	}

	spacing := Config.Float32("ui.spacing")
	center := matreex.NewElement()
	center.Translate(
		-(spacing*float32(width)/2  - spacing/2),
		-(spacing*float32(height)/2 - spacing/2),
		-(spacing*float32(length)/2 - spacing/2),
		nil,
	);
	camMat := cam.MatElement()
	camMat.AddChild(center)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			for z := 0; z < length; z++ {
				mat := &ledTransforms[x*height*length + y*length + z]
				mat.LoadIdentity()
				mat.Translate(
					float32(x) * spacing,
					float32(y) * spacing,
					float32(z) * spacing,
					nil,
				)
				center.AddChild(mat)
			}
		}
	}
	projection.AddChild(camMat)

	util.RunGLAsync(InitGL)
	util.Check(irix.OpenWindow(1280, 768, false))
	irix.WindowTitle(title)
	irix.Main(UpdateDisplay)
}

func InitGL() {
	bg := Config.Float32("ui.background")
	gl.ClearColor(bg, bg, bg, 1.0)

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
	projection.LoadProjection(
		cam.GetFovy(),
		irix.WindowAspect(),
		cam.GetZNear(),
		cam.GetZFar(),
	)
	projection.Update()
	cam.UpdateLogic(delta)

	colorUniform := ledShader.Uniform["color_led"]
	for i, mat := range ledTransforms {
		r := frontBuffer[i*3]
		g := frontBuffer[i*3 + 1]
		b := frontBuffer[i*3 + 2]
		if !showOff && (r==0 && g==0 && b==0) {
			continue
		}
		ledShader.MatModviewProj(&mat.Abs)
		colorUniform.Uniform3f(r, g, b)
		sphere.Render()
	}
}

type CubeCamera struct {
	RotX    float32
	RotY    float32
	Zoom    float32
	mat     matreex.Element
	inverse xmath.Matrix4
}

func NewCubeCamera() *CubeCamera {
	return &CubeCamera{
		Zoom: -160,
		mat:  *matreex.NewElement(),
	}
}

func (cam *CubeCamera) UpdateLogic(delta float32) {
	cam.mat.LoadIdentity()
	cam.mat.Translate(0, 0, cam.Zoom, nil)
	cam.mat.Rotate(cam.RotX, 0, 1, 0, nil)
	cam.mat.Rotate(cam.RotY, 1, 0, 0, nil)
	cam.mat.Invert(&cam.inverse)
}

func (cam *CubeCamera) Unproject(v *xmath.Vector4) (ret xmath.Vector4) {
	cam.inverse.TransformV(v, &ret)
	return
}

func (cam *CubeCamera) MatElement() *matreex.Element {
	return &cam.mat
}

func (cam *CubeCamera) GetFovy() float32 {
	return Config.Float32("ui.fovy")
}

func (cam *CubeCamera) GetZNear() float32 {
	return 1.0
}

func (cam *CubeCamera) GetZFar() float32 {
	return 640
}
