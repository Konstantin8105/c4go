#-----------#
# Example 1 #
#-----------#

#=========
| POINTS |
=========#
9 # number of points #

# Nodes which define the boundary #
0:  0.0  0.0    0.25    1
1:  5.0  0.0    0.25    2
2:  5.0  2.0    0.25    2
3:  4.0  3.0    0.25    3
4:  0.0  3.0    0.25    3

# Nodes which define the hole #
5:  1.0  1.0    0.1    4
6:  1.0  2.0    0.1    4
7:  2.0  2.0    0.1    4
8:  2.0  1.0    0.1    4

#===========
| SEGMENTS |
===========#
9 # Number of segments #

# Boundary segments #
0:  0  1    1
1:  1  2    2
2:  2  3    2
3:  3  4    3
4:  4  0    3

# Hole segments #
5:  5  6    4
6:  6  7    4
7:  7  8    4
8:  8  5    4
