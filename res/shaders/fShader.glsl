#version 410
out vec4 frag_colour;
uniform vec3 color;
void main() {
	float depth = 0.0 + gl_FragCoord.z;
	frag_colour = vec4(color, 1.0);
}