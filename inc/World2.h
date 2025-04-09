#ifndef WORLD2_H
#define WORLD2_H
#pragma once

#include <cstdint>

struct Cube {
    float x, y, z;
    uint8_t type;
};

struct Platform {
    float x, y, z;
    float width, height, depth;
};

struct World2 {
    Platform platform;
    Cube cubes[10][10];
};

#endif //WORLD2_H
