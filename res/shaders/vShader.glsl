#version 410
in vec3 vp;
uniform vec3 pos;
void main() {
	gl_Position = vec4(vp + pos, 1.0);
}