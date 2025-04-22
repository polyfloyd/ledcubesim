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
