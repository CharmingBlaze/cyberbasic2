#include "bullet_wrapper.h"
#include <stdlib.h>
#include <string.h>
#include <math.h>

// Simplified physics implementation for demonstration
// In a real implementation, this would interface with Bullet Physics

struct cb_physics_body {
    cb_body_type type;
    cb_shape_type shape;
    cb_vector3 size;
    float mass;
    cb_transform transform;
    cb_vector3 linear_velocity;
    cb_vector3 angular_velocity;
    float friction;
    float restitution;
    bool active;
};

struct cb_physics_world {
    cb_vector3 gravity;
    cb_physics_body_t* bodies;
    int body_count;
    int body_capacity;
};

struct cb_constraint {
    int type; // 0=point, 1=hinge, 2=slider
    cb_physics_body_t body_a;
    cb_physics_body_t body_b;
    cb_vector3 pivot_a;
    cb_vector3 pivot_b;
    cb_vector3 axis_a;
    cb_vector3 axis_b;
};

// World management
cb_physics_world_t cb_physics_create_world(cb_vector3 gravity) {
    cb_physics_world_t world = malloc(sizeof(struct cb_physics_world));
    world->gravity = gravity;
    world->bodies = malloc(sizeof(cb_physics_body_t) * 100);
    world->body_count = 0;
    world->body_capacity = 100;
    return world;
}

void cb_physics_destroy_world(cb_physics_world_t world) {
    if (world) {
        for (int i = 0; i < world->body_count; i++) {
            free(world->bodies[i]);
        }
        free(world->bodies);
        free(world);
    }
}

void cb_physics_step_simulation(cb_physics_world_t world, float time_step) {
    if (!world) return;
    
    // Simple physics simulation
    for (int i = 0; i < world->body_count; i++) {
        cb_physics_body_t body = world->bodies[i];
        if (!body || !body->active || body->type == CB_BODY_STATIC) continue;
        
        // Apply gravity
        body->linear_velocity.y += world->gravity.y * time_step;
        
        // Update position
        body->transform.position.x += body->linear_velocity.x * time_step;
        body->transform.position.y += body->linear_velocity.y * time_step;
        body->transform.position.z += body->linear_velocity.z * time_step;
        
        // Simple ground collision
        if (body->transform.position.y < body->size.y * 0.5f) {
            body->transform.position.y = body->size.y * 0.5f;
            body->linear_velocity.y *= -body->restitution;
            body->linear_velocity.x *= (1.0f - body->friction);
            body->linear_velocity.z *= (1.0f - body->friction);
        }
    }
}

// Body management
cb_physics_body_t cb_physics_create_body(cb_physics_world_t world, cb_body_type type, cb_shape_type shape, cb_vector3 size, float mass) {
    if (!world || world->body_count >= world->body_capacity) return NULL;
    
    cb_physics_body_t body = malloc(sizeof(struct cb_physics_body));
    body->type = type;
    body->shape = shape;
    body->size = size;
    body->mass = mass;
    body->transform.position = (cb_vector3){ 0, 0, 0 };
    body->transform.rotation = cb_physics_quaternion_identity();
    body->linear_velocity = (cb_vector3){ 0, 0, 0 };
    body->angular_velocity = (cb_vector3){ 0, 0, 0 };
    body->friction = 0.5f;
    body->restitution = 0.1f;
    body->active = true;
    
    world->bodies[world->body_count++] = body;
    return body;
}

void cb_physics_destroy_body(cb_physics_world_t world, cb_physics_body_t body) {
    if (!world || !body) return;
    
    // Find and remove body from world
    for (int i = 0; i < world->body_count; i++) {
        if (world->bodies[i] == body) {
            free(body);
            // Shift remaining bodies
            for (int j = i; j < world->body_count - 1; j++) {
                world->bodies[j] = world->bodies[j + 1];
            }
            world->body_count--;
            break;
        }
    }
}

void cb_physics_set_transform(cb_physics_body_t body, cb_transform transform) {
    if (body) {
        body->transform = transform;
    }
}

cb_transform cb_physics_get_transform(cb_physics_body_t body) {
    if (body) {
        return body->transform;
    }
    return (cb_transform){ {0, 0, 0}, {0, 0, 0, 1} };
}

void cb_physics_set_position(cb_physics_body_t body, cb_vector3 position) {
    if (body) {
        body->transform.position = position;
    }
}

cb_vector3 cb_physics_get_position(cb_physics_body_t body) {
    if (body) {
        return body->transform.position;
    }
    return (cb_vector3){ 0, 0, 0 };
}

void cb_physics_set_rotation(cb_physics_body_t body, cb_quaternion rotation) {
    if (body) {
        body->transform.rotation = rotation;
    }
}

cb_quaternion cb_physics_get_rotation(cb_physics_body_t body) {
    if (body) {
        return body->transform.rotation;
    }
    return cb_physics_quaternion_identity();
}

// Velocity and forces
void cb_physics_set_linear_velocity(cb_physics_body_t body, cb_vector3 velocity) {
    if (body) {
        body->linear_velocity = velocity;
    }
}

cb_vector3 cb_physics_get_linear_velocity(cb_physics_body_t body) {
    if (body) {
        return body->linear_velocity;
    }
    return (cb_vector3){ 0, 0, 0 };
}

void cb_physics_set_angular_velocity(cb_physics_body_t body, cb_vector3 velocity) {
    if (body) {
        body->angular_velocity = velocity;
    }
}

cb_vector3 cb_physics_get_angular_velocity(cb_physics_body_t body) {
    if (body) {
        return body->angular_velocity;
    }
    return (cb_vector3){ 0, 0, 0 };
}

void cb_physics_apply_central_force(cb_physics_body_t body, cb_vector3 force) {
    if (body && body->type == CB_BODY_DYNAMIC) {
        // F = ma, so a = F/m
        body->linear_velocity.x += (force.x / body->mass) * 0.016f; // Assuming 60 FPS
        body->linear_velocity.y += (force.y / body->mass) * 0.016f;
        body->linear_velocity.z += (force.z / body->mass) * 0.016f;
    }
}

void cb_physics_apply_force(cb_physics_body_t body, cb_vector3 force, cb_vector3 relative_position) {
    // Simplified - just apply as central force
    cb_physics_apply_central_force(body, force);
}

void cb_physics_apply_impulse(cb_physics_body_t body, cb_vector3 impulse, cb_vector3 relative_position) {
    if (body && body->type == CB_BODY_DYNAMIC) {
        // Impulse changes velocity directly: Δv = J/m
        body->linear_velocity.x += impulse.x / body->mass;
        body->linear_velocity.y += impulse.y / body->mass;
        body->linear_velocity.z += impulse.z / body->mass;
    }
}

void cb_physics_apply_torque(cb_physics_body_t body, cb_vector3 torque) {
    if (body && body->type == CB_BODY_DYNAMIC) {
        // Simplified torque application
        body->angular_velocity.x += torque.x * 0.016f;
        body->angular_velocity.y += torque.y * 0.016f;
        body->angular_velocity.z += torque.z * 0.016f;
    }
}

void cb_physics_apply_torque_impulse(cb_physics_body_t body, cb_vector3 torque_impulse) {
    if (body && body->type == CB_BODY_DYNAMIC) {
        body->angular_velocity.x += torque_impulse.x;
        body->angular_velocity.y += torque_impulse.y;
        body->angular_velocity.z += torque_impulse.z;
    }
}

// Mass and inertia
void cb_physics_set_mass(cb_physics_body_t body, float mass) {
    if (body) {
        body->mass = mass;
    }
}

float cb_physics_get_mass(cb_physics_body_t body) {
    if (body) {
        return body->mass;
    }
    return 0.0f;
}

void cb_physics_set_friction(cb_physics_body_t body, float friction) {
    if (body) {
        body->friction = friction;
    }
}

float cb_physics_get_friction(cb_physics_body_t body) {
    if (body) {
        return body->friction;
    }
    return 0.0f;
}

void cb_physics_set_restitution(cb_physics_body_t body, float restitution) {
    if (body) {
        body->restitution = restitution;
    }
}

float cb_physics_get_restitution(cb_physics_body_t body) {
    if (body) {
        return body->restitution;
    }
    return 0.0f;
}

// Constraints (simplified implementation)
cb_constraint_t cb_physics_create_point_constraint(cb_physics_body_t body_a, cb_physics_body_t body_b, cb_vector3 pivot_a, cb_vector3 pivot_b) {
    cb_constraint_t constraint = malloc(sizeof(struct cb_constraint));
    constraint->type = 0;
    constraint->body_a = body_a;
    constraint->body_b = body_b;
    constraint->pivot_a = pivot_a;
    constraint->pivot_b = pivot_b;
    return constraint;
}

cb_constraint_t cb_physics_create_hinge_constraint(cb_physics_body_t body_a, cb_physics_body_t body_b, cb_vector3 pivot_a, cb_vector3 pivot_b, cb_vector3 axis_a, cb_vector3 axis_b) {
    cb_constraint_t constraint = malloc(sizeof(struct cb_constraint));
    constraint->type = 1;
    constraint->body_a = body_a;
    constraint->body_b = body_b;
    constraint->pivot_a = pivot_a;
    constraint->pivot_b = pivot_b;
    constraint->axis_a = axis_a;
    constraint->axis_b = axis_b;
    return constraint;
}

cb_constraint_t cb_physics_create_slider_constraint(cb_physics_body_t body_a, cb_physics_body_t body_b, cb_vector3 pivot_a, cb_vector3 pivot_b, cb_vector3 axis_a, cb_vector3 axis_b) {
    cb_constraint_t constraint = malloc(sizeof(struct cb_constraint));
    constraint->type = 2;
    constraint->body_a = body_a;
    constraint->body_b = body_b;
    constraint->pivot_a = pivot_a;
    constraint->pivot_b = pivot_b;
    constraint->axis_a = axis_a;
    constraint->axis_b = axis_b;
    return constraint;
}

void cb_physics_destroy_constraint(cb_constraint_t constraint) {
    if (constraint) {
        free(constraint);
    }
}

// Raycasting
cb_ray_cast_result cb_physics_ray_cast(cb_physics_world_t world, cb_vector3 start, cb_vector3 end) {
    cb_ray_cast_result result = { false, {0, 0, 0}, {0, 0, 0}, 0.0f, NULL };
    
    if (!world) return result;
    
    // Simple raycast implementation
    cb_vector3 direction = {
        end.x - start.x,
        end.y - start.y,
        end.z - start.z
    };
    float distance = sqrtf(direction.x * direction.x + direction.y * direction.y + direction.z * direction.z);
    
    if (distance > 0.0f) {
        direction.x /= distance;
        direction.y /= distance;
        direction.z /= distance;
    }
    
    // Check intersection with bodies
    for (int i = 0; i < world->body_count; i++) {
        cb_physics_body_t body = world->bodies[i];
        if (!body || !body->active) continue;
        
        // Simple box intersection
        cb_vector3 min = {
            body->transform.position.x - body->size.x * 0.5f,
            body->transform.position.y - body->size.y * 0.5f,
            body->transform.position.z - body->size.z * 0.5f
        };
        cb_vector3 max = {
            body->transform.position.x + body->size.x * 0.5f,
            body->transform.position.y + body->size.y * 0.5f,
            body->transform.position.z + body->size.z * 0.5f
        };
        
        // Check if ray intersects with box (simplified)
        if (start.x >= min.x && start.x <= max.x &&
            start.y >= min.y && start.y <= max.y &&
            start.z >= min.z && start.z <= max.z) {
            result.hit = true;
            result.hit_point = start;
            result.hit_normal = (cb_vector3){ 0, 1, 0 };
            result.hit_fraction = 0.0f;
            result.hit_body = body;
            break;
        }
    }
    
    return result;
}

// Collision detection
cb_collision_result cb_physics_check_collision(cb_physics_body_t body_a, cb_physics_body_t body_b) {
    cb_collision_result result = { false, {0, 0, 0}, {0, 0, 0}, 0.0f };
    
    if (!body_a || !body_b || !body_a->active || !body_b->active) return result;
    
    // Simple AABB collision
    cb_vector3 min_a = {
        body_a->transform.position.x - body_a->size.x * 0.5f,
        body_a->transform.position.y - body_a->size.y * 0.5f,
        body_a->transform.position.z - body_a->size.z * 0.5f
    };
    cb_vector3 max_a = {
        body_a->transform.position.x + body_a->size.x * 0.5f,
        body_a->transform.position.y + body_a->size.y * 0.5f,
        body_a->transform.position.z + body_a->size.z * 0.5f
    };
    
    cb_vector3 min_b = {
        body_b->transform.position.x - body_b->size.x * 0.5f,
        body_b->transform.position.y - body_b->size.y * 0.5f,
        body_b->transform.position.z - body_b->size.z * 0.5f
    };
    cb_vector3 max_b = {
        body_b->transform.position.x + body_b->size.x * 0.5f,
        body_b->transform.position.y + body_b->size.y * 0.5f,
        body_b->transform.position.z + body_b->size.z * 0.5f
    };
    
    bool overlap = (max_a.x >= min_b.x && min_a.x <= max_b.x) &&
                   (max_a.y >= min_b.y && min_a.y <= max_b.y) &&
                   (max_a.z >= min_b.z && min_a.z <= max_b.z);
    
    if (overlap) {
        result.colliding = true;
        result.contact_point = (cb_vector3){ 
            (max_a.x + min_a.x + max_b.x + min_b.x) * 0.25f,
            (max_a.y + min_a.y + max_b.y + min_b.y) * 0.25f,
            (max_a.z + min_a.z + max_b.z + min_b.z) * 0.25f
        };
        result.contact_normal = (cb_vector3){ 0, 1, 0 }; // Simplified
        result.penetration_depth = 0.1f; // Simplified
    }
    
    return result;
}

// Utility functions
cb_vector3 cb_physics_vector3_create(float x, float y, float z) {
    return (cb_vector3){ x, y, z };
}

cb_quaternion cb_physics_quaternion_identity(void) {
    return (cb_quaternion){ 0, 0, 0, 1 };
}

cb_quaternion cb_physics_quaternion_from_euler(float yaw, float pitch, float roll) {
    // Simplified quaternion from Euler angles
    float cy = cosf(yaw * 0.5f);
    float sy = sinf(yaw * 0.5f);
    float cp = cosf(pitch * 0.5f);
    float sp = sinf(pitch * 0.5f);
    float cr = cosf(roll * 0.5f);
    float sr = sinf(roll * 0.5f);
    
    cb_quaternion q;
    q.w = cr * cp * cy + sr * sp * sy;
    q.x = sr * cp * cy - cr * sp * sy;
    q.y = cr * sp * cy + sr * cp * sy;
    q.z = cr * cp * sy - sr * sp * cy;
    
    return q;
}

cb_quaternion cb_physics_quaternion_multiply(cb_quaternion q1, cb_quaternion q2) {
    cb_quaternion result;
    result.w = q1.w * q2.w - q1.x * q2.x - q1.y * q2.y - q1.z * q2.z;
    result.x = q1.w * q2.x + q1.x * q2.w + q1.y * q2.z - q1.z * q2.y;
    result.y = q1.w * q2.y - q1.x * q2.z + q1.y * q2.w + q1.z * q2.x;
    result.z = q1.w * q2.z + q1.x * q2.y - q1.y * q2.x + q1.z * q2.w;
    return result;
}

cb_vector3 cb_physics_quaternion_rotate_vector(cb_quaternion q, cb_vector3 v) {
    // Simplified quaternion-vector rotation
    cb_vector3 result = v;
    // In a real implementation, this would be proper quaternion rotation
    return result;
}
