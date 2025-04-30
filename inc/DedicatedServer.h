// DedicatedServer.h
#ifndef DEDICATEDSERVER_H
#define DEDICATEDSERVER_H

#include <enet6/enet.h>
#include <iostream>
#include <thread>
#include <atomic>

#include "MessageType.h"
#include "World.h"

class DedicatedServer {
public:
    DedicatedServer(uint16_t port, size_t maxClients);
    ~DedicatedServer();

    void start();
    void stop();

    void broadcastMessage(MessageType messageType, const std::string& data) const;
    void broadcastWorld(const World &world) const;

    void processInput(const char *data, size_t dataLength);

private:
    void serverLoop();

    ENetAddress _address;
    ENetHost* _server = nullptr;
    std::atomic<bool> _running{false};
    std::thread _serverThread;
};

#endif // DEDICATEDSERVER_H