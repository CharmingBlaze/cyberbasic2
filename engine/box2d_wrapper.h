#ifndef BOX2D_WRAPPER_H
#define BOX2D_WRAPPER_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Vector2 structure for 2D physics
typedef struct {
    float x, y;
} cb_vector2;

// 2D Transform structure
typedef struct {
    cb_vector2 position;
    float angle; // in radians
} cb_transform2d;

// 2D body types
typedef enum {
    CB_BODY_2D_STATIC,
    CB_BODY_2D_DYNAMIC,
    CB_BODY_2D_KINEMATIC
} cb_body_2d_type;

// 2D collision shapes
typedef enum {
    CB_SHAPE_2D_BOX,
    CB_SHAPE_2D_CIRCLE,
    CB_SHAPE_2D_EDGE,
    CB_SHAPE_2D_POLYGON,
    CB_SHAPE_2D_CHAIN
} cb_shape_2d_type;

// 2D physics body handle
typedef struct cb_physics_body_2d* cb_physics_body_2d_t;

// 2D physics world handle
typedef struct cb_physics_world_2d* cb_physics_world_2d_t;

// World management
cb_physics_world_2d_t cb_physics_2d_create_world(cb_vector2 gravity);
void cb_physics_2d_destroy_world(cb_physics_world_2d_t world);
void cb_physics_2d_step_simulation(cb_physics_world_2d_t world, float time_step, int velocity_iterations, int position_iterations);

// Body management
cb_physics_body_2d_t cb_physics_2d_create_body(cb_physics_world_2d_t world, cb_body_2d_type type, cb_shape_2d_type shape, cb_vector2 size, float density);
void cb_physics_2d_destroy_body(cb_physics_world_2d_t world, cb_physics_body_2d_t body);
void cb_physics_2d_set_transform(cb_physics_body_2d_t body, cb_transform2d transform);
cb_transform2d cb_physics_2d_get_transform(cb_physics_body_2d_t body);
void cb_physics_2d_set_position(cb_physics_body_2d_t body, cb_vector2 position);
cb_vector2 cb_physics_2d_get_position(cb_physics_body_2d_t body);
void cb_physics_2d_set_angle(cb_physics_body_2d_t body, float angle);
float cb_physics_2d_get_angle(cb_physics_body_2d_t body);

// Velocity and forces
void cb_physics_2d_set_linear_velocity(cb_physics_body_2d_t body, cb_vector2 velocity);
cb_vector2 cb_physics_2d_get_linear_velocity(cb_physics_body_2d_t body);
void cb_physics_2d_set_angular_velocity(cb_physics_body_2d_t body, float velocity);
float cb_physics_2d_get_angular_velocity(cb_physics_body_2d_t body);
void cb_physics_2d_apply_force(cb_physics_body_2d_t body, cb_vector2 force, cb_vector2 point);
void cb_physics_2d_apply_force_to_center(cb_physics_body_2d_t body, cb_vector2 force);
void cb_physics_2d_apply_linear_impulse(cb_physics_body_2d_t body, cb_vector2 impulse, cb_vector2 point);
void cb_physics_2d_apply_linear_impulse_to_center(cb_physics_body_2d_t body, cb_vector2 impulse);
void cb_physics_2d_apply_torque(cb_physics_body_2d_t body, float torque);
void cb_physics_2d_apply_angular_impulse(cb_physics_body_2d_t body, float impulse);

// Mass and properties
void cb_physics_2d_set_density(cb_physics_body_2d_t body, float density);
float cb_physics_2d_get_density(cb_physics_body_2d_t body);
void cb_physics_2d_set_friction(cb_physics_body_2d_t body, float friction);
float cb_physics_2d_get_friction(cb_physics_body_2d_t body);
void cb_physics_2d_set_restitution(cb_physics_body_2d_t body, float restitution);
float cb_physics_2d_get_restitution(cb_physics_body_2d_t body);
void cb_physics_2d_set_gravity_scale(cb_physics_body_2d_t body, float scale);
float cb_physics_2d_get_gravity_scale(cb_physics_body_2d_t body);

// 2D Constraints and joints
typedef struct cb_constraint_2d* cb_constraint_2d_t;

cb_constraint_2d_t cb_physics_2d_create_revolute_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor);
cb_constraint_2d_t cb_physics_2d_create_prismatic_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, cb_vector2 axis);
cb_constraint_2d_t cb_physics_2d_create_distance_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor_a, cb_vector2 anchor_b, float length);
cb_constraint_2d_t cb_physics_2d_create_pulley_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 ground_anchor_a, cb_vector2 ground_anchor_b, cb_vector2 anchor_a, cb_vector2 anchor_b, float ratio);
cb_constraint_2d_t cb_physics_2d_create_mouse_joint(cb_physics_body_2d_t body, cb_vector2 target);
cb_constraint_2d_t cb_physics_2d_create_gear_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_constraint_2d_t joint_a, cb_constraint_2d_t joint_b, float ratio);
cb_constraint_2d_t cb_physics_2d_create_wheel_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, cb_vector2 axis, float damping);
cb_constraint_2d_t cb_physics_2d_create_weld_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, float angle);
cb_constraint_2d_t cb_physics_2d_create_friction_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, float max_force, float max_torque);

void cb_physics_2d_destroy_joint(cb_constraint_2d_t joint);

// Joint properties
void cb_physics_2d_joint_set_motor_speed(cb_constraint_2d_t joint, float speed);
void cb_physics_2d_joint_set_max_motor_force(cb_constraint_2d_t joint, float force);
void cb_physics_2d_joint_set_limits(cb_constraint_2d_t joint, float lower, float upper);
void cb_physics_2d_joint_set_frequency(cb_constraint_2d_t joint, float hz);
void cb_physics_2d_joint_set_damping(cb_constraint_2d_t joint, float damping);

// 2D Raycasting
typedef struct {
    bool hit;
    cb_vector2 point;
    cb_vector2 normal;
    float fraction;
    cb_physics_body_2d_t body;
} cb_ray_cast_2d_result;

cb_ray_cast_2d_result cb_physics_2d_ray_cast(cb_physics_world_2d_t world, cb_vector2 start, cb_vector2 end);

// 2D Collision detection and queries
typedef struct {
    bool colliding;
    cb_vector2 contact_points[2];
    cb_vector2 contact_normal;
    int contact_count;
    float separation;
} cb_collision_2d_result;

cb_collision_2d_result cb_physics_2d_check_collision(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b);

// AABB queries
typedef struct {
    cb_vector2 lower_bound;
    cb_vector2 upper_bound;
} cb_aabb_2d;

typedef struct {
    cb_physics_body_2d_t* bodies;
    int count;
} cb_query_2d_result;

cb_query_2d_result cb_physics_2d_query_aabb(cb_physics_world_2d_t world, cb_aabb_2d aabb);
cb_physics_body_2d_t cb_physics_2d_query_point(cb_physics_world_2d_t world, cb_vector2 point);

// Shape casting
typedef struct {
    bool hit;
    cb_vector2 point;
    cb_vector2 normal;
    float fraction;
    cb_physics_body_2d_t body;
} cb_shape_cast_2d_result;

cb_shape_cast_2d_result cb_physics_2d_shape_cast(cb_physics_world_2d_t world, cb_shape_2d_type shape, cb_vector2 shape_size, cb_transform2d transform, cb_vector2 translation);

// World properties
void cb_physics_2d_set_gravity(cb_physics_world_2d_t world, cb_vector2 gravity);
cb_vector2 cb_physics_2d_get_gravity(cb_physics_world_2d_t world);
void cb_physics_2d_set_allow_sleeping(cb_physics_world_2d_t world, bool allow);
bool cb_physics_2d_get_allow_sleeping(cb_physics_world_2d_t world);

// Body filtering and categories
typedef uint16_t cb_filter_category;
void cb_physics_2d_set_filter_category(cb_physics_body_2d_t body, cb_filter_category category);
void cb_physics_2d_set_filter_mask(cb_physics_body_2d_t body, cb_filter_category mask);
void cb_physics_2d_set_filter_group_index(cb_physics_body_2d_t body, int group_index);

// Debug drawing
typedef struct {
    void (*draw_circle)(cb_vector2 center, float radius, cb_vector2 color);
    void (*draw_segment)(cb_vector2 p1, cb_vector2 p2, cb_vector2 color);
    void (*draw_polygon)(cb_vector2* vertices, int vertex_count, cb_vector2 color);
    void (*draw_solid_polygon)(cb_vector2* vertices, int vertex_count, cb_vector2 color);
} cb_debug_draw_2d;

void cb_physics_2d_set_debug_draw(cb_physics_world_2d_t world, cb_debug_draw_2d* debug_draw);
void cb_physics_2d_draw_debug_data(cb_physics_world_2d_t world);

// Utility functions
cb_vector2 cb_physics_2d_vector_create(float x, float y);
cb_vector2 cb_physics_2d_vector_add(cb_vector2 a, cb_vector2 b);
cb_vector2 cb_physics_2d_vector_subtract(cb_vector2 a, cb_vector2 b);
cb_vector2 cb_physics_2d_vector_multiply(cb_vector2 v, float scalar);
cb_vector2 cb_physics_2d_vector_normalize(cb_vector2 v);
float cb_physics_2d_vector_length(cb_vector2 v);
float cb_physics_2d_vector_length_squared(cb_vector2 v);
float cb_physics_2d_vector_dot(cb_vector2 a, cb_vector2 b);
float cb_physics_2d_vector_cross(cb_vector2 a, cb_vector2 b);
cb_vector2 cb_physics_2d_vector_cross_float(float s, cb_vector2 a);
cb_vector2 cb_physics_2d_vector_cross_vector(cb_vector2 a, float s);

// Common filter categories
#define CB_FILTER_2D_CATEGORY_1   0x0001
#define CB_FILTER_2D_CATEGORY_2   0x0002
#define CB_FILTER_2D_CATEGORY_3   0x0004
#define CB_FILTER_2D_CATEGORY_4   0x0008
#define CB_FILTER_2D_CATEGORY_5   0x0010
#define CB_FILTER_2D_CATEGORY_6   0x0020
#define CB_FILTER_2D_CATEGORY_7   0x0040
#define CB_FILTER_2D_CATEGORY_8   0x0080
#define CB_FILTER_2D_CATEGORY_9   0x0100
#define CB_FILTER_2D_CATEGORY_10  0x0200
#define CB_FILTER_2D_CATEGORY_11  0x0400
#define CB_FILTER_2D_CATEGORY_12  0x0800
#define CB_FILTER_2D_CATEGORY_13  0x1000
#define CB_FILTER_2D_CATEGORY_14  0x2000
#define CB_FILTER_2D_CATEGORY_15  0x4000
#define CB_FILTER_2D_CATEGORY_16  0x8000

#define CB_FILTER_2D_ALL_CATEGORIES 0xFFFF
#define CB_FILTER_2D_ALL_MASK       0xFFFF

#ifdef __cplusplus
}
#endif

#endif // BOX2D_WRAPPER_H
