#ifndef WORLDSERIALIZER_H
#define WORLDSERIALIZER_H
#pragma once

#include "World2.h"
#include <string>

bool saveWorldToFile(const World2& world, const std::string& filename);
bool loadWorldFromFile(World2& world, const std::string& filename);

#endif //WORLDSERIALIZER_H
