#include "tests.h"

class person {
public:
    float name;
    int number;
};

class Rectangle {
    int width, height;
  public:
    Rectangle ();
    Rectangle (int,int);
    void set_values (int,int);
    int area() {return width*height;}
};

Rectangle::Rectangle () {
  width = 5;
  height = 5;
}

void Rectangle::set_values (int x, int y) {
  width = x;
  height = y;
}

Rectangle::Rectangle (int a, int b) {
  width = a;
  height = b;
}

class Circle {
    int radius;
  public:
    Circle(int r) { radius = r; }
    int circum() {return 2*radius*3;}
};

int main()
{
    plan(7);

	person obj;
	obj.name = 2.3;
	obj.number = 4;

	is_eq(obj.name,2.3);
	is_eq(obj.number,4);

  Rectangle rect;
  rect.set_values (3,4);
  is_eq(rect.area(),12);

Rectangle rectb (5,6);
is_eq(rectb.area(),30);

Rectangle rectc;
is_eq(rectc.area(),25);

  Circle foo (10);   // functional form
  is_eq(foo.circum(), 6*10);
  Circle bar = 20;   // assignment init.
  is_eq(bar.circum(), 6*20);

    done_testing();
}
