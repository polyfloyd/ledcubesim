package main

import (
	xmath   "polyfloyd/irix/math"
	matreex "polyfloyd/irix/math/matreex"
)

type CubeCamera struct {
	RotX    float32
	RotY    float32
	Zoom    float32
	mat     matreex.Element
	inverse xmath.Matrix4
}

func NewCubeCamera() *CubeCamera {
	return &CubeCamera{
		Zoom: -80,
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
	return 320
}
