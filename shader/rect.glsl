#version 330

layout (location = 0) in vec2 in_vertex;
uniform mat3 projModel;

void main() {
    vec3 vertex = projModel * vec3(in_vertex, 1);
    gl_Position = vec4(vertex.xy, 0, 1);
}
