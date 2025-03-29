#include <iostream>
#include <vector>
#include <string>
#include <cstring>
#include <cassert>
#include <cerrno>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/ip.h>

constexpr size_t kMaxMsgSize = 32 << 20;  // 32 MB

void die(const std::string &msg) {
    std::cerr << "[Error] " << msg << ": " << strerror(errno) << std::endl;
    exit(EXIT_FAILURE);
}

bool readFull(int fd, uint8_t *buf, size_t n) {
    while (n > 0) {
        ssize_t rv = read(fd, buf, n);
        if (rv <= 0) return false;
        n -= rv;
        buf += rv;
    }
    return true;
}

bool writeAll(int fd, const uint8_t *buf, size_t n) {
    while (n > 0) {
        ssize_t rv = write(fd, buf, n);
        if (rv <= 0) return false;
        n -= rv;
        buf += rv;
    }
    return true;
}

bool sendRequest(int fd, const std::string &text) {
    if (text.size() > kMaxMsgSize) return false;

    uint32_t len = text.size();
    std::vector<uint8_t> buffer(4 + len);
    memcpy(buffer.data(), &len, 4);
    memcpy(buffer.data() + 4, text.data(), len);

    return writeAll(fd, buffer.data(), buffer.size());
}

bool readResponse(int fd) {
    uint32_t len;
    if (!readFull(fd, reinterpret_cast<uint8_t *>(&len), 4) || len > kMaxMsgSize) {
        std::cerr << "Invalid response length." << std::endl;
        return false;
    }

    std::vector<uint8_t> buffer(len);
    if (!readFull(fd, buffer.data(), len)) {
        std::cerr << "Failed to read response." << std::endl;
        return false;
    }

    std::cout << "Response (" << len << " bytes): "
              << std::string(buffer.begin(), buffer.begin() + std::min<size_t>(100, len))
              << std::endl;
    return true;
}

int main() {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) die("socket()");

    sockaddr_in addr{};
    addr.sin_family = AF_INET;
    addr.sin_port = htons(1234);
    addr.sin_addr.s_addr = htonl(INADDR_LOOPBACK);

    if (connect(fd, reinterpret_cast<sockaddr *>(&addr), sizeof(addr)) < 0)
        die("connect()");

    std::vector<std::string> requests = {std::string(kMaxMsgSize, 'z')};
    for (const auto &req : requests) {
        if (!sendRequest(fd, req)) {
            std::cerr << "Failed to send request." << std::endl;
            break;
        }
    }
    for ([[maybe_unused]] const auto &req : requests) {
        if (!readResponse(fd)) break;
    }

    close(fd);
    return 0;
}
