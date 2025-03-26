#include <enet.h>
#include <iostream>

const char* SERVER_IP = "127.0.0.1"; // Remplacez par l'IP du serveur
const int SERVER_PORT = 1234;

int main() {
    if (enet_initialize() != 0) {
        std::cerr << "Erreur lors de l'initialisation d'ENet.\n";
        return EXIT_FAILURE;
    }

    ENetHost* client = enet_host_create(NULL, 1, 2, 0, 0);
    if (!client) {
        std::cerr << "Erreur lors de la création du client.\n";
        return EXIT_FAILURE;
    }

    ENetAddress address;
    ENetPeer* peer;
    enet_address_set_host(&address, SERVER_IP);
    address.port = SERVER_PORT;

    peer = enet_host_connect(client, &address, 2, 0);
    if (!peer) {
        std::cerr << "Échec de connexion au serveur.\n";
        return EXIT_FAILURE;
    }

    std::cout << "Connexion au serveur en cours...\n";

    while (true) {
        ENetEvent event;
        while (enet_host_service(client, &event, 1000) > 0) {
            switch (event.type) {
                case ENET_EVENT_TYPE_CONNECT:
                    std::cout << "Connecté au serveur !\n";
                break;
                case ENET_EVENT_TYPE_RECEIVE:
                    std::cout << "Message reçu : " << event.packet->data << "\n";
                enet_packet_destroy(event.packet);
                break;
                case ENET_EVENT_TYPE_DISCONNECT:
                    std::cout << "Déconnecté du serveur.\n";
                return 0;
            }
        }
    }

    enet_host_destroy(client);
    enet_deinitialize();
    return EXIT_SUCCESS;

    //interpolation
    void updatePosition(const Vector3& newPos) {
        playerPos = playerPos.lerp(newPos, 0.1f); // 0.1f est le facteur de lissage
    }

}
