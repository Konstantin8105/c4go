#!/bin/bash

set -e

go build

mkdir -p ./testdata/



# prepare variables
	export C4GO_DIR=$GOPATH/src/github.com/Konstantin8105/c4go
	export C4GO=$C4GO_DIR/c4go
	export TEMP_FOLDER="./testdata/gtk"
	export CFILE=$TEMP_FOLDER/gtk.c
	export GOFILE=$TEMP_FOLDER/gtk.go

# prepare C code
    if [ ! -d $TEMP_FOLDER ]; then
		mkdir -p $TEMP_FOLDER
		touch $CFILE 
		echo '
#include <gtk/gtk.h>

static void
activate (GtkApplication* app,
          gpointer        user_data)
{
  GtkWidget *window;

  window = gtk_application_window_new (app);
  gtk_window_set_title (GTK_WINDOW (window), "Window");
  gtk_window_set_default_size (GTK_WINDOW (window), 200, 200);
  gtk_widget_show_all (window);
}

int
main (int    argc,
      char **argv)
{
  GtkApplication *app;
  int status;

  app = gtk_application_new ("org.gtk.example", G_APPLICATION_FLAGS_NONE);
  g_signal_connect (app, "activate", G_CALLBACK (activate), NULL);
  status = g_application_run (G_APPLICATION (app), argc, argv);
  g_object_unref (app);

  return status;
}' > $CFILE
	fi

# remove go files from last transpilation
	echo "***** remove go files"
	rm -f $TEMP_FOLDER/*.go
	rm -f $TEMP_FOLDER/*.app

# transpilation 
export FILES=`ls $TEMP_FOLDER/*.c | tr "\n" " "`
$C4GO transpile  -s           \
	             -clang-flag="-pthread -I/usr/include/gtk-3.0 -I/usr/include/at-spi2-atk/2.0 -I/usr/include/at-spi-2.0 -I/usr/include/dbus-1.0 -I/usr/lib/x86_64-linux-gnu/dbus-1.0/include -I/usr/include/gtk-3.0 -I/usr/include/gio-unix-2.0/ -I/usr/include/cairo -I/usr/include/pango-1.0 -I/usr/include/harfbuzz -I/usr/include/pango-1.0 -I/usr/include/atk-1.0 -I/usr/include/cairo -I/usr/include/pixman-1 -I/usr/include/freetype2 -I/usr/include/libpng16 -I/usr/include/freetype2 -I/usr/include/libpng16 -I/usr/include/gdk-pixbuf-2.0 -I/usr/include/libpng16 -I/usr/include/glib-2.0 -I/usr/lib/x86_64-linux-gnu/glib-2.0/include -lgtk-3 -lgdk-3 -lpangocairo-1.0 -lpango-1.0 -latk-1.0 -lcairo-gobject -lcairo -lgdk_pixbuf-2.0 -lgio-2.0 -lgobject-2.0 -lglib-2.0"  \
	             -o="$GOFILE" \
				 $CFILE

echo "Calculate warnings : $TEMP_FOLDER"
# show warnings comments in Go source
	WARNINGS=`cat $GOFILE | grep "^// Warning" | sort | uniq | wc -l`
	echo "		After transpiling : $WARNINGS warnings."
# show amount error from `go build`:
	WARNINGS_GO=`go build -o $TEMP_FOLDER/gtk.app -gcflags="-e" $FILE 2>&1 | wc -l`
	echo "		Go build : $WARNINGS_GO warnings"
# amount unsafe
	UNSAFE=`cat $GOFILE | grep "unsafe\." | wc -l`
	echo "		Unsafe   : $UNSAFE"
# amount Go code lines
	LINES=`wc $GOFILE`
	echo "(lines,words,bytes)	 : $GOLINES"
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

