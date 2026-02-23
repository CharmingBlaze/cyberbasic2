# Makefile for CyberBasic Engine
# Requires Raylib and Bullet Physics libraries

# Compiler settings
CC = gcc
CFLAGS = -Wall -Wextra -std=c99 -O2 -fPIC
LDFLAGS = -shared

# Library paths and includes (adjust these paths based on your installation)
RAYLIB_INCLUDE = -I/usr/local/include
RAYLIB_LIB = -L/usr/local/lib -lraylib

BULLET_INCLUDE = -I/usr/local/include/bullet
BULLET_LIB = -L/usr/local/lib -lBulletDynamics -lBulletCollision -lLinearMath

BOX2D_INCLUDE = -I/usr/local/include
BOX2D_LIB = -L/usr/local/lib -lBox2D

# Platform-specific settings
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    LDFLAGS += -lm -lpthread -ldl -lrt -lX11
endif
ifeq ($(UNAME_S),Darwin)
    LDFLAGS += -framework OpenGL -framework Cocoa -framework IOKit -framework CoreAudio -framework CoreVideo
endif

# Source files
SOURCES = raylib_wrapper.c bullet_wrapper.c box2d_wrapper.c engine_core.c
OBJECTS = $(SOURCES:.c=.o)
HEADERS = raylib_wrapper.h bullet_wrapper.h box2d_wrapper.h

# Target library
TARGET_LIB = libcyberbasic_engine.a
TARGET_SO = libcyberbasic_engine.so

# Default target
all: $(TARGET_LIB) $(TARGET_SO)

# Static library
$(TARGET_LIB): $(OBJECTS)
	ar rcs $@ $^

# Shared library
$(TARGET_SO): $(OBJECTS)
	$(CC) $(LDFLAGS) -o $@ $^ $(RAYLIB_LIB) $(BULLET_LIB) $(BOX2D_LIB)

# Object files
%.o: %.c $(HEADERS)
	$(CC) $(CFLAGS) $(RAYLIB_INCLUDE) $(BULLET_INCLUDE) $(BOX2D_INCLUDE) -c $< -o $@

# Clean
clean:
	rm -f $(OBJECTS) $(TARGET_LIB) $(TARGET_SO)

# Install (system-wide - requires sudo)
install: $(TARGET_LIB) $(TARGET_SO)
	cp $(TARGET_LIB) /usr/local/lib/
	cp $(TARGET_SO) /usr/local/lib/
	cp $(HEADERS) /usr/local/include/
	ldconfig

# Uninstall (system-wide - requires sudo)
uninstall:
	rm -f /usr/local/lib/$(TARGET_LIB)
	rm -f /usr/local/lib/$(TARGET_SO)
	rm -f /usr/local/include/raylib_wrapper.h
	rm -f /usr/local/include/bullet_wrapper.h
	rm -f /usr/local/include/box2d_wrapper.h
	ldconfig

# Test build (simple test program)
test: test_engine.c
	$(CC) $(CFLAGS) $(RAYLIB_INCLUDE) $(BULLET_INCLUDE) $(BOX2D_INCLUDE) -o test_engine test_engine.c $(TARGET_LIB) $(RAYLIB_LIB) $(BULLET_LIB) $(BOX2D_LIB)

.PHONY: all clean install uninstall test
