

#include "Client.h"


Client::Client() : client(nullptr), peer(nullptr)
{

}

bool Client::initialize() {
    if (enet_initialize() != 0) {
        std::cerr << "An error occurred while initializing ENet." << std::endl;
        return false;
    }
    return true;
}

bool Client::connectToServer(const std::string& host, enet_uint16 port) {
    ENetAddress address;
    enet_address_set_host(&address, ENET_ADDRESS_TYPE_ANY, host.c_str());
    address.port = port;

    client = enet_host_create(address.type, nullptr, 1, 2, 0, 0);
    if (!client) {
        std::cerr << "Failed to create ENet client host." << std::endl;
        return false;
    }

    peer = enet_host_connect(client, &address, 2, 0);
    if (!peer) {
        std::cerr << "No available peers for initiating an ENet connection." << std::endl;
        return false;
    }

    ENetEvent event;
    if (enet_host_service(client, &event, 5000) > 0 && event.type == ENET_EVENT_TYPE_CONNECT) {
        std::cout << "Connection to " << host << ":" << port << " succeeded." << std::endl;
        return true;
    } else {
        enet_peer_reset(peer);
        std::cerr << "Connection to " << host << ":" << port << " failed." << std::endl;
        return false;
    }
}

// not sure about that
// void Client::synchronizeTime() {
//     auto start = std::chrono::high_resolution_clock::now();
//
//     // Send a time synchronization request
//     ENetPacket* packet = enet_packet_create("TIME_SYNC", strlen("TIME_SYNC") + 1, ENET_PACKET_FLAG_RELIABLE);
//     enet_peer_send(peer, 0, packet);
//
//     // Wait for the server's response
//     ENetEvent event;
//     if (enet_host_service(client, &event, 5000) > 0 && event.type == ENET_EVENT_TYPE_RECEIVE) {
//         auto end = std::chrono::high_resolution_clock::now();
//         auto rtt = std::chrono::duration_cast<std::chrono::milliseconds>(end - start).count();
//
//         std::string serverTime((char*)event.packet->data);
//         std::cout << "Server time: " << serverTime << ", RTT: " << rtt << " ms" << std::endl;
//
//         enet_packet_destroy(event.packet);
//     } else {
//         std::cerr << "Time synchronization failed." << std::endl;
//     }
// }

// void Client::receiveSnapshot(std::shared_ptr<World> world) {
//     ENetEvent event;
//     while (enet_host_service(client, &event, 1000) > 0) {
//         switch (event.type) {
//             case ENET_EVENT_TYPE_RECEIVE:
//                 std::string snapshot(static_cast<char *>(event.packet->data), event.packet->dataLength);
//                 world->deserialize(snapshot);
//                 enet_packet_destroy(event.packet);
//                 break;
//             case ENET_EVENT_TYPE_DISCONNECT:
//                 std::cout << "Disconnected from server." << std::endl;
//                 return;
//             default:
//                 break;
//         }
//     }
// }

void Client::cleanup() {
    if (client) {
        enet_host_destroy(client);
    }
    enet_deinitialize();
}