cmake_minimum_required(VERSION 3.30)
project(MagnumBootstrap)

set(CMAKE_MODULE_PATH "${PROJECT_SOURCE_DIR}/modules/" ${CMAKE_MODULE_PATH})

# Activer C++ 17
set(CMAKE_CXX_STANDARD 17)

add_subdirectory("externals/enet6")

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
        src/server.cpp
        src/test_enet6.cpp
        src/WorldSerializer.cpp
        inc/WorldSerializer.h
        inc/World2.h
        src/World2.cpp)



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
        enet6
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
        enet6
)


# Add exe server
add_executable(Server src/server.cpp)
target_link_libraries(Server PRIVATE enet6)

# Add exe client
add_executable(Client src/client.cpp)
target_link_libraries(Client PRIVATE enet6)


add_executable(TestEnet6 src/test_enet6.cpp)
target_link_libraries(TestEnet6 PRIVATE ${CMAKE_SOURCE_DIR}/externals/enet6/build/libenet.a)
target_include_directories(Server PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/externals/enet6/include)
target_include_directories(Client PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/externals/enet6/include)
target_include_directories(TestEnet6 PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/externals/enet6/include)
target_include_directories(MagnumBootstrap PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/externals/enet6/include)


