// Client.cpp
#include "Client.h"

#include <sstream>

#include "MessageType.h"

Client::Client(){
    if (enet_initialize() != 0) {
        throw std::runtime_error("Failed to initialize ENet.");
    }

    // Create the ENet client host without resolving any address yet
    _client = enet_host_create(ENET_ADDRESS_TYPE_ANY, nullptr, 1, 2, 0, 0);
    if (!_client) {
        throw std::runtime_error("Failed to create ENet client host.");
    }
}

Client::~Client() {
    disconnect();
    if (_client) {
        enet_host_destroy(_client);
    }
    enet_deinitialize();
}

void Client::connect(const std::string& host, const uint16_t port) {
    if (_connected) return;

    // Resolve the address and set the port
    enet_address_set_host(&_address, ENET_ADDRESS_TYPE_ANY, host.c_str());
    _address.port = port;

    // Create a peer and initiate the connection
    _peer = enet_host_connect(_client, &_address, 2, 0);
    if (!_peer) {
        throw std::runtime_error("No available peers for initiating an ENet connection.");
    }

    // Wait for the connection to succeed
    ENetEvent event;
    if (enet_host_service(_client, &event, 5000) > 0 && event.type == ENET_EVENT_TYPE_CONNECT) {
        _connected = true;
        _clientThread = std::thread(&Client::clientLoop, this);
        std::cout << "Connected to " << host << ":" << port << std::endl;
    } else {
        enet_peer_reset(_peer);
        _peer = nullptr;
        throw std::runtime_error("Connection to " + host + ":" + std::to_string(port) + " failed.");
    }
}

void Client::disconnect() {
    if (!_connected) return;

    _connected = false;
    if (_clientThread.joinable()) {
        _clientThread.join();
    }

    if (_peer) {
        enet_peer_disconnect(_peer, 0);
        _peer = nullptr;
    }
    std::cout << "Disconnected from server." << std::endl;
}

void Client::sendMessage(const std::string& message) {
    if (!_connected || !_peer) return;

    ENetPacket* packet = enet_packet_create(message.c_str(), message.size() + 1, ENET_PACKET_FLAG_RELIABLE);
    enet_peer_send(_peer, 0, packet);
}

void Client::clientLoop() {
    while (_connected) {
        ENetEvent event;
        while (enet_host_service(_client, &event, 1000) > 0) {
            switch (event.type) {
                case ENET_EVENT_TYPE_RECEIVE: {
                    std::string message(reinterpret_cast<char*>(event.packet->data), event.packet->dataLength);
                    size_t delimiter = message.find(':');
                    if (delimiter != std::string::npos) {
                        int messageTypeInt = std::stoi(message.substr(0, delimiter));
                        MessageType messageType = static_cast<MessageType>(messageTypeInt);
                        std::string data = message.substr(delimiter + 1);

                        switch (messageType) {
                            case MessageType::CREATE_GAMEOBJECT:
                                if (_world) {
                                    std::istringstream iss(data);
                                    auto newObject = _world->createGameObject();
                                    newObject->deserialize(iss);
                                }
                                break;

                            case MessageType::UPDATE_WORLD:
                                if (_world) {
                                    updateWorldFromServer(data);
                                }
                            default:
                                break;
                        }
                    }
                    enet_packet_destroy(event.packet);
                    break;
                }
                default:
                    break;
            }
        }
    }
}

void Client::updateWorldFromServer(const std::string& serializedData) const {
    if (!_world) return;

    std::istringstream iss(serializedData);
    try {
        _world->deserialize(iss);
    }
    catch (const std::exception& e) {
        std::cerr << "Failed to update world from server: " << e.what() << std::endl;
    }
}

void Client::sendInput(MessageType messageType, const int keyCode) {
    if (!_connected || !_peer) return;

    // Serialize the message type and key code into a binary stream
    std::ostringstream oss(std::ios::binary);
    int messageTypeInt = static_cast<int>(messageType);
    oss.write(reinterpret_cast<const char*>(&messageTypeInt), sizeof(messageTypeInt));
    oss.write(reinterpret_cast<const char*>(&keyCode), sizeof(keyCode));

    // Create and send the packet
    const std::string& binaryData = oss.str();
    ENetPacket* packet = enet_packet_create(binaryData.data(), binaryData.size(), ENET_PACKET_FLAG_RELIABLE);
    enet_peer_send(_peer, 0, packet);
}