void if_1()
{
    int a = 5;
    int b = 2;
    int c = 4;
    if (a > b) {
        return;
    } else if (c <= a) {
        a = 0;
    }
    (void)(a);
    (void)(b);
    (void)(c);

    int w = 2 > 1 ? -1 : 5;
    int r;
    r = 2 > 1 ? -1 : 5;
    r = (2 > 1) ? -1 : 5;
    r = (w > 1) ? -1 : 5;
    r = w > 1 ? -1 : 5;
    r = (w > 1) + (r == 4) ? -1 : 5;
    if (w > 0) {
        r = 3;
    }
}
