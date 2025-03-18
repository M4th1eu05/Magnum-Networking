include(FetchContent)

if(WIN32)
    set(CMAKE_PREFIX_PATH ${CMAKE_CURRENT_SOURCE_DIR}/glfw-3.4.bin.WIN64 ${CMAKE_PREFIX_PATH})
endif()

set(IMGUI_DIR ${CMAKE_CURRENT_SOURCE_DIR}/imgui)

set(MAGNUM_WITH_GLFWAPPLICATION ON CACHE BOOL "" FORCE)
set(MAGNUM_WITH_BULLET ON CACHE BOOL "" FORCE)
set(MAGNUM_WITH_IMGUI ON CACHE BOOL "" FORCE)

add_subdirectory(bullet)

add_subdirectory(corrade)
add_subdirectory(magnum)
add_subdirectory(magnum-integration)
add_subdirectory(nlohmannjson)
