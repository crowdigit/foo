#version 330

in vec2 texCoord;
out vec4 out_color;

uniform sampler2D sampler;

void main() {
    out_color = texture(sampler, texCoord);
}
