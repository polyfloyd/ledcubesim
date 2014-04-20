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

const (
	CUBE_TOTAL_VOXELS = CUBE_WIDTH * CUBE_LENGTH * CUBE_HEIGHT
)

var (
	cam       = NewCubeCamera()
	sphere    *mesh.Mesh
	ledShader *shader.Program

	projection    = matreex.NewElement()
	ledTransforms = make([]matreex.Element, CUBE_TOTAL_VOXELS)

	shouldSwapBuffer = false
	DisplayBackBuffer = make([]float32, CUBE_TOTAL_VOXELS * 3)
	frontBuffer       = make([]float32, CUBE_TOTAL_VOXELS * 3)
)

func SwapDisplayBuffer() {
	shouldSwapBuffer = true
}

func StartDisplay(title string) {
	irix.UseVSync(true)

	input.OnMouseScroll(func(dx, dy float64) {
		cam.Zoom += float32(dy) * ZOOM_ACCELERATION
	})
	input.OnMouseDrag(glfw.MouseButtonLeft, func(x, y float64) {
		cam.RotX += float32(x / 10)
		cam.RotY += float32(y / 10)
	})


	colors := [][]float32{
		{ 0, 0, 1 },
		{ 0, 1, 0 },
		{ 0, 1, 1 },
		{ 1, 0, 0 },
		{ 0, 0, 0 },
	}
	for i := 0; i < CUBE_TOTAL_VOXELS*3; i += 3 {
		c := colors[(i/3) % len(colors)]
		frontBuffer[i+0] = c[0]
		frontBuffer[i+1] = c[1]
		frontBuffer[i+2] = c[2]
	}


	center := matreex.NewElement()
	center.Translate(
		-(1 + LED_DISTANCE*CUBE_WIDTH /2),
		-(1 + LED_DISTANCE*CUBE_HEIGHT/2),
		-(1 + LED_DISTANCE*CUBE_LENGTH/2),
		nil,
	);
	camMat := cam.MatElement()
	camMat.AddChild(center)
	for x := 0; x < CUBE_WIDTH; x++ {
		for y := 0; y < CUBE_HEIGHT; y++ {
			for z := 0; z < CUBE_LENGTH; z++ {
				mat := &ledTransforms[x*CUBE_HEIGHT*CUBE_LENGTH + y*CUBE_LENGTH + z]
				mat.LoadIdentity()
				mat.Translate(
					float32(x * LED_DISTANCE),
					float32(y * LED_DISTANCE),
					float32(z * LED_DISTANCE),
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
	gl.ClearColor(BACKGROUND, BACKGROUND, BACKGROUND, 1.0)

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
	if shouldSwapBuffer {
		t := DisplayBackBuffer
		DisplayBackBuffer = frontBuffer
		frontBuffer = t
		shouldSwapBuffer = false
	}
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
		if !RENDER_OFF && (r==0 && g==0 && b==0) {
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
	return FOVY
}

func (cam *CubeCamera) GetZNear() float32 {
	return 1.0
}

func (cam *CubeCamera) GetZFar() float32 {
	return 640
}
