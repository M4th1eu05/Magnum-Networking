//
// Created by Tarook on 18/03/2025.
//

#include "World.h"

World::World() {
    _bWorld.setGravity({0.0f, -10.0f, 0.0f});
    _bWorld.setDebugDrawer(&_debugDraw);
}

std::shared_ptr<GameObject> World::createGameObject(std::shared_ptr<GameObject> parent) {
    std::shared_ptr<GameObject> object;
    if (!parent) {
        object = std::make_shared<GameObject>(&_scene);
    }
    else {
        object = std::make_shared<GameObject>(parent.get());
    }

    addObject(object);
    return object;
}

void World::addObject(const std::shared_ptr<GameObject> &object) {
    _objects.push_back(object);
}

void World::removeObject(std::shared_ptr<GameObject> object) {
    std::erase(_objects, object);
}

void World::update() {
    for (std::shared_ptr<GameObject>& object : _objects) {
        object->update();
    }

    _bWorld.stepSimulation(physicsTimeStep);
}