#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>

int main() {
    struct sockaddr_in peer;
    peer.sin_family = AF_INET;
    peer.sin_port = htons(3035);
    peer.sin_addr.s_addr = inet_addr("127.0.0.1");

    int s = socket(AF_INET, SOCK_STREAM, 0);
    if (s < 0) {
        fprintf(stderr, "ошибка создания сокета\n");
        return 1;
    }

    int conn = connect(s, (struct sockaddr *) &peer, sizeof(peer));
    if (conn < 0) {
        fprintf(stderr, "ошибка подключения\n");
        return 1;
    }
    char buff[100];
    while (1) {
        char string[100];
        scanf("%s", string);



        int err = send(s, string, 100, 0);
        if (err < 0) {
            fprintf(stderr, "ошибка отправки сообщения\n");
            return 1;
        }

        err = recv(s, buff, 100, 0);
        if (err < 0) {
            fprintf(stderr, "ошибка получения ответа\n");
            return 1;
        }

        printf("%s\n", buff);

        if (string[0] == '!') {
            break;
        }
    }

    int err = shutdown(s, 1);
    if (err < 0) {
        fprintf(stderr, "ошибка отправки FIN-сегмента\n");
        return 1;
    }


    err = close(s);
    if (err < 0) {
        fprintf(stderr, "ошибка закрытия дескриптора\n");
        return 1;
    }

    return 0;
}
