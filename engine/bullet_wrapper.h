#ifndef BULLET_WRAPPER_H
#define BULLET_WRAPPER_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Vector3 structure for physics
typedef struct {
    float x, y, z;
} cb_vector3;

// Quaternion structure for rotation
typedef struct {
    float x, y, z, w;
} cb_quaternion;

// Transform structure
typedef struct {
    cb_vector3 position;
    cb_quaternion rotation;
} cb_transform;

// Physics body types
typedef enum {
    CB_BODY_STATIC,
    CB_BODY_DYNAMIC,
    CB_BODY_KINEMATIC
} cb_body_type;

// Collision shapes
typedef enum {
    CB_SHAPE_BOX,
    CB_SHAPE_SPHERE,
    CB_SHAPE_PLANE,
    CB_SHAPE_CYLINDER,
    CB_SHAPE_CAPSULE,
    CB_SHAPE_MESH
} cb_shape_type;

// Physics body handle
typedef struct cb_physics_body* cb_physics_body_t;

// Physics world handle
typedef struct cb_physics_world* cb_physics_world_t;

// World management
cb_physics_world_t cb_physics_create_world(cb_vector3 gravity);
void cb_physics_destroy_world(cb_physics_world_t world);
void cb_physics_step_simulation(cb_physics_world_t world, float time_step);

// Body management
cb_physics_body_t cb_physics_create_body(cb_physics_world_t world, cb_body_type type, cb_shape_type shape, cb_vector3 size, float mass);
void cb_physics_destroy_body(cb_physics_world_t world, cb_physics_body_t body);
void cb_physics_set_transform(cb_physics_body_t body, cb_transform transform);
cb_transform cb_physics_get_transform(cb_physics_body_t body);
void cb_physics_set_position(cb_physics_body_t body, cb_vector3 position);
cb_vector3 cb_physics_get_position(cb_physics_body_t body);
void cb_physics_set_rotation(cb_physics_body_t body, cb_quaternion rotation);
cb_quaternion cb_physics_get_rotation(cb_physics_body_t body);

// Velocity and forces
void cb_physics_set_linear_velocity(cb_physics_body_t body, cb_vector3 velocity);
cb_vector3 cb_physics_get_linear_velocity(cb_physics_body_t body);
void cb_physics_set_angular_velocity(cb_physics_body_t body, cb_vector3 velocity);
cb_vector3 cb_physics_get_angular_velocity(cb_physics_body_t body);
void cb_physics_apply_central_force(cb_physics_body_t body, cb_vector3 force);
void cb_physics_apply_force(cb_physics_body_t body, cb_vector3 force, cb_vector3 relative_position);
void cb_physics_apply_impulse(cb_physics_body_t body, cb_vector3 impulse, cb_vector3 relative_position);
void cb_physics_apply_torque(cb_physics_body_t body, cb_vector3 torque);
void cb_physics_apply_torque_impulse(cb_physics_body_t body, cb_vector3 torque);

// Mass and inertia
void cb_physics_set_mass(cb_physics_body_t body, float mass);
float cb_physics_get_mass(cb_physics_body_t body);
void cb_physics_set_friction(cb_physics_body_t body, float friction);
float cb_physics_get_friction(cb_physics_body_t body);
void cb_physics_set_restitution(cb_physics_body_t body, float restitution);
float cb_physics_get_restitution(cb_physics_body_t body);

// Constraints and joints
typedef struct cb_constraint* cb_constraint_t;

cb_constraint_t cb_physics_create_point_constraint(cb_physics_body_t body_a, cb_physics_body_t body_b, cb_vector3 pivot_a, cb_vector3 pivot_b);
cb_constraint_t cb_physics_create_hinge_constraint(cb_physics_body_t body_a, cb_physics_body_t body_b, cb_vector3 pivot_a, cb_vector3 pivot_b, cb_vector3 axis_a, cb_vector3 axis_b);
cb_constraint_t cb_physics_create_slider_constraint(cb_physics_body_t body_a, cb_physics_body_t body_b, cb_vector3 pivot_a, cb_vector3 pivot_b, cb_vector3 axis_a, cb_vector3 axis_b);
void cb_physics_destroy_constraint(cb_constraint_t constraint);

// Raycasting
typedef struct {
    bool hit;
    cb_vector3 hit_point;
    cb_vector3 hit_normal;
    float hit_fraction;
    cb_physics_body_t hit_body;
} cb_ray_cast_result;

cb_ray_cast_result cb_physics_ray_cast(cb_physics_world_t world, cb_vector3 start, cb_vector3 end);

// Collision detection
typedef struct {
    bool colliding;
    cb_vector3 contact_point;
    cb_vector3 contact_normal;
    float penetration_depth;
} cb_collision_result;

cb_collision_result cb_physics_check_collision(cb_physics_body_t body_a, cb_physics_body_t body_b);

// Utility functions
cb_vector3 cb_physics_vector3_create(float x, float y, float z);
cb_quaternion cb_physics_quaternion_identity(void);
cb_quaternion cb_physics_quaternion_from_euler(float yaw, float pitch, float roll);
cb_quaternion cb_physics_quaternion_multiply(cb_quaternion q1, cb_quaternion q2);
cb_vector3 cb_physics_quaternion_rotate_vector(cb_quaternion q, cb_vector3 v);

#ifdef __cplusplus
}
#endif

#endif // BULLET_WRAPPER_H
