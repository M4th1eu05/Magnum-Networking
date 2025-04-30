// Client.h
#ifndef CLIENT_H
#define CLIENT_H

#include <enet6/enet.h>
#include <iostream>
#include <string>
#include <thread>
#include <atomic>
#include <memory>

#include "MessageType.h"
#include "World.h"

class Client {
public:
    // pass world
    Client();

    ~Client();

    void connect(const std::string &host, uint16_t port);
    void disconnect();
    void sendMessage(const std::string& message);

    void updateWorldFromServer(const std::string &serializedData) const;

    void sendInput(MessageType messageType, int keyCode);

    void setWorld(const std::shared_ptr<World> &world) {
        _world = world;
    }

private:
    void clientLoop();

    ENetAddress _address;
    ENetHost* _client = nullptr;
    ENetPeer* _peer = nullptr;
    std::atomic<bool> _connected{false};
    std::thread _clientThread;

    std::shared_ptr<World> _world;
};

#endif // CLIENT_H