#include "raylib_wrapper.h"
#include "raylib.h"

// Graphics initialization
bool cb_init_window(int width, int height, const char* title) {
    InitWindow(width, height, title);
    return IsWindowReady();
}

void cb_close_window(void) {
    CloseWindow();
}

bool cb_window_should_close(void) {
    return WindowShouldClose();
}

void cb_begin_drawing(void) {
    BeginDrawing();
}

void cb_end_drawing(void) {
    EndDrawing();
}

void cb_clear_background(int r, int g, int b) {
    ClearBackground((Color){ r, g, b, 255 });
}

// Texture/Image operations
cb_image cb_load_image(const char* filename) {
    Image img = LoadImage(filename);
    return (cb_image){ 
        img.data, 
        img.width, 
        img.height, 
        img.mipmaps, 
        img.format 
    };
}

void cb_unload_image(cb_image image) {
    Image img = { image.data, image.width, image.height, image.mipmaps, image.format };
    UnloadImage(img);
}

cb_texture cb_load_texture_from_image(cb_image image) {
    Image img = { image.data, image.width, image.height, image.mipmaps, image.format };
    Texture2D tex = LoadTextureFromImage(img);
    return (cb_texture){ 
        tex.id, 
        tex.width, 
        tex.height, 
        tex.mipmaps, 
        tex.format 
    };
}

void cb_unload_texture(cb_texture texture) {
    Texture2D tex = { texture.id, texture.width, texture.height, texture.mipmaps, texture.format };
    UnloadTexture(tex);
}

void cb_draw_texture(cb_texture texture, int x, int y, int tint) {
    Texture2D tex = { texture.id, texture.width, texture.height, texture.mipmaps, texture.format };
    DrawTexture(tex, x, y, (Color){ tint, tint, tint, 255 });
}

// 2D Drawing
void cb_draw_rectangle(int x, int y, int width, int height, int r, int g, int b, int a) {
    DrawRectangle(x, y, width, height, (Color){ r, g, b, a });
}

void cb_draw_circle(int x, int y, int radius, int r, int g, int b, int a) {
    DrawCircle(x, y, radius, (Color){ r, g, b, a });
}

void cb_draw_text(const char* text, int x, int y, int fontSize, int r, int g, int b, int a) {
    DrawText(text, x, y, fontSize, (Color){ r, g, b, a });
}

// 3D Operations
cb_camera3d cb_create_camera(cb_vector3 position, cb_vector3 target, cb_vector3 up, float fovy) {
    Camera3D camera = {
        { position.x, position.y, position.z },
        { target.x, target.y, target.z },
        { up.x, up.y, up.z },
        fovy,
        CAMERA_PERSPECTIVE
    };
    return (cb_camera3d){
        { camera.position.x, camera.position.y, camera.position.z },
        { camera.target.x, camera.target.y, camera.target.z },
        { camera.up.x, camera.up.y, camera.up.z },
        camera.fovy,
        camera.type
    };
}

void cb_update_camera(cb_camera3d* camera) {
    Camera3D cam = {
        { camera->position.x, camera->position.y, camera->position.z },
        { camera->target.x, camera->target.y, camera->target.z },
        { camera->up.x, camera->up.y, camera->up.z },
        camera->fovy,
        camera->type
    };
    UpdateCamera(&cam);
    
    // Update the wrapper struct
    camera->position = (cb_vector3){ cam.position.x, cam.position.y, cam.position.z };
    camera->target = (cb_vector3){ cam.target.x, cam.target.y, cam.target.z };
    camera->up = (cb_vector3){ cam.up.x, cam.up.y, cam.up.z };
}

void cb_begin_mode_3d(cb_camera3d camera) {
    Camera3D cam = {
        { camera.position.x, camera.position.y, camera.position.z },
        { camera.target.x, camera.target.y, camera.target.z },
        { camera.up.x, camera.up.y, camera.up.z },
        camera.fovy,
        camera.type
    };
    BeginMode3D(cam);
}

void cb_end_mode_3d(void) {
    EndMode3D();
}

// Model operations
cb_model cb_load_model(const char* filename) {
    Model model = LoadModel(filename);
    return (cb_model){
        model.meshes,
        model.meshCount,
        model.materials,
        model.materialCount,
        { 
            { model.transform.translation.x, model.transform.translation.y, model.transform.translation.z },
            { model.transform.rotation.x, model.transform.rotation.y, model.transform.rotation.z },
            { model.transform.scale.x, model.transform.scale.y, model.transform.scale.z }
        }
    };
}

void cb_unload_model(cb_model model) {
    Model m = { 
        model.meshes, 
        model.meshCount, 
        model.materials, 
        model.materialCount,
        {
            { model.transform.position.x, model.transform.position.y, model.transform.position.z },
            { model.transform.rotation.x, model.transform.rotation.y, model.transform.rotation.z },
            { model.transform.scale.x, model.transform.scale.y, model.transform.scale.z }
        }
    };
    UnloadModel(m);
}

void cb_draw_model(cb_model model, cb_vector3 position, float scale, int tint) {
    Model m = { 
        model.meshes, 
        model.meshCount, 
        model.materials, 
        model.materialCount,
        {
            { model.transform.position.x, model.transform.position.y, model.transform.position.z },
            { model.transform.rotation.x, model.transform.rotation.y, model.transform.rotation.z },
            { model.transform.scale.x, model.transform.scale.y, model.transform.scale.z }
        }
    };
    DrawModel(m, (Vector3){ position.x, position.y, position.z }, scale, (Color){ tint, tint, tint, 255 });
}

void cb_draw_cube(cb_vector3 position, float size, int r, int g, int b, int a) {
    DrawCube((Vector3){ position.x, position.y, position.z }, size, size, size, (Color){ r, g, b, a });
}

void cb_draw_sphere(cb_vector3 position, float radius, int r, int g, int b, int a) {
    DrawSphere((Vector3){ position.x, position.y, position.z }, radius, (Color){ r, g, b, a });
}

// Input handling
bool cb_is_key_pressed(int key) {
    return IsKeyPressed(key);
}

bool cb_is_key_down(int key) {
    return IsKeyDown(key);
}

bool cb_is_key_released(int key) {
    return IsKeyReleased(key);
}

bool cb_is_key_up(int key) {
    return IsKeyUp(key);
}

bool cb_is_mouse_button_pressed(int button) {
    return IsMouseButtonPressed(button);
}

bool cb_is_mouse_button_down(int button) {
    return IsMouseButtonDown(button);
}

bool cb_is_mouse_button_released(int button) {
    return IsMouseButtonReleased(button);
}

bool cb_is_mouse_button_up(int button) {
    return IsMouseButtonUp(button);
}

cb_vector3 cb_get_mouse_position(void) {
    Vector2 pos = GetMousePosition();
    return (cb_vector3){ pos.x, pos.y, 0.0f };
}

// Audio operations
cb_sound cb_load_sound(const char* filename) {
    Sound sound = LoadSound(filename);
    return (cb_sound){ sound.stream, sound.sampleCount };
}

void cb_unload_sound(cb_sound sound) {
    Sound s = { sound.stream, sound.sampleCount };
    UnloadSound(s);
}

void cb_play_sound(cb_sound sound) {
    Sound s = { sound.stream, sound.sampleCount };
    PlaySound(s);
}

void cb_stop_sound(cb_sound sound) {
    Sound s = { sound.stream, sound.sampleCount };
    StopSound(s);
}

void cb_set_sound_volume(cb_sound sound, float volume) {
    Sound s = { sound.stream, sound.sampleCount };
    SetSoundVolume(s, volume);
}

cb_music cb_load_music_stream(const char* filename) {
    Music music = LoadMusicStream(filename);
    return (cb_music){ music.stream, music.looping };
}

void cb_unload_music_stream(cb_music music) {
    Music m = { music.stream, music.loop };
    UnloadMusicStream(m);
}

void cb_play_music(cb_music music) {
    Music m = { music.stream, music.loop };
    PlayMusicStream(m);
}

void cb_stop_music(cb_music music) {
    Music m = { music.stream, music.loop };
    StopMusicStream(m);
}

void cb_set_music_volume(cb_music music, float volume) {
    Music m = { music.stream, music.loop };
    SetMusicVolume(m, volume);
}

// Utility functions
void cb_set_target_fps(int fps) {
    SetTargetFPS(fps);
}

int cb_get_fps(void) {
    return GetFPS();
}

float cb_get_frame_time(void) {
    return GetFrameTime();
}
