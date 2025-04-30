//
// Created by Tarook on 18/03/2025.
//

#include "World.h"

#include <iostream>

#include "DedicatedServer.h"
#include "MessageType.h"
#include "Magnum/Timeline.h"

World::World(Timeline& timeline)
    : _timeline(timeline){
    _bWorld.setGravity({0.0f, -10.0f, 0.0f});
    _bWorld.setDebugDrawer(&_debugDraw);
}

std::shared_ptr<GameObject> World::createGameObject(const std::shared_ptr<GameObject> &parent) {
    std::shared_ptr<GameObject> object;
    if (!parent) {
        object = std::make_shared<GameObject>(&_scene, shared_from_this());
    }
    else {
        object = std::make_shared<GameObject>(parent.get(), shared_from_this());
    }

    addObject(object);

    // Notify the server about the new GameObject
    if (_server) {
        std::ostringstream oss;
        object->serialize(oss);
        const std::string serializedData = oss.str();
        _server->broadcastMessage(MessageType::CREATE_GAMEOBJECT, serializedData);
    }

    return object;
}

void World::addObject(const std::shared_ptr<GameObject> &object) {
    _objects.push_back(object);
}

void World::removeObject(std::shared_ptr<GameObject> object) {
    std::erase(_objects, object);
}

void World::update() {
    //std::cout << "World update called" << std::endl;
    for (std::shared_ptr<GameObject>& object : _objects) {
        object->update();
    }

    _bWorld.stepSimulation(_timeline.previousFrameDuration(), 5);
}