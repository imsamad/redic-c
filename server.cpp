#include <arpa/inet.h>
#include <bits/stdc++.h>
#include <cerrno>
#include <errno.h>
#include <fcntl.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <vector>

using namespace std;

int MAX_MSG_SIZE = 1024;

void msg(const char *msg) { printf("msg: errorno - %d, %s", errno, msg); }

void die(const char *msg_) {
  msg(msg_);
  abort();
}

void set_nb_flags(int fd) {
  errno = 0;

  int flags = fcntl(fd, F_GETFL, 0);
  if (errno) {
    die("fcntl");
  }

  flags |= O_NONBLOCK;

  (void)fcntl(fd, F_SETFL, flags);

  if (errno) {
    die("set flags");
  }
}

struct Conn {
  int fd;
  bool want_read = false;
  bool want_write = false;
  bool want_close = false;
  vector<uint8_t> incoming;
  vector<uint8_t> outgoing;
};

Conn *handle_accept(int fd) {
  struct sockaddr_in client_addr = {};
  socklen_t socklen = sizeof(client_addr);
  int conn_fd = accept(fd, (struct sockaddr *)&client_addr, &socklen);

  if (conn_fd < 0) {
    msg("connection failed");
    return nullptr;
  }

  uint32_t ip = ntohl(client_addr.sin_addr.s_addr);
  printf("new client %u.%u.%u.%u:%u", (ip >> 24) & 255, (ip >> 16) & 255,
         (ip >> 8) & 255, ip & 255, ntohs(client_addr.sin_port));

  set_nb_flags(conn_fd);

  Conn *conn = new Conn();
  conn->fd = conn_fd;
  conn->want_read = true;

  return conn;
}

void handle_read(Conn *conn) {}
void handle_write(Conn *conn) {}

int main() {

  int fd = socket(AF_INET, SOCK_STREAM, 0);

  if (fd < 0) {
    die("socket");
  }

  struct sockaddr_in addr = {};
  addr.sin_family = AF_INET;
  addr.sin_addr.s_addr = htonl(0);
  addr.sin_port = htons(1234);

  if (bind(fd, (const struct sockaddr *)&addr, sizeof(addr)) == -1) {
    die("bind");
  }

  set_nb_flags(fd);

  if (listen(fd, SOMAXCONN)) {
    die("listen");
  }
  vector<Conn *> fd2conn;
  vector<struct pollfd> poll_args;

  while (1) {
    poll_args.clear();
    struct pollfd pfd = {fd, POLLIN, 0};
    poll_args.push_back(pfd);

    for (Conn *conn : fd2conn) {
      if (!conn) {
        continue;
      }

      struct pollfd pfd = {conn->fd, POLLERR, 0};

      if (conn->want_read) {
        pfd.events |= POLLIN;
      }
      if (conn->want_write) {
        pfd.events |= POLLOUT;
      }
      poll_args.push_back(pfd);
    }

    int rv = poll(poll_args.data(), (nfds_t)poll_args.size(), -1);

    if (rv < 0 && errno == EINTR) {
      continue;
    }
    if (rv < 0) {
      die("poll");
    }

    if (poll_args[0].revents) {
      if (Conn *conn = handle_accept(fd)) {
        if (fd2conn.size() < conn->fd) {
          fd2conn.resize(conn->fd + 1);
        }
        assert(!fd2conn[conn->fd]);
        fd2conn[conn->fd] = conn;
      }
    }

    for (int i = 1; i < poll_args.size(); i++) {
      uint32_t ready = poll_args[i].revents;

      if (ready == 0) {
        continue;
      }

      Conn *conn = fd2conn[poll_args[i].fd];

      if (ready & POLLIN) {
        assert(conn->want_read);
        handle_read(conn);
      }
      if (ready & POLLOUT) {
        assert(conn->want_write);
        handle_write(conn);
      }
      if ((ready & POLLERR) || conn->want_close) {
        assert(conn->want_close);
        (void)close(conn->fd);
        fd2conn[poll_args[i].fd] = nullptr;
        delete conn;
      }
    }
  }
}
