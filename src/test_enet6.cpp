#include <enet6/enet.h>
#include <iostream>

int main() {
    if (enet_initialize() != 0) {
        std::cerr << "Erreur : Impossible d'initialiser Enet6 !" << std::endl;
        return EXIT_FAILURE;
    }

    std::cout << "Enet6 est bien initialisé !" << std::endl;
    enet_deinitialize();
    return EXIT_SUCCESS;
}