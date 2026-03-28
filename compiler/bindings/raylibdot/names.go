package raylibdot

// modelNames: raylib_3d model, animation, material, billboard (not camera/lights).
var modelNames = []string{
	"LoadModel", "LoadModelAnimated", "LoadModelFromMesh", "UnloadModel", "LoadCube",
	"SetModelColor", "RotateModel", "DrawModelSimple", "DrawModel",
	"SetModelPosition", "SetModelRotation", "SetModelScale", "DrawModelWithState",
	"DrawModelEx", "DrawModelWires", "IsModelValid", "GetModelBoundingBox",
	"DrawModelWiresEx", "DrawModelPoints", "DrawModelPointsEx",
	"DrawBillboard", "DrawBillboardRec", "SetModelMeshMaterial", "DrawBillboardPro",
	"LoadModelAnimations", "GetModelAnimationId", "PlayModelAnimation",
	"SetModelTexture", "SetMaterialTexture", "SetMaterialColor", "SetModelShader",
	"SetMaterialFloat", "SetMaterialVector",
	"UpdateModelAnimation", "UpdateModelAnimationBones",
	"UnloadModelAnimation", "UnloadModelAnimations", "IsModelAnimationValid",
	"GetModelAnimationFrameCount", "CreateModelAnimState", "UpdateModelAnimState",
	"SetModelAnimStateFrame", "GetModelAnimStateFrame", "DestroyModelAnimState",
}

// shapes3dNames: primitive 3D draws from raylib_3d.go.
var shapes3dNames = []string{
	"DrawCube", "DrawCubeWires", "DrawSphere", "DrawSphereWires", "DrawPlane",
	"DrawLine3D", "DrawPoint3D", "DrawCircle3D", "DrawCubeV",
	"DrawCylinder", "DrawCylinderWires", "DrawText3D", "DrawRay",
	"DrawTriangle3D", "DrawTriangleStrip3D", "DrawCubeWiresV", "DrawSphereEx",
	"DrawCylinderEx", "DrawCylinderWiresEx", "DrawCapsule", "DrawCapsuleWires",
	"DrawBoundingBox",
}

// meshNames: raylib_mesh.go (+ shared material helpers).
var meshNames = []string{
	"MeshCreate", "MeshUpdate", "MeshSetVertices", "MeshSetNormals", "MeshSetUVs", "MeshSetIndices",
	"GenMeshPoly", "GenMeshPlane", "GenMeshCube", "GenMeshSphere", "GenMeshHemiSphere",
	"GenMeshCylinder", "GenMeshCone", "GenMeshTorus", "GenMeshKnot", "GenMeshHeightmap", "GenMeshCubicmap",
	"UploadMesh", "UnloadMesh", "GetMeshBoundingBox", "ExportMesh",
	"DrawMesh", "DrawMeshMatrix", "UpdateMeshBuffer", "DrawMeshInstanced",
	"LoadMaterialDefault", "IsMaterialValid", "UnloadMaterial", "SetMaterialTexture",
	"LoadMaterials", "GetMaterialIdFromLoad", "GetRayCollisionMesh", "GetRayCollisionModel",
}

// imageNames: raylib_images.go (all registered).
var imageNames = []string{
	"LoadImage", "LoadImageRaw", "LoadImageAnim", "GetLoadImageAnimFrames", "LoadImageAnimFromMemory",
	"LoadImageFromMemory", "LoadImageFromTexture", "LoadImageFromScreen", "IsImageValid", "UnloadImage",
	"ExportImage", "SaveImage", "ExportImageToMemory", "ExportImageAsCode",
	"GenImageColor", "GenImageGradientLinear", "GenImageGradientRadial", "GenImageGradientSquare",
	"GenImageChecked", "GenImageWhiteNoise", "GenImagePerlinNoise", "GenImageCellular", "GenImageText",
	"ImageCopy", "ImageFromImage", "ImageFromChannel", "ImageText", "ImageTextEx",
	"ImageFormat", "ImageToPOT", "ImageCrop", "ImageAlphaCrop", "ImageAlphaClear", "ImageAlphaMask",
	"ImageAlphaPremultiply", "ImageBlurGaussian", "ImageKernelConvolution",
	"ImageResize", "ImageResizeNN", "ImageResizeCanvas", "ImageMipmaps", "ImageDither",
	"ImageFlipVertical", "ImageFlipHorizontal", "ImageRotate", "ImageRotateCW", "ImageRotateCCW",
	"ImageColorTint", "ImageColorInvert", "ImageColorGrayscale", "ImageColorContrast",
	"ImageColorBrightness", "ImageColorReplace",
	"LoadImageColors", "UnloadImageColors", "GetLoadedImageColor", "GetImageColor",
	"ImageClearBackground", "ImageDrawPixel", "ImageDrawPixelV", "ImageDrawLine", "ImageDrawLineV", "ImageDrawLineEx",
	"ImageDrawCircle", "ImageDrawCircleV", "ImageDrawCircleLines", "ImageDrawCircleLinesV",
	"ImageDrawRectangle", "ImageDrawRectangleV", "ImageDrawRectangleRec", "ImageDrawRectangleLines",
	"ImageDrawTriangle", "ImageDrawTriangleEx", "ImageDrawTriangleLines",
	"ImageDrawTriangleFan", "ImageDrawTriangleStrip", "ImageDraw", "ImageDrawText", "ImageDrawTextEx",
}

// fontNames: raylib_fonts.go.
var fontNames = []string{
	"GetFontDefault", "LoadFont", "LoadFontEx", "DrawTextExFont", "MeasureTextEx", "UnloadFont",
	"LoadFontFromImage", "LoadFontFromMemory", "IsFontValid", "LoadFontData", "GenImageFontAtlas",
	"UnloadFontData", "ExportFontAsCode", "DrawTextCodepoint", "DrawTextCodepoints",
	"GetGlyphIndex", "GetGlyphInfo", "GetGlyphAtlasRec",
}

// rlaudioNames: raylib_audio.go.
var rlaudioNames = []string{
	"InitAudioDevice", "CloseAudioDevice", "IsAudioDeviceReady",
	"LoadSound", "PlaySound", "StopSound", "SetSoundVolume", "UnloadSound",
	"LoadMusicStream", "PlayMusicStream", "UpdateMusicStream", "StopMusicStream",
	"SetMusicVolume", "UnloadMusicStream", "SetMasterVolume", "GetMasterVolume",
	"PauseSound", "ResumeSound", "IsSoundPlaying", "SetSoundPitch", "SetSoundPan",
	"PauseMusicStream", "ResumeMusicStream", "IsMusicStreamPlaying",
	"LoadMusic", "PlayMusic", "PauseMusic", "ResumeMusic", "IsMusicPlaying",
	"SeekMusicStream", "SetMusicPitch", "SetMusicPan", "GetMusicTimeLength", "GetMusicTimePlayed", "IsMusicValid",
	"LoadMusicStreamFromMemory",
	"LoadWave", "LoadWaveFromMemory", "IsWaveValid", "UnloadWave", "ExportWave", "WaveCopy", "WaveCrop", "WaveFormat",
	"LoadWaveSamples", "UnloadWaveSamples", "ExportWaveAsCode",
	"LoadSoundFromWave", "LoadSoundAlias", "IsSoundValid", "UpdateSound", "UnloadSoundAlias",
	"LoadAudioStream", "IsAudioStreamValid", "UnloadAudioStream", "UpdateAudioStream", "IsAudioStreamProcessed",
	"PlayAudioStream", "PauseAudioStream", "ResumeAudioStream", "IsAudioStreamPlaying", "StopAudioStream",
	"SetAudioStreamVolume", "SetAudioStreamPitch", "SetAudioStreamPan",
	"SetAudioStreamBufferSizeDefault", "SetAudioStreamCallback",
	"AttachAudioStreamProcessor", "DetachAudioStreamProcessor",
	"AttachAudioMixedProcessor", "DetachAudioMixedProcessor",
}
