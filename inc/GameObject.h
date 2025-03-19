//
// Created by Tarook on 18/03/2025.
//

#ifndef GAMEOBJECT_H
#define GAMEOBJECT_H


#include <BaseComponent.h>
#include <memory>
#include <Magnum/SceneGraph/MatrixTransformation3D.h>
#include <Magnum/SceneGraph/Scene.h>

#include "ILifeCycle.h"

using namespace Magnum;

typedef SceneGraph::Object<SceneGraph::MatrixTransformation3D> Object3D;

class GameObject : public Object3D, public ILifeCycle, public std::enable_shared_from_this<GameObject> {
public:
    explicit GameObject(Object3D* parent = nullptr) : Object3D(parent) {}

    void start() override;
    void update() override;
    void stop() override;

    void removeComponent(const std::shared_ptr<BaseComponent> &component);

    template<typename T, typename... Args>
    std::shared_ptr<T> addComponent(Args&&... args) {
        static_assert(std::is_base_of_v<BaseComponent, T>, "T must be derived from BaseComponent");
        auto component = std::make_shared<T>(std::forward<Args>(args)..., shared_from_this());
        addComponentInternal(component);
        return component;
    }

private:
    void addComponentInternal(const std::shared_ptr<BaseComponent> &component);

private:
    std::vector<std::shared_ptr<BaseComponent>> components;
};


#endif //GAMEOBJECT_H
