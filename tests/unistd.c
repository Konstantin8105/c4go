#include <unistd.h> 
#include "tests.h"

#define MSGSIZE 16 
char* msg1 = "hello, world #1"; 
char* msg2 = "hello, world #2"; 
char* msg3 = "hello, world #3"; 
  
int main() 
{ 
    plan(3);

    char inbuf[MSGSIZE]; 
    int p[2]; 
  
    if (pipe(p) < 0) 
        exit(1); 
  
    write(p[1], msg1, MSGSIZE); 
    write(p[1], msg2, MSGSIZE); 
    write(p[1], msg3, MSGSIZE); 

    read(p[0], inbuf, MSGSIZE); 
	is_streq(inbuf, msg1);

    read(p[0], inbuf, MSGSIZE); 
	is_streq(inbuf, msg2);

    read(p[0], inbuf, MSGSIZE); 
	is_streq(inbuf, msg3);

    done_testing();
} 
