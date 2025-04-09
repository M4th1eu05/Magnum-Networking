#include <enet6/enet.h>
#include <iostream>
#include "World2.h"
#include "WorldSerializer.h"


int main() {
    if (enet_initialize() != 0) {
        std::cerr << "Erreur : Impossible d'initialiser Enet6 !" << std::endl;
        return EXIT_FAILURE;
    }

    std::cout << "Enet6 est bien initialisé !" << std::endl;
    enet_deinitialize();
    return EXIT_SUCCESS;

    World2 world;

    // Setup de base pour test
    world.platform = {0.0f, 0.0f, 0.0f, 20.0f, 1.0f, 20.0f};
    for (int i = 0; i < 10; ++i)
        for (int j = 0; j < 10; ++j)
            world.cubes[i][j] = {i * 1.0f, 1.0f, j * 1.0f, 1}; // type 1 = par défaut

    if (saveWorldToFile(world, "world.dat"))
        std::cout << "World saved!\n";

    World2 loaded;
    if (loadWorldFromFile(loaded, "world.dat"))
        std::cout << "World loaded!\n";

    return 0;
}