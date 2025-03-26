cmake_minimum_required(VERSION 3.30)
project(MagnumBootstrap)

set(CMAKE_MODULE_PATH "${PROJECT_SOURCE_DIR}/modules/" ${CMAKE_MODULE_PATH})

# Activer C++ 17
set(CMAKE_CXX_STANDARD 17)

include(FetchContent)

# Récupérer ENet depuis GitHub
FetchContent_Declare(
        enet
        GIT_REPOSITORY https://github.com/lsalzman/enet.git
        GIT_TAG master  # Utilise la dernière version disponible
)

# Charger et intégrer ENet
FetchContent_MakeAvailable(enet)

# Add exe server
add_executable(Server src/server.cpp)
target_link_libraries(Server PRIVATE enet)

# Add exe client
add_executable(Client src/client.cpp)
target_link_libraries(Client PRIVATE enet)

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

add_executable(MagnumBootstrap
        src/main.cpp
        src/World.cpp
        src/GameObject.cpp
        src/Components/BaseComponent.cpp
        src/Components/Rigidbody.cpp
        src/client.cpp
        src/server.cpp)



target_include_directories(MagnumBootstrap PUBLIC inc inc/Components)
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