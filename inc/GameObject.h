//
// Created by Tarook on 18/03/2025.
//

#ifndef GAMEOBJECT_H
#define GAMEOBJECT_H


#include <Components/BaseComponent.h>
#include <memory>
#include <Magnum/SceneGraph/MatrixTransformation3D.h>
#include <Magnum/SceneGraph/Scene.h>

#include "ILifeCycle.h"
#include "Components/ISerializable.h"
#include "Magnum/Math/Quaternion.h"

class World;
using namespace Magnum;

typedef SceneGraph::Object<SceneGraph::MatrixTransformation3D> Object3D;

class GameObject : public Object3D, public ILifeCycle, public ISerializable, public std::enable_shared_from_this<GameObject> {
public:
    explicit GameObject(Object3D* parent = nullptr, const std::shared_ptr<World> &world = nullptr) : Object3D(parent), _world(world) {}

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

    template<typename T>
    std::shared_ptr<T> getComponent() const {
        static_assert(std::is_base_of_v<BaseComponent, T>, "T must be derived from BaseComponent");
        for (const auto& component : components) {
            if (auto casted = std::dynamic_pointer_cast<T>(component)) {
                return casted;
            }
        }
        return nullptr;
    }

    std::shared_ptr<World> getWorld() const {
        return _world;
    }

    template<typename T>
    std::vector<std::shared_ptr<T>> getComponents() const {
        static_assert(std::is_base_of_v<BaseComponent, T>, "T must be derived from BaseComponent");
        std::vector<std::shared_ptr<T>> result;
        for (const auto& component : components) {
            if (auto casted = std::dynamic_pointer_cast<T>(component)) {
                result.push_back(casted);
            }
        }
        return result;
    }

    void serialize(std::ostream& os) const override {
        // Serialize transformation data
        Matrix4 transformationMatrix = this->transformationMatrix();
        const Vector3 position = transformationMatrix.translation();
        const Vector3 scale = transformationMatrix.scaling();
        const Quaternion rotation = Quaternion::fromMatrix(transformationMatrix.rotationScaling());

        // Write position, rotation, and scale to the stream
        os.write(reinterpret_cast<const char*>(&position), sizeof(position));
        os.write(reinterpret_cast<const char*>(&rotation), sizeof(rotation));
        os.write(reinterpret_cast<const char*>(&scale), sizeof(scale));

        // Serialize GameObject data
        size_t componentCount = components.size();
        os.write(reinterpret_cast<const char*>(&componentCount), sizeof(componentCount));
        for (const auto& component : components) {
            component->serialize(os);
        }
    }

    void deserialize(std::istream& is) override {
        // Deserialize transformation data
        Vector3 position;
        Quaternion rotation;
        Vector3 scale;

        is.read(reinterpret_cast<char*>(&position), sizeof(position));
        is.read(reinterpret_cast<char*>(&rotation), sizeof(rotation));
        is.read(reinterpret_cast<char*>(&scale), sizeof(scale));

        // Apply the deserialized transformation to the GameObject
        this->setTransformation(Matrix4::from(rotation.toMatrix(), position) * Matrix4::scaling(scale));

        // Deserialize GameObject components
        size_t componentCount;
        is.read(reinterpret_cast<char*>(&componentCount), sizeof(componentCount));
        for (size_t i = 0; i < componentCount; ++i) {
            // Assuming a factory method to create components from the stream
            //auto component = createComponentFromStream(is);
            //addComponentInternal(component);
        }
    }

    /*
    std::shared_ptr<BaseComponent> createComponentFromStream(std::istream& is) {
        // Read the component type identifier
        std::string componentType;
        size_t typeLength;
        is.read(reinterpret_cast<char*>(&typeLength), sizeof(typeLength));
        componentType.resize(typeLength);
        is.read(&componentType[0], typeLength);

        // Use a registry to find the appropriate constructor
        static const std::unordered_map<std::string, std::function<std::shared_ptr<BaseComponent>(std::istream&)>> componentRegistry = {
            {"Rigidbody", [](std::istream& is) {
                // Deserialize Rigidbody-specific parameters
                Float mass;
                is.read(reinterpret_cast<char*>(&mass), sizeof(mass));

                // Deserialize collision shape (assuming a factory or placeholder for deserialization)
                btCollisionShape* bShape = deserializeCollisionShape(is);

                // Deserialize world (assuming a global or accessible context for the world)
                auto world = getWorldFromContext();

                // Deserialize GameObject (assuming the GameObject is already created and passed)
                auto gameObject = getGameObjectFromContext();

                // Create and return the Rigidbody
                auto component = std::make_shared<Rigidbody>(mass, bShape, world, gameObject);
                component->deserialize(is);
                return component;
            }},
            // Add other components here
        };
        auto it = componentRegistry.find(componentType);
        if (it == componentRegistry.end()) {
            throw std::runtime_error("Unknown component type: " + componentType);
        }

        return it->second(is);
    }
    */

private:
    void addComponentInternal(const std::shared_ptr<BaseComponent> &component);

private:
    std::vector<std::shared_ptr<BaseComponent>> components;
protected:
    std::shared_ptr<World> _world;
};


#endif //GAMEOBJECT_H
