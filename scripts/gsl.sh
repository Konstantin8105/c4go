#!/bin/bash

set -e

go build

mkdir -p ./testdata/



# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export TEMP_FOLDER="./testdata/gsl"
	export CFILE=$TEMP_FOLDER/gsl.c
	export GOFILE=$TEMP_FOLDER/gsl.go

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		touch $CFILE 
		echo '
#include <stdio.h>
#include <gsl/gsl_errno.h>
#include <gsl/gsl_matrix.h>
#include <gsl/gsl_odeiv2.h>

int
func (double t, const double y[], double f[],
      void *params)
{
  (void)(t); /* avoid unused parameter warning */
  double mu = *(double *)params;
  f[0] = y[1];
  f[1] = -y[0] - mu*y[1]*(y[0]*y[0] - 1);
  return GSL_SUCCESS;
}

int
jac (double t, const double y[], double *dfdy,
     double dfdt[], void *params)
{
  (void)(t); /* avoid unused parameter warning */
  double mu = *(double *)params;
  gsl_matrix_view dfdy_mat
    = gsl_matrix_view_array (dfdy, 2, 2);
  gsl_matrix * m = &dfdy_mat.matrix;
  gsl_matrix_set (m, 0, 0, 0.0);
  gsl_matrix_set (m, 0, 1, 1.0);
  gsl_matrix_set (m, 1, 0, -2.0*mu*y[0]*y[1] - 1.0);
  gsl_matrix_set (m, 1, 1, -mu*(y[0]*y[0] - 1.0));
  dfdt[0] = 0.0;
  dfdt[1] = 0.0;
  return GSL_SUCCESS;
}

int
main (void)
{
  double mu = 10;
  gsl_odeiv2_system sys = {func, jac, 2, &mu};

  gsl_odeiv2_driver * d =
    gsl_odeiv2_driver_alloc_y_new (&sys, gsl_odeiv2_step_rk8pd,
                                  1e-6, 1e-6, 0.0);
  int i;
  double t = 0.0, t1 = 100.0;
  double y[2] = { 1.0, 0.0 };

  for (i = 1; i <= 100; i++)
    {
      double ti = i * t1 / 100.0;
      int status = gsl_odeiv2_driver_apply (d, &t, ti, y);

      if (status != GSL_SUCCESS)
        {
          printf ("error, return value=%d\n", status);
          break;
        }

      printf ("%.5e %.5e %.5e\n", t, y[0], y[1]);
    }

  gsl_odeiv2_driver_free (d);
  return 0;
}
' > $CFILE
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app
# ast 
if [ "$1" == "-a" ]; then
	$C4GO ast		-clang-flag="-lgsl  -lgslcblas"  \
					 $CFILE
fi

# transpilation 
$C4GO transpile  -s           \
	             -clang-flag="-lgsl  -lgslcblas"  \
	             -o="$GOFILE" \
				 $CFILE

echo "Calculate warnings : $TEMP_FOLDER"
# show warnings comments in Go source
	WARNINGS=`cat $GOFILE | grep "^// Warning" | sort | uniq | wc -l`
	echo "		After transpiling : $WARNINGS warnings."
# show amount error from `go build`:
	WARNINGS_GO=`go build -o $TEMP_FOLDER/gsl.app -gcflags="-e" $FILE 2>&1 | wc -l`
	echo "		Go build : $WARNINGS_GO warnings"
# amount unsafe
	UNSAFE=`cat $GOFILE | grep "unsafe\." | wc -l`
	echo "		Unsafe   : $UNSAFE"
# amount Go code lines
	LINES=`wc $GOFILE`
	echo "(lines,words,bytes)	 : $LINES"
# defers
	DEFER=`cat $GOFILE | grep "defer func" | wc -l`
	echo "defer func           	 : $DEFER"

# Arguments menu
echo "    -s for show detail of Go build errors"
if [ "$1" == "-s" ]; then
	# show go build warnings	
		for f in $TEMP_FOLDER/*.go ; do
			# iteration by Go files
				echo "	file : $f"
			# c4go warnings
				cat $f | grep "^// Warning" | sort | uniq
			# show amount error from `go build`:
				go build -o $f.app -gcflags="-e" $f 2>&1 | sort 
		done
fi

