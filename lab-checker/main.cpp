#include "include/tcp-server/Server.hpp"
#include <iostream>

int main() {
    try {
        Server server(1234);
        server.run();
    } catch (const std::exception& e) {
        std::cerr << "Server error: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}