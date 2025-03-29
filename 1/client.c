#include <netinet/in.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <sys/types.h>
#include <unistd.h>

void die(const char *msg) {
  fprintf(stderr, "[%d] %s\n", errno, msg);
  abort();
}

void msg(const char *msg) {
    fprintf(stderr, "%s\n", msg);
}

int main() {
    int fd = socket(AF_INET, SOCK_STREAM, 0);

    if (fd == -1) {
        die("socket");
    }

    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_port = ntohs(1234);
    addr.sin_addr.s_addr = htonl(INADDR_LOOPBACK);

    if (connect(fd, (const struct sockaddr *)&addr,sizeof(addr)) == 1) {
        die("connect");
    }

    for(int i = 0; 1; i++) {
        char msg[] = "hello";
        write(fd, msg,sizeof(msg));

        char rbuf[64] = {};
        ssize_t n = read(fd, rbuf, sizeof(rbuf) - 1);

        if (n < 0) {
            die("read");
        }

        printf("server say: %s \n", rbuf);
    }

    // close(fd);

    return 0;
}
