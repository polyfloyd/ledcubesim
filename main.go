package main

import (
	glfw    "github.com/go-gl/glfw3"
	gl      "github.com/polyfloyd/gl"
	irix    "polyfloyd/irix"
	input   "polyfloyd/irix/input"
	matreex "polyfloyd/irix/math/matreex"
	mesh    "polyfloyd/irix/res/mesh"
	shader  "polyfloyd/irix/res/shader"
	util    "polyfloyd/irix/util"
)

const (
	CUBE_WIDTH  = 16
	CUBE_LENGTH = 16
	CUBE_HEIGHT = 16

	BACKGROUND   = 0.12
	LED_DISTANCE = 8

	ZOOM_ACCELERATION = 8
	FOVY              = 45.0
)

const (
	CUBE_TOTAL_VOXELS = CUBE_WIDTH * CUBE_LENGTH * CUBE_HEIGHT
)

var (
	cam       = NewCubeCamera()
	sphere    *mesh.Mesh
	ledShader *shader.Program

	projection    = matreex.NewElement()
	center        = matreex.NewElement()
	ledTransforms = make([]matreex.Element, CUBE_TOTAL_VOXELS)

	frontBuffer = make([]float32, CUBE_TOTAL_VOXELS * 3)

	matrixUniformLocation gl.UniformLocation
	colorUniformLocation  gl.UniformLocation
)

func main() {
	irix.MustGLExt("GL_ARB_vertex_buffer_object")

	irix.FrameCap = 60
	irix.ShowFPS = true

	input.OnKeyPress(glfw.KeyF4, func(_ glfw.ModifierKey) {
		irix.Wireframe = !irix.Wireframe
	})

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
		{ 1, 0, 1 },
	}
	for i := 0; i < CUBE_TOTAL_VOXELS*3; i += 3 {
		c := colors[(i/3) % len(colors)]
		frontBuffer[i+0] = c[0]
		frontBuffer[i+1] = c[1]
		frontBuffer[i+2] = c[2]
	}


	camMat := cam.MatElement()
	center.Translate(
		-(1 + LED_DISTANCE*CUBE_WIDTH /2),
		-(1 + LED_DISTANCE*CUBE_LENGTH/2),
		-(1 + LED_DISTANCE*CUBE_HEIGHT/2),
		nil,
	);
	camMat.AddChild(center)
	for x := 0; x < CUBE_WIDTH; x++ {
		for y := 0; y < CUBE_LENGTH; y++ {
			for z := 0; z < CUBE_HEIGHT; z++ {
				mat := &ledTransforms[x*CUBE_LENGTH*CUBE_HEIGHT + y*CUBE_HEIGHT + z]
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
	irix.WindowTitle("A Simulator For LED Cubes")
	irix.Main(Cycle)
}

func InitGL() {
	gl.ClearColor(BACKGROUND, BACKGROUND, BACKGROUND, 1.0)

	sphereBuilder := mesh.GenIcosahedron(2)
	m, err := mesh.Build(sphereBuilder)
	util.Check(err)
	sphere = m[0]
	sphere.Load()

	vert, err := shader.CreateVertexObject(`
		#version 330 core

		{{.vert_position  }}
		{{.vert_normal    }}
		{{.vert_tex2      }}
		{{.vert_color     }}
		uniform mat4 mat_modviewproj;

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

		in vec3 frag_normal;
		in vec2 frag_tex2;
		in vec3 frag_color;
		uniform vec3 color_led;

		out vec3 color;

		vec3 LIGHT_VEC   = normalize(vec3(1, 1, 1));
		vec3 LIGHT_COLOR = vec3(0.5, 0.5, 0.5);

		void main() {
			float cosTheta = clamp(dot(frag_normal, LIGHT_VEC), 0, 1);
			color = color_led + LIGHT_COLOR * cosTheta;
		}
	`)
	util.Check(err)
	util.Check(frag.Load())
	ledShader, err = shader.Link(true, vert, frag)
	util.Check(err)

	matrixUniformLocation = ledShader.UniLoc("mat_modviewproj")
	colorUniformLocation  = ledShader.UniLoc("color_led")

	sphere.Enable()
	ledShader.Enable()
}

func Cycle(delta float32) {
	projection.LoadProjection(
		cam.GetFovy(),
		irix.WindowAspect(),
		cam.GetZNear(),
		cam.GetZFar(),
	)
	projection.Update()
	cam.UpdateLogic(delta)

	c := 0
	for _, mat := range ledTransforms {
		matrixUniformLocation.UniformMatrix4f(false, (*[16]float32)(&mat.Abs))
		colorUniformLocation.Uniform3f(
			frontBuffer[c],
			frontBuffer[c + 1],
			frontBuffer[c + 2],
		)
		c += 3
		sphere.RenderGL()
	}
}
