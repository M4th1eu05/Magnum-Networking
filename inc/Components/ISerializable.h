//
// Created by Tarook on 02/04/2025.
//

#ifndef ISERIALIZABLE_H
#define ISERIALIZABLE_H

#include <iostream>

class ISerializable {
public:
    virtual ~ISerializable() = default;
    virtual void serialize(std::ostream& os) const = 0;
    virtual void deserialize(std::istream& is) = 0;
};

#endif //ISERIALIZABLE_H
