//
// Created by Tarook on 18/03/2025.
//

#ifndef BASECOMPONENTS_H
#define BASECOMPONENTS_H
#include <memory>

#include "ILifeCycle.h"


class GameObject;

class BaseComponent : public ILifeCycle {
public:
    explicit BaseComponent(std::shared_ptr<GameObject> gameObject) : gameObject(gameObject) {}

    std::shared_ptr<GameObject> gameObject = nullptr;
};


#endif //BASECOMPONENTS_H
