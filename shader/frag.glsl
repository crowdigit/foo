#version 330

in vec2 texCoord;
out vec4 out_color;

uniform vec3 color;
uniform sampler2D sampler;

void main() {
    out_color = texture(sampler, texCoord);
    // out_color = vec4(color, 0.0);
}
