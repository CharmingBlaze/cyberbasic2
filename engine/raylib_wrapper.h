#ifndef RAYLIB_WRAPPER_H
#define RAYLIB_WRAPPER_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Graphics initialization
bool cb_init_window(int width, int height, const char* title);
void cb_close_window(void);
bool cb_window_should_close(void);
void cb_begin_drawing(void);
void cb_end_drawing(void);
void cb_clear_background(int r, int g, int b);

// Texture/Image operations
typedef struct {
    unsigned int id;
    int width;
    int height;
    int mipmaps;
    int format;
} cb_texture;

typedef struct {
    void* data;
    int width;
    int height;
    int mipmaps;
    int format;
} cb_image;

cb_image cb_load_image(const char* filename);
void cb_unload_image(cb_image image);
cb_texture cb_load_texture_from_image(cb_image image);
void cb_unload_texture(cb_texture texture);
void cb_draw_texture(cb_texture texture, int x, int y, int tint);

// 2D Drawing
void cb_draw_rectangle(int x, int y, int width, int height, int r, int g, int b, int a);
void cb_draw_circle(int x, int y, int radius, int r, int g, int b, int a);
void cb_draw_text(const char* text, int x, int y, int fontSize, int r, int g, int b, int a);

// 3D Operations
typedef struct {
    float x, y, z;
} cb_vector3;

typedef struct {
    cb_vector3 position;
    cb_vector3 target;
    cb_vector3 up;
    float fovy;
    int type;
} cb_camera3d;

typedef struct {
    cb_vector3 position;
    cb_vector3 size;
    cb_vector3 rotation;
} cb_bounding_box;

cb_camera3d cb_create_camera(cb_vector3 position, cb_vector3 target, cb_vector3 up, float fovy);
void cb_update_camera(cb_camera3d* camera);
void cb_begin_mode_3d(cb_camera3d camera);
void cb_end_mode_3d(void);

// Model operations
typedef struct {
    cb_vector3 position;
    cb_vector3 rotation;
    cb_vector3 scale;
} cb_transform;

typedef struct {
    void* meshes;
    int meshCount;
    void* materials;
    int materialCount;
    cb_transform transform;
} cb_model;

cb_model cb_load_model(const char* filename);
void cb_unload_model(cb_model model);
void cb_draw_model(cb_model model, cb_vector3 position, float scale, int tint);
void cb_draw_cube(cb_vector3 position, float size, int r, int g, int b, int a);
void cb_draw_sphere(cb_vector3 position, float radius, int r, int g, int b, int a);

// Input handling
bool cb_is_key_pressed(int key);
bool cb_is_key_down(int key);
bool cb_is_key_released(int key);
bool cb_is_key_up(int key);
bool cb_is_mouse_button_pressed(int button);
bool cb_is_mouse_button_down(int button);
bool cb_is_mouse_button_released(int button);
bool cb_is_mouse_button_up(int button);
cb_vector3 cb_get_mouse_position(void);

// Audio operations
typedef struct {
    void* stream;
    unsigned int sampleCount;
} cb_sound;

typedef struct {
    void* stream;
    bool loop;
} cb_music;

cb_sound cb_load_sound(const char* filename);
void cb_unload_sound(cb_sound sound);
void cb_play_sound(cb_sound sound);
void cb_stop_sound(cb_sound sound);
void cb_set_sound_volume(cb_sound sound, float volume);
cb_music cb_load_music_stream(const char* filename);
void cb_unload_music_stream(cb_music music);
void cb_play_music(cb_music music);
void cb_stop_music(cb_music music);
void cb_set_music_volume(cb_music music, float volume);

// Utility functions
void cb_set_target_fps(int fps);
int cb_get_fps(void);
float cb_get_frame_time(void);

#ifdef __cplusplus
}
#endif

#endif // RAYLIB_WRAPPER_H
