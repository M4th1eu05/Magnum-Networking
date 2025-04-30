#include "../inc/DedicatedServer.h"

#include <sstream>

#include "World.h"

DedicatedServer::DedicatedServer(const uint16_t port, const size_t maxClients) {
    if (enet_initialize() != 0) {
        throw std::runtime_error("Failed to initialize ENet6.");
    }

    // Configure the address for dual-stack (IPv4 and IPv6)
    enet_address_build_any(&_address, ENET_ADDRESS_TYPE_IPV6);
    _address.port = port;

    // Create the server host
    _server = enet_host_create(ENET_ADDRESS_TYPE_ANY, &_address, maxClients, 2, 0, 0);
    if (!_server) {
        throw std::runtime_error("Failed to create ENet6 server host.");
    }

    std::cout << "Server initialized on port " << port << std::endl;
    // log address
    char addressString[65];
    enet_address_get_host(&_address, addressString, sizeof(addressString));
    std::cout << "Server address: " << addressString << ":" << port << std::endl;
}

DedicatedServer::~DedicatedServer() {
    stop();
    if (_server) {
        enet_host_destroy(_server);
    }
    enet_deinitialize();
}

void DedicatedServer::start() {
    if (_running) return;

    _running = true;
    _serverThread = std::thread(&DedicatedServer::serverLoop, this);
    std::cout << "Server started." << std::endl;
}

void DedicatedServer::stop() {
    if (!_running) return;

    _running = false;
    if (_serverThread.joinable()) {
        _serverThread.join();
    }
    std::cout << "Server stopped." << std::endl;
}

void DedicatedServer::serverLoop() {
    while (_running) {
        ENetEvent event;
        while (enet_host_service(_server, &event, 1000) > 0) {
            std::string message;
            size_t delimiter;

            switch (event.type) {
                case ENET_EVENT_TYPE_CONNECT:
                    std::cout << "Client connected." << std::endl;
                    break;

                case ENET_EVENT_TYPE_RECEIVE:
                    std::cout << "Message received: "
                              << reinterpret_cast<char*>(event.packet->data) << std::endl;

                    message = std::string(reinterpret_cast<char*>(event.packet->data), event.packet->dataLength);
                    delimiter = message.find(':');
                    if (delimiter != std::string::npos) {
                        int messageTypeInt = std::stoi(message.substr(0, delimiter));
                        MessageType messageType = static_cast<MessageType>(messageTypeInt);
                        std::string data = message.substr(delimiter + 1);

                        if (messageType == MessageType::CLIENT_INPUT) {
                            std::cout << "Received input: " << data << std::endl;
                            // Add logic to handle the input
                        }
                    }
                    enet_packet_destroy(event.packet);
                    break;

                case ENET_EVENT_TYPE_DISCONNECT:
                    std::cout << "Client disconnected." << std::endl;
                    break;

                case ENET_EVENT_TYPE_DISCONNECT_TIMEOUT:
                    std::cout << "Client disconnected due to timeout." << std::endl;
                    break;

                default:
                    break;
            }
        }
    }
}

void DedicatedServer::broadcastMessage(MessageType messageType, const std::string& data) const {
    if (!_server) return;

    // Serialize the message type as an integer
    int messageTypeInt = static_cast<int>(messageType);
    const std::string message = std::to_string(messageTypeInt) + ":" + data;

    // Create an ENet packet
    ENetPacket* packet = enet_packet_create(message.data(), message.size(), ENET_PACKET_FLAG_RELIABLE);

    // Broadcast the packet to all connected peers
    enet_host_broadcast(_server, 0, packet);
    enet_host_flush(_server);
}

void DedicatedServer::broadcastWorld(const World& world) const {
    if (!_server) return;

    // Serialize the world
    std::ostringstream oss;
    world.serialize(oss);
    const std::string serializedData = oss.str();

    // Serialize the message type as an integer
    int messageTypeInt = static_cast<int>(MessageType::UPDATE_WORLD);
    const std::string message = std::to_string(messageTypeInt) + ":" + serializedData;

    // Create an ENet packet
    ENetPacket* packet = enet_packet_create(message.data(), message.size(), ENET_PACKET_FLAG_RELIABLE);

    // Broadcast the packet to all connected peers
    enet_host_broadcast(_server, 0, packet);
    enet_host_flush(_server);
}

void DedicatedServer::processInput(const char* data, const size_t dataLength) {
    std::istringstream iss(std::string(data, dataLength), std::ios::binary);

    // Deserialize the message type and key code
    int messageTypeInt;
    int keyCode;
    iss.read(reinterpret_cast<char*>(&messageTypeInt), sizeof(messageTypeInt));
    iss.read(reinterpret_cast<char*>(&keyCode), sizeof(keyCode));

    MessageType messageType = static_cast<MessageType>(messageTypeInt);
    if (messageType == MessageType::CLIENT_INPUT) {
        std::cout << "Received input: KeyCode=" << keyCode << std::endl;
        // Handle the input
    }
}