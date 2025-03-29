#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <unistd.h>
#include <string.h>

void die(const char *msg) {
  fprintf(stderr, "[%d] %s\n", errno, msg);
  abort();
}
void msg(const char *msg) {
    fprintf(stderr, "%s\n", msg);
}
void do_something(int connfd)
{
    char rbuf[64] = {};

    ssize_t n = read(connfd, rbuf, sizeof(rbuf) - 1);
    if(n < 0) {
        msg("read() error");
        return;
    }
    fprintf(stderr, "client says: %s \n", rbuf);
    char wbuff[] = "world";
    write(connfd, wbuff, strlen(wbuff));
}

int main() {
  int fd = socket(AF_INET, SOCK_STREAM, 0);

  if (fd < 0) {
    die("socket");
  }
  int val = 1;
  setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &val, sizeof(val));

  struct sockaddr_in addr = {};
  addr.sin_family = AF_INET;
  addr.sin_port = ntohs(1234);
  addr.sin_addr.s_addr = ntohl(0);

  if (bind(fd, (const struct sockaddr *)&addr, sizeof(addr)) == -1) {
    die("listen");
  }

  if (listen(fd, SOMAXCONN) == -1) {
    die("listen");
  }

  while(1) {
      struct sockaddr_in client_addr = {};

      socklen_t socklen = sizeof(client_addr);

      int conn_fd = accept(fd,(struct sockaddr *)&client_addr,&socklen);
      printf("conn_fd %d",conn_fd);
      if (conn_fd == -1) {
          continue;
      }
      while(1) {
          do_something(conn_fd);
      }
      // close(conn_fd);
  }
  return 0;
}
