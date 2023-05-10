#version 330

layout (location = 0) in vec2 in_vertex;
uniform mat3 ProjModel;

void main() {
    vec3 vertex = ProjModel * vec3(in_vertex, 0);
    gl_Position = vec4(vertex, 1);
}
