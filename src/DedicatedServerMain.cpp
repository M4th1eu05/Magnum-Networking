//
// Created by Tarook on 29/04/2025.
//
#include "../inc/DedicatedServer.h"
#include <iostream>

int main() {
    try {
        const uint16_t port = 1234; // Example port
        const size_t maxClients = 32; // Example max clients

        DedicatedServer server(port, maxClients);
        server.start();

        std::cout << "Press Enter to stop the server..." << std::endl;
        std::cin.get();

        server.stop();
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}