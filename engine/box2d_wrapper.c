#include "box2d_wrapper.h"
#include <stdlib.h>
#include <string.h>
#include <math.h>

// Simplified 2D physics implementation for demonstration
// In a real implementation, this would interface with Box2D

struct cb_physics_body_2d {
    cb_body_2d_type type;
    cb_shape_2d_type shape;
    cb_vector2 size;
    float density;
    cb_transform2d transform;
    cb_vector2 linear_velocity;
    float angular_velocity;
    float friction;
    float restitution;
    float gravity_scale;
    bool active;
    cb_filter_category category;
    cb_filter_category mask;
    int group_index;
};

struct cb_physics_world_2d {
    cb_vector2 gravity;
    cb_physics_body_2d_t* bodies;
    int body_count;
    int body_capacity;
    bool allow_sleeping;
    cb_debug_draw_2d* debug_draw;
};

struct cb_constraint_2d {
    int type; // 0=revolute, 1=prismatic, 2=distance, 3=pulley, 4=mouse, 5=gear, 6=wheel, 7=weld, 8=friction
    cb_physics_body_2d_t body_a;
    cb_physics_body_2d_t body_b;
    cb_vector2 anchor_a;
    cb_vector2 anchor_b;
    cb_vector2 axis;
    float length;
    float ratio;
    float max_force;
    float max_torque;
    float motor_speed;
    float lower_limit;
    float upper_limit;
    float frequency;
    float damping;
    cb_constraint_2d_t joint_a;
    cb_constraint_2d_t joint_b;
};

// World management
cb_physics_world_2d_t cb_physics_2d_create_world(cb_vector2 gravity) {
    cb_physics_world_2d_t world = malloc(sizeof(struct cb_physics_world_2d));
    world->gravity = gravity;
    world->bodies = malloc(sizeof(cb_physics_body_2d_t) * 100);
    world->body_count = 0;
    world->body_capacity = 100;
    world->allow_sleeping = true;
    world->debug_draw = NULL;
    return world;
}

void cb_physics_2d_destroy_world(cb_physics_world_2d_t world) {
    if (world) {
        for (int i = 0; i < world->body_count; i++) {
            free(world->bodies[i]);
        }
        free(world->bodies);
        free(world);
    }
}

void cb_physics_2d_step_simulation(cb_physics_world_2d_t world, float time_step, int velocity_iterations, int position_iterations) {
    if (!world) return;
    
    // Simple physics simulation
    for (int i = 0; i < world->body_count; i++) {
        cb_physics_body_2d_t body = world->bodies[i];
        if (!body || !body->active || body->type == CB_BODY_2D_STATIC) continue;
        
        // Apply gravity
        cb_vector2 gravity_force = {
            world->gravity.x * body->density * body->gravity_scale,
            world->gravity.y * body->density * body->gravity_scale
        };
        
        body->linear_velocity.x += gravity_force.x * time_step;
        body->linear_velocity.y += gravity_force.y * time_step;
        
        // Update position
        body->transform.position.x += body->linear_velocity.x * time_step;
        body->transform.position.y += body->linear_velocity.y * time_step;
        body->transform.angle += body->angular_velocity * time_step;
        
        // Simple ground collision
        if (body->shape == CB_SHAPE_2D_BOX) {
            float half_height = body->size.y * 0.5f;
            if (body->transform.position.y < half_height) {
                body->transform.position.y = half_height;
                body->linear_velocity.y *= -body->restitution;
                body->linear_velocity.x *= (1.0f - body->friction);
                body->angular_velocity *= (1.0f - body->friction);
            }
        } else if (body->shape == CB_SHAPE_2D_CIRCLE) {
            float radius = body->size.x * 0.5f;
            if (body->transform.position.y < radius) {
                body->transform.position.y = radius;
                body->linear_velocity.y *= -body->restitution;
                body->linear_velocity.x *= (1.0f - body->friction);
                body->angular_velocity *= (1.0f - body->friction);
            }
        }
    }
}

// Body management
cb_physics_body_2d_t cb_physics_2d_create_body(cb_physics_world_2d_t world, cb_body_2d_type type, cb_shape_2d_type shape, cb_vector2 size, float density) {
    if (!world || world->body_count >= world->body_capacity) return NULL;
    
    cb_physics_body_2d_t body = malloc(sizeof(struct cb_physics_body_2d));
    body->type = type;
    body->shape = shape;
    body->size = size;
    body->density = density;
    body->transform.position = (cb_vector2){ 0, 0 };
    body->transform.angle = 0.0f;
    body->linear_velocity = (cb_vector2){ 0, 0 };
    body->angular_velocity = 0.0f;
    body->friction = 0.5f;
    body->restitution = 0.1f;
    body->gravity_scale = 1.0f;
    body->active = true;
    body->category = CB_FILTER_2D_CATEGORY_1;
    body->mask = CB_FILTER_2D_ALL_MASK;
    body->group_index = 0;
    
    world->bodies[world->body_count++] = body;
    return body;
}

void cb_physics_2d_destroy_body(cb_physics_world_2d_t world, cb_physics_body_2d_t body) {
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

void cb_physics_2d_set_transform(cb_physics_body_2d_t body, cb_transform2d transform) {
    if (body) {
        body->transform = transform;
    }
}

cb_transform2d cb_physics_2d_get_transform(cb_physics_body_2d_t body) {
    if (body) {
        return body->transform;
    }
    return (cb_transform2d){ {0, 0}, 0.0f };
}

void cb_physics_2d_set_position(cb_physics_body_2d_t body, cb_vector2 position) {
    if (body) {
        body->transform.position = position;
    }
}

cb_vector2 cb_physics_2d_get_position(cb_physics_body_2d_t body) {
    if (body) {
        return body->transform.position;
    }
    return (cb_vector2){ 0, 0 };
}

void cb_physics_2d_set_angle(cb_physics_body_2d_t body, float angle) {
    if (body) {
        body->transform.angle = angle;
    }
}

float cb_physics_2d_get_angle(cb_physics_body_2d_t body) {
    if (body) {
        return body->transform.angle;
    }
    return 0.0f;
}

// Velocity and forces
void cb_physics_2d_set_linear_velocity(cb_physics_body_2d_t body, cb_vector2 velocity) {
    if (body) {
        body->linear_velocity = velocity;
    }
}

cb_vector2 cb_physics_2d_get_linear_velocity(cb_physics_body_2d_t body) {
    if (body) {
        return body->linear_velocity;
    }
    return (cb_vector2){ 0, 0 };
}

void cb_physics_2d_set_angular_velocity(cb_physics_body_2d_t body, float velocity) {
    if (body) {
        body->angular_velocity = velocity;
    }
}

float cb_physics_2d_get_angular_velocity(cb_physics_body_2d_t body) {
    if (body) {
        return body->angular_velocity;
    }
    return 0.0f;
}

void cb_physics_2d_apply_force(cb_physics_body_2d_t body, cb_vector2 force, cb_vector2 point) {
    // Simplified - just apply as central force
    cb_physics_2d_apply_force_to_center(body, force);
}

void cb_physics_2d_apply_force_to_center(cb_physics_body_2d_t body, cb_vector2 force) {
    if (body && body->type == CB_BODY_2D_DYNAMIC) {
        // F = ma, so a = F/m
        float mass = body->density;
        if (mass > 0.0f) {
            body->linear_velocity.x += (force.x / mass) * 0.016f; // Assuming 60 FPS
            body->linear_velocity.y += (force.y / mass) * 0.016f;
        }
    }
}

void cb_physics_2d_apply_linear_impulse(cb_physics_body_2d_t body, cb_vector2 impulse, cb_vector2 point) {
    // Simplified - just apply as central impulse
    cb_physics_2d_apply_linear_impulse_to_center(body, impulse);
}

void cb_physics_2d_apply_linear_impulse_to_center(cb_physics_body_2d_t body, cb_vector2 impulse) {
    if (body && body->type == CB_BODY_2D_DYNAMIC) {
        float mass = body->density;
        if (mass > 0.0f) {
            // Impulse changes velocity directly: Δv = J/m
            body->linear_velocity.x += impulse.x / mass;
            body->linear_velocity.y += impulse.y / mass;
        }
    }
}

void cb_physics_2d_apply_torque(cb_physics_body_2d_t body, float torque) {
    if (body && body->type == CB_BODY_2D_DYNAMIC) {
        // Simplified torque application
        body->angular_velocity += torque * 0.016f;
    }
}

void cb_physics_2d_apply_angular_impulse(cb_physics_body_2d_t body, float impulse) {
    if (body && body->type == CB_BODY_2D_DYNAMIC) {
        body->angular_velocity += impulse;
    }
}

// Mass and properties
void cb_physics_2d_set_density(cb_physics_body_2d_t body, float density) {
    if (body) {
        body->density = density;
    }
}

float cb_physics_2d_get_density(cb_physics_body_2d_t body) {
    if (body) {
        return body->density;
    }
    return 0.0f;
}

void cb_physics_2d_set_friction(cb_physics_body_2d_t body, float friction) {
    if (body) {
        body->friction = friction;
    }
}

float cb_physics_2d_get_friction(cb_physics_body_2d_t body) {
    if (body) {
        return body->friction;
    }
    return 0.0f;
}

void cb_physics_2d_set_restitution(cb_physics_body_2d_t body, float restitution) {
    if (body) {
        body->restitution = restitution;
    }
}

float cb_physics_2d_get_restitution(cb_physics_body_2d_t body) {
    if (body) {
        return body->restitution;
    }
    return 0.0f;
}

void cb_physics_2d_set_gravity_scale(cb_physics_body_2d_t body, float scale) {
    if (body) {
        body->gravity_scale = scale;
    }
}

float cb_physics_2d_get_gravity_scale(cb_physics_body_2d_t body) {
    if (body) {
        return body->gravity_scale;
    }
    return 1.0f;
}

// Constraints (simplified implementation)
cb_constraint_2d_t cb_physics_2d_create_revolute_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 0;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->anchor_a = anchor;
    joint->anchor_b = anchor;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_prismatic_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, cb_vector2 axis) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 1;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->anchor_a = anchor;
    joint->anchor_b = anchor;
    joint->axis = axis;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_distance_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor_a, cb_vector2 anchor_b, float length) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 2;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->anchor_a = anchor_a;
    joint->anchor_b = anchor_b;
    joint->length = length;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_pulley_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 ground_anchor_a, cb_vector2 ground_anchor_b, cb_vector2 anchor_a, cb_vector2 anchor_b, float ratio) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 3;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->anchor_a = anchor_a;
    joint->anchor_b = anchor_b;
    joint->ratio = ratio;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_mouse_joint(cb_physics_body_2d_t body, cb_vector2 target) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 4;
    joint->body_a = body;
    joint->body_b = NULL;
    joint->anchor_a = target;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_gear_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_constraint_2d_t joint_a, cb_constraint_2d_t joint_b, float ratio) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 5;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->joint_a = joint_a;
    joint->joint_b = joint_b;
    joint->ratio = ratio;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_wheel_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, cb_vector2 axis, float damping) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 6;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->anchor_a = anchor;
    joint->axis = axis;
    joint->damping = damping;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_weld_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, float angle) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 7;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->anchor_a = anchor;
    joint->anchor_b = anchor;
    return joint;
}

cb_constraint_2d_t cb_physics_2d_create_friction_joint(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b, cb_vector2 anchor, float max_force, float max_torque) {
    cb_constraint_2d_t joint = malloc(sizeof(struct cb_constraint_2d));
    joint->type = 8;
    joint->body_a = body_a;
    joint->body_b = body_b;
    joint->anchor_a = anchor;
    joint->anchor_b = anchor;
    joint->max_force = max_force;
    joint->max_torque = max_torque;
    return joint;
}

void cb_physics_2d_destroy_joint(cb_constraint_2d_t joint) {
    if (joint) {
        free(joint);
    }
}

// Joint properties
void cb_physics_2d_joint_set_motor_speed(cb_constraint_2d_t joint, float speed) {
    if (joint) {
        joint->motor_speed = speed;
    }
}

void cb_physics_2d_joint_set_max_motor_force(cb_constraint_2d_t joint, float force) {
    if (joint) {
        joint->max_force = force;
    }
}

void cb_physics_2d_joint_set_limits(cb_constraint_2d_t joint, float lower, float upper) {
    if (joint) {
        joint->lower_limit = lower;
        joint->upper_limit = upper;
    }
}

void cb_physics_2d_joint_set_frequency(cb_constraint_2d_t joint, float hz) {
    if (joint) {
        joint->frequency = hz;
    }
}

void cb_physics_2d_joint_set_damping(cb_constraint_2d_t joint, float damping) {
    if (joint) {
        joint->damping = damping;
    }
}

// 2D Raycasting
cb_ray_cast_2d_result cb_physics_2d_ray_cast(cb_physics_world_2d_t world, cb_vector2 start, cb_vector2 end) {
    cb_ray_cast_2d_result result = { false, {0, 0}, {0, 0}, 0.0f, NULL };
    
    if (!world) return result;
    
    // Simple raycast implementation
    cb_vector2 direction = {
        end.x - start.x,
        end.y - start.y
    };
    float distance = sqrtf(direction.x * direction.x + direction.y * direction.y);
    
    if (distance > 0.0f) {
        direction.x /= distance;
        direction.y /= distance;
    }
    
    // Check intersection with bodies
    for (int i = 0; i < world->body_count; i++) {
        cb_physics_body_2d_t body = world->bodies[i];
        if (!body || !body->active) continue;
        
        // Simple AABB intersection
        cb_vector2 min = {
            body->transform.position.x - body->size.x * 0.5f,
            body->transform.position.y - body->size.y * 0.5f
        };
        cb_vector2 max = {
            body->transform.position.x + body->size.x * 0.5f,
            body->transform.position.y + body->size.y * 0.5f
        };
        
        // Check if ray intersects with box (simplified)
        if (start.x >= min.x && start.x <= max.x &&
            start.y >= min.y && start.y <= max.y) {
            result.hit = true;
            result.point = start;
            result.normal = (cb_vector2){ 0, 1 };
            result.fraction = 0.0f;
            result.body = body;
            break;
        }
    }
    
    return result;
}

// 2D Collision detection
cb_collision_2d_result cb_physics_2d_check_collision(cb_physics_body_2d_t body_a, cb_physics_body_2d_t body_b) {
    cb_collision_2d_result result = { false, {{0, 0}, {0, 0}}, {0, 0}, 0, 0.0f };
    
    if (!body_a || !body_b || !body_a->active || !body_b->active) return result;
    
    // Simple AABB collision
    cb_vector2 min_a = {
        body_a->transform.position.x - body_a->size.x * 0.5f,
        body_a->transform.position.y - body_a->size.y * 0.5f
    };
    cb_vector2 max_a = {
        body_a->transform.position.x + body_a->size.x * 0.5f,
        body_a->transform.position.y + body_a->size.y * 0.5f
    };
    
    cb_vector2 min_b = {
        body_b->transform.position.x - body_b->size.x * 0.5f,
        body_b->transform.position.y - body_b->size.y * 0.5f
    };
    cb_vector2 max_b = {
        body_b->transform.position.x + body_b->size.x * 0.5f,
        body_b->transform.position.y + body_b->size.y * 0.5f
    };
    
    bool overlap = (max_a.x >= min_b.x && min_a.x <= max_b.x) &&
                   (max_a.y >= min_b.y && min_a.y <= max_b.y);
    
    if (overlap) {
        result.colliding = true;
        result.contact_points[0] = (cb_vector2){ 
            (max_a.x + min_a.x + max_b.x + min_b.x) * 0.25f,
            (max_a.y + min_a.y + max_b.y + min_b.y) * 0.25f
        };
        result.contact_normal = (cb_vector2){ 0, 1 }; // Simplified
        result.contact_count = 1;
        result.separation = 0.1f; // Simplified
    }
    
    return result;
}

// AABB queries
cb_query_2d_result cb_physics_2d_query_aabb(cb_physics_world_2d_t world, cb_aabb_2d aabb) {
    cb_query_2d_result result = { NULL, 0 };
    
    if (!world) return result;
    
    // Count bodies in AABB
    int count = 0;
    for (int i = 0; i < world->body_count; i++) {
        cb_physics_body_2d_t body = world->bodies[i];
        if (!body || !body->active) continue;
        
        cb_vector2 min = {
            body->transform.position.x - body->size.x * 0.5f,
            body->transform.position.y - body->size.y * 0.5f
        };
        cb_vector2 max = {
            body->transform.position.x + body->size.x * 0.5f,
            body->transform.position.y + body->size.y * 0.5f
        };
        
        if (max.x >= aabb.lower_bound.x && min.x <= aabb.upper_bound.x &&
            max.y >= aabb.lower_bound.y && min.y <= aabb.upper_bound.y) {
            count++;
        }
    }
    
    // Allocate and fill result
    if (count > 0) {
        result.bodies = malloc(sizeof(cb_physics_body_2d_t) * count);
        result.count = count;
        
        int index = 0;
        for (int i = 0; i < world->body_count; i++) {
            cb_physics_body_2d_t body = world->bodies[i];
            if (!body || !body->active) continue;
            
            cb_vector2 min = {
                body->transform.position.x - body->size.x * 0.5f,
                body->transform.position.y - body->size.y * 0.5f
            };
            cb_vector2 max = {
                body->transform.position.x + body->size.x * 0.5f,
                body->transform.position.y + body->size.y * 0.5f
            };
            
            if (max.x >= aabb.lower_bound.x && min.x <= aabb.upper_bound.x &&
                max.y >= aabb.lower_bound.y && min.y <= aabb.upper_bound.y) {
                result.bodies[index++] = body;
            }
        }
    }
    
    return result;
}

cb_physics_body_2d_t cb_physics_2d_query_point(cb_physics_world_2d_t world, cb_vector2 point) {
    if (!world) return NULL;
    
    for (int i = 0; i < world->body_count; i++) {
        cb_physics_body_2d_t body = world->bodies[i];
        if (!body || !body->active) continue;
        
        cb_vector2 min = {
            body->transform.position.x - body->size.x * 0.5f,
            body->transform.position.y - body->size.y * 0.5f
        };
        cb_vector2 max = {
            body->transform.position.x + body->size.x * 0.5f,
            body->transform.position.y + body->size.y * 0.5f
        };
        
        if (point.x >= min.x && point.x <= max.x &&
            point.y >= min.y && point.y <= max.y) {
            return body;
        }
    }
    
    return NULL;
}

// Shape casting
cb_shape_cast_2d_result cb_physics_2d_shape_cast(cb_physics_world_2d_t world, cb_shape_2d_type shape, cb_vector2 shape_size, cb_transform2d transform, cb_vector2 translation) {
    cb_shape_cast_2d_result result = { false, {0, 0}, {0, 0}, 0.0f, NULL };
    
    if (!world) return result;
    
    // Simplified shape casting - just use raycast from center
    cb_vector2 start = transform.position;
    cb_vector2 end = {
        transform.position.x + translation.x,
        transform.position.y + translation.y
    };
    
    cb_ray_cast_2d_result ray_result = cb_physics_2d_ray_cast(world, start, end);
    
    if (ray_result.hit) {
        result.hit = true;
        result.point = ray_result.point;
        result.normal = ray_result.normal;
        result.fraction = ray_result.fraction;
        result.body = ray_result.body;
    }
    
    return result;
}

// World properties
void cb_physics_2d_set_gravity(cb_physics_world_2d_t world, cb_vector2 gravity) {
    if (world) {
        world->gravity = gravity;
    }
}

cb_vector2 cb_physics_2d_get_gravity(cb_physics_world_2d_t world) {
    if (world) {
        return world->gravity;
    }
    return (cb_vector2){ 0, -9.81f };
}

void cb_physics_2d_set_allow_sleeping(cb_physics_world_2d_t world, bool allow) {
    if (world) {
        world->allow_sleeping = allow;
    }
}

bool cb_physics_2d_get_allow_sleeping(cb_physics_world_2d_t world) {
    if (world) {
        return world->allow_sleeping;
    }
    return true;
}

// Body filtering
void cb_physics_2d_set_filter_category(cb_physics_body_2d_t body, cb_filter_category category) {
    if (body) {
        body->category = category;
    }
}

void cb_physics_2d_set_filter_mask(cb_physics_body_2d_t body, cb_filter_category mask) {
    if (body) {
        body->mask = mask;
    }
}

void cb_physics_2d_set_filter_group_index(cb_physics_body_2d_t body, int group_index) {
    if (body) {
        body->group_index = group_index;
    }
}

// Debug drawing
void cb_physics_2d_set_debug_draw(cb_physics_world_2d_t world, cb_debug_draw_2d* debug_draw) {
    if (world) {
        world->debug_draw = debug_draw;
    }
}

void cb_physics_2d_draw_debug_data(cb_physics_world_2d_t world) {
    if (!world || !world->debug_draw) return;
    
    for (int i = 0; i < world->body_count; i++) {
        cb_physics_body_2d_t body = world->bodies[i];
        if (!body || !body->active) continue;
        
        cb_vector2 color = (cb_vector2){ 255, 255 }; // White
        
        if (body->shape == CB_SHAPE_2D_BOX) {
            cb_vector2 vertices[4];
            float half_width = body->size.x * 0.5f;
            float half_height = body->size.y * 0.5f;
            
            vertices[0] = (cb_vector2){ -half_width, -half_height };
            vertices[1] = (cb_vector2){ half_width, -half_height };
            vertices[2] = (cb_vector2){ half_width, half_height };
            vertices[3] = (cb_vector2){ -half_width, half_height };
            
            // Rotate vertices
            float cos_a = cosf(body->transform.angle);
            float sin_a = sinf(body->transform.angle);
            
            for (int j = 0; j < 4; j++) {
                float x = vertices[j].x * cos_a - vertices[j].y * sin_a + body->transform.position.x;
                float y = vertices[j].x * sin_a + vertices[j].y * cos_a + body->transform.position.y;
                vertices[j] = (cb_vector2){ x, y };
            }
            
            if (world->debug_draw->draw_polygon) {
                world->debug_draw->draw_polygon(vertices, 4, color);
            }
        } else if (body->shape == CB_SHAPE_2D_CIRCLE) {
            if (world->debug_draw->draw_circle) {
                world->debug_draw->draw_circle(body->transform.position, body->size.x * 0.5f, color);
            }
        }
    }
}

// Utility functions
cb_vector2 cb_physics_2d_vector_create(float x, float y) {
    return (cb_vector2){ x, y };
}

cb_vector2 cb_physics_2d_vector_add(cb_vector2 a, cb_vector2 b) {
    return (cb_vector2){ a.x + b.x, a.y + b.y };
}

cb_vector2 cb_physics_2d_vector_subtract(cb_vector2 a, cb_vector2 b) {
    return (cb_vector2){ a.x - b.x, a.y - b.y };
}

cb_vector2 cb_physics_2d_vector_multiply(cb_vector2 v, float scalar) {
    return (cb_vector2){ v.x * scalar, v.y * scalar };
}

cb_vector2 cb_physics_2d_vector_normalize(cb_vector2 v) {
    float length = sqrtf(v.x * v.x + v.y * v.y);
    if (length > 0.0f) {
        return (cb_vector2){ v.x / length, v.y / length };
    }
    return (cb_vector2){ 0, 0 };
}

float cb_physics_2d_vector_length(cb_vector2 v) {
    return sqrtf(v.x * v.x + v.y * v.y);
}

float cb_physics_2d_vector_length_squared(cb_vector2 v) {
    return v.x * v.x + v.y * v.y;
}

float cb_physics_2d_vector_dot(cb_vector2 a, cb_vector2 b) {
    return a.x * b.x + a.y * b.y;
}

float cb_physics_2d_vector_cross(cb_vector2 a, cb_vector2 b) {
    return a.x * b.y - a.y * b.x;
}

cb_vector2 cb_physics_2d_vector_cross_float(float s, cb_vector2 a) {
    return (cb_vector2){ -s * a.y, s * a.x };
}

cb_vector2 cb_physics_2d_vector_cross_vector(cb_vector2 a, float s) {
    return (cb_vector2){ s * a.y, -s * a.x };
}
