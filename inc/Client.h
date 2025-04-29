// Client.h
#ifndef CLIENT_H
#define CLIENT_H

#include <enet6/enet.h>
#include <iostream>
#include <string>
#include <thread>
#include <atomic>

class Client {
public:
    Client();

    ~Client();

    void connect(const std::string &host, uint16_t port);
    void disconnect();
    void sendMessage(const std::string& message);

private:
    void clientLoop();

    ENetAddress _address;
    ENetHost* _client = nullptr;
    ENetPeer* _peer = nullptr;
    std::atomic<bool> _connected{false};
    std::thread _clientThread;
};

#endif // CLIENT_H