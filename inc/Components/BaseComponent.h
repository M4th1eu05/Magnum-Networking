//
// Created by Tarook on 18/03/2025.
//

#ifndef BASECOMPONENTS_H
#define BASECOMPONENTS_H
#include <memory>

#include "ILifeCycle.h"
#include "ISerializable.h"


class GameObject;

class BaseComponent : public ILifeCycle, public ISerializable {
public:
    explicit BaseComponent(std::shared_ptr<GameObject> gameObject) : gameObject(gameObject) {}

    std::shared_ptr<GameObject> gameObject = nullptr;

    void serialize(std::ostream& os) const override;

    void deserialize(std::istream& is) override;
};


#endif //BASECOMPONENTS_H
