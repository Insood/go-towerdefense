#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

void main()
{
    vec2 uv = fragTexCoord;
    float lineWidth = 0.035;

    float edgeX = 1.0 - step(lineWidth, uv.x) + step(1.0 - lineWidth, uv.x);
    float edgeY = 1.0 - step(lineWidth, uv.y) + step(1.0 - lineWidth, uv.y);
    float gridLine = clamp(max(edgeX, edgeY), 0.0, 1.0);

    vec3 fillColor = vec3(0.53, 0.91, 0.53);
    vec3 lineColor = vec3(0.35, 0.40, 0.46);
    vec3 color = mix(fillColor, lineColor, gridLine);

    finalColor = vec4(color, 1.0) * fragColor;
}
