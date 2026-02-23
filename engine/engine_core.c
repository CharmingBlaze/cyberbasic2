#include "raylib_wrapper.h"
#include "bullet_wrapper.h"
#include <stdio.h>
#include <stdlib.h>

// Engine state structure
typedef struct {
    int screen_width;
    int screen_height;
    const char* title;
    int target_fps;
    bool running;
    
    // Physics world
    cb_physics_world_t physics_world;
    
    // Camera
    cb_camera3d camera;
    bool camera_3d_mode;
    
    // Audio state
    cb_music current_music;
    bool music_playing;
} cb_engine_state;

static cb_engine_state g_engine = {0};

// Engine initialization and management
bool cb_engine_init(int width, int height, const char* title) {
    g_engine.screen_width = width;
    g_engine.screen_height = height;
    g_engine.title = title;
    g_engine.target_fps = 60;
    g_engine.running = false;
    g_engine.camera_3d_mode = false;
    g_engine.music_playing = false;
    
    // Initialize graphics
    if (!cb_init_window(width, height, title)) {
        printf("Failed to initialize window\n");
        return false;
    }
    
    // Initialize physics
    g_engine.physics_world = cb_physics_create_world((cb_vector3){0.0f, -9.81f, 0.0f});
    if (!g_engine.physics_world) {
        printf("Failed to create physics world\n");
        cb_close_window();
        return false;
    }
    
    // Set up default camera
    g_engine.camera = cb_create_camera(
        (cb_vector3){10.0f, 10.0f, 10.0f},
        (cb_vector3){0.0f, 0.0f, 0.0f},
        (cb_vector3){0.0f, 1.0f, 0.0f},
        45.0f
    );
    
    cb_set_target_fps(g_engine.target_fps);
    g_engine.running = true;
    
    printf("Engine initialized successfully\n");
    return true;
}

void cb_engine_shutdown(void) {
    if (g_engine.music_playing) {
        cb_stop_music(g_engine.current_music);
        cb_unload_music_stream(g_engine.current_music);
    }
    
    if (g_engine.physics_world) {
        cb_physics_destroy_world(g_engine.physics_world);
    }
    
    cb_close_window();
    g_engine.running = false;
    printf("Engine shutdown complete\n");
}

bool cb_engine_is_running(void) {
    return g_engine.running && !cb_window_should_close();
}

void cb_engine_begin_frame(void) {
    cb_begin_drawing();
    cb_clear_background(135, 206, 235); // Sky blue background
}

void cb_engine_end_frame(void) {
    cb_end_drawing();
}

void cb_engine_update_physics(float delta_time) {
    if (g_engine.physics_world) {
        cb_physics_step_simulation(g_engine.physics_world, delta_time);
    }
}

// Graphics functions
void cb_engine_begin_3d_mode(void) {
    cb_begin_mode_3d(g_engine.camera);
    g_engine.camera_3d_mode = true;
}

void cb_engine_end_3d_mode(void) {
    cb_end_mode_3d();
    g_engine.camera_3d_mode = false;
}

void cb_engine_set_camera_position(float x, float y, float z) {
    g_engine.camera.position = (cb_vector3){x, y, z};
}

void cb_engine_set_camera_target(float x, float y, float z) {
    g_engine.camera.target = (cb_vector3){x, y, z};
}

// Physics functions
cb_physics_body_t cb_engine_create_physics_body(int body_type, int shape_type, float x, float y, float z, float size_x, float size_y, float size_z, float mass) {
    if (!g_engine.physics_world) return NULL;
    
    cb_body_type type = (cb_body_type)body_type;
    cb_shape_type shape = (cb_shape_type)shape_type;
    cb_vector3 size = {size_x, size_y, size_z};
    
    cb_physics_body_t body = cb_physics_create_body(g_engine.physics_world, type, shape, size, mass);
    if (body) {
        cb_physics_set_position(body, (cb_vector3){x, y, z});
    }
    
    return body;
}

void cb_engine_set_body_position(cb_physics_body_t body, float x, float y, float z) {
    if (body) {
        cb_physics_set_position(body, (cb_vector3){x, y, z});
    }
}

void cb_engine_set_body_velocity(cb_physics_body_t body, float vx, float vy, float vz) {
    if (body) {
        cb_physics_set_linear_velocity(body, (cb_vector3){vx, vy, vz});
    }
}

void cb_engine_apply_force(cb_physics_body_t body, float fx, float fy, float fz) {
    if (body) {
        cb_physics_apply_central_force(body, (cb_vector3){fx, fy, fz});
    }
}

// Raycasting
bool cb_engine_ray_cast(float start_x, float start_y, float start_z, float dir_x, float dir_y, float dir_z, float max_distance, float* hit_x, float* hit_y, float* hit_z) {
    if (!g_engine.physics_world) return false;
    
    cb_vector3 start = {start_x, start_y, start_z};
    cb_vector3 end = {
        start_x + dir_x * max_distance,
        start_y + dir_y * max_distance,
        start_z + dir_z * max_distance
    };
    
    cb_ray_cast_result result = cb_physics_ray_cast(g_engine.physics_world, start, end);
    
    if (result.hit && hit_x && hit_y && hit_z) {
        *hit_x = result.hit_point.x;
        *hit_y = result.hit_point.y;
        *hit_z = result.hit_point.z;
    }
    
    return result.hit;
}

// Audio functions
bool cb_engine_load_music(const char* filename) {
    if (g_engine.music_playing) {
        cb_stop_music(g_engine.current_music);
        cb_unload_music_stream(g_engine.current_music);
        g_engine.music_playing = false;
    }
    
    g_engine.current_music = cb_load_music_stream(filename);
    return g_engine.current_music.stream != NULL;
}

void cb_engine_play_music(void) {
    if (!g_engine.music_playing && g_engine.current_music.stream) {
        cb_play_music(g_engine.current_music);
        g_engine.music_playing = true;
    }
}

void cb_engine_stop_music(void) {
    if (g_engine.music_playing) {
        cb_stop_music(g_engine.current_music);
        g_engine.music_playing = false;
    }
}

void cb_engine_set_music_volume(float volume) {
    if (g_engine.current_music.stream) {
        cb_set_music_volume(g_engine.current_music, volume);
    }
}

// Utility functions
int cb_engine_get_fps(void) {
    return cb_get_fps();
}

float cb_engine_get_frame_time(void) {
    return cb_get_frame_time();
}

// Input functions
bool cb_engine_is_key_pressed(int key) {
    return cb_is_key_pressed(key);
}

bool cb_engine_is_key_down(int key) {
    return cb_is_key_down(key);
}

bool cb_engine_is_key_released(int key) {
    return cb_is_key_released(key);
}

bool cb_engine_is_key_up(int key) {
    return cb_is_key_up(key);
}

void cb_engine_get_mouse_position(float* x, float* y) {
    cb_vector3 pos = cb_get_mouse_position();
    if (x) *x = pos.x;
    if (y) *y = pos.y;
}

// Drawing helpers
void cb_engine_draw_grid(float size, int slices) {
    // Simple grid drawing using Raylib
    for (int i = -slices; i <= slices; i++) {
        float pos = i * size;
        // Draw lines in X direction
        cb_draw_cube((cb_vector3){pos, 0, 0}, 0.1f, 100, 100, 100, 255);
        // Draw lines in Z direction  
        cb_draw_cube((cb_vector3){0, 0, pos}, 0.1f, 100, 100, 100, 255);
    }
}

void cb_engine_draw_axes(void) {
    // X axis - red
    cb_draw_cube((cb_vector3){5.0f, 0.0f, 0.0f}, 10.0f, 255, 0, 0, 255);
    // Y axis - green
    cb_draw_cube((cb_vector3){0.0f, 5.0f, 0.0f}, 10.0f, 0, 255, 0, 255);
    // Z axis - blue
    cb_draw_cube((cb_vector3){0.0f, 0.0f, 5.0f}, 10.0f, 0, 0, 255, 255);
}

// Debug information
void cb_engine_draw_debug_info(void) {
    char fps_text[32];
    sprintf(fps_text, "FPS: %d", cb_engine_get_fps());
    cb_draw_text(fps_text, 10, 10, 20, 255, 255, 255, 255);
    
    char frame_time_text[32];
    sprintf(frame_time_text, "Frame Time: %.3f ms", cb_engine_get_frame_time() * 1000.0f);
    cb_draw_text(frame_time_text, 10, 35, 20, 255, 255, 255, 255);
    
    if (g_engine.physics_world) {
        cb_draw_text("Physics: Active", 10, 60, 20, 0, 255, 0, 255);
    } else {
        cb_draw_text("Physics: Inactive", 10, 60, 20, 255, 0, 0, 255);
    }
    
    if (g_engine.camera_3d_mode) {
        cb_draw_text("Mode: 3D", 10, 85, 20, 255, 255, 255, 255);
    } else {
        cb_draw_text("Mode: 2D", 10, 85, 20, 255, 255, 255, 255);
    }
}
