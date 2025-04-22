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
