#include "tests.h"

//--------------------------
class person {
public:
    float name;
    int number;
};

void simple()
{
	person obj;
	obj.name = 2.3;
	obj.number = 4;

	is_eq(obj.name,2.3);
	is_eq(obj.number,4);
}

/*
//--------------------------
class Rectangle {
    int width, height;
  public:
    Rectangle ();
    int area() {return width*height;}
};

Rectangle::Rectangle () {
  width = 4;
  height = 5;
}

void defaulf_constructor()
{
	Rectangle rec;
	is_eq(rec.area(),20);
}

//--------------------------
class Rectangle2 {
    int width, height;
  public:
    Rectangle2 (int,int);
    int area() {return width*height;}
};

Rectangle2::Rectangle2 (int a, int b) {
  width = a;
  height = b;
}

void constructor()
{
	Rectangle2 rec(2,3);
	is_eq(rec.area(),6);
}

//--------------------------
class Rectangle3 {
    int width, height;
  public:
    void set_values (int,int);
    int area() {return width*height;}
};

void Rectangle3::set_values (int x, int y) {
  width = x;
  height = y;
}

void setter()
{
	Rectangle3 rec;
	rec.set_values(4,5);
	is_eq(rec.area(),20);
}
*/

int main()
{
    plan(2);

	simple();
	// defaulf_constructor();
	// constructor();
	// setter();

    done_testing();
}
