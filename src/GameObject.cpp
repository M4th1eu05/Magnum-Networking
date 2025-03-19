//
// Created by Tarook on 18/03/2025.
//

#include "GameObject.h"

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
    if (!component->gameObject) {
        components.push_back(component);
        component->gameObject = shared_from_this();
    }
}

void GameObject::removeComponent(const std::shared_ptr<BaseComponent> &component) {
    component->gameObject = nullptr;
    std::erase(components, component);
}
