#include <stdio.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>

int main() {
    int conn_num = 5;

    struct sockaddr_in addr = {
            AF_INET,
            htons(3035),
            {inet_addr("127.0.0.1")}
    };

    int s = socket(AF_INET, SOCK_STREAM, 0);
    if (s < 0) {
        fprintf(stderr, "ошибка создания сокета");
        return 1;
    }

    int err = bind(s, (struct sockaddr *) &addr, sizeof(addr));
    if (err < 0) {
        fprintf(stderr, "ошибка bind");
        return 1;
    }

    err = listen(s, 32);
    if (err < 0) {
        fprintf(stderr, "ошибка установки прослушивания");
        return 1;
    }

    struct sockaddr_in client_addr;
    int client_addr_size = sizeof(client_addr);

    int i = 0;
    char buff[100];
    buff[0] = 'e';
    buff[1] = 'c';
    buff[2] = 'h';
    buff[3] = 'o';
    buff[4] = ' ';

    while (i < conn_num) {
        int client = accept(s, (struct sockaddr *) &client_addr, &client_addr_size);
        if (client < 0) {
            fprintf(stderr, "ошибка приема подключения");
            return 1;
        }

        while (1) {
            for (int i = 5; i < 100; ++i)
                buff[i] = 0;

            err = recv(client, buff + 5, 100, 0);
            if (err < 0) {
                fprintf(stderr, "ошибка получения сообщения");
                return 1;
            }

            printf("%s from %s:%hu\n", buff + 5, inet_ntoa(client_addr.sin_addr), ntohs(client_addr.sin_port));

            if ((buff + 5)[0] == '!') {
                err = send(client, "close", 5, 0);
                if (err < 0) {
                    fprintf(stderr, "ошибка передачи сообщения");
                    return 1;
                }
                break;
            } else {
                int size = 0;
                for (int i = 5; i<100; ++i)
                    if (buff[i] == 0)
                        break;
                    else
                        size++;

                for (int i = 5; i<5+size/2; ++i){
                    char c = buff[i];
                    buff[i] = buff[5+size - (i-5+1)];
                    buff[5+size - (i-5+1)] = c;
                }
                err = send(client, buff, 100, 0);
                if (err < 0) {
                    fprintf(stderr, "ошибка передачи сообщения");
                    return 1;
                }
            }
        }

        err = shutdown(client, 1);
        if (err < 0) {
            fprintf(stderr, "ошибка получения сообщения");
            return 1;
        }
        ++i;
    }

    err = close(s);
    if (err < 0) {
        fprintf(stderr, "ошибка закрытия дескриптора");
        return 1;
    }

    return 0;
}
