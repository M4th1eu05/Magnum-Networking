//
// Created by mathi on 29/04/2025.
//

#ifndef CLIENT_H
#define CLIENT_H
#include <string>

#include "../cmake-build-debug/_deps/enet6-src/include/enet6/types.h"

class Client {
public:
    Client();
    ~Client();

    bool initialize();
    bool connectToServer(const std::string& host, enet_uint16 port);
    void cleanup();

private:
    ENetHost* client;
    ENetPeer* peer;
};

#endif //CLIENT_H
