cmake_minimum_required(VERSION 3.30)
project(MagnumBootstrap)

set(CMAKE_MODULE_PATH "${PROJECT_SOURCE_DIR}/modules/" ${CMAKE_MODULE_PATH})

add_subdirectory(externals)
find_package(Magnum REQUIRED
        GL
        GlfwApplication
        Primitives
        Shaders
        Trade
        MeshTools
        Primitives
        SceneGraph)
find_package(MagnumIntegration REQUIRED ImGui Bullet)
find_package(Bullet REQUIRED Dynamics)

set_directory_properties(PROPERTIES CORRADE_USE_PEDANTIC_FLAGS ON)

include(FetchContent)

# Fetch enet6
FetchContent_Declare(
        enet6
        GIT_REPOSITORY https://github.com/SirLynix/enet6.git
        GIT_TAG main # Replace with a specific tag or commit if needed
)

FetchContent_MakeAvailable(enet6)

add_executable(MagnumBootstrap
        src/main.cpp
        src/World.cpp
        src/GameObject.cpp
        src/Components/BaseComponent.cpp
        inc/ILifeCycle.h
        src/Components/Rigidbody.cpp
        inc/Components/Rigidbody.h
        inc/Components/ISerializable.h
        src/Components/Collider.cpp
        inc/Components/Collider.h
)


target_include_directories(MagnumBootstrap PUBLIC inc inc/Components)
target_include_directories(MagnumBootstrap PUBLIC inc inc ${enet6_SOURCE_DIR}/include)


set_property(TARGET MagnumBootstrap PROPERTY CXX_STANDARD 20)
target_link_libraries(MagnumBootstrap PRIVATE
        Magnum::Application
        Magnum::GL
        Magnum::Magnum
        Magnum::MeshTools
        Magnum::Primitives
        Magnum::SceneGraph
        Magnum::Shaders
        Magnum::Trade
        MagnumIntegration::Bullet
        Bullet::Dynamics
        MagnumIntegration::ImGui
        BulletDynamics
        BulletCollision
        LinearMath
        Bullet3Common
        nlohmann_json
)


add_executable(BulletExample src/BulletExample.cpp)
set_property(TARGET BulletExample PROPERTY CXX_STANDARD 20)
target_link_libraries(BulletExample PRIVATE
        Magnum::Application
        Magnum::GL
        Magnum::Magnum
        Magnum::MeshTools
        Magnum::Primitives
        Magnum::SceneGraph
        Magnum::Shaders
        Magnum::Trade
        MagnumIntegration::Bullet
        Bullet::Dynamics
        MagnumIntegration::ImGui
        BulletDynamics
        BulletCollision
        LinearMath
        Bullet3Common
        nlohmann_json
)

# Link enet6 to your target
target_link_libraries(MagnumBootstrap PRIVATE enet)
target_link_libraries(BulletExample PRIVATE enet)

add_executable(DedicatedServer
        src/DedicatedServerMain.cpp
        src/DedicatedServer.cpp
        inc/DedicatedServer.h
        src/World.cpp
        src/GameObject.cpp
        src/Components/BaseComponent.cpp
        inc/ILifeCycle.h
        src/Components/Rigidbody.cpp
        inc/Components/Rigidbody.h
        inc/Components/ISerializable.h
        src/Components/Collider.cpp
        inc/Components/Collider.h
        inc/MessageType.h
)

target_include_directories(DedicatedServer PUBLIC inc ${enet6_SOURCE_DIR}/include)

target_link_libraries(DedicatedServer PRIVATE enet)
set_property(TARGET DedicatedServer PROPERTY CXX_STANDARD 20)
target_link_libraries(DedicatedServer PRIVATE
        Magnum::Application
        Magnum::GL
        Magnum::Magnum
        Magnum::MeshTools
        Magnum::Primitives
        Magnum::SceneGraph
        Magnum::Shaders
        Magnum::Trade
        MagnumIntegration::Bullet
        Bullet::Dynamics
        MagnumIntegration::ImGui
        BulletDynamics
        BulletCollision
        LinearMath
        Bullet3Common
        nlohmann_json
)


add_executable(Client
        src/ClientMain.cpp
        src/Client.cpp
        inc/Client.h
        src/World.cpp
        src/GameObject.cpp
        src/Components/BaseComponent.cpp
        inc/ILifeCycle.h
        src/Components/Rigidbody.cpp
        inc/Components/Rigidbody.h
        inc/Components/ISerializable.h
        src/Components/Collider.cpp
        inc/Components/Collider.h
        inc/MessageType.h
        src/DedicatedServer
)

target_include_directories(Client PUBLIC inc ${enet6_SOURCE_DIR}/include)

target_link_libraries(Client PRIVATE enet)

set_property(TARGET Client PROPERTY CXX_STANDARD 20)
target_link_libraries(Client PRIVATE
        Magnum::Application
        Magnum::GL
        Magnum::Magnum
        Magnum::MeshTools
        Magnum::Primitives
        Magnum::SceneGraph
        Magnum::Shaders
        Magnum::Trade
        MagnumIntegration::Bullet
        Bullet::Dynamics
        MagnumIntegration::ImGui
        BulletDynamics
        BulletCollision
        LinearMath
        Bullet3Common
        nlohmann_json
)