/*%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%        o88��88o  o88��88o  o88��88o  888  888       %%%%%%%%%%%%
%%%%%%%%%%%%        888  888  ���  888  888  ���  888  888       %%%%%%%%%%%%
%%%%%%%%%%%%        888  888       888  888       888  888       %%%%%%%%%%%%
%%%%%%%%%%%%        8888888�  o8888888  �888888o  888  888       %%%%%%%%%%%%
%%%%%%%%%%%%        888       888  888       888  888  888       %%%%%%%%%%%%
%%%%%%%%%%%%        888  ooo  888  888  ooo  888  �8888888       %%%%%%%%%%%%
%%%%%%%%%%%%        �88oo88�  �8888�88  �88oo88�       888       %%%%%%%%%%%%
%%%%%%%%%%%%                                      ooo  888       %%%%%%%%%%%%
%%%%%%%%%%%%                                      �888888�       %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                         888         %%%%%%%%%%%%
%%%%%%%%%%%%                                         888         %%%%%%%%%%%%
%%%%%%%%%%%%                                         888         %%%%%%%%%%%%
%%%%%%%%%%%%      8888888o888o   o88��88o  o88��88o  888��88o    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  88o  888  888  888  ���  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  888  888  888       888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  8888888�  �888888o  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  888            888  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  888  ooo  ooo  888  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%      888  888  888  �88oo88�  �88oo88�  888  888    %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%                                                     %%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%                        %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%  Author: Bojan NICENO  %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%                        %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%% niceno@univ.trieste.it %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%                        %%%%%%%%%%%%%%%%%%%%%%%%%%%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%*/
#define GRAPHICS OFF

#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifndef max
#define max(a,b)  (((a) > (b)) ? (a) : (b))
#endif
#ifndef min
#define min(a,b)  (((a) < (b)) ? (a) : (b))
#endif
#ifndef PI
#define PI    3.14159265359
#endif

#define SMALL 1e-30
#define GREAT 1e+30

#define ON      0 
#define OFF    -1       /* element is switched off */
#define WAIT   -2       /* node is waiting (it is a corner node) */
#define A       3
#define D       4
#define W       5

#define MAX_NODES 200000

/*-----------------------------+
|  definitions for the chains  |
+-----------------------------*/
#define CLOSED 0
#define OPEN   1
#define INSIDE 2


struct ele
 {
  int i,  j,  k;
  int ei, ej, ek;
  int si, sj, sk;

  int mark;             /* is it off (ON or OFF) */
  int state;            /* is it (D)one, (A)ctive or (W)aiting */
  int material;

  double xv, yv, xin, yin, R, r, Det;

  int new_numb;         /* used for renumeration */
 }
elem[MAX_NODES*2];


struct sid
 {
  int ea, eb;           /* left and right element */
  int a, b, c, d;       /* left, right, start and end point */

  int mark;             /* is it off, is on the boundary */

  double s;

  int new_numb;         /* used for renumeration */
 }
side[MAX_NODES*3];


struct nod
 {
  double x, y, F;
			
  double sumx, sumy;
  int    Nne;

  int mark;             /* is it off */

  int next;             /* next node in the boundary chain */
  int chain;            /* on which chains is the node */
  int inserted;

  int new_numb;         /* used for renumeration */
 }
node[MAX_NODES], point[MAX_NODES/2];


 struct seg
  {int n0, n1;
   int N; int chain; int bound; int mark;}
 *segment;
 
 struct chai
  {int s0, s1, type;}
 *chain;


int Ne, Nn, Ns, Nc;             /* number of: elements, nodes, sides */
int ugly;                       /* mora li biti globalna ??? */


/*=========================================================================*/
double area(struct nod *na, struct nod *nb, struct nod *nc)
{
 return 0.5 * (   ((*nb).x-(*na).x)*((*nc).y-(*na).y) 
		- ((*nb).y-(*na).y)*((*nc).x-(*na).x));
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
double dist(struct nod *na, struct nod *nb)
{
 return sqrt(   ((*nb).x-(*na).x)*((*nb).x-(*na).x)
	      + ((*nb).y-(*na).y)*((*nb).y-(*na).y) );
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
in_elem(struct nod *n)
{
 int e;
 
 for(e=0; e<Ne; e++)    /* This must search through all elements ?? */
  {
   if(    area(n, &node[elem[e].i], &node[elem[e].j]) >= 0.0
       && area(n, &node[elem[e].j], &node[elem[e].k]) >= 0.0
       && area(n, &node[elem[e].k], &node[elem[e].i]) >= 0.0 )
   
   break;
  }
 return e;
}
/*-in_elem-----------------------------------------------------------------*/


/*=========================================================================*/
bowyer(int n, int spac)
{
 int e, i, s, swap;
 struct nod vor;

 do  
  { 
   swap=0;
   for(s=0; s<Ns; s++)
   if(side[s].mark==0)
/* if( !( (node[side[s].c].inserted>1 && node[side[s].d].bound==OFF && side[s].s<(node[side[s].c].F+node[side[s].d].F) ) ||
	  (node[side[s].d].inserted>1 && node[side[s].c].bound==OFF && side[s].s<(node[side[s].c].F+node[side[s].d].F) ) ) ) */
    {
     if(side[s].a==n)
      {e=side[s].eb; 
       if(e!=OFF)
	{vor.x=elem[e].xv; 
	 vor.y=elem[e].yv;
	 if( dist(&vor, &node[n]) < elem[e].R )
	  {swap_side(s); swap=1;}}}
   
     else if(side[s].b==n)
      {e=side[s].ea; 
       if(e!=OFF)
	{vor.x=elem[e].xv; 
	 vor.y=elem[e].yv;
	 if( dist(&vor, &node[n]) < elem[e].R )
	  {swap_side(s); swap=1;}}}
    }
  }
 while(swap==1);

}
/*-bowyer------------------------------------------------------------------*/


/*=========================================================================*/
circles(int e)
/*---------------------------------------------------+
|  This function calculates radii of inscribed and   |
|  circumscribed circle for a given element (int e)  |
+---------------------------------------------------*/
{
 double x, y, xi, yi, xj, yj, xk, yk, xij, yij, xjk, yjk, num, den;
 double si, sj, sk, O;

 xi=node[elem[e].i].x; yi=node[elem[e].i].y;
 xj=node[elem[e].j].x; yj=node[elem[e].j].y;
 xk=node[elem[e].k].x; yk=node[elem[e].k].y;
   
 xij=0.5*(xi+xj); yij=0.5*(yi+yj);
 xjk=0.5*(xj+xk); yjk=0.5*(yj+yk);

 num = (xij-xjk)*(xj-xi) + (yij-yjk)*(yj-yi);
 den = (xj -xi) *(yk-yj) - (xk -xj) *(yj-yi);

 if(den>0)
  {
   elem[e].xv = x = xjk + num/den*(yk-yj);
   elem[e].yv = y = yjk - num/den*(xk-xj);

   elem[e].R  = sqrt( (xi-x)*(xi-x) + (yi-y)*(yi-y) );
  }

 si=side[elem[e].si].s;
 sj=side[elem[e].sj].s;
 sk=side[elem[e].sk].s;
 O =si+sj+sk;
 elem[e].Det = xi*(yj-yk) - xj*(yi-yk) + xk*(yi-yj);

 elem[e].xin = ( xi*si + xj*sj + xk*sk ) / O;
 elem[e].yin = ( yi*si + yj*sj + yk*sk ) / O;

 elem[e].r   = elem[e].Det / O;
}
/*-circles-----------------------------------------------------------------*/


/*=========================================================================*/
spacing(int e, int n)
/*----------------------------------------------------------------+
|  This function calculates the value of the spacing function in  |
|  a new node 'n' which is inserted in element 'e' by a linear    |
|  approximation from the values of the spacing function in the   |
|  elements nodes.                                                |
+----------------------------------------------------------------*/
{
 double dxji, dxki, dyji, dyki, dx_i, dy_i, det, a, b;

 dxji = node[elem[e].j].x - node[elem[e].i].x;
 dyji = node[elem[e].j].y - node[elem[e].i].y;
 dxki = node[elem[e].k].x - node[elem[e].i].x;
 dyki = node[elem[e].k].y - node[elem[e].i].y;
 dx_i = node[n].x - node[elem[e].i].x;
 dy_i = node[n].y - node[elem[e].i].y;

 det = dxji*dyki - dxki*dyji;

 a = (+ dyki*dx_i - dxki*dy_i)/det;
 b = (- dyji*dx_i + dxji*dy_i)/det;

 node[n].F = node[elem[e].i].F + 
	     a*(node[elem[e].j].F - node[elem[e].i].F) +
	     b*(node[elem[e].k].F - node[elem[e].i].F);
}
/*-spacing-----------------------------------------------------------------*/


/*=========================================================================*/
insert_node(double x, double y, int spac,
	 int prev_n, int prev_s_mark, int mark, int next_s_mark, int next_n)
{
 int    i,j,k,en, n, e,ei,ej,ek, s,si,sj,sk;
 double sx, sy;

 Nn++;          /* one new node */
 
 node[Nn-1].x = x;
 node[Nn-1].y = y;
 node[Nn-1].mark = mark;

/* find the element which contains new node */ 
 e = in_elem(&node[Nn-1]);

/* calculate the spacing function in the new node */
 if(spac==ON)
   spacing(e, Nn-1);

 i =elem[e].i;  j =elem[e].j;  k =elem[e].k;
 ei=elem[e].ei; ej=elem[e].ej; ek=elem[e].ek; 
 si=elem[e].si; sj=elem[e].sj; sk=elem[e].sk; 
 
 Ne+=2;
 Ns+=3;

/*---------------+
|  new elements  |
+---------------*/ 
 elem[Ne-2].i=Nn-1;  elem[Ne-2].j=k;     elem[Ne-2].k=i;
 elem[Ne-1].i=Nn-1;  elem[Ne-1].j=i;     elem[Ne-1].k=j; 
 
 elem[Ne-2].ei=ej;   elem[Ne-2].ej=Ne-1; elem[Ne-2].ek=e;
 elem[Ne-1].ei=ek;   elem[Ne-1].ej=e;    elem[Ne-1].ek=Ne-2;
 
 elem[Ne-2].si=sj;   elem[Ne-2].sj=Ns-2; elem[Ne-2].sk=Ns-3;
 elem[Ne-1].si=sk;   elem[Ne-1].sj=Ns-1; elem[Ne-1].sk=Ns-2;
 
/*------------+ 
|  new sides  |
+------------*/ 
 side[Ns-3].c =k;    side[Ns-3].d =Nn-1;     /* c-d */
 side[Ns-3].a =j;    side[Ns-3].b =i;        /* a-b */
 side[Ns-3].ea=e;    side[Ns-3].eb=Ne-2;
 
 side[Ns-2].c =i;    side[Ns-2].d =Nn-1;     /* c-d */
 side[Ns-2].a =k;    side[Ns-2].b =j;        /* a-b */
 side[Ns-2].ea=Ne-2; side[Ns-2].eb=Ne-1;
 
 side[Ns-1].c =j;    side[Ns-1].d =Nn-1;     /* c-d */
 side[Ns-1].a =i;    side[Ns-1].b =k;        /* a-b */
 side[Ns-1].ea=Ne-1; side[Ns-1].eb=e;       

 for(s=1; s<=3; s++)
  {sx = node[side[Ns-s].c].x - node[side[Ns-s].d].x;
   sy = node[side[Ns-s].c].y - node[side[Ns-s].d].y;
   side[Ns-s].s = sqrt(sx*sx+sy*sy);}

 elem[e].i  = Nn-1;
 elem[e].ej = Ne-2;
 elem[e].ek = Ne-1;
 elem[e].sj = Ns-3;
 elem[e].sk = Ns-1;

 if(side[si].a==i) {side[si].a=Nn-1; side[si].ea=e;}
 if(side[si].b==i) {side[si].b=Nn-1; side[si].eb=e;}
 
 if(side[sj].a==j) {side[sj].a=Nn-1; side[sj].ea=Ne-2;}
 if(side[sj].b==j) {side[sj].b=Nn-1; side[sj].eb=Ne-2;}
 
 if(side[sk].a==k) {side[sk].a=Nn-1; side[sk].ea=Ne-1;} 
 if(side[sk].b==k) {side[sk].b=Nn-1; side[sk].eb=Ne-1;} 

 if(ej!=-1)
  {if(elem[ej].ei==e) {elem[ej].ei=Ne-2;}
   if(elem[ej].ej==e) {elem[ej].ej=Ne-2;}
   if(elem[ej].ek==e) {elem[ej].ek=Ne-2;}}

 if(ek!=-1)
  {if(elem[ek].ei==e) {elem[ek].ei=Ne-1;}
   if(elem[ek].ej==e) {elem[ek].ej=Ne-1;}
   if(elem[ek].ek==e) {elem[ek].ek=Ne-1;}}

/* Find circumenters for two new elements, 
   and for the one who's segment has changed */
 circles(e);
 circles(Ne-2);
 circles(Ne-1);

 bowyer(Nn-1, spac);

/*-------------------------------------------------+
|  NEW ! Insert boundary conditions for the sides  |
+-------------------------------------------------*/
 for(s=3; s<Ns; s++)
  {
   if(side[s].c==prev_n && side[s].d==Nn-1)  side[s].mark=prev_s_mark;
   if(side[s].d==prev_n && side[s].c==Nn-1)  side[s].mark=prev_s_mark;
   if(side[s].c==next_n && side[s].d==Nn-1)  side[s].mark=next_s_mark;
   if(side[s].d==next_n && side[s].c==Nn-1)  side[s].mark=next_s_mark;
  }

 return e;
}
/*-insert_node-------------------------------------------------------------*/


/*=========================================================================*/
swap_side(int s)
{
 int    a, b, c, d, ea, eb, eac, ead, ebc, ebd, sad, sac, sbc, sbd;
 double sx, sy;
 
 ea=side[s].ea; 
 eb=side[s].eb;
 a=side[s].a; b=side[s].b; c=side[s].c; d=side[s].d;

 if(elem[ea].ei==eb) {ead=elem[ea].ej; eac=elem[ea].ek; 
		      sad=elem[ea].sj; sac=elem[ea].sk;}
 
 if(elem[ea].ej==eb) {ead=elem[ea].ek; eac=elem[ea].ei; 
		      sad=elem[ea].sk; sac=elem[ea].si;}   
 
 if(elem[ea].ek==eb) {ead=elem[ea].ei; eac=elem[ea].ej;
		      sad=elem[ea].si; sac=elem[ea].sj;}

 if(elem[eb].ei==ea) {ebc=elem[eb].ej; ebd=elem[eb].ek;
		      sbc=elem[eb].sj; sbd=elem[eb].sk;}

 if(elem[eb].ej==ea) {ebc=elem[eb].ek; ebd=elem[eb].ei;
		      sbc=elem[eb].sk; sbd=elem[eb].si;}
 
 if(elem[eb].ek==ea) {ebc=elem[eb].ei; ebd=elem[eb].ej;
		      sbc=elem[eb].si; sbd=elem[eb].sj;}

 elem[ea].i =a;   elem[ea].j =b;   elem[ea].k =d;  
 elem[ea].ei=ebd; elem[ea].ej=ead; elem[ea].ek=eb;  
 elem[ea].si=sbd; elem[ea].sj=sad; elem[ea].sk=s;  
  
 elem[eb].i =a;   elem[eb].j =c;   elem[eb].k =b;  
 elem[eb].ei=ebc; elem[eb].ej=ea;  elem[eb].ek=eac;  
 elem[eb].si=sbc; elem[eb].sj=s;   elem[eb].sk=sac;  

 if(eac!=-1)
  {
   if(elem[eac].ei==ea) elem[eac].ei=eb;
   if(elem[eac].ej==ea) elem[eac].ej=eb;
   if(elem[eac].ek==ea) elem[eac].ek=eb; 
  }
 
 if(ebd!=-1)
  {
   if(elem[ebd].ei==eb) elem[ebd].ei=ea;
   if(elem[ebd].ej==eb) elem[ebd].ej=ea;
   if(elem[ebd].ek==eb) elem[ebd].ek=ea; 
  }
 
 if(side[sad].ea==ea) {side[sad].a=b;}
 if(side[sad].eb==ea) {side[sad].b=b;}

 if(side[sbc].ea==eb) {side[sbc].a=a;}
 if(side[sbc].eb==eb) {side[sbc].b=a;}

 if(side[sbd].ea==eb) {side[sbd].ea=ea; side[sbd].a=a;}
 if(side[sbd].eb==eb) {side[sbd].eb=ea; side[sbd].b=a;}
 
 if(a<b)
  {side[s].c=a; side[s].d=b; side[s].a=d; side[s].b=c;
   side[s].ea=ea; side[s].eb=eb;}
 else 
  {side[s].c=b; side[s].d=a; side[s].a=c; side[s].b=d;
   side[s].ea=eb; side[s].eb=ea;}

 sx = node[side[s].c].x - node[side[s].d].x;
 sy = node[side[s].c].y - node[side[s].d].y;
 side[s].s = sqrt(sx*sx+sy*sy);

 if(side[sac].ea==ea) {side[sac].ea=eb; side[sac].a=b;}
 if(side[sac].eb==ea) {side[sac].eb=eb; side[sac].b=b;}
 
 if(side[sad].ea==ea) {side[sad].a=b;}
 if(side[sad].eb==ea) {side[sad].b=b;}

 if(side[sbc].ea==eb) {side[sbc].a=a;}
 if(side[sbc].eb==eb) {side[sbc].b=a;}

 if(side[sbd].ea==eb) {side[sbd].ea=ea; side[sbd].a=a;}
 if(side[sbd].eb==eb) {side[sbd].eb=ea; side[sbd].b=a;}

 circles(ea);
 circles(eb);
}
/*-swap_side---------------------------------------------------------------*/


/*=========================================================================*/
void erase()
{
 int s, n, e;

 int a, b, c, d, ea, eb;

/*--------------------------+
|                           |
|  Negative area check for  |
|  elimination of elements  |
|                           |
+--------------------------*/
 for(e=0; e<Ne; e++)
   if( (node[elem[e].i].chain==node[elem[e].j].chain) &&
       (node[elem[e].j].chain==node[elem[e].k].chain) &&
       (chain[node[elem[e].i].chain].type==CLOSED) )
  {
   a = min( min(elem[e].i, elem[e].j), elem[e].k );
   c = max( max(elem[e].i, elem[e].j), elem[e].k );
   b = elem[e].i+elem[e].j+elem[e].k - a - c;

   if(a<3)
     elem[e].mark=OFF;

   else if(area(&node[a], &node[b], &node[c]) < 0.0)
     elem[e].mark=OFF;
  }

 for(e=0; e<Ne; e++)
  {if(elem[elem[e].ei].mark==OFF) elem[e].ei=OFF;
   if(elem[elem[e].ej].mark==OFF) elem[e].ej=OFF;
   if(elem[elem[e].ek].mark==OFF) elem[e].ek=OFF;}

/*-----------------------+
|                        |
|  Elimination of sides  |
|                        |
+-----------------------*/
 for(s=0; s< 3; s++)
   side[s].mark=OFF;

 for(s=3; s<Ns; s++)
   if( (elem[side[s].ea].mark==OFF) && (elem[side[s].eb].mark==OFF) )
     side[s].mark=OFF;

 for(s=3; s<Ns; s++)
   if(side[s].mark!=OFF)
    {
     if(elem[side[s].ea].mark==OFF) {side[s].ea=OFF; side[s].a=OFF;}
     if(elem[side[s].eb].mark==OFF) {side[s].eb=OFF; side[s].b=OFF;}
    }

/*-----------------------+
|                        |
|  Elimination of nodes  |
|                        |
+-----------------------*/
 for(n=0; n< 3; n++)
   node[n].mark=OFF;

}
/*-erase-------------------------------------------------------------------*/


/*=========================================================================*/
diamond()
{
 int    ea, eb, eac, ead, ebc, ebd, s;
 
 for(s=0; s<Ns; s++)
   if(side[s].mark!=OFF)
    {
     ea=side[s].ea;
     eb=side[s].eb;

     if(elem[ea].ei==eb) {ead=elem[ea].ej; eac=elem[ea].ek;}
     if(elem[ea].ej==eb) {ead=elem[ea].ek; eac=elem[ea].ei;}   
     if(elem[ea].ek==eb) {ead=elem[ea].ei; eac=elem[ea].ej;}
     if(elem[eb].ei==ea) {ebc=elem[eb].ej; ebd=elem[eb].ek;}
     if(elem[eb].ej==ea) {ebc=elem[eb].ek; ebd=elem[eb].ei;}
     if(elem[eb].ek==ea) {ebc=elem[eb].ei; ebd=elem[eb].ej;}

     if( (eac==OFF || elem[eac].state==D) &&
	 (ebc==OFF || elem[ebc].state==D) &&
	 (ead==OFF || elem[ead].state==D) &&
	 (ebd==OFF || elem[ebd].state==D) )
      {
       elem[ea].state=D;
       elem[eb].state=D;
      }
    }
}
/*-diamond-----------------------------------------------------------------*/


/*=========================================================================*/
classify()
/*----------------------------------------------------------+
|  This function searches through all elements every time.  |
|  Some optimisation will definitely bee needed             |
|                                                           |
|  But it also must me noted, that this function defines    |
|  the strategy for insertion of new nodes                  |
|                                                           |
|  It's MUCH MUCH better when the ugliest element is found  |
|  as one with highest ratio of R/r !!! (before it was      |
|  element with greater R)                                  |
+----------------------------------------------------------*/
{
 int e, ei, ej, ek,si,sj,sk;
 double ratio=-GREAT, F;

 ugly=OFF;

 for(e=0; e<Ne; e++)
   if(elem[e].mark!=OFF)
    {
     ei=elem[e].ei; ej=elem[e].ej; ek=elem[e].ek;

     F=(node[elem[e].i].F + node[elem[e].j].F + node[elem[e].k].F)/3.0;

     elem[e].state=W;

/*--------------------------+
|  0.577 is ideal triangle  |
+--------------------------*/
     if(elem[e].R < 0.700*F) elem[e].state=D; /* 0.0866; 0.07 */

/*------------------------+
|  even this is possible  |
+------------------------*/
     if(ei!=OFF && ej!=OFF && ek!=OFF)
       if(elem[ei].state==D && elem[ej].state==D && elem[ek].state==D)
	 elem[e].state=D;
    }

/*--------------------------------------+
|  Diamond check. Is it so important ?  |
+--------------------------------------*/
   diamond();   

/*------------------------------------------------+
|  First part of the trick:                       |
|    search through the elements on the boundary  |
+------------------------------------------------*/
 for(e=0; e<Ne; e++)
   if(elem[e].mark!=OFF && elem[e].state!=D)
    {
     si=elem[e].si; sj=elem[e].sj; sk=elem[e].sk;

     if(side[si].mark!=0) elem[e].state=A;
     if(side[sj].mark!=0) elem[e].state=A;
     if(side[sk].mark!=0) elem[e].state=A;
  
     if(elem[e].state==A && elem[e].R/elem[e].r > ratio)
      {ratio=max(ratio, elem[e].R/elem[e].r);
       ugly=e;}
    }

/*-------------------------------------------------+
|  Second part of the trick:                       |
|    if non-acceptable element on the boundary is  |
|    found, ignore the elements inside the domain  |
+-------------------------------------------------*/
 if(ugly==OFF)
   for(e=0; e<Ne; e++)
     if(elem[e].mark!=OFF)
      {
       if(elem[e].state!=D)
	{
	 ei=elem[e].ei; ej=elem[e].ej; ek=elem[e].ek;

	 if(ei!=OFF)
	   if(elem[ei].state==D) elem[e].state=A;
  
	 if(ej!=OFF)
	   if(elem[ej].state==D) elem[e].state=A;
  
	 if(ek!=OFF)
	   if(elem[ek].state==D) elem[e].state=A;
  
	 if(elem[e].state==A && elem[e].R/elem[e].r > ratio)
	  {ratio=max(ratio, elem[e].R/elem[e].r);
	   ugly=e;}
	}
      }

}
/*-classify----------------------------------------------------------------*/


/*=========================================================================*/
new_node()
/*---------------------------------------------------+
|  This function is very important.                  |
|  It determines the position of the inserted node.  |
+---------------------------------------------------*/
{
 int    s=OFF, n, e;
 double xM, yM, xCa, yCa, p, px, py, q, qx, qy, rhoM, rho_M, d;

 struct nod Ca;

/*-------------------------------------------------------------------------+
|  It's obvious that elements which are near the boundary, will come into  |
|  play first.                                                             |
|                                                                          |
|  However, some attention has to be payed for the case when two accepted  |
|  elements surround the ugly one                                          |
|                                                                          |
|  What if new points falls outside the domain                             |
+-------------------------------------------------------------------------*/
 if(elem[elem[ugly].ei].state==D)    {s=elem[ugly].si; n=elem[ugly].i;}
 if(elem[elem[ugly].ej].state==D)    {s=elem[ugly].sj; n=elem[ugly].j;}
 if(elem[elem[ugly].ek].state==D)    {s=elem[ugly].sk; n=elem[ugly].k;}
 if(side[elem[ugly].si].mark > 0)    {s=elem[ugly].si; n=elem[ugly].i;}
 if(side[elem[ugly].sj].mark > 0)    {s=elem[ugly].sj; n=elem[ugly].j;}
 if(side[elem[ugly].sk].mark > 0)    {s=elem[ugly].sk; n=elem[ugly].k;}
 if(s==OFF) return;

 xM  = 0.5*(node[side[s].c].x + node[side[s].d].x);
 yM  = 0.5*(node[side[s].c].y + node[side[s].d].y);

 Ca.x = elem[ugly].xv;
 Ca.y = elem[ugly].yv;

 p  = 0.5*side[s].s;    /* not checked */

 qx = Ca.x-xM;
 qy = Ca.y-yM;
 q  = sqrt(qx*qx+qy*qy);

 rhoM = 0.577 *  0.5*(node[side[s].c].F + node[side[s].d].F);

 rho_M = min( max( rhoM, p), 0.5*(p*p+q*q)/q );

 if(rho_M < p) d=rho_M;
 else          d=rho_M+sqrt(rho_M*rho_M-p*p); 

/*---------------------------------------------------------------------+
|  The following line check can the new point fall outside the domain. |
|  However, I can't remember how it works, but I believe that it is    |
|  still a weak point of the code, particulary when there are lines    |
|  inside the domain                                                   |
+---------------------------------------------------------------------*/

 if( area(&node[side[s].c], &node[side[s].d], &Ca) *
     area(&node[side[s].c], &node[side[s].d], &node[n]) > 0.0 )
   insert_node(xM + d*qx/q,  yM + d*qy/q, ON, OFF, 0, 0, 0, OFF);
/*
 else
  {
   node[n].x = xM - d*qx/q;
   node[n].y = yM - d*qy/q;
   node[n].mark=6;   
   for(e=0; e<Ne; e++) 
     if(elem[e].i==n || elem[e].j==n || elem[e].k==n)
       circles(e);
  }
*/
 return;
}
/*-new_node----------------------------------------------------------------*/


/*=========================================================================*/
neighbours() 
/*--------------------------------------------------------------+
|  Counting the elements which surround each node.              |
|  It is important for the two functions: 'relax' and 'smooth'  |
+--------------------------------------------------------------*/
{ 
 int s;
 
 for(s=0; s<Ns; s++)
   if(side[s].mark==0)
    {
     if(node[side[s].c].mark==0)
       node[side[s].c].Nne++;
       
     if(node[side[s].d].mark==0)
       node[side[s].d].Nne++;
    }
}
/*-neighbours--------------------------------------------------------------*/


/*=========================================================================*/
materials()
{
 int e, c, mater, iter, over, s;
 int ei, ej, ek, si, sj, sk;

 for(e=0; e<Ne; e++)
   if(elem[e].mark!=OFF)   
     elem[e].material=OFF;

 for(c=0; c<Nc; c++)
  {
   if(point[c].inserted==0)
    {
     elem[in_elem(&point[c])].material=point[c].mark;
     mater=ON;
    }
  }

 if(mater==ON)
  {
   for(;;) 
    {      
     over=ON;

     for(e=0; e<Ne; e++)
       if(elem[e].mark!=OFF && elem[e].material==OFF)
	{
	 ei=elem[e].ei;
	 ej=elem[e].ej;
	 ek=elem[e].ek;

	 si=elem[e].si;
	 sj=elem[e].sj;
	 sk=elem[e].sk;

   
	 if(ei!=OFF)
	   if(elem[ei].material!=OFF && side[si].mark==0)
	    {
	     elem[e].material=elem[ei].material;
	     over=OFF;
	    }

	 if(ej!=OFF)
	   if(elem[ej].material!=OFF && side[sj].mark==0)
	    {
	     elem[e].material=elem[ej].material;
	     over=OFF;
	    }

	 if(ek!=OFF)
	   if(elem[ek].material!=OFF && side[sk].mark==0)
	    {
	     elem[e].material=elem[ek].material;
	     over=OFF;
	    }

	}

     if(over==ON) break;

    } /* for(iter) */

  }
}
/*-materials---------------------------------------------------------------*/


/*=========================================================================*/
relax()
{
 int s, T, E;
 
 for(T=6; T>=3; T--)
   for(s=0; s<Ns; s++)
     if(side[s].mark==0)
       if( (node[side[s].a].mark==0) &&
	   (node[side[s].b].mark==0) &&
	   (node[side[s].c].mark==0) &&
	   (node[side[s].d].mark==0) )
      {
       E =   node[side[s].c].Nne + node[side[s].d].Nne 
	   - node[side[s].a].Nne - node[side[s].b].Nne;

       if(E==T) 
	{node[side[s].a].Nne++; node[side[s].b].Nne++; 
	 node[side[s].c].Nne--; node[side[s].d].Nne--;  
	 swap_side(s);}
      }

}
/*-relax-------------------------------------------------------------------*/


/*=========================================================================*/
int smooth()
{
 int it, s, n, e;
 
 for(it=0; it<10; it++)
  {    
   for(s=0; s<Ns; s++)
     if(side[s].mark==0)
      {
       if(node[side[s].c].mark==0)
	{node[side[s].c].sumx += node[side[s].d].x;
	 node[side[s].c].sumy += node[side[s].d].y;}
       
       if(node[side[s].d].mark==0)
	{node[side[s].d].sumx += node[side[s].c].x;
	 node[side[s].d].sumy += node[side[s].c].y;}
      }
    
   for(n=0; n<Nn; n++)
     if(node[n].mark==0)
      {node[n].x=node[n].sumx/node[n].Nne; node[n].sumx=0.0;
       node[n].y=node[n].sumy/node[n].Nne; node[n].sumy=0.0;}
  }
       
 for(e=0; e<Ne; e++)
   if(elem[e].mark!=OFF)
     circles(e);

 return 0;
}
/*-smooth------------------------------------------------------------------*/


/*=========================================================================*/
renum()
{
 int n, o, s, e, e2, c, d, i, j, k;
 int new_elem=0, new_node=0, new_side=0, next_e, next_s, lowest;

 for(n=0; n<Nn; n++) node[n].new_numb=OFF;
 for(e=0; e<Ne; e++) elem[e].new_numb=OFF;
 for(s=0; s<Ns; s++) side[s].new_numb=OFF;

/*-------------------------------+
|  Searching the first element.  |
|  It is the first which is ON   |
+-------------------------------*/
 for(e=0; e<Ne; e++)
   if(elem[e].mark!=OFF)
     break;

/*----------------------------------------------------------+
|  Assigning numbers 0 and 1 to the nodes of first element  |
+----------------------------------------------------------*/
 node[elem[e].i].new_numb  = new_node; new_node++;
 node[elem[e].j].new_numb  = new_node; new_node++;

/*%%%%%%%%%%%%%%%%%%%%%%%%%
%                         %
%  Renumeration of nodes  %
%                         % 
%%%%%%%%%%%%%%%%%%%%%%%%%*/
 do
  {
   lowest = Nn+Nn;
   next_e = OFF;

   for(e=0; e<Ne; e++)
     if(elem[e].mark!=OFF && elem[e].new_numb==OFF)
      {
       i=node[elem[e].i].new_numb;
       j=node[elem[e].j].new_numb;
       k=node[elem[e].k].new_numb;

       if( i+j+k+2 == abs(i) + abs(j) + abs(k) )
	{
	 if( (i==OFF) && (j+k) < lowest) {next_e=e; lowest=j+k;}
	 if( (j==OFF) && (i+k) < lowest) {next_e=e; lowest=i+k;}
	 if( (k==OFF) && (i+j) < lowest) {next_e=e; lowest=i+j;}
	}
      }

   if(next_e!=OFF)
    {
     i=node[elem[next_e].i].new_numb;
     j=node[elem[next_e].j].new_numb;
     k=node[elem[next_e].k].new_numb;

/*----------------------------------+
|  Assign a new number to the node  |
+----------------------------------*/
     if(i==OFF) {node[elem[next_e].i].new_numb = new_node; new_node++;}
     if(j==OFF) {node[elem[next_e].j].new_numb = new_node; new_node++;}
     if(k==OFF) {node[elem[next_e].k].new_numb = new_node; new_node++;}
    }
  }
 while(next_e != OFF);

/*%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%                             %
%  Renumeration of triangles  %
%                             %
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%*/
 do
  {
   lowest = Nn+Nn+Nn;
   next_e = OFF;

   for(e=0; e<Ne; e++)
     if(elem[e].mark!=OFF && elem[e].new_numb==OFF)
      {
       i=node[elem[e].i].new_numb;
       j=node[elem[e].j].new_numb;
       k=node[elem[e].k].new_numb;

       if( (i+j+k) < lowest )
	{
	 next_e=e;
	 lowest=i+j+k;
	}
      }

   if(next_e!=OFF)
    {
     elem[next_e].new_numb=new_elem; new_elem++;
    }
  }
 while(next_e != OFF);



/*%%%%%%%%%%%%%%%%%%%%%%%%%
%                         %
%  Renumeration of sides  %
%                         %
%%%%%%%%%%%%%%%%%%%%%%%%%*/
 do
  {
   lowest = Nn+Nn;
   next_s = OFF;

   for(s=0; s<Ns; s++)
     if(side[s].mark!=OFF && side[s].new_numb==OFF)
      {
       c=node[side[s].c].new_numb;
       d=node[side[s].d].new_numb;

       if( (c+d) < lowest)
	{
	 lowest=c+d; next_s=s;
	}
      }

   if(next_s!=OFF)
    {
     side[next_s].new_numb=new_side;
     new_side++;
    }
      
  }
 while(next_s != OFF);

}
/*-renum-------------------------------------------------------------------*/


char name[128]; int len;


/*=========================================================================*/
load_i(FILE *in, int *numb)
{
 char dum, dummy[128];

 for(;;)
  {fscanf(in,"%s", dummy);
   if(dummy[0]=='#' && strlen(dummy)>1 && dummy[strlen(dummy)-1]=='#') {}
   else if(dummy[0]=='#') {do{fscanf(in,"%c", &dum);} while(dum!='#');}
   else                   {*numb=atoi(dummy); break;} }
}

load_d(FILE *in, double *numb)
{
 char dum, dummy[128];

 for(;;)
  {fscanf(in,"%s", dummy);
   if(dummy[0]=='#' && strlen(dummy)>1 && dummy[strlen(dummy)-1]=='#') {}
   else if(dummy[0]=='#') {do{fscanf(in,"%c", &dum);} while(dum!='#');}
   else                   {*numb=atof(dummy); break;} }
}

load_s(FILE *in, char *string)
{
 char dum, dummy[128];

 for(;;)
  {fscanf(in,"%s", dummy);
   if(dummy[0]=='#' && strlen(dummy)>1 && dummy[strlen(dummy)-1]=='#') {}
   else if(dummy[0]=='#') {do{fscanf(in,"%c", &dum);} while(dum!='#');}
   else                   {strcpy(string, dummy); break;} }
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
load()
{
 int  c, n, s, Fl, M, N0, chains, bound;
 char dummy[128];
 double xmax=-GREAT, xmin=+GREAT, ymax=-GREAT, ymin=+GREAT, xt, yt, gab;

 FILE *in;

 int m;
 double xO, yO, xN, yN, xC, yC, L, Lx, Ly, dLm, ddL, L_tot;
 
 int *inserted;

/*----------+
|           |
|  Loading  |
|           |
+----------*/
 if((in=fopen(name, "r"))==NULL)
  {fprintf(stderr, "Cannot load file %s !\n", name);
   return 1;}

 load_i(in, &Nc);
 inserted=(int *) calloc(Nc, sizeof(int)); 
 for(n=0; n<Nc; n++)
  {
   load_s(in, dummy);     
   load_d(in, &point[n].x);
   load_d(in, &point[n].y);
   load_d(in, &point[n].F);
   load_i(in, &point[n].mark); 

   xmax=max(xmax, point[n].x); ymax=max(ymax, point[n].y);
   xmin=min(xmin, point[n].x); ymin=min(ymin, point[n].y);

   point[n].inserted=0; /* it is only loaded */
  }

 load_i(in, &Fl);
 segment=(struct seg *) calloc(Fl+1, sizeof(struct seg));
 chain  =(struct chai *) calloc(Fl+1, sizeof(struct chai)); /* approximation */
 segment[Fl].n0=-1;
 segment[Fl].n1=-1;

 for(s=0; s<Fl; s++)
  {
   load_s(in, dummy);
   load_i(in, &segment[s].n0);
   load_i(in, &segment[s].n1);
   load_i(in, &segment[s].mark);
  }
 fclose(in);

/*----------------------+
   counting the chains
+----------------------*/
 chains=0;
 chain[chains].s0=0;
 for(s=0; s<Fl; s++)
  {
   point[segment[s].n0].inserted++;
   point[segment[s].n1].inserted++;

   segment[s].chain           = chains;

   if(segment[s].n1!=segment[s+1].n0)
    {chain[chains].s1=s;
     chains++;
     chain[chains].s0=s+1;}
  }

/*-------------------------------------+
   counting the nodes on each segment
+-------------------------------------*/
 for(s=0; s<Fl; s++)
  {
   xO=point[segment[s].n0].x; yO=point[segment[s].n0].y;
   xN=point[segment[s].n1].x; yN=point[segment[s].n1].y; 

   Lx=(xN-xO); Ly=(yN-yO); L=sqrt(Lx*Lx+Ly*Ly);

   if( (point[segment[s].n0].F+point[segment[s].n1].F > L ) &&
       (segment[s].n0 != segment[s].n1) )
    {point[segment[s].n0].F = min(point[segment[s].n0].F,L);
     point[segment[s].n1].F = min(point[segment[s].n1].F,L);}
  }

/*-------------------------------------+
   counting the nodes on each segment
+-------------------------------------*/
 for(s=0; s<Fl; s++)
  {
   xO=point[segment[s].n0].x; yO=point[segment[s].n0].y;
   xN=point[segment[s].n1].x; yN=point[segment[s].n1].y; 

   Lx=(xN-xO); Ly=(yN-yO); L=sqrt(Lx*Lx+Ly*Ly);

   if(point[segment[s].n1].F+point[segment[s].n0].F<=L)
    {dLm=0.5*(point[segment[s].n0].F+point[segment[s].n1].F);
     segment[s].N=ceil(L/dLm);}
   else
     segment[s].N=1;
  }


 for(n=0; n<chains; n++)
  {
   if( segment[chain[n].s0].n0 == segment[chain[n].s1].n1 )
    {chain[n].type=CLOSED;}

   if( segment[chain[n].s0].n0 != segment[chain[n].s1].n1 )
    {chain[n].type=OPEN;}

   if( (point[segment[chain[n].s0].n0].inserted==1) &&
       (point[segment[chain[n].s1].n1].inserted==1) )
    {chain[n].type=INSIDE;}
  }

/*------------+
|             |
|  Inserting  |
|             |
+------------*/
 xt = 0.5*(xmax+xmin);
 yt = 0.5*(ymax+ymin);

 gab=max((xmax-xmin),(ymax-ymin));
 
 Nn = 3;
 node[2].x = xt;                node[2].y = yt + 2.8*gab; 
 node[0].x = xt - 2.0*gab;      node[0].y = yt - 1.4*gab; 
 node[1].x = xt + 2.0*gab;      node[1].y = yt - 1.4*gab; 
 node[2].inserted=2;
 node[1].inserted=2;
 node[0].inserted=2;
/*
 node[2].type=;
 node[1].type=;
 node[0].type=;
*/
 node[2].next=1;
 node[1].next=0;
 node[0].next=2;

 Ne=1;
 elem[0].i =0;  elem[0].j = 1; elem[0].k = 2;
 elem[0].ei=-1; elem[0].ej=-1; elem[0].ek=-1;
 elem[0].si= 1; elem[0].sj= 2; elem[0].sk= 0;
 
 Ns=3;
 side[0].c=0; side[0].d=1; side[0].a=2; side[0].b=-1; 
 side[1].c=1; side[1].d=2; side[1].a=0; side[1].b=-1; 
 side[2].c=0; side[2].d=2; side[2].a=-1; side[2].b=1;  
 side[0].ea= 0; side[0].eb=-1;
 side[1].ea= 0; side[1].eb=-1; 
 side[2].ea=-1; side[2].eb= 0;


 for(n=0; n<Nc; n++)
   point[n].new_numb=OFF;

 for(c=0; c<chains; c++)
  {
   for(s=chain[c].s0; s<=chain[c].s1; s++)
    {
     xO=point[segment[s].n0].x; yO=point[segment[s].n0].y;
     xN=point[segment[s].n1].x; yN=point[segment[s].n1].y; 

/*===============
*  first point  *
===============*/
     if( point[segment[s].n0].new_numb == OFF )
      {
       if(s==chain[c].s0) /* first segment in the chain */
	 insert_node(xO, yO, OFF,
	 OFF,  OFF, point[segment[s].n0].mark, OFF, OFF);

       else if(s==chain[c].s1 && segment[s].N==1)
	 insert_node(xO, yO, OFF,
	 Nn-1, segment[s-1].mark,
	 point[segment[s].n0].mark,
	 segment[s].mark, point[segment[chain[c].s0].n0].new_numb);

       else
	{
	 insert_node(xO, yO, OFF,
	 Nn-1, segment[s-1].mark, point[segment[s].n0].mark, OFF, OFF);
	}

       node[Nn-1].next     = Nn;     /* Nn-1 is index of inserted node */
       node[Nn-1].chain    = segment[s].chain;
       node[Nn-1].F        = point[segment[s].n0].F;
       point[segment[s].n0].new_numb=Nn-1;
      }

     Lx=(xN-xO);  Ly=(yN-yO);  L=sqrt(Lx*Lx+Ly*Ly);
     dLm=L/segment[s].N;

     if(point[segment[s].n0].F + point[segment[s].n1].F <= L)
      { 
       if(point[segment[s].n0].F > point[segment[s].n1].F)
	{M=-segment[s].N/2; ddL=(point[segment[s].n1].F-dLm)/M;}
       else
	{M=+segment[s].N/2; ddL=(dLm-point[segment[s].n0].F)/M;}
      }

/*=================
*  middle points  *
=================*/
     L_tot=0;
     if(point[segment[s].n0].F + point[segment[s].n1].F <= L)
       for(m=1; m<abs(segment[s].N); m++)
	{
	 L_tot+=(dLm-M*ddL);
  
	 if(point[segment[s].n0].F > point[segment[s].n1].F)
	  {M++; if(M==0 && segment[s].N%2==0) M++;}
	 else
	  {M--; if(M==0 && segment[s].N%2==0) M--;}

	 if(s==chain[c].s1 && m==(abs(segment[s].N)-1))
	  {insert_node(xO+Lx/L*L_tot, yO+Ly/L*L_tot, OFF,
	   Nn-1, segment[s].mark, segment[s].mark, segment[s].mark, point[segment[s].n1].new_numb);
	   node[Nn-1].next = Nn;}
	 
	 else if(m==1)
	  {insert_node(xO+Lx/L*L_tot, yO+Ly/L*L_tot, OFF,
	   point[segment[s].n0].new_numb, segment[s].mark, segment[s].mark, OFF, OFF);
	   node[Nn-1].next = Nn;}

	 else
	  {insert_node(xO+Lx/L*L_tot, yO+Ly/L*L_tot, OFF,
	   Nn-1, segment[s].mark, segment[s].mark, OFF, OFF);
	   node[Nn-1].next = Nn;}

	 node[Nn-1].chain    = segment[s].chain;
	 node[Nn-1].F        = 0.5*(node[Nn-2].F + (dLm-M*ddL));
	}

/*==============
*  last point  * -> just for the inside chains
==============*/
     if( (point[segment[s].n1].new_numb == OFF) && (s==chain[c].s1) )
      {
       insert_node(xN, yN, OFF,
       Nn-1, segment[s].mark, point[segment[s].n1].mark, OFF, OFF);
       node[Nn-1].next     = OFF;
       node[Nn-1].chain    = segment[s].chain;
       node[Nn-1].F        = point[segment[s].n1].F;
      }

     if( chain[c].type==CLOSED && s==chain[c].s1)
       node[Nn-1].next     = point[segment[chain[c].s0].n0].new_numb;

     if( chain[c].type==OPEN && s==chain[c].s1)
       node[Nn-1].next     = OFF;
    }
  }

 free(segment);
 free(inserted);

 return 0;
}
/*-load--------------------------------------------------------------------*/


/*=========================================================================*/
save()
{
 int  e, s, n, r_Nn=0, r_Ns=0, r_Ne=0;

 struct nod *r_node;
 struct ele *r_elem;
 struct sid *r_side;

 FILE *out;
 
 r_node=(struct nod *) calloc(Nn, sizeof(struct nod));
 r_elem=(struct ele *) calloc(Ne, sizeof(struct ele));
 r_side=(struct sid *) calloc(Ns, sizeof(struct sid));
 if(r_side==NULL)
  {fprintf(stderr, "Sorry, cannot allocate enough memory !\n");
   return 1;}

 for(n=0; n<Nn; n++)
   if(node[n].mark!=OFF && node[n].new_numb!=OFF)
    {
     r_Nn++;
     r_node[node[n].new_numb].x    = node[n].x;
     r_node[node[n].new_numb].y    = node[n].y;
     r_node[node[n].new_numb].mark = node[n].mark;
    }

 for(e=0; e<Ne; e++)
   if(elem[e].mark!=OFF && elem[e].new_numb!=OFF)
    {
     r_Ne++;
     r_elem[elem[e].new_numb].i  = node[elem[e].i].new_numb;
     r_elem[elem[e].new_numb].j  = node[elem[e].j].new_numb;
     r_elem[elem[e].new_numb].k  = node[elem[e].k].new_numb;
     r_elem[elem[e].new_numb].si = side[elem[e].si].new_numb;
     r_elem[elem[e].new_numb].sj = side[elem[e].sj].new_numb;
     r_elem[elem[e].new_numb].sk = side[elem[e].sk].new_numb;
     r_elem[elem[e].new_numb].xv = elem[e].xv;
     r_elem[elem[e].new_numb].yv = elem[e].yv;
     r_elem[elem[e].new_numb].material = elem[e].material;

     if(elem[e].ei != -1)
       r_elem[elem[e].new_numb].ei = elem[elem[e].ei].new_numb;
     else
       r_elem[elem[e].new_numb].ei = -1;

     if(elem[e].ej != -1)
       r_elem[elem[e].new_numb].ej = elem[elem[e].ej].new_numb;
     else
       r_elem[elem[e].new_numb].ej = -1;

     if(elem[e].ek != -1)
       r_elem[elem[e].new_numb].ek = elem[elem[e].ek].new_numb;
     else
       r_elem[elem[e].new_numb].ek = -1;
    }

 for(s=0; s<Ns; s++)
   if(side[s].mark!=OFF && side[s].new_numb!=OFF)
    {
     r_Ns++;
     r_side[side[s].new_numb].c    = node[side[s].c].new_numb;
     r_side[side[s].new_numb].d    = node[side[s].d].new_numb;
     r_side[side[s].new_numb].mark = side[s].mark;

     if(side[s].a != OFF)
      {r_side[side[s].new_numb].a  = node[side[s].a].new_numb;
       r_side[side[s].new_numb].ea = elem[side[s].ea].new_numb;}
     else
      {r_side[side[s].new_numb].a  = OFF;
       r_side[side[s].new_numb].ea = OFF;}

     if(side[s].b != OFF)
      {r_side[side[s].new_numb].b  = node[side[s].b].new_numb;
       r_side[side[s].new_numb].eb = elem[side[s].eb].new_numb;}
     else
      {r_side[side[s].new_numb].b  = OFF;
       r_side[side[s].new_numb].eb = OFF;}
    }

/*------------+
|             |
|  Node data  |
|             |
+------------*/
 name[len-1] = 'n';

 if((out=fopen(name, "w"))==NULL)
  {fprintf(stderr, "Cannot save file %s !\n", name);
   return 1;}
 
 fprintf(out, "%d\n", r_Nn);
 for(n=0; n<r_Nn; n++)
   fprintf(out, "%4d:  %18.15e %18.15e  %d\n",
		 n, r_node[n].x, r_node[n].y, r_node[n].mark);
 fprintf(out, "----------------------------------------------------------\n");
 fprintf(out, "   n:  x                      y                       mark\n");

 fclose(out);

/*---------------+
|                |
|  Element data  |
|                |
+---------------*/
 name[len-1] = 'e';

 if((out=fopen(name, "w"))==NULL)
  {fprintf(stderr, "Cannot save file %s !\n", name);
   return 1;}

 fprintf(out, "%d\n", r_Ne);
 for(e=0; e<r_Ne; e++)
   fprintf(out, "%4d: %4d %4d %4d  %4d %4d %4d  %4d %4d %4d  %18.15e %18.15e  %4d\n",
		 e, r_elem[e].i,  r_elem[e].j,  r_elem[e].k,
		    r_elem[e].ei, r_elem[e].ej, r_elem[e].ek,
		    r_elem[e].si, r_elem[e].sj, r_elem[e].sk,
		    r_elem[e].xv, r_elem[e].yv,
		    r_elem[e].material);
 fprintf(out, "---------------------------------------------------");
 fprintf(out, "-------------------------------------------------------\n");
 fprintf(out, "   e:   i,   j,   k,   ei,  ej,  ek,   si,  sj,  sk");  
 fprintf(out, "   xV,                    yV                       sign\n");  

 fclose(out);

/*------------+
|             |
|  Side data  |
|             |
+------------*/
 name[len-1] = 's';

 if((out=fopen(name, "w"))==NULL)
  {fprintf(stderr, "Cannot save file %s !\n", name);
   return 1;}
 
 fprintf(out, "%d\n", r_Ns);
 for(s=0; s<r_Ns; s++)
   fprintf(out, "%4d:  %4d %4d %4d %4d  %d\n",
		 s, r_side[s].c, r_side[s].d, r_side[s].ea, r_side[s].eb, r_side[s].mark);
 fprintf(out, "--------------------------------\n");
 fprintf(out, "   s:    c    d   ea   eb   mark\n");

 fclose(out);

 return 0;
}
/*-save--------------------------------------------------------------------*/


FILE *dxf_file;
char dxf_name[128];

/*=========================================================================*/
start_dxf()
{
 if((dxf_file=fopen(dxf_name,"w"))==NULL)
  {
   printf("A file '%s' cannot be opened for output ! \n\n", dxf_name);
   return 1;
  }
 else
  {
   fprintf(dxf_file, "0\n");
   fprintf(dxf_file, "SECTION\n");
   fprintf(dxf_file, "2\n");
   fprintf(dxf_file, "ENTITIES\n");
  }

 return 0;
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
line_dxf(double x1, double y1, double z1, 
	 double x2, double y2, double z2, 
	 char *layer)
{
 fprintf(dxf_file, "0\n");
 fprintf(dxf_file, "LINE\n");
 fprintf(dxf_file, "8\n");
 fprintf(dxf_file, "%s\n", layer);
 fprintf(dxf_file, "10\n");
 fprintf(dxf_file, "%lf\n", x1);
 fprintf(dxf_file, "20\n");
 fprintf(dxf_file, "%lf\n", y1);
 fprintf(dxf_file, "30\n");
 fprintf(dxf_file, "%lf\n", z1);
 fprintf(dxf_file, "11\n");
 fprintf(dxf_file, "%lf\n", x2);
 fprintf(dxf_file, "21\n");
 fprintf(dxf_file, "%lf\n", y2);
 fprintf(dxf_file, "31\n");
 fprintf(dxf_file, "%lf\n", z2);

 return 0;
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
end_dxf()
{
 fprintf(dxf_file, "0\n");
 fprintf(dxf_file, "ENDSEC\n");
 fprintf(dxf_file, "0\n");
 fprintf(dxf_file, "EOF\n");
 fclose(dxf_file);

 return 0;
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
draw_dxf()
{
 int    e, n, s, ei, ej, ek, ea, eb;
 double x, y, xc, yc, xd, yd, xa, ya, xb, yb;
 char   numb[128];

/*----------------+
|  Draw boundary  |
+----------------*/
 for(s=0; s<Ns; s++)
   if(side[s].mark>0) /* It means, side is on the boundary */
    {
     xc=node[side[s].c].x; yc=node[side[s].c].y;
     xd=node[side[s].d].x; yd=node[side[s].d].y;
     line_dxf(xc, yc, 0, xd, yd, 0, "boundary");
    }
 
/*----------------+
|  Draw Delaunay  |
+----------------*/
 for(s=0; s<Ns; s++)
   if(side[s].mark==0) /* It means: side is in the domain */
    {
     xc=node[side[s].c].x; yc=node[side[s].c].y;
     xd=node[side[s].d].x; yd=node[side[s].d].y;
     line_dxf(xc, yc, 0, xd, yd, 0, "delaunay");
    }

/*---------------+
|  Draw Voronoi  |
+---------------*/
 for(s=0; s<Ns; s++)
   if(side[s].mark!=OFF)
    {
     if((ea=side[s].ea)!=OFF)
      {xa=elem[ea].xv;
       ya=elem[ea].yv;}
     else
      {xa=0.5*(node[side[s].c].x+node[side[s].d].x);
       ya=0.5*(node[side[s].c].y+node[side[s].d].y);}
     
     if((eb=side[s].eb)!=OFF)
      {xb=elem[eb].xv;
       yb=elem[eb].yv;}
     else
      {xb=0.5*(node[side[s].c].x+node[side[s].d].x);
       yb=0.5*(node[side[s].c].y+node[side[s].d].y);}
     
     line_dxf(xa, ya, 0, xb, yb, 0, "voronoi");
    }
}
/*-draw_dxf---------------------------------------------------------------*/

FILE *fig_file;
char fig_name[128];

/*=========================================================================*/
start_fig()
{
 if((fig_file=fopen(fig_name,"w"))==NULL)
  {
   printf("A file '%s' cannot be opened for output ! \n\n", fig_name);
   return 1;
  }
 else
  {
   fprintf(fig_file, "#FIG 3.1\n");
   fprintf(fig_file, "Landscape\n");
   fprintf(fig_file, "Center\n");
   fprintf(fig_file, "Metric\n");
   fprintf(fig_file, "1200 2\n");
  }

 return 0;
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
line_fig(int x1, int y1, 
	 int x2, int y2, 
	 int style, int width, int color, float le)
{
 fprintf(fig_file, "2 ");
 fprintf(fig_file, "1 ");
 fprintf(fig_file, "%d ", style); /* 0 - solid, 1 - dashed, 2 - dotted */
 fprintf(fig_file, "%d ", width);
 fprintf(fig_file, "%d ", color); /* pen color */
 fprintf(fig_file, "7 ");         /* fill color 0 - black, 7 - white */ 
 fprintf(fig_file, "0 ");         /* depth */ 
 fprintf(fig_file, "0 ");         /* ? */  
 fprintf(fig_file, "-1 ");        /* fill style, -1 - no fill */
 fprintf(fig_file, "%5.3f ", le); /* lenght for dashes */ 
 fprintf(fig_file, "0 ");         /* join style */  
 fprintf(fig_file, "0 ");         /* cap style */ 
 fprintf(fig_file, "-1 ");        /* ? */
 fprintf(fig_file, "0 ");         /* forward arrow */
 fprintf(fig_file, "0 ");         /* backward arrow */
 fprintf(fig_file, "2\n");        /* number of points */  

 fprintf(fig_file, "         ");
 fprintf(fig_file, "%d ",  x1);
 fprintf(fig_file, "%d ",  y1);
 fprintf(fig_file, "%d ",  x2);
 fprintf(fig_file, "%d\n", y2);

 return 0;
}
/*-------------------------------------------------------------------------*/


/*=========================================================================*/
end_fig()
{
 fclose(fig_file);

 return 0;
}
/*-------------------------------------------------------------------------*/


/*===========================================================================
 Let's say that drawing area is 20 x 20 cm. One cm in xfig is 450 poins.
 It means that drawing area is 9000 x 9000 points.
---------------------------------------------------------------------------*/
draw_fig()
{
 int    e, n, s, ei, ej, ek, ea, eb;
 double x, y, xc, yc, xd, yd, xa, ya, xb, yb,
	xmax=-GREAT, xmin=+GREAT, ymax=-GREAT, ymin=+GREAT, scl;
 char   numb[128];

 for(n=0; n<Nn; n++)
   if(node[n].mark!=OFF)
    {
     xmin=min(xmin, node[n].x); ymin=min(ymin, node[n].y);
     xmax=max(xmax, node[n].x); ymax=max(ymax, node[n].y);
    }
 scl =min( 9000.0/(ymax-ymin+SMALL), 9000.0/(xmax-xmin+SMALL) );

/*----------------+
|  Draw boundary  |
+----------------*/
 for(s=0; s<Ns; s++)
   if(side[s].mark>0) /* It means, side is on the boundary */
    {
     xc=node[side[s].c].x; yc=node[side[s].c].y;
     xd=node[side[s].d].x; yd=node[side[s].d].y;
     line_fig ( 450+(int)floor(scl*xc), 450+(int)floor(scl*yc), 
		450+(int)floor(scl*xd), 450+(int)floor(scl*yd), 
		0, 3, 0, 0.000);
    }
 
/*----------------+
|  Draw Delaunay  |
+----------------*/
 for(s=0; s<Ns; s++)
   if(side[s].mark==0) /* It means: side is in the domain */
    {
     xc=node[side[s].c].x; yc=node[side[s].c].y;
     xd=node[side[s].d].x; yd=node[side[s].d].y;
     line_fig( 450+(int)floor(scl*xc), 450+(int)floor(scl*yc), 
	       450+(int)floor(scl*xd), 450+(int)floor(scl*yd),
	       0, 1, 1, 0.000);
    }

/*---------------+
|  Draw Voronoi  |
+---------------*/
 for(s=0; s<Ns; s++)
   if(side[s].mark!=OFF)
    {
     if((ea=side[s].ea)!=OFF)
      {xa=elem[ea].xv;
       ya=elem[ea].yv;}
     else
      {xa=0.5*(node[side[s].c].x+node[side[s].d].x);
       ya=0.5*(node[side[s].c].y+node[side[s].d].y);}
     
     if((eb=side[s].eb)!=OFF)
      {xb=elem[eb].xv;
       yb=elem[eb].yv;}
     else
      {xb=0.5*(node[side[s].c].x+node[side[s].d].x);
       yb=0.5*(node[side[s].c].y+node[side[s].d].y);}
     
     line_fig( 450+(int)floor(scl*xa), 450+(int)floor(scl*ya), 
	       450+(int)floor(scl*xb), 450+(int)floor(scl*yb),
	       0, 1, 4, 0.000);
    }
}
/*-draw_fig---------------------------------------------------------------*/


#if GRAPHICS == ON
/*--------------------------------------------------------------------------+
|          NONSTANDARD MS-DOS GRAPHICAL FUNCTIONS   (WATCOM 10.0a)          |
+--------------------------------------------------------------------------*/
#include <graph.h>
#define B_CIRCLE(x1, y1, x2, y2)  _ellipse(_GBORDER,       x1, y1, x2, y2)
#define I_CIRCLE(x1, y1, x2, y2)  _ellipse(_GFILLINTERIOR, x1, y1, x2, y2)
#define CLEARSCREEN               _clearscreen(_GCLEARSCREEN)
#define CLOSEGRAPH                _setvideomode(_DEFAULTMODE)
#define COLOR(b)                  _setcolor(b)
#define FILL(x1, y1, color)       _floodfill(x1, y1, color)
#define GRTEXT(x1, y1, string)    _moveto(x1, y1); _outgtext(string)
#define LINE(x1, y1, x2, y2)      _moveto(x1, y1); _lineto(x2,y2)
#define OPEN_SVGA                 _setvideomode(_SVRES256COLOR)
#define OPEN_VGA                  _setvideomode(_VRES16COLOR)
#define RECTANGLE(x1, y1, x2, y2) _rectangle(_GBORDER, x1, y1, x2, y2)
/*------------------------------------------------------------------------*/


/*=========================================================================*/
draw(int mesh, int voronoi, int marks, int fill)
{
 int    e, n, s, X0, Y0, ei, ej, ek, ea, eb;
 double scl, x, y, xc, yc, xd, yd, x1, y1, x2, y2,
	xmax=-GREAT, xmin=+GREAT, ymax=-GREAT, ymin=+GREAT;
 char   numb[128];

 for(n=0; n<Nn; n++)
   if(node[n].mark!=OFF)
    {
     xmin=min(xmin, node[n].x); ymin=min(ymin, node[n].y);
     xmax=max(xmax, node[n].x); ymax=max(ymax, node[n].y);
    }
 scl =min( 380.0/(ymax-ymin+SMALL), 700.0/(xmax-xmin+SMALL) );
 
 X0 = 50;
 Y0 = 550;
 
 if(mesh!=OFF)
  {
   if(mesh==0) COLOR(7);
   else        COLOR(mesh);

   for(s=0; s<Ns; s++)
     if(side[s].mark!=OFF)
      {
       xc=node[side[s].c].x; yc=node[side[s].c].y;
       xd=node[side[s].d].x; yd=node[side[s].d].y;

/* if(side[s].new_numb==OFF) COLOR(13); */

       LINE(xc*scl + X0, -yc*scl + Y0, xd*scl + X0, -yd*scl + Y0);
      }

   if(fill==ON && mesh==ON)
     for(e=0; e<Ne; e++)
       if(elem[e].mark!=OFF)
        {
         x=0.333333333*(node[elem[e].i].x+node[elem[e].j].x+node[elem[e].k].x);
         y=0.333333333*(node[elem[e].i].y+node[elem[e].j].y+node[elem[e].k].y);
         COLOR(16-elem[e].material);
         FILL(x*scl+X0, -y*scl+Y0, 7);
        }
  } /* mesh==ON */

 for(s=0; s<Ns; s++)
   if(side[s].mark>0) /* It means, side is on the boundary */
    {
     if(marks==ON) COLOR(16-side[s].mark);
     else          COLOR(15);

     xc=node[side[s].c].x; yc=node[side[s].c].y;
     xd=node[side[s].d].x; yd=node[side[s].d].y;
   
     LINE(xc*scl + X0, -yc*scl + Y0, xd*scl + X0, -yd*scl + Y0);
    }
 
 if(marks==ON)
   for(n=0; n<Nn; n++)
    {
     if(node[n].mark>0)  /* node is on the boundary */
      {
       COLOR(16-node[n].mark);
       x=node[n].x; y=node[n].y;
       I_CIRCLE(x*scl+X0+2, -y*scl+Y0+2, x*scl+X0-2, -y*scl+Y0-2);
      }
    }

 if(voronoi!=OFF)
  {
   if(voronoi==ON) COLOR(7);
   else            COLOR(voronoi);

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
     
       LINE(x1*scl + X0, -y1*scl + Y0, x2*scl + X0, -y2*scl + Y0);
      }
  }
}
/*-draw-------------------------------------------------------------------*/
#endif


int main(int argc, char *argv[])
{
 int arg, ans, d=ON, r=ON, s=ON, dxf=OFF, fig=OFF, m=ON, g=ON, exa=OFF, Nn0;

 if(argc<2)
  {printf("\n*********************************************************");
   printf("\n****************                        *****************");
   printf("\n****************   PROGRAM:  EasyMesh   *****************");
   printf("\n****************                        *****************");
   printf("\n****************      version 1.4       *****************");
   printf("\n****************                        *****************");
   printf("\n****************  Author: Bojan NICENO  *****************");
   printf("\n**************** niceno@univ.trieste.it *****************");
   printf("\n****************                        *****************");
   printf("\n*********************************************************");
   printf("\n\nUsage:  EasyMesh  <NAME>  [<options>]");
   printf("\n\n***************");
   printf("\n*** OPTIONS ***");
   printf("\n***************");
   printf("\n\nValid options are:");
   printf("\n   -d        don't triangulate domain");
   printf("\n   -g        without graphic output");
   printf("\n   -m        without messages");
   printf("\n   -r        without relaxation");
   printf("\n   -s        without Laplacian smoothing");
   printf("\n   +dxf      create drawing in DXF format");
   printf("\n   +fig      create drawing in fig format");
   printf("\n   +example  create example input file");
   printf("\n\n*************");
   printf("\n*** INPUT ***");
   printf("\n*************");
   printf("\n\nInput file (NAME.d) has the following format");
   printf("\n  first line:          <Nbp>");
   printf("\n  following Nbp lines: <point:> <x> <y> <spacing> <marker>");
   printf("\n  one line:            <Nbs>");
   printf("\n  following Nbs lines: <segment:> <start_point> <end_point> <marker>");
   printf("\n\n  where:");
   printf("\n    Nbn     is the number of points defining the boundary");
   printf("\n    Nbp     is the number of sides defining the boundary");
   printf("\n    marker  is the boundary condition marker");
   printf("\n\nNote: Input file has to end with the extension: .d !");
   printf("\n\n**************");
   printf("\n*** OUTPUT ***");
   printf("\n**************");
   printf("\n\nEasyMesh produces the following three output files:");
   printf("\n  NAME.n");
   printf("\n  NAME.e");
   printf("\n  NAME.s");
   printf("\n\nNode file (NAME.n) has the following format:");
   printf("\n  first line:         <Nn>");
   printf("\n  following Nn lines: <node:> <x> <y> <marker> ");
   printf("\n  last two lines:     comments inserted by the program ");
   printf("\n\n  where:");
   printf("\n  Nn      is the number of nodes");
   printf("\n  x, y    are the node coordinates");
   printf("\n  marker  is the node boundary marker");
   printf("\n\nElement file (NAME.e) has the following format:");
   printf("\n  first line          <Ne> ");
   printf("\n  following Ne lines: <element:> <i> <j> <k> <ei> <ej> <ek> <si> <sj> <sk> <xV> <yV> <marker> ");
   printf("\n  last two lines:     comments inserted by the program ");
   printf("\n\n  where:");
   printf("\n    Ne          is the number of elements");
   printf("\n    i,   j,  k  are the nodes belonging to the element, ");
   printf("\n    ei, ej, ek  are the neighbouring elements,");
   printf("\n    si, sj, sk  are the element sides. ");
   printf("\n    xV, yV      are the coordinates of the element circumcenter");
   printf("\n    marker      is the side boundary marker");
   printf("\n\nSide file (NAME.s) has the following format:");
   printf("\n  first line:         <Ns> ");
   printf("\n  following Ns lines: <side:> <c> <d> <ea> <eb> <marker> "); 
   printf("\n  last two lines:     comments inserted by the program ");
   printf("\n\n  where:");
   printf("\n    Ns      is the number of sides");
   printf("\n    c,  d   are the starting and ending node of the side,");
   printf("\n    ea, eb  are the elements on the left and on the right of the side.");
   printf("\n\nNote: If eb equals to -1, it means that the right element does not exists,");
   printf("\n      i.e. the side is on the boundary !");  
   printf("\n\n");
   
   return 0;}
 else
  {if(strcmp(argv[1], "+example")==0) exa=ON;
   strcpy(name,     argv[1]);
   len=strlen(name);
   if(name[len-2]=='.')
     if(name[len-1]=='d' || name[len-1]=='D' )
       name[len-2]='\0';
   strcpy(dxf_name, name); strcat(dxf_name, ".dxf"); 
   strcpy(fig_name, name); strcat(fig_name, ".fig");}

/*-----------------------+
|  command line options  |
+-----------------------*/
 for(arg=2; arg<argc; arg++)
  {
   if(strcmp(argv[arg],"-d")      ==0) {d=OFF; r=OFF; s=OFF;}
   if(strcmp(argv[arg],"+dxf")    ==0) dxf=ON;
   if(strcmp(argv[arg],"+fig")    ==0) fig=ON;
   if(strcmp(argv[arg],"+example")==0) exa=ON;
   if(strcmp(argv[arg],"-g")      ==0) g=OFF;
   if(strcmp(argv[arg],"-r")      ==0) r=OFF;
   if(strcmp(argv[arg],"-s")      ==0) s=OFF;
   if(strcmp(argv[arg],"-m")      ==0) m=OFF;
  }

 if(exa==ON)
  {
   FILE *out;
   if( (out=fopen("example.d", "w"))==NULL )
    {printf("Can't open file 'example.d' for output !");
     return 1;}
   else
    {
     fprintf(out, "\n#*******************************************************");
     fprintf(out, "\n");
     fprintf(out, "\n  Everything enclosed inside the cashes is a comment.");   
     fprintf(out, "\n");
     fprintf(out, "\n  This input file is used to generate the triangular");
     fprintf(out, "\n  mesh over the following domain:");
     fprintf(out, "\n");
     fprintf(out, "\n     3--------------------2");
     fprintf(out, "\n     |                    |");
     fprintf(out, "\n     |    5----6          |");
     fprintf(out, "\n     |    |    |          |");
     fprintf(out, "\n     |    |    |          |");
     fprintf(out, "\n     |    4----7          |");
     fprintf(out, "\n     |                    |");
     fprintf(out, "\n     0--------------------1");
     fprintf(out, "\n");
     fprintf(out, "\n");
     fprintf(out, "\n  Run EasyMesh with the command:");
     fprintf(out, "\n");
     fprintf(out, "\n  > EasyMesh example +fig");
     fprintf(out, "\n ");
     fprintf(out, "\n  if you want to see the results with xfig,");
     fprintf(out, "\n  or with");
     fprintf(out, "\n");
     fprintf(out, "\n  > EasyMesh example +dxf");
     fprintf(out, "\n");
     fprintf(out, "\n  if you want to see the results with some tool that");
     fprintf(out, "\n  suports Autodesk DXF format");
     fprintf(out, "\n");
     fprintf(out, "\n*******************************************************#");
     fprintf(out, "\n");
     fprintf(out, "\n#*********");
     fprintf(out, "\n  POINTS");
     fprintf(out, "\n*********#");
     fprintf(out, "\n8 # number of points defining the boundary #");
     fprintf(out, "\n");
     fprintf(out, "\n# rectangular domain #");
     fprintf(out, "\n#-------+-----+-----+---------+--------#");
     fprintf(out, "\n# point |  x  |  y  | spacing | marker #");
     fprintf(out, "\n#-------+-----+-----+---------+--------#"); 
     fprintf(out, "\n   0:     0.0   0.0    0.05       1");
     fprintf(out, "\n   1:     2.0   0.0    0.25       2");
     fprintf(out, "\n   2:     2.0   1.6    0.25       2");
     fprintf(out, "\n   3:     0.0   1.6    0.1        1");
     fprintf(out, "\n");
     fprintf(out, "\n# square hole #");
     fprintf(out, "\n#-------+-----+-----+---------+--------#");
     fprintf(out, "\n# point |  x  |  y  | spacing | marker #");
     fprintf(out, "\n#-------+-----+-----+---------+--------#"); 
     fprintf(out, "\n   4:     0.5   0.5    0.05       3");
     fprintf(out, "\n   5:     0.5   1.1    0.08       3");
     fprintf(out, "\n   6:     1.1   1.1    0.2        3");
     fprintf(out, "\n   7:     1.1   0.5    0.2        3");
     fprintf(out, "\n");
     fprintf(out, "\n#***********");
     fprintf(out, "\n  SEGMENTS");
     fprintf(out, "\n***********#");
     fprintf(out, "\n8 # number of segments #");
     fprintf(out, "\n");
     fprintf(out, "\n# domain #");
     fprintf(out, "\n#---------+-----+-----+--------#");
     fprintf(out, "\n# segment | 1st | 2nd | marker #");
     fprintf(out, "\n#---------+-----+-----+--------#");    
     fprintf(out, "\n     0:      0     1      1");
     fprintf(out, "\n     1:      1     2      2");
     fprintf(out, "\n     2:      2     3      1");
     fprintf(out, "\n     3:      3     0      1");
     fprintf(out, "\n");
     fprintf(out, "\n# hole #");
     fprintf(out, "\n#---------+-----+-----+--------#");
     fprintf(out, "\n# segment | 1st | 2nd | marker #");
     fprintf(out, "\n#---------+-----+-----+--------#");
     fprintf(out, "\n     4:      4     5      3");
     fprintf(out, "\n     5:      5     6      3");
     fprintf(out, "\n     6:      6     7      3");
     fprintf(out, "\n     7:      7     4      3\n");
     printf("\nThe file 'example.d' created !\n");
     fclose(out);
     return 0;
    }
  }

 strcat(name, ".d");
 len=strlen(name);


 printf("\nLoading the input file: %s \n", name); fflush(stdout);
 if(load()!=0)
   return 1;
 erase();
 classify();

 if(m==ON)
   printf("Working...\n"); fflush(stdout);

 if(d==ON)
   do
    {
     Nn0=Nn;
     new_node();
     classify();
     if(Nn==MAX_NODES-1) break;
     if(Nn==Nn0) break;
    }
   while(ugly!=OFF);

 neighbours();
 
 if(r==ON || s==ON)
   if(m==ON)
     printf("Improving the grid quality...\n"); fflush(stdout);
 if(r==ON)           relax();
 if(r==ON || s==ON) smooth();

 if(m==ON)
   printf("Renumerating nodes, elements and sides...\n"); fflush(stdout);
 renum();

 if(m==ON)
   printf("Processing material marks... \n"); fflush(stdout);
 materials();

#if GRAPHICS == ON
 if(g==ON)
  {
   do
    {
     printf("\n****************************");
     printf("\n**   Enter your choice:   **");
     printf("\n****************************\n\n");
     printf("0. Exit\n");
     printf("1. Delaunay\n");
     printf("2. Voronoi\n");
     printf("3. Delaunay and Voronoi\n");
     printf("4. Materials\n");
     printf("5. Boundary condition marks only\n");
     printf("6. Boundary condition marks with Delaunay\n");
     printf("7. Boundary condition marks with Voronoi\n\n");

     printf("-> "); scanf("%d", &ans);

     switch(ans)
      {
       case 1: OPEN_SVGA; draw(11,  OFF, OFF, OFF); printf("Press any key to continue !"); fflush(stdout); getch(); CLOSEGRAPH; break;
       case 2: OPEN_SVGA; draw(OFF, 12,  OFF, OFF); printf("Press any key to continue !"); fflush(stdout); getch(); CLOSEGRAPH; break;
       case 3: OPEN_SVGA; draw(11,  12,  OFF, OFF); printf("Press any key to continue !"); fflush(stdout); getch(); CLOSEGRAPH; break;
       case 4: OPEN_SVGA; draw(ON,  OFF, OFF, ON ); printf("Press any key to continue !"); fflush(stdout); getch(); CLOSEGRAPH; break;
       case 5: OPEN_SVGA; draw(OFF, OFF, ON,  OFF); printf("Press any key to continue !"); fflush(stdout); getch(); CLOSEGRAPH; break;
       case 6: OPEN_SVGA; draw(ON,  OFF, ON,  OFF); printf("Press any key to continue !"); fflush(stdout); getch(); CLOSEGRAPH; break;
       case 7: OPEN_SVGA; draw(OFF, ON,  ON,  OFF); printf("Press any key to continue !"); fflush(stdout); getch(); CLOSEGRAPH; break;
      }
    }
   while(ans!=0);
  }
#endif

 save();

 if(dxf==ON)
  {start_dxf(); draw_dxf(); end_dxf();}

 if(fig==ON)
  {start_fig(); draw_fig(); end_fig();}

 return 1;
}

