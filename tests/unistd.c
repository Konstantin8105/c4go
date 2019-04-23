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
    int fd = STDIN_FILENO;
    char data[128];
    for (int i = 0; i < 128; i++) {
        data[i] = '\x00';
    }
    printf("data = %s\n", data);
    is_streq(data, "");
    ssize_t s = read(fd, data, 128);
    if (s < 0) {
        fail("not good");
        write(2, "An error occurred in the read.\n", 31);
    }
    printf("data = %s\n", data);
    is_streq(data, "7");
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
//

void test_fstat()
{
	char * filename = "./tests/stdio.txt";

    int file=0;
        if((file=open(filename,O_RDONLY)) < -1)
            return;
 
    struct stat fileStat;
    if(fstat(file,&fileStat) < 0)    
        return;
 
    printf("Information for %s\n",filename);
    printf("---------------------------\n");
    printf("File Size: \t\t%d bytes\n",fileStat.st_size);
    printf("Number of Links: \t%d\n",fileStat.st_nlink);
    printf("File inode: \t\t%d\n",fileStat.st_ino);
 
    printf("File Permissions: \t");
    printf( (S_ISDIR(fileStat.st_mode)) ? "d" : "-");
    printf( (fileStat.st_mode & S_IRUSR) ? "r" : "-");
    printf( (fileStat.st_mode & S_IWUSR) ? "w" : "-");
    printf( (fileStat.st_mode & S_IXUSR) ? "x" : "-");
    printf( (fileStat.st_mode & S_IRGRP) ? "r" : "-");
    printf( (fileStat.st_mode & S_IWGRP) ? "w" : "-");
    printf( (fileStat.st_mode & S_IXGRP) ? "x" : "-");
    printf( (fileStat.st_mode & S_IROTH) ? "r" : "-");
    printf( (fileStat.st_mode & S_IWOTH) ? "w" : "-");
    printf( (fileStat.st_mode & S_IXOTH) ? "x" : "-");
    printf("\n\n");
 
    printf("The file %s a symbolic link\n\n", (S_ISLNK(fileStat.st_mode)) ? "is" : "is not");
}

int main()
{
    plan(10);

    START_TEST(write);
    START_TEST(read);
    START_TEST(read_file);
    // START_TEST(struct);
	START_TEST(fstat);

    done_testing();
}
