#include "WorldSerializer.h"
#include <fstream>

bool saveWorldToFile(const World2& world, const std::string& filename) {
    std::ofstream out(filename, std::ios::binary);
    if (!out) return false;

    // Save Platform
    out.write(reinterpret_cast<const char*>(&world.platform), sizeof(Platform));

    // Save cubes
    for (int i = 0; i < 10; ++i)
        out.write(reinterpret_cast<const char*>(&world.cubes[i]), sizeof(Cube) * 10);

    return true;
}

bool loadWorldFromFile(World2& world, const std::string& filename) {
    std::ifstream in(filename, std::ios::binary);
    if (!in) return false;

    // Load platform
    in.read(reinterpret_cast<char*>(&world.platform), sizeof(Platform));

    // Load cubes
    for (int i = 0; i < 10; ++i)
        in.read(reinterpret_cast<char*>(&world.cubes[i]), sizeof(Cube) * 10);

    return true;
}

