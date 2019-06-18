/*%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                888                                  %%%%%%%%%%%%
%%%%%%%%%%%%                888                                  %%%%%%%%%%%%
%%%%%%%%%%%%      o88°°88o  888°°88o  o88°°88o  888       888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  °°°  888  888  888  888  888       888    %%%%%%%%%%%%
%%%%%%%%%%%%      888       888  888  888  888  888  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      °888888o  888  888  888  888  888  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%           888  888  888  888  888  888  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      ooo  888  888  888  888  888  °88  888  88°    %%%%%%%%%%%%
%%%%%%%%%%%%      °88oo88°  888  888  °88oo88°   °888° °888°     %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                         888         %%%%%%%%%%%%
%%%%%%%%%%%%                                         888         %%%%%%%%%%%%
%%%%%%%%%%%%                                         888         %%%%%%%%%%%%
%%%%%%%%%%%%      8888888o888o   o88°°88o  o88°°88o  888°°88o    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  88o  888  888  888  °°°  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  888  888  888       888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  8888888°  °888888o  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  888            888  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  888  ooo  ooo  888  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  °88oo88°  °88oo88°  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%                        %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%  Author: Bojan NICENO  %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%                        %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%% niceno@univ.trieste.it %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%                        %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%*/
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/Xos.h>
#include <X11/Xatom.h>

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#ifndef max
#define max(a,b)  (((a) > (b)) ? (a) : (b))
#endif
#ifndef min
#define min(a,b)  (((a) < (b)) ? (a) : (b))
#endif
#define SMALL 1e-30
#define GREAT 1e+30

/*%%%%%%%%%%%%%%%%%%%%%%%%%%%%
   Definitions for the mesh 
%%%%%%%%%%%%%%%%%%%%%%%%%%%%*/
int Nn, Ne, Ns;
double xmax, ymax, xmin, ymin, scl, X0, Y0;

struct ele
 {
  int i,  j,  k;
  int mark;
  double xv, yv;
 }
* elem;


struct sid
 {
  int ea, eb;           /* left and right element */
  int a, b, c, d;       /* left, right, start and end point */
  int mark;             /* is it off, is on the boundary */
 }
* side;


struct nod
 {
  double x, y;
  int mark;
 }
* node;
/*%%%%%%%%%%%%%%%%%%%%%%%%%%*/

#define MAIN_WDTH     750
#define MAIN_HGHT     500
#define BUTTON_HEIGHT  18
#define BUTTON_WIDTH   80
#define NBUTTONS       10
#define ON              0
#define OFF            -1

#define DELAUNAY  0
#define VORONOI   1
#define NODES     2
#define ELEMENTS  3
#define MATERIALS 4
#define BOUNDARY  5
#define ZOOM      6
#define MOVE      7
#define FIT       8
#define QUIT      9


Display *display;
Window  main_win;
Window  draw_win;
int main_wdth, main_hght, draw_wdth, draw_hght;  

int         scr_num;
static char *prog_name, *file_name;
GC          gc_BoW, gc_WoB, gc_XOR, gc_THICK, gc_DASHED, gc_numb;
XFontStruct *text_font, *numb_font; 

int main_wdth=MAIN_WDTH, main_hght=MAIN_HGHT;

struct butt_data
 {
  int    x0, y0, hght, wdth, border;
  char   *caption;
  int    pressed;
 }

butt_data[NBUTTONS]=
 {
  { MAIN_WDTH-BUTTON_WIDTH-20,  10, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Delaunay",  ON},
  { MAIN_WDTH-BUTTON_WIDTH-20,  40, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Voronoi",   OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20,  70, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Nodes",     OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20, 100, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Elements",  OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20, 130, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Materials", OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20, 160, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Boundary",  OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20, 190, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Zoom",      OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20, 220, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Move",      OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20, 250, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Fit",       OFF},
  { MAIN_WDTH-BUTTON_WIDTH-20, 280, BUTTON_HEIGHT, BUTTON_WIDTH, 2, "Quit",      OFF},
 };

Window button[NBUTTONS];

/*========================================================================*/
getGC(Window win)
{
 unsigned long valuemask=0;
 XGCValues     values;

 int         dash_offset=0;
 static char dash_list[2]={1,1};
 int         list_len=2;
  
/* Normal, thin lines */
 gc_BoW   = XCreateGC(display, win, valuemask, &values);
 gc_WoB   = XCreateGC(display, win, valuemask, &values);

 XSetFont(display, gc_BoW,  text_font->fid);
 XSetFont(display, gc_WoB,  text_font->fid);

 XSetForeground(display, gc_BoW, BlackPixel(display, scr_num));
 XSetForeground(display, gc_WoB, WhitePixel(display, scr_num));

 XSetLineAttributes(display, gc_BoW, 0, LineSolid, CapRound, JoinRound);
 XSetLineAttributes(display, gc_WoB, 0, LineSolid, CapRound, JoinRound);

/* Thick lines */
 gc_THICK = XCreateGC(display, win, valuemask, &values);
 XSetForeground(display, gc_THICK, BlackPixel(display, scr_num));
 XSetLineAttributes(display, gc_THICK, 3, LineSolid, CapRound, JoinRound);

/* Dashed lines */
 gc_DASHED = XCreateGC(display, win, valuemask, &values);
 XSetForeground(display, gc_DASHED, BlackPixel(display, scr_num));
 XSetLineAttributes(display, gc_DASHED, 0, LineOnOffDash, CapRound, JoinRound);
 XSetDashes(display, gc_DASHED, dash_offset, dash_list, list_len);

/* numbers */
 gc_numb = XCreateGC(display, win, valuemask, &values);
 XSetFont(display, gc_numb, numb_font->fid);
 XSetForeground(display, gc_numb, BlackPixel(display, scr_num));

/* Invisible lines */
 gc_XOR = XCreateGC(display, win, 0, NULL);
 XSetFunction(display, gc_XOR, GXxor);
 XSetForeground(display, gc_XOR, WhitePixel(display, scr_num));
}

/*========================================================================*/
load_fonts()
{
 if( (text_font = XLoadQueryFont(display, "-*-helvetica-bold-r-normal--12-*")) == NULL )
  {
   (void) fprintf(stderr, "%s: Cannot open font\n", prog_name);
   exit(-1);
  } 

 if( (numb_font = XLoadQueryFont(display, "-*-helvetica-bold-r-normal--10-*")) == NULL )
  {
   (void) fprintf(stderr, "%s: Cannot open font\n", prog_name);
   exit(-1);
  } 
}

/*========================================================================*/
draw(Window win, GC gc, int win_x_dim, int win_y_dim)
{
 int x_0, y_0, x_dim, y_dim;

 x_0 = win_x_dim/5;
 y_0 = win_y_dim/5;
 
 x_dim = 3*win_x_dim/5; 
 y_dim = 3*win_y_dim/5; 

 XDrawRectangle(display, win, gc, x_0, y_0, x_dim, y_dim);
}

/*========================================================================*/
place_text(Window win, GC gc, XFontStruct *text_font, 
           int win_x_dim, int win_y_dim, char *string)
{
 int width, height; /* string height and width */

 width  = XTextWidth(text_font, string, strlen(string)); 
 height = text_font->ascent + text_font->descent;

 XDrawString(display, win, gc,
             (win_x_dim-width)/2, (win_y_dim+height)/2,
             string, strlen(string));  
}


/*========================================================================*/
create_buttons(Window parent)
{
 int b;
 
 unsigned long        valuemask = CWWinGravity;
 unsigned long        border, background;
 XSetWindowAttributes attr;
 
 attr.win_gravity = NorthEastGravity;

 for(b=0; b<NBUTTONS; b++)
  {
   if(butt_data[b].pressed==OFF)
    {border    =BlackPixel(display, scr_num);
     background=WhitePixel(display, scr_num);}
   else
    {border    =WhitePixel(display, scr_num);
     background=BlackPixel(display, scr_num);}

   button[b] = XCreateSimpleWindow(display, parent,
                                   butt_data[b].x0,   butt_data[b].y0, 
                                   butt_data[b].wdth, butt_data[b].hght, 
                                   butt_data[b].border,
                                   border, background);

   XSelectInput(display, button[b], ExposureMask | ButtonPressMask);
   XChangeWindowAttributes(display, button[b], valuemask, &attr);
   XMapWindow(display, button[b]);
  }
}

/*========================================================================*/
write_on_button(int b)
{
 int width, height; /* string height and width */

 unsigned long        valuemask = CWBackPixel | CWBorderPixel;
 XSetWindowAttributes attr;
 GC gc;

 height = text_font->ascent + text_font->descent;
 width  = XTextWidth(text_font, butt_data[b].caption, strlen(butt_data[b].caption)); 

 if(butt_data[b].pressed==ON)  
  {attr.background_pixel = BlackPixel(display, scr_num);
   attr.border_pixel     = WhitePixel(display, scr_num);}

 if(butt_data[b].pressed==OFF) 
  {attr.background_pixel = WhitePixel(display, scr_num);
   attr.border_pixel     = BlackPixel(display, scr_num);}

 XChangeWindowAttributes(display, button[b], valuemask, &attr);
 XClearWindow(display, button[b]); 
/* XFlush(display); */

 if(butt_data[b].pressed==ON)  gc=gc_WoB;
 if(butt_data[b].pressed==OFF) gc=gc_BoW;

 XDrawString(display, button[b], gc,
            (butt_data[b].wdth-width)/2, (butt_data[b].hght+height)/2,
             butt_data[b].caption, strlen(butt_data[b].caption));
}

/*=========================================================================*/
draw_mesh(Window win)
{
 int    e, n, s, ei, ej, ek, ea, eb;
 double x, y, xc, yc, xd, yd, x1, y1, x2, y2;
 char   numb[80];
 int f_hght, f_wdth;

 f_hght = 9; /*numb_font->ascent + numb_font->descent;*/

/***********************
*  Draw Delaunay Mesh  *
***********************/
 if(butt_data[DELAUNAY].pressed==ON)
   for(s=0; s<Ns; s++)
     if(side[s].mark!=OFF) 
      {
       xc=node[side[s].c].x; yc=node[side[s].c].y;
       xd=node[side[s].d].x; yd=node[side[s].d].y;
        
       XDrawLine(display, win, gc_BoW, 
                 (int)(xc*scl + X0), (int)(-yc*scl + Y0), 
                 (int)(xd*scl + X0), (int)(-yd*scl + Y0));
      }

/**********************
*  Draw Voronoi Mesh  *
**********************/
 if(butt_data[VORONOI].pressed==ON)
   for(s=0; s<Ns; s++)
     if(side[s].mark!=OFF)
      {
       if((ea=side[s].ea)!=OFF)
        {x1=elem[ea].xv;
         y1=elem[ea].yv;}
       else
        {x1=0.5*(node[side[s].c].x+node[side[s].d].x);
         y1=0.5*(node[side[s].c].y+node[side[s].d].y);}
     
       if((eb=side[s].eb)!=OFF)
        {x2=elem[eb].xv;
         y2=elem[eb].yv;}
       else
        {x2=0.5*(node[side[s].c].x+node[side[s].d].x);
         y2=0.5*(node[side[s].c].y+node[side[s].d].y);}

       XDrawLine(display, win, gc_DASHED, 
                 (int)(x1*scl + X0), (int)(-y1*scl + Y0), 
                 (int)(x2*scl + X0), (int)(-y2*scl + Y0));
      }

 for(s=0; s<Ns; s++)
   if(side[s].mark>0) /* It means, side is on the boundary */
    {
     xc=node[side[s].c].x; yc=node[side[s].c].y;
     xd=node[side[s].d].x; yd=node[side[s].d].y;

     XDrawLine(display, win, gc_THICK, 
	       (int)(xc*scl + X0), (int)(-yc*scl + Y0), 
               (int)(xd*scl + X0), (int)(-yd*scl + Y0));
    }

 if(butt_data[MATERIALS].pressed==ON || butt_data[ELEMENTS].pressed==ON)
   for(e=0; e<Ne; e++)
    {
     x=0.333*(node[elem[e].i].x+node[elem[e].j].x+node[elem[e].k].x);
     y=0.333*(node[elem[e].i].y+node[elem[e].j].y+node[elem[e].k].y);

     if(butt_data[MATERIALS].pressed==ON) sprintf(numb, "%d", elem[e].mark);
     if(butt_data[ELEMENTS].pressed ==ON) sprintf(numb, "%d", e);

     f_wdth  = XTextWidth(numb_font, numb, strlen(numb)); 

     XClearArea(display, win, 
                (int)(x*scl+X0-f_wdth/2-1), (int)(-y*scl+Y0-1-f_hght/2), 
                 f_wdth+2, f_hght+1, False);

     XDrawString(display, win, gc_numb, 
                 (int)(x*scl+X0-f_wdth/2), (int)(-y*scl+Y0+f_hght/2), 
                  numb, strlen(numb)); 
    }

 if(butt_data[NODES].pressed==ON)
   for(n=0; n<Nn; n++)
    {
     x=node[n].x; y=node[n].y;

     sprintf(numb, "%d", n);
     f_wdth  = XTextWidth(numb_font, numb, strlen(numb)); 

     XClearArea(display, win, 
                (int)(x*scl+X0-f_wdth/2-1), (int)(-y*scl+Y0-5-f_hght), 
                 f_wdth+2, f_hght+1, False);

     XDrawString(display, win, gc_numb, 
                 (int)(x*scl+X0-f_wdth/2), (int)(-y*scl+Y0-5), 
                  numb, strlen(numb)); 
    }

/*****************************
*  Draw Boundary Conditions  *
*****************************/
 if(butt_data[BOUNDARY].pressed==ON)
  {
   for(s=0; s<Ns; s++)
     if(side[s].mark>0) /* It means, side is on the boundary */
      {
       x = 0.5*(node[side[s].c].x + node[side[s].d].x); 
       y = 0.5*(node[side[s].c].y + node[side[s].d].y);

       sprintf(numb, "%d", side[s].mark);
       f_wdth  = XTextWidth(numb_font, numb, strlen(numb)); 

       XClearArea(display, win, 
                  (int)(x*scl+X0-f_wdth/2-1), (int)(-y*scl+Y0-1-f_hght/2), 
                   f_wdth+2, f_hght+1, False);

       XDrawString(display, win, gc_numb, 
                   (int)(x*scl+X0-f_wdth/2), (int)(-y*scl+Y0+f_hght/2), 
                    numb, strlen(numb)); 
      }

   for(n=0; n<Nn; n++)
     if(node[n].mark>0)
      {
       x=node[n].x; y=node[n].y;

       sprintf(numb, "%d", node[n].mark);
       f_wdth  = XTextWidth(numb_font, numb, strlen(numb)); 

       XClearArea(display, win, 
                  (int)(x*scl+X0-f_wdth/2-1), (int)(-y*scl+Y0-1-f_hght/2), 
                   f_wdth+2, f_hght+1, False);

       XDrawString(display, win, gc_numb, 
                   (int)(x*scl+X0-f_wdth/2), (int)(-y*scl+Y0+f_hght/2), 
                    numb, strlen(numb)); 
      }
  }
}
/*-draw_mesh--------------------------------------------------------------*/


/*========================================================================*/
load_mesh()
{
 int n, s, e, len;
 int d1, d2, d3, d4, d5, d6;
 char dummy[80];
 FILE *in;

 strcat(file_name, ".n");
 len=strlen(file_name);
 
/*--------+
|  Nodes  |
+--------*/
 in=fopen(file_name, "r");
 if(in==NULL) 
  {fprintf(stderr, "%s: cannot open file: %s\n\n", prog_name, file_name); 
   fflush(stdout);
   exit(-1);}
   
 fscanf(in, "%d", &Nn);
 node=(struct nod *) calloc(Nn, sizeof(struct nod)); 
 if(node==NULL) 
  {fprintf(stderr, "%s: cannot allocate enough memory\n\n", prog_name); 
   fflush(stdout);
   exit(-1);}

 for(n=0; n<Nn; n++)
  {
   fscanf(in, "%s %lf %lf %d", dummy, &node[n].x, &node[n].y, &node[n].mark);
   xmax=max(xmax, node[n].x); ymax=max(ymax, node[n].y);
   xmin=min(xmin, node[n].x); ymin=min(ymin, node[n].y);
  }
 fclose(in);

/*-----------+
|  Elements  |
+-----------*/
 file_name[len-1]='e';
 in=fopen(file_name, "r");
 if(in==NULL)
  {fprintf(stderr, "%s: cannot open file: %s\n\n", prog_name, file_name); 
   fflush(stdout);
   exit(-1);}

 fscanf(in, "%d", &Ne);

 elem=(struct ele *) calloc(Ne, sizeof(struct ele)); 
 if(elem==NULL) 
  {fprintf(stderr, "%s: cannot allocate enough memory\n\n", prog_name); 
   fflush(stdout);
   exit(-1);}

 for(e=0; e<Ne; e++)
  {
   fscanf(in, "%s %d %d %d %d %d %d %d %d %d %lf %lf %d", 
               dummy, &elem[e].i, &elem[e].j, &elem[e].k,
                      &d1, &d2, &d3, &d4, &d5, &d6,
                      &elem[e].xv, &elem[e].yv, &elem[e].mark);
  }
 fclose(in);

/*--------+
|  Sides  |
+--------*/
 file_name[len-1]='s';
 in=fopen(file_name, "r");
 if(in==NULL)
  {fprintf(stderr, "%s: cannot open file: %s\n\n", prog_name, file_name); 
   fflush(stdout);
   exit(-1);}

 fscanf(in, "%d", &Ns);

 side=(struct sid *) calloc(Ns, sizeof(struct sid)); 
 if(node==NULL) 
  {fprintf(stderr, "%s: cannot allocate enough memory\n\n", prog_name); 
   fflush(stdout);
   exit(-1);}

 for(s=0; s<Ns; s++)
  {
   fscanf(in, dummy);     
   fscanf(in, "%s %d %d %d %d %d", 
                      dummy, &side[s].c, &side[s].d, &side[s].ea, &side[s].eb,
                             &side[s].mark);
  }
 fclose(in);

}
/*------------------------------------------------------------------------*/


/*========================================================================*/
init(int argc, char **argv)
{
 XWindowAttributes main_win_attr;
 Pixmap ico_pixm, back_pixm;
 
/*----------------------+
|  Connect to X server  |
+----------------------*/
  {
   char *disp_name=NULL; /* ako nije definirano od korisnika mora biti NULL */

   if( (display=XOpenDisplay(disp_name)) == NULL ) 
    {
     fprintf(stderr, "%s: cannot connect to X server %s\n",
                      prog_name, XDisplayName(disp_name));
     exit(-1);
    }
  }

 scr_num = DefaultScreen(display);

/*---------------------------------------+
|  Creating the main application window  |
+---------------------------------------*/
 main_win = XCreateSimpleWindow(display, RootWindow(display, scr_num),
                                0, 0, MAIN_WDTH, MAIN_HGHT, 4,
                                BlackPixel(display, scr_num), 
                                WhitePixel(display, scr_num) );
 XGetWindowAttributes(display, main_win, &main_win_attr);
 main_wdth=main_win_attr.width;
 main_hght=main_win_attr.height;

/*------------------------------------+
|  Preparing an icon an a background  |
+------------------------------------*/
 {
#define icon_width 57
#define icon_height 57
static unsigned char icon_bits[] = {
   0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x00, 0x00, 0x00, 0x00, 0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f, 0x00,
   0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0xe4, 0xff, 0xd7, 0xff,
   0xaf, 0xff, 0x4f, 0x00, 0xd4, 0xff, 0xbb, 0xff, 0x77, 0xff, 0x57, 0x00,
   0xb4, 0xff, 0xbb, 0xff, 0x7b, 0xff, 0x5b, 0x00, 0x74, 0xff, 0x7d, 0xff,
   0xfb, 0xfe, 0x5d, 0x00, 0xf4, 0xfc, 0xfd, 0xfe, 0xfd, 0xfe, 0x5e, 0x00,
   0xf4, 0xfb, 0xfe, 0xfe, 0xfe, 0x7d, 0x5f, 0x00, 0xf4, 0xf7, 0xfe, 0x7d,
   0xff, 0xbd, 0x5f, 0x00, 0xf4, 0x6f, 0xff, 0x7b, 0xff, 0xdb, 0x5f, 0x00,
   0xf4, 0x5f, 0xff, 0xbb, 0xff, 0xeb, 0x5f, 0x00, 0xf4, 0x1f, 0xfc, 0xd7,
   0xff, 0xf0, 0x5f, 0x00, 0xf4, 0xa7, 0x03, 0xd6, 0x00, 0xc7, 0x5f, 0x00,
   0xf4, 0x79, 0xff, 0x01, 0xff, 0xbb, 0x5f, 0x00, 0x74, 0x7e, 0xff, 0xd7,
   0xff, 0x7b, 0x5e, 0x00, 0x94, 0x7f, 0xff, 0xbb, 0xff, 0xfb, 0x59, 0x00,
   0xe4, 0x7f, 0xff, 0xbb, 0xff, 0xfd, 0x57, 0x00, 0xf4, 0xff, 0xfe, 0x7d,
   0xff, 0xfd, 0x4f, 0x00, 0xe4, 0xff, 0xfe, 0xfe, 0xfe, 0xfd, 0x5f, 0x00,
   0x94, 0xff, 0x7e, 0xff, 0xfd, 0xfe, 0x47, 0x00, 0x74, 0xfe, 0x7e, 0xff,
   0xfd, 0xfe, 0x59, 0x00, 0xf4, 0xfd, 0xbd, 0xff, 0xfb, 0x7e, 0x5e, 0x00,
   0xf4, 0xf3, 0xdd, 0xff, 0x77, 0x9f, 0x5f, 0x00, 0xf4, 0xcf, 0xed, 0xff,
   0x6f, 0xe7, 0x5f, 0x00, 0xf4, 0x3f, 0xed, 0xff, 0x6f, 0xf9, 0x5f, 0x00,
   0xf4, 0xff, 0xf0, 0xff, 0x1f, 0xfe, 0x5f, 0x00, 0xf4, 0xff, 0x03, 0x00,
   0x80, 0xff, 0x5f, 0x00, 0xf4, 0xff, 0xf0, 0xff, 0x1f, 0xfe, 0x5f, 0x00,
   0xf4, 0x3f, 0xf5, 0xff, 0x5f, 0xf9, 0x5f, 0x00, 0xf4, 0xcf, 0xed, 0xff,
   0x6f, 0xe7, 0x5f, 0x00, 0xf4, 0xf3, 0xdd, 0xff, 0x77, 0xdf, 0x5f, 0x00,
   0xf4, 0xfd, 0xbd, 0xff, 0x7b, 0x3f, 0x5f, 0x00, 0x74, 0xfe, 0xbe, 0xff,
   0xfb, 0xfe, 0x5c, 0x00, 0x94, 0xff, 0x7e, 0xff, 0xfd, 0xfe, 0x53, 0x00,
   0xe4, 0xff, 0xfe, 0xfe, 0xfe, 0xfe, 0x4f, 0x00, 0xf4, 0xff, 0xfe, 0xfe,
   0xfe, 0xfe, 0x5f, 0x00, 0xe4, 0x7f, 0xff, 0x7d, 0xff, 0xfd, 0x47, 0x00,
   0x94, 0x7f, 0xff, 0xbb, 0xff, 0xfd, 0x59, 0x00, 0x74, 0x7e, 0xff, 0xd7,
   0xff, 0x7d, 0x5e, 0x00, 0xf4, 0x79, 0xff, 0xd7, 0xff, 0x9d, 0x5f, 0x00,
   0xf4, 0xa7, 0x1f, 0x00, 0xf8, 0xe3, 0x5f, 0x00, 0xf4, 0x1f, 0xe0, 0xd7,
   0x07, 0xf8, 0x5f, 0x00, 0xf4, 0x9f, 0xff, 0xbb, 0xff, 0xf3, 0x5f, 0x00,
   0xf4, 0x6f, 0xff, 0xbb, 0xff, 0xed, 0x5f, 0x00, 0xf4, 0x77, 0xff, 0x7d,
   0xff, 0xdd, 0x5f, 0x00, 0xf4, 0xfb, 0xfe, 0xfe, 0xfe, 0xbe, 0x5f, 0x00,
   0xf4, 0xfc, 0x7e, 0xff, 0xfd, 0x7e, 0x5f, 0x00, 0x74, 0xff, 0x7d, 0xff,
   0x7d, 0xff, 0x5c, 0x00, 0xb4, 0xff, 0xbd, 0xff, 0x7b, 0xff, 0x5b, 0x00,
   0xd4, 0xff, 0xdb, 0xff, 0xb7, 0xff, 0x57, 0x00, 0xe4, 0xff, 0xeb, 0xff,
   0xaf, 0xff, 0x4f, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00,
   0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f, 0x00, 0x00, 0x00, 0x00, 0x00,
   0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00};

  #define back_width 4
  #define back_height 4
  static unsigned char back_bits[] = {0x05, 0x0a, 0x05, 0x0a};
 
  ico_pixm = XCreateBitmapFromData(display, main_win,
                                    icon_bits, icon_width, icon_height);

  back_pixm = XCreatePixmapFromBitmapData(display, main_win, 
                back_bits, back_width, back_height,
                BlackPixel(display, scr_num), WhitePixel(display, scr_num),
                main_win_attr.depth);
 }

/*----------------------------+
|  Set the window background  |
+----------------------------*/
 XSetWindowBackgroundPixmap(display, main_win, back_pixm);

/*------------------------+
|  Set window properties  |
+------------------------*/
 {
  XSizeHints    size_hints;
  XWMHints      wm_hints;
  XClassHint    class_hints;
  XTextProperty winname, iconame;

  char *win_name="ShowMesh 1.0";
  char *ico_name="ShowMesh";

  size_hints.flags      = PPosition | PSize | PMinSize;   
  size_hints.min_width  = 50; /* 300 */
  size_hints.min_height = NBUTTONS*(BUTTON_HEIGHT+12)+10;  
  
  wm_hints.initial_state = NormalState;  /* Normal or Iconified     */
  wm_hints.input         = True;         /* It needs keyboard input */
  wm_hints.icon_pixmap   = ico_pixm;
  wm_hints.flags         = StateHint | IconPixmapHint | InputHint;
  
  class_hints.res_name  = prog_name;
  class_hints.res_class = "ShowMesh";

  if( XStringListToTextProperty(&win_name,1,&winname) == 0)
   {fprintf(stderr, "%s: structure allocation for window name failed\n",
                     prog_name);
    exit(-1);}

  if( XStringListToTextProperty(&ico_name,1,&iconame) == 0)
   {fprintf(stderr, "%s: structure allocation for icon name failed\n",
                     prog_name);
    exit(-1);}

  XSetWMProperties(display, main_win, 
                   &winname, &iconame,
                   argv, argc,
                   &size_hints, &wm_hints, &class_hints);

 }

/*---------------------+
|  Select input types  |
+---------------------*/
 XSelectInput(display, main_win, ExposureMask | StructureNotifyMask );
              
 load_fonts();
 getGC(main_win);

/*---------------------+
|  Display the window  |
+---------------------*/ 
 XMapWindow(display, main_win);

/*--------------------+
|  Create the button  |
+--------------------*/
 create_buttons(main_win);

/*----------------------------+
|  Create the drawing window  |
+----------------------------*/
  draw_wdth = main_wdth-BUTTON_WIDTH-40;
  draw_hght = main_hght-20;
  draw_win = XCreateSimpleWindow(display, main_win,
                                 10, 10, draw_wdth, draw_hght, 3, 
                                 BlackPixel(display, scr_num), 
                                 WhitePixel(display, scr_num) );

 XSelectInput(display, draw_win, ExposureMask | PointerMotionMask | ButtonPressMask);
 XMapWindow(display, draw_win);

 scl =min( (0.9*(double)draw_hght)/(ymax-ymin+SMALL), 
           (0.9*(double)draw_wdth)/(xmax-xmin+SMALL) );
 X0 = (int)(xmax-xmin)*scl*0.05;
 Y0 = draw_hght-(int)(ymax-ymin)*scl*0.05;

}
/*------------------------------------------------------------------------*/


/*========================================================================*/
void main(int argc, char *argv[])
{
 XEvent      report;

 if(argc!=2)
  {printf("\n*********************************************************");
   printf("\n****************                        *****************");
   printf("\n****************   PROGRAM:  ShowMesh   *****************");
   printf("\n****************                        *****************");
   printf("\n****************      version  1.0      *****************");
   printf("\n****************                        *****************");
   printf("\n****************  Author: Bojan NICENO  *****************");
   printf("\n**************** niceno@univ.trieste.it *****************");
   printf("\n****************                        *****************");
   printf("\n*********************************************************");
   printf("\n\nUsage:  ShowMesh  <NAME>");
   printf("\n\nShowMesh uses the following three input files:");
   printf("\n  NAME.n");
   printf("\n  NAME.e");
   printf("\n  NAME.s");
   printf("\nThese files are created with EasyMesh.\n\n");
   exit(-1);}
 
/*-------------------------------+
|  Copy the name of the program  |
+-------------------------------*/
 prog_name=argv[0];
 file_name=argv[1];

 load_mesh();
 
 init(argc, argv);

/*============#
#  MAIN LOOP  #
#============*/
 {
  Window        root, child; /* for XQuerryPointer */
  unsigned int  mouse_butt; 
  int           b, x0=OFF, y0, x_root, y_root, x_new, y_new, x_old, y_old;
  double        x0_fiz, y0_fiz, scl_new; 
 
   while(1)
    {
     XNextEvent(display, &report);
     switch(report.type)
      {
/******************
*  Expose Window  *
******************/
       case Expose: 
         if(report.xany.window==draw_win) 
           draw_mesh(draw_win);        
         for(b=0; b<NBUTTONS; b++)
           if(report.xany.window==button[b]) 
             write_on_button(b);
         break;

/******************
*  Resize Window  *
******************/
       case ConfigureNotify: 
         main_wdth = report.xconfigure.width;
         main_hght = report.xconfigure.height;
         draw_wdth = main_wdth-BUTTON_WIDTH-40;
         draw_hght = main_hght-20;
         XResizeWindow(display, draw_win, draw_wdth, draw_hght);  
         break;

/*****************
*  Button Press  *
*****************/
       case ButtonPress: 
         if(report.xany.window==button[QUIT])
          {XUnloadFont(display, text_font->fid);
           XFreeGC(display, gc_WoB); 
           XFreeGC(display, gc_BoW); 
           XCloseDisplay(display); 
           exit(1);}

         for(b=0; b<NBUTTONS; b++)
          {if(report.xany.window==button[b])
            {if(butt_data[b].pressed==ON) 
              {butt_data[b].pressed=OFF; write_on_button(b);}
             else                         
              {butt_data[b].pressed=ON;  write_on_button(b);}
          }}

         if(report.xany.window==button[MATERIALS] && butt_data[MATERIALS].pressed==ON)
          {butt_data[ELEMENTS].pressed=OFF; write_on_button(ELEMENTS);}

         if(report.xany.window==button[ELEMENTS] && butt_data[ELEMENTS].pressed==ON)
          {butt_data[MATERIALS].pressed=OFF; write_on_button(MATERIALS);}

         if(butt_data[DELAUNAY].pressed==OFF) 
          {butt_data[ELEMENTS].pressed=OFF; butt_data[MATERIALS].pressed=OFF;
           write_on_button(ELEMENTS);       write_on_button(MATERIALS);} 

         if(butt_data[DELAUNAY].pressed==OFF && butt_data[VORONOI].pressed==OFF) 
          {butt_data[NODES].pressed=OFF; write_on_button(NODES);}
        
         if(report.xany.window==button[FIT])
          {butt_data[FIT].pressed=OFF;
           scl =min( (0.9*(double)draw_hght)/(ymax-ymin+SMALL), 
                     (0.9*(double)draw_wdth)/(xmax-xmin+SMALL) );
           X0 = (int)(xmax-xmin)*scl*0.05;
           Y0 = draw_hght-(int)(ymax-ymin)*scl*0.05;
           XClearWindow(display, draw_win); 
           draw_mesh(draw_win);
           write_on_button(FIT);}

         if(report.xany.window>=button[DELAUNAY] && report.xany.window<=button[BOUNDARY])
          {XClearWindow(display, draw_win); 
           draw_mesh(draw_win);}

         if(report.xany.window==draw_win)
          {
           if(butt_data[MOVE].pressed==ON)
            {
             if(x0==OFF) {x0=report.xmotion.x; y0=report.xmotion.y;}
             else {butt_data[MOVE].pressed=OFF; 
                   write_on_button(MOVE);
                   X0+=(x_new-x0); Y0+=(y_new-y0);   x0=OFF;
                   XClearWindow(display, draw_win); draw_mesh(draw_win);} 
            }  
           if(butt_data[ZOOM].pressed==ON)
            {
             if(x0==OFF) {x0=report.xmotion.x; y0=report.xmotion.y;}
             else {butt_data[ZOOM].pressed=OFF;
                   write_on_button(ZOOM);  
                   x0_fiz = (min(x0, x_new)-X0)/scl; 
                   y0_fiz = (min(y0, y_new)-Y0)/scl; 
                   if(x0!=x_new && y0!=y_new)
                     scl_new=( min( (double)draw_wdth/abs(x0-x_new), 
                                    (double)draw_hght/abs(y0-y_new) ) )*scl;
                   if( max(scl_new*xmax, scl_new*ymax) < 32768 )
                    {scl=scl_new;
                     X0=-x0_fiz*scl; Y0=-y0_fiz*scl;}   
                   x0=OFF;
                   XClearWindow(display, draw_win); draw_mesh(draw_win);} 
            }  
          }  
         break;

/*****************
*  Mouse Motion  *
*****************/
       case MotionNotify:
         if(report.xany.window==draw_win)
          {
           if(butt_data[MOVE].pressed==ON)
            {x_old=x_new; y_old=y_new;
             XQueryPointer(display, report.xmotion.window,
                           &root, &child, &x_root, &y_root, &x_new, &y_new, 
                           &mouse_butt); 
             if(x0!=OFF) 
              {XDrawLine(display, draw_win, gc_XOR, x0, y0, x_old, y_old);
               XDrawLine(display, draw_win, gc_XOR, x0, y0, x_new, y_new);}}

           if(butt_data[ZOOM].pressed==ON)
            {x_old=x_new; y_old=y_new;
             XQueryPointer(display, report.xmotion.window,
                           &root, &child, &x_root, &y_root, &x_new, &y_new, 
                           &mouse_butt); 
             if(x0!=OFF) 
              {XDrawRectangle(display, draw_win, gc_XOR, 
               min(x0, x_old), min(y0, y_old), abs(x_old-x0), abs(y_old-y0));
               XDrawRectangle(display, draw_win, gc_XOR, 
               min(x0, x_new), min(y0, y_new), abs(x_new-x0), abs(y_new-y0));}}
          }
         break;
      } /* end switch */
    } /* end while */
  } /* end of block */

}