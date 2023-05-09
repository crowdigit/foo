#version 330

layout (location = 0) in vec2 in_vertex;
uniform mat3 ProjModel;

void main() {
    gl_Position = vec4(ProjModel * vec3(in_vertex, 1.0), 1.0);
}
