package shadersys

// Embedded GLSL for shader.pbr / shader.toon / shader.dissolve (OpenGL 3.3).
// Vertex layout matches raylib's default 3D pipeline (vertexPosition, vertexTexCoord, vertexNormal, vertexColor).

const presetVertexShader = `#version 330
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;
out vec2 fragTexCoord;
out vec4 fragColor;
out vec3 fragWorldPos;
out vec3 fragWorldNor;
uniform mat4 matProjection;
uniform mat4 matView;
uniform mat4 matModel;
void main() {
  fragTexCoord = vertexTexCoord;
  fragColor = vertexColor;
  mat3 nmat = transpose(inverse(mat3(matModel)));
  fragWorldNor = normalize(nmat * vertexNormal);
  vec4 w = matModel * vec4(vertexPosition, 1.0);
  fragWorldPos = w.xyz;
  gl_Position = matProjection * matView * w;
}
`

const presetPBRFragment = `#version 330
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragWorldPos;
in vec3 fragWorldNor;
uniform sampler2D texture0;
out vec4 finalColor;
void main() {
  vec3 L = normalize(vec3(1.0, 2.0, 1.0));
  float ndl = max(dot(normalize(fragWorldNor), L), 0.0);
  vec4 tex = texture(texture0, fragTexCoord);
  if (tex.a < 0.01) discard;
  finalColor = tex * fragColor * vec4(vec3(0.12 + 0.88 * ndl), 1.0);
}
`

const presetToonFragment = `#version 330
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragWorldPos;
in vec3 fragWorldNor;
uniform sampler2D texture0;
out vec4 finalColor;
void main() {
  vec3 L = normalize(vec3(1.0, 2.0, 1.0));
  float ndl = dot(normalize(fragWorldNor), L);
  float shade = ndl > 0.55 ? 1.0 : (ndl > 0.2 ? 0.5 : 0.22);
  vec4 tex = texture(texture0, fragTexCoord);
  if (tex.a < 0.01) discard;
  finalColor = tex * fragColor * vec4(vec3(shade), 1.0);
}
`

const presetDissolveFragment = `#version 330
in vec2 fragTexCoord;
in vec4 fragColor;
in vec3 fragWorldPos;
in vec3 fragWorldNor;
uniform sampler2D texture0;
uniform float dissolve;
out vec4 finalColor;
void main() {
  vec4 tex = texture(texture0, fragTexCoord);
  if (tex.a < 0.01) discard;
  float edge = (fragTexCoord.x + fragTexCoord.y) * 0.5;
  if (edge > dissolve) discard;
  finalColor = tex * fragColor;
}
`
