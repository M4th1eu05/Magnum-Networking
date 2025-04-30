//
// Created by mathi on 29/04/2025.
//

#ifndef APIHANDLER_H
#define APIHANDLER_H
#include <functional>
#include <string>

#include "cpr/cpr.h"


class APIHandler {
public:

    bool POST(const std::string &url, const std::string &postData, const std::function<void(std::string)> &callback);
    bool GET(const std::string &url, const std::function<void(std::string)>& callback);

private:

};

#endif //APIHANDLER_H
