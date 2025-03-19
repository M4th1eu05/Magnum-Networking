//
// Created by Tarook on 18/03/2025.
//

#include "GameObject.h"

#include <iostream>


void GameObject::start() {
    for (std::shared_ptr<BaseComponent>& component : components) {
        component->start();
    }
}

void GameObject::update()
{

    for (std::shared_ptr<BaseComponent>& component : components) {
        component->update();
    }
}

void GameObject::stop() {
    for (std::shared_ptr<BaseComponent>& component : components) {
        component->stop();
    }
}

void GameObject::addComponentInternal(const std::shared_ptr<BaseComponent> &component) {
    components.push_back(component);
}

void GameObject::removeComponent(const std::shared_ptr<BaseComponent> &component) {
    component->gameObject = nullptr;
    std::erase(components, component);
}
