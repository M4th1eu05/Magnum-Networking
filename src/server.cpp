#include <enet/enet.h>
#include <iostream>

const int MAX_PLAYERS = 4;
const int PORT = 1234;

int main() {
    if (enet_initialize() != 0) {
        std::cerr << "Erreur lors de l'initialisation d'ENet.\n";
        return EXIT_FAILURE;
    }

    ENetAddress address;
    ENetHost* server;

    address.host = ENET_HOST_ANY;
    address.port = PORT;

    server = enet_host_create(&address, MAX_PLAYERS, 2, 0, 0);
    if (!server) {
        std::cerr << "Erreur lors de la création du serveur.\n";
        return EXIT_FAILURE;
    }

    std::cout << "Serveur démarré sur le port " << PORT << ".\n";

    while (true) {
        ENetEvent event;
        while (enet_host_service(server, &event, 1000) > 0) {
            switch (event.type) {
                case ENET_EVENT_TYPE_CONNECT:
                    std::cout << "Un joueur s'est connecté.\n";
                break;
                case ENET_EVENT_TYPE_RECEIVE:
                    std::cout << "Message reçu : " << event.packet->data << "\n";
                enet_packet_destroy(event.packet);
                break;
                case ENET_EVENT_TYPE_DISCONNECT:
                    std::cout << "Un joueur s'est déconnecté.\n";
                break;
            }
        }
    }

    enet_host_destroy(server);
    enet_deinitialize();
    return EXIT_SUCCESS;

    void sendInput(ENetPeer* peer, const std::string& input) {
        ENetPacket* packet = enet_packet_create(input.c_str(),
                                                input.size() + 1,
                                                ENET_PACKET_FLAG_RELIABLE);
        enet_peer_send(peer, 0, packet);
        enet_host_flush(peer->host);
    }

    case ENET_EVENT_TYPE_RECEIVE:
        std::string input(reinterpret_cast<char*>(event.packet->data));
    std::cout << "Input reçu du joueur : " << input << "\n";
    enet_packet_destroy(event.packet);
    break;

    //interpolation
    struct Vector3 {
        float x, y, z;
        Vector3 lerp(const Vector3& target, float alpha) {
            return {x + alpha * (target.x - x),
                    y + alpha * (target.y - y),
                    z + alpha * (target.z - z)};
        }
    };


}
