package main

import (
	glfw    "github.com/go-gl/glfw3"
	gl      "github.com/polyfloyd/gl"
	irix    "polyfloyd/irix"
	input   "polyfloyd/irix/input"
	// xmath   "polyfloyd/irix/math"
	matreex "polyfloyd/irix/math/matreex"
	// res     "polyfloyd/irix/res"
	mesh    "polyfloyd/irix/res/mesh"
	shader  "polyfloyd/irix/res/shader"
	uni     "polyfloyd/irix/res/uniform"
	util    "polyfloyd/irix/util"
)

const (
	BACKGROUND_R = 0.2
	BACKGROUND_G = 0.2
	BACKGROUND_B = 0.2

	CUBE_WIDTH  = 16
	CUBE_LENGTH = 16
	CUBE_HEIGHT = 16

	LED_DISTANCE = 8

	ZOOM_ACCELERATION = 2
)

var (
	cam       *CubeCamera
	sphere    *mesh.Mesh
	ledShader *shader.Program

	projection    = matreex.NewElement()
	center        = matreex.NewElement()
	ledTransforms = make([]matreex.Element, CUBE_WIDTH * CUBE_LENGTH * CUBE_HEIGHT)
)

func main() {

	irix.MustGLExt("GL_ARB_vertex_buffer_object")

	irix.FrameCap = 60
	irix.ShowFPS = true

	input.OnKeyPress(glfw.KeyF4, func(_ glfw.ModifierKey) {
		irix.Wireframe = !irix.Wireframe
	})

	cam = NewCubeCamera()
	input.OnMouseScroll(func(dx, dy float64) {
		cam.Zoom += float32(dy) * ZOOM_ACCELERATION
	})
	input.OnMouseDrag(glfw.MouseButtonLeft, func(x, y float64) {
		cam.RotX += float32(x / 10)
		cam.RotY += float32(y / 10)
	})

	err := irix.OpenWindow(1280, 768, false)
	util.Check(err)

	util.RunGLAsync(func() {
		gl.ClearColor(BACKGROUND_R, BACKGROUND_G, BACKGROUND_B, 1.0)

		var err error

		sphereBuilder := mesh.GenIcosahedron(2)
		sphere, _ = sphereBuilder.Build(nil)
		sphere.Load()

		// shader.PrintSource = true
		vert, err := shader.CreateVertexObject(`
			#version 430 core

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
			#version 430 core

			{{.mtl_ambient}}
			{{.mtl_diffuse}}
			uniform sampler2D tex;

			in vec3 frag_normal;
			in vec2 frag_tex2;
			in vec3 frag_color;

			out vec3 color;

			vec3 LIGHT_VEC = normalize(vec3(1, 1, 1));
			vec3 LIGHT_COLOR = vec3(0, 0.6, 1);

			void main() {
				float cosTheta = clamp(dot(frag_normal, LIGHT_VEC), 0, 1);
				color =
					mtl_ambient * mtl_diffuse +
					(mtl_diffuse + LIGHT_COLOR) * cosTheta;
			}
		`)
		util.Check(err)
		util.Check(frag.Load())
		ledShader, err = shader.Link(true, vert, frag)
		util.Check(err)

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

	})

	irix.WindowTitle("A Simulator For LED Cubes")
	irix.Main(Cycle)
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

	ledShader.Enable()

	for _, mat := range ledTransforms {
		uni.ModelviewProjection(&mat.Abs)
		sphere.RenderGL()
	}

	ledShader.Disable()
}
