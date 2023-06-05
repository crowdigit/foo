#version 330

layout (location = 0) in vec2 in_vertex;
uniform mat4 projModel;

void main() {
    gl_Position = projModel * vec4(in_vertex, 0, 1);
}
