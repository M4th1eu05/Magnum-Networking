//
// Created by Tarook on 18/03/2025.
//

#ifndef LIFECYCLE_H
#define LIFECYCLE_H

class ILifeCycle {
public:
    virtual void start() = 0;
    virtual void update() = 0;
    virtual void stop() = 0;
};
#endif //LIFECYCLE_H
