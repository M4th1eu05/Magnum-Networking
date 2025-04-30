//
// Created by Tarook on 18/03/2025.
//

#ifndef LIFECYCLE_H
#define LIFECYCLE_H

class ILifeCycle {
public:
    virtual void start() {};
    virtual void update() {};
    virtual void stop() {};
    virtual void destroy() {};
};
#endif //LIFECYCLE_H
