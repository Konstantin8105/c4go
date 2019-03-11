#include "tests.h"
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

#define MSGSIZE 16
char* msg1 = "hello, world #1";
char* msg2 = "hello, world #2";
char* msg3 = "hello, world #3";

void test_write()
{
    char inbuf[MSGSIZE];
    for (int i = 0; i < MSGSIZE; i++) {
        inbuf[i] = '\x00';
    }
    int p[2] = { 0, 0 };

    if (pipe(p) < 0)
        exit(1);

    is_true(p[0] != p[1]);

    write(p[1], msg1, MSGSIZE);
    write(p[1], msg2, MSGSIZE);
    write(p[1], msg3, MSGSIZE);

    read(p[0], inbuf, MSGSIZE);
    is_streq(inbuf, msg1);

    read(p[0], inbuf, MSGSIZE);
    is_streq(inbuf, msg2);

    read(p[0], inbuf, MSGSIZE);
    is_streq(inbuf, msg3);
}

void test_read()
{
    char data[128];
    for (int i = 0; i < 128; i++) {
        data[i] = '\x00';
    }
    printf("data = %s\n", data);
    ssize_t s = read(STDIN_FILENO, data, 128);
    if (s < 0) {
        fail("not good");
        write(2, "An error occurred in the read.\n", 31);
    }
    printf("data = %s\n", data);
    is_true(strlen(data) > 0);

    diag("wrong read");
    ssize_t sw = read(272727272, data, 128);
    is_eq(sw, -1);
}

void test_read_file()
{
    int fd;
    ssize_t size_read;
    char buffer[200];
    fd = open("./tests/stdio.txt", O_RDONLY);
    size_read = read(fd, buffer, sizeof(buffer));
    if (size_read == -1) {
        fail("not ok");
        return;
    }
    is_true(size_read > 0);
    is_true(strlen(buffer) > 0);
    buffer[size_read] = '\x00';
    printf("buffer = `%s`\n", buffer);
    close(fd);
}

typedef struct {
    int left;
    int right;
} pair_t;

// void test_struct()
// {
// {// write
// pair_t p;
// p.left = 10;
// p.right = 20;
// write(1, &p, sizeof(pair_t));
// }
// { // read
// pair_t p;
// read(0, &p, sizeof(pair_t));
// printf("left: %d right: %d\n",p.left, p.right);
// }
// }

int main()
{
    plan(8);

    START_TEST(write);
    START_TEST(read);
    START_TEST(read_file);
    // START_TEST(struct);

    done_testing();
}
