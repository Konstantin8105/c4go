package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	file, err := os.Open("/tmp/gcc.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var cList map[string]bool = map[string]bool{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "ARG") {
			continue
		}
		index := strings.Index(line, "-c")
		if index < 0 {
			continue
		}
		line = line[index+len("-c "):]
		index = strings.Index(line, " ")
		if index < 0 {
			continue
		}
		line = line[:index]
		if !strings.HasSuffix(strings.ToLower(line), ".c") {
			continue
		}

		folder := "/tmp/GSL/gsl-2.4/"
		var fileList []string
		// find all C source files
		err = filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
			if strings.HasSuffix(strings.ToLower(f.Name()), ".c") {
				if strings.HasSuffix(path, "/"+line) {
					fileList = append(fileList, path)
				}
			}
			return nil
		})
		if err != nil {
			err = fmt.Errorf("Cannot walk: %v", err)
			return
		}

		for _, f := range fileList {
			cList[f] = true
		}

		fmt.Println("line = ", line, fileList)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// sorting list
	var sortList []string
	for k := range cList {
		sortList = append(sortList, k)
	}
	sort.Strings(sortList)

	f, err := os.Create("./travis/gsl.list")
	if err != nil {
		log.Fatal(err)
	}
	for _, k := range sortList {
		fmt.Printf("%s ", k)
		f.WriteString(fmt.Sprintf("%s\n", k))
	}
	f.Close()

	fg, err := os.Open("./travis/gsl.list")
	if err != nil {
		log.Fatal(err)
	}
	var list []string
	scannerG := bufio.NewScanner(fg)
	for scannerG.Scan() {
		line := scannerG.Text()
		cmd := exec.Command("c4go", "transpile",
			"-clang-flag=-DHAVE_CONFIG_H",
			"-clang-flag=-I/tmp/GSL/gsl-2.4/",
			"-o="+line[:len(line)-2]+".go",
			line)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Fail: ", line, err)
		} else {
			list = append(list, line)
			fmt.Println("Ok: ", line)
		}
	}
	fmt.Println(list)
	fg.Close()

	folder := "/tmp/GSL/gsl-2.4/"
	// find all C source files
	err = filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".go") {
			l, err := getLogs(path)
			fmt.Println(">>", path)
			for _, t := range l {
				fmt.Println(t)
			}
			fmt.Println(err)
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("Cannot walk: %v", err)
		return
	}
}

func getLogs(goFile string) (logs []string, err error) {
	file, err := os.Open(goFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// ignore
		// Warning (*ast.TranslationUnitDecl):  :0 :cannot transpileRecordDecl `__WAIT_STATUS`. could not determine the size of type `union __WAIT_STATUS` for that reason: Cannot determine sizeof : |union __WAIT_STATUS|. err = Cannot canculate `union` sizeof for `string`. Cannot determine sizeof : |union wait *|. err = error in union
		if strings.Contains(line, "union __WAIT_STATUS") {
			continue
		}

		if strings.Contains(line, "//") && strings.Contains(line, "AST") {
			logs = append(logs, line)
		}
		if strings.HasPrefix(line, "// Warning") {
			logs = append(logs, line)
		}
	}

	err = scanner.Err()
	return
}

// c4go transpile -clang-flag="-DHAVE_CONFIG_H" -clang-flag="-I/tmp/GSL/gsl-2.4/" /tmp/GSL/gsl-2.4/blas/blas.c /tmp/GSL/gsl-2.4/block/block.c /tmp/GSL/gsl-2.4/block/file.c /tmp/GSL/gsl-2.4/block/init.c /tmp/GSL/gsl-2.4/bspline/bspline.c /tmp/GSL/gsl-2.4/bspline/greville.c /tmp/GSL/gsl-2.4/cblas/caxpy.c /tmp/GSL/gsl-2.4/cblas/ccopy.c /tmp/GSL/gsl-2.4/cblas/cdotc_sub.c /tmp/GSL/gsl-2.4/cblas/cdotu_sub.c /tmp/GSL/gsl-2.4/cblas/cgbmv.c /tmp/GSL/gsl-2.4/cblas/cgemm.c /tmp/GSL/gsl-2.4/cblas/cgemv.c /tmp/GSL/gsl-2.4/cblas/cgerc.c /tmp/GSL/gsl-2.4/cblas/cgeru.c /tmp/GSL/gsl-2.4/cblas/chbmv.c /tmp/GSL/gsl-2.4/cblas/chemm.c /tmp/GSL/gsl-2.4/cblas/chemv.c /tmp/GSL/gsl-2.4/cblas/cher.c /tmp/GSL/gsl-2.4/cblas/cher2.c /tmp/GSL/gsl-2.4/cblas/cher2k.c /tmp/GSL/gsl-2.4/cblas/cherk.c /tmp/GSL/gsl-2.4/cblas/chpmv.c /tmp/GSL/gsl-2.4/cblas/chpr.c /tmp/GSL/gsl-2.4/cblas/chpr2.c /tmp/GSL/gsl-2.4/cblas/cscal.c /tmp/GSL/gsl-2.4/cblas/csscal.c /tmp/GSL/gsl-2.4/cblas/cswap.c /tmp/GSL/gsl-2.4/cblas/csymm.c /tmp/GSL/gsl-2.4/cblas/csyr2k.c /tmp/GSL/gsl-2.4/cblas/csyrk.c /tmp/GSL/gsl-2.4/cblas/ctbmv.c /tmp/GSL/gsl-2.4/cblas/ctbsv.c /tmp/GSL/gsl-2.4/cblas/ctpmv.c /tmp/GSL/gsl-2.4/cblas/ctpsv.c /tmp/GSL/gsl-2.4/cblas/ctrmm.c /tmp/GSL/gsl-2.4/cblas/ctrmv.c /tmp/GSL/gsl-2.4/cblas/ctrsm.c /tmp/GSL/gsl-2.4/cblas/ctrsv.c /tmp/GSL/gsl-2.4/cblas/dasum.c /tmp/GSL/gsl-2.4/cblas/daxpy.c /tmp/GSL/gsl-2.4/cblas/dcopy.c /tmp/GSL/gsl-2.4/cblas/ddot.c /tmp/GSL/gsl-2.4/cblas/dgbmv.c /tmp/GSL/gsl-2.4/cblas/dgemm.c /tmp/GSL/gsl-2.4/cblas/dgemv.c /tmp/GSL/gsl-2.4/cblas/dger.c /tmp/GSL/gsl-2.4/cblas/dnrm2.c /tmp/GSL/gsl-2.4/cblas/drot.c /tmp/GSL/gsl-2.4/cblas/drotg.c /tmp/GSL/gsl-2.4/cblas/drotm.c /tmp/GSL/gsl-2.4/cblas/drotmg.c /tmp/GSL/gsl-2.4/cblas/dsbmv.c /tmp/GSL/gsl-2.4/cblas/dscal.c /tmp/GSL/gsl-2.4/cblas/dsdot.c /tmp/GSL/gsl-2.4/cblas/dspmv.c /tmp/GSL/gsl-2.4/cblas/dspr.c /tmp/GSL/gsl-2.4/cblas/dspr2.c /tmp/GSL/gsl-2.4/cblas/dswap.c /tmp/GSL/gsl-2.4/cblas/dsymm.c /tmp/GSL/gsl-2.4/cblas/dsymv.c /tmp/GSL/gsl-2.4/cblas/dsyr.c /tmp/GSL/gsl-2.4/cblas/dsyr2.c /tmp/GSL/gsl-2.4/cblas/dsyr2k.c /tmp/GSL/gsl-2.4/cblas/dsyrk.c /tmp/GSL/gsl-2.4/cblas/dtbmv.c /tmp/GSL/gsl-2.4/cblas/dtbsv.c /tmp/GSL/gsl-2.4/cblas/dtpmv.c /tmp/GSL/gsl-2.4/cblas/dtpsv.c /tmp/GSL/gsl-2.4/cblas/dtrmm.c /tmp/GSL/gsl-2.4/cblas/dtrmv.c /tmp/GSL/gsl-2.4/cblas/dtrsm.c /tmp/GSL/gsl-2.4/cblas/dtrsv.c /tmp/GSL/gsl-2.4/cblas/dzasum.c /tmp/GSL/gsl-2.4/cblas/dznrm2.c /tmp/GSL/gsl-2.4/cblas/hypot.c /tmp/GSL/gsl-2.4/cblas/icamax.c /tmp/GSL/gsl-2.4/cblas/idamax.c /tmp/GSL/gsl-2.4/cblas/isamax.c /tmp/GSL/gsl-2.4/cblas/izamax.c /tmp/GSL/gsl-2.4/cblas/sasum.c /tmp/GSL/gsl-2.4/cblas/saxpy.c /tmp/GSL/gsl-2.4/cblas/scasum.c /tmp/GSL/gsl-2.4/cblas/scnrm2.c /tmp/GSL/gsl-2.4/cblas/scopy.c /tmp/GSL/gsl-2.4/cblas/sdot.c /tmp/GSL/gsl-2.4/cblas/sdsdot.c /tmp/GSL/gsl-2.4/cblas/sgbmv.c /tmp/GSL/gsl-2.4/cblas/sgemm.c /tmp/GSL/gsl-2.4/cblas/sgemv.c /tmp/GSL/gsl-2.4/cblas/sger.c /tmp/GSL/gsl-2.4/cblas/snrm2.c /tmp/GSL/gsl-2.4/cblas/srot.c /tmp/GSL/gsl-2.4/cblas/srotg.c /tmp/GSL/gsl-2.4/cblas/srotm.c /tmp/GSL/gsl-2.4/cblas/srotmg.c /tmp/GSL/gsl-2.4/cblas/ssbmv.c /tmp/GSL/gsl-2.4/cblas/sscal.c /tmp/GSL/gsl-2.4/cblas/sspmv.c /tmp/GSL/gsl-2.4/cblas/sspr.c /tmp/GSL/gsl-2.4/cblas/sspr2.c /tmp/GSL/gsl-2.4/cblas/sswap.c /tmp/GSL/gsl-2.4/cblas/ssymm.c /tmp/GSL/gsl-2.4/cblas/ssymv.c /tmp/GSL/gsl-2.4/cblas/ssyr.c /tmp/GSL/gsl-2.4/cblas/ssyr2.c /tmp/GSL/gsl-2.4/cblas/ssyr2k.c /tmp/GSL/gsl-2.4/cblas/ssyrk.c /tmp/GSL/gsl-2.4/cblas/stbmv.c /tmp/GSL/gsl-2.4/cblas/stbsv.c /tmp/GSL/gsl-2.4/cblas/stpmv.c /tmp/GSL/gsl-2.4/cblas/stpsv.c /tmp/GSL/gsl-2.4/cblas/strmm.c /tmp/GSL/gsl-2.4/cblas/strmv.c /tmp/GSL/gsl-2.4/cblas/strsm.c /tmp/GSL/gsl-2.4/cblas/strsv.c /tmp/GSL/gsl-2.4/cblas/xerbla.c /tmp/GSL/gsl-2.4/cblas/zaxpy.c /tmp/GSL/gsl-2.4/cblas/zcopy.c /tmp/GSL/gsl-2.4/cblas/zdotc_sub.c /tmp/GSL/gsl-2.4/cblas/zdotu_sub.c /tmp/GSL/gsl-2.4/cblas/zdscal.c /tmp/GSL/gsl-2.4/cblas/zgbmv.c /tmp/GSL/gsl-2.4/cblas/zgemm.c /tmp/GSL/gsl-2.4/cblas/zgemv.c /tmp/GSL/gsl-2.4/cblas/zgerc.c /tmp/GSL/gsl-2.4/cblas/zgeru.c /tmp/GSL/gsl-2.4/cblas/zhbmv.c /tmp/GSL/gsl-2.4/cblas/zhemm.c /tmp/GSL/gsl-2.4/cblas/zhemv.c /tmp/GSL/gsl-2.4/cblas/zher.c /tmp/GSL/gsl-2.4/cblas/zher2.c /tmp/GSL/gsl-2.4/cblas/zher2k.c /tmp/GSL/gsl-2.4/cblas/zherk.c /tmp/GSL/gsl-2.4/cblas/zhpmv.c /tmp/GSL/gsl-2.4/cblas/zhpr.c /tmp/GSL/gsl-2.4/cblas/zhpr2.c /tmp/GSL/gsl-2.4/cblas/zscal.c /tmp/GSL/gsl-2.4/cblas/zswap.c /tmp/GSL/gsl-2.4/cblas/zsymm.c /tmp/GSL/gsl-2.4/cblas/zsyr2k.c /tmp/GSL/gsl-2.4/cblas/zsyrk.c /tmp/GSL/gsl-2.4/cblas/ztbmv.c /tmp/GSL/gsl-2.4/cblas/ztbsv.c /tmp/GSL/gsl-2.4/cblas/ztpmv.c /tmp/GSL/gsl-2.4/cblas/ztpsv.c /tmp/GSL/gsl-2.4/cblas/ztrmm.c /tmp/GSL/gsl-2.4/cblas/ztrmv.c /tmp/GSL/gsl-2.4/cblas/ztrsm.c /tmp/GSL/gsl-2.4/cblas/ztrsv.c /tmp/GSL/gsl-2.4/cdf/beta.c /tmp/GSL/gsl-2.4/cdf/betainv.c /tmp/GSL/gsl-2.4/cdf/binomial.c /tmp/GSL/gsl-2.4/cdf/cauchy.c /tmp/GSL/gsl-2.4/cdf/cauchyinv.c /tmp/GSL/gsl-2.4/cdf/chisq.c /tmp/GSL/gsl-2.4/cdf/chisqinv.c /tmp/GSL/gsl-2.4/cdf/exponential.c /tmp/GSL/gsl-2.4/cdf/exponentialinv.c /tmp/GSL/gsl-2.4/cdf/exppow.c /tmp/GSL/gsl-2.4/cdf/fdist.c /tmp/GSL/gsl-2.4/cdf/fdistinv.c /tmp/GSL/gsl-2.4/cdf/flat.c /tmp/GSL/gsl-2.4/cdf/flatinv.c /tmp/GSL/gsl-2.4/cdf/gamma.c /tmp/GSL/gsl-2.4/cdf/gammainv.c /tmp/GSL/gsl-2.4/cdf/gauss.c /tmp/GSL/gsl-2.4/cdf/gaussinv.c /tmp/GSL/gsl-2.4/cdf/geometric.c /tmp/GSL/gsl-2.4/cdf/gumbel1.c /tmp/GSL/gsl-2.4/cdf/gumbel1inv.c /tmp/GSL/gsl-2.4/cdf/gumbel2.c /tmp/GSL/gsl-2.4/cdf/gumbel2inv.c /tmp/GSL/gsl-2.4/cdf/hypergeometric.c /tmp/GSL/gsl-2.4/cdf/laplace.c /tmp/GSL/gsl-2.4/cdf/laplaceinv.c /tmp/GSL/gsl-2.4/cdf/logistic.c /tmp/GSL/gsl-2.4/cdf/logisticinv.c /tmp/GSL/gsl-2.4/cdf/lognormal.c /tmp/GSL/gsl-2.4/cdf/lognormalinv.c /tmp/GSL/gsl-2.4/cdf/nbinomial.c /tmp/GSL/gsl-2.4/cdf/pareto.c /tmp/GSL/gsl-2.4/cdf/paretoinv.c /tmp/GSL/gsl-2.4/cdf/pascal.c /tmp/GSL/gsl-2.4/cdf/poisson.c /tmp/GSL/gsl-2.4/cdf/rayleigh.c /tmp/GSL/gsl-2.4/cdf/rayleighinv.c /tmp/GSL/gsl-2.4/cdf/tdist.c /tmp/GSL/gsl-2.4/cdf/tdistinv.c /tmp/GSL/gsl-2.4/cdf/weibull.c /tmp/GSL/gsl-2.4/cdf/weibullinv.c /tmp/GSL/gsl-2.4/cheb/deriv.c /tmp/GSL/gsl-2.4/cheb/eval.c /tmp/GSL/gsl-2.4/cheb/init.c /tmp/GSL/gsl-2.4/cheb/integ.c /tmp/GSL/gsl-2.4/combination/combination.c /tmp/GSL/gsl-2.4/combination/file.c /tmp/GSL/gsl-2.4/combination/init.c /tmp/GSL/gsl-2.4/combination/inline.c /tmp/GSL/gsl-2.4/complex/inline.c /tmp/GSL/gsl-2.4/complex/math.c /tmp/GSL/gsl-2.4/deriv/deriv.c /tmp/GSL/gsl-2.4/dht/dht.c /tmp/GSL/gsl-2.4/diff/diff.c /tmp/GSL/gsl-2.4/doc/examples/blas.c /tmp/GSL/gsl-2.4/doc/examples/block.c /tmp/GSL/gsl-2.4/doc/examples/bspline.c /tmp/GSL/gsl-2.4/doc/examples/combination.c /tmp/GSL/gsl-2.4/doc/examples/diff.c /tmp/GSL/gsl-2.4/doc/examples/dwt.c /tmp/GSL/gsl-2.4/doc/examples/fft.c /tmp/GSL/gsl-2.4/doc/examples/interp.c /tmp/GSL/gsl-2.4/doc/examples/interp2d.c /tmp/GSL/gsl-2.4/doc/examples/matrix.c /tmp/GSL/gsl-2.4/doc/examples/multiset.c /tmp/GSL/gsl-2.4/doc/examples/poisson.c /tmp/GSL/gsl-2.4/doc/examples/qrng.c /tmp/GSL/gsl-2.4/doc/examples/rng.c /tmp/GSL/gsl-2.4/doc/examples/rquantile.c /tmp/GSL/gsl-2.4/doc/examples/rstat.c /tmp/GSL/gsl-2.4/doc/examples/siman.c /tmp/GSL/gsl-2.4/doc/examples/spmatrix.c /tmp/GSL/gsl-2.4/doc/examples/stat.c /tmp/GSL/gsl-2.4/doc/examples/vector.c /tmp/GSL/gsl-2.4/eigen/francis.c /tmp/GSL/gsl-2.4/eigen/gen.c /tmp/GSL/gsl-2.4/eigen/genherm.c /tmp/GSL/gsl-2.4/eigen/genhermv.c /tmp/GSL/gsl-2.4/eigen/gensymm.c /tmp/GSL/gsl-2.4/eigen/gensymmv.c /tmp/GSL/gsl-2.4/eigen/genv.c /tmp/GSL/gsl-2.4/eigen/herm.c /tmp/GSL/gsl-2.4/eigen/hermv.c /tmp/GSL/gsl-2.4/eigen/jacobi.c /tmp/GSL/gsl-2.4/eigen/nonsymm.c /tmp/GSL/gsl-2.4/eigen/nonsymmv.c /tmp/GSL/gsl-2.4/eigen/schur.c /tmp/GSL/gsl-2.4/eigen/sort.c /tmp/GSL/gsl-2.4/eigen/symm.c /tmp/GSL/gsl-2.4/eigen/symmv.c /tmp/GSL/gsl-2.4/err/error.c /tmp/GSL/gsl-2.4/err/message.c /tmp/GSL/gsl-2.4/err/stream.c /tmp/GSL/gsl-2.4/err/strerror.c /tmp/GSL/gsl-2.4/fft/dft.c /tmp/GSL/gsl-2.4/fft/fft.c /tmp/GSL/gsl-2.4/fit/linear.c /tmp/GSL/gsl-2.4/histogram/add.c /tmp/GSL/gsl-2.4/histogram/add2d.c /tmp/GSL/gsl-2.4/histogram/calloc_range.c /tmp/GSL/gsl-2.4/histogram/calloc_range2d.c /tmp/GSL/gsl-2.4/histogram/copy.c /tmp/GSL/gsl-2.4/histogram/copy2d.c /tmp/GSL/gsl-2.4/histogram/file.c /tmp/GSL/gsl-2.4/histogram/file2d.c /tmp/GSL/gsl-2.4/histogram/get.c /tmp/GSL/gsl-2.4/histogram/get2d.c /tmp/GSL/gsl-2.4/histogram/init.c /tmp/GSL/gsl-2.4/histogram/init2d.c /tmp/GSL/gsl-2.4/histogram/maxval.c /tmp/GSL/gsl-2.4/histogram/maxval2d.c /tmp/GSL/gsl-2.4/histogram/oper.c /tmp/GSL/gsl-2.4/histogram/oper2d.c /tmp/GSL/gsl-2.4/histogram/params.c /tmp/GSL/gsl-2.4/histogram/params2d.c /tmp/GSL/gsl-2.4/histogram/pdf.c /tmp/GSL/gsl-2.4/histogram/pdf2d.c /tmp/GSL/gsl-2.4/histogram/reset.c /tmp/GSL/gsl-2.4/histogram/reset2d.c /tmp/GSL/gsl-2.4/histogram/stat.c /tmp/GSL/gsl-2.4/histogram/stat2d.c /tmp/GSL/gsl-2.4/ieee-utils/env.c /tmp/GSL/gsl-2.4/ieee-utils/fp.c /tmp/GSL/gsl-2.4/ieee-utils/make_rep.c /tmp/GSL/gsl-2.4/ieee-utils/print.c /tmp/GSL/gsl-2.4/ieee-utils/read.c /tmp/GSL/gsl-2.4/integration/chebyshev.c /tmp/GSL/gsl-2.4/integration/chebyshev2.c /tmp/GSL/gsl-2.4/integration/cquad.c /tmp/GSL/gsl-2.4/integration/exponential.c /tmp/GSL/gsl-2.4/integration/fixed.c /tmp/GSL/gsl-2.4/integration/gegenbauer.c /tmp/GSL/gsl-2.4/integration/glfixed.c /tmp/GSL/gsl-2.4/integration/hermite.c /tmp/GSL/gsl-2.4/integration/jacobi.c /tmp/GSL/gsl-2.4/integration/laguerre.c /tmp/GSL/gsl-2.4/integration/legendre.c /tmp/GSL/gsl-2.4/integration/qag.c /tmp/GSL/gsl-2.4/integration/qagp.c /tmp/GSL/gsl-2.4/integration/qags.c /tmp/GSL/gsl-2.4/integration/qawc.c /tmp/GSL/gsl-2.4/integration/qawf.c /tmp/GSL/gsl-2.4/integration/qawo.c /tmp/GSL/gsl-2.4/integration/qaws.c /tmp/GSL/gsl-2.4/integration/qcheb.c /tmp/GSL/gsl-2.4/integration/qk.c /tmp/GSL/gsl-2.4/integration/qk15.c /tmp/GSL/gsl-2.4/integration/qk21.c /tmp/GSL/gsl-2.4/integration/qk31.c /tmp/GSL/gsl-2.4/integration/qk41.c /tmp/GSL/gsl-2.4/integration/qk51.c /tmp/GSL/gsl-2.4/integration/qk61.c /tmp/GSL/gsl-2.4/integration/qmomo.c /tmp/GSL/gsl-2.4/integration/qmomof.c /tmp/GSL/gsl-2.4/integration/qng.c /tmp/GSL/gsl-2.4/integration/rational.c /tmp/GSL/gsl-2.4/integration/workspace.c /tmp/GSL/gsl-2.4/interpolation/accel.c /tmp/GSL/gsl-2.4/interpolation/bicubic.c /tmp/GSL/gsl-2.4/interpolation/bilinear.c /tmp/GSL/gsl-2.4/interpolation/cspline.c /tmp/GSL/gsl-2.4/interpolation/inline.c /tmp/GSL/gsl-2.4/interpolation/interp.c /tmp/GSL/gsl-2.4/interpolation/interp2d.c /tmp/GSL/gsl-2.4/interpolation/linear.c /tmp/GSL/gsl-2.4/interpolation/poly.c /tmp/GSL/gsl-2.4/interpolation/spline.c /tmp/GSL/gsl-2.4/interpolation/spline2d.c /tmp/GSL/gsl-2.4/interpolation/steffen.c /tmp/GSL/gsl-2.4/linalg/balance.c /tmp/GSL/gsl-2.4/linalg/balancemat.c /tmp/GSL/gsl-2.4/linalg/bidiag.c /tmp/GSL/gsl-2.4/linalg/cholesky.c /tmp/GSL/gsl-2.4/linalg/choleskyc.c /tmp/GSL/gsl-2.4/linalg/cod.c /tmp/GSL/gsl-2.4/linalg/condest.c /tmp/GSL/gsl-2.4/linalg/exponential.c /tmp/GSL/gsl-2.4/linalg/hermtd.c /tmp/GSL/gsl-2.4/linalg/hessenberg.c /tmp/GSL/gsl-2.4/linalg/hesstri.c /tmp/GSL/gsl-2.4/linalg/hh.c /tmp/GSL/gsl-2.4/linalg/householder.c /tmp/GSL/gsl-2.4/linalg/householdercomplex.c /tmp/GSL/gsl-2.4/linalg/inline.c /tmp/GSL/gsl-2.4/linalg/invtri.c /tmp/GSL/gsl-2.4/linalg/lq.c /tmp/GSL/gsl-2.4/linalg/lu.c /tmp/GSL/gsl-2.4/linalg/luc.c /tmp/GSL/gsl-2.4/linalg/mcholesky.c /tmp/GSL/gsl-2.4/linalg/multiply.c /tmp/GSL/gsl-2.4/linalg/pcholesky.c /tmp/GSL/gsl-2.4/linalg/ptlq.c /tmp/GSL/gsl-2.4/linalg/qr.c /tmp/GSL/gsl-2.4/linalg/qrpt.c /tmp/GSL/gsl-2.4/linalg/svd.c /tmp/GSL/gsl-2.4/linalg/symmtd.c /tmp/GSL/gsl-2.4/linalg/tridiag.c /tmp/GSL/gsl-2.4/matrix/copy.c /tmp/GSL/gsl-2.4/matrix/file.c /tmp/GSL/gsl-2.4/matrix/getset.c /tmp/GSL/gsl-2.4/matrix/init.c /tmp/GSL/gsl-2.4/matrix/matrix.c /tmp/GSL/gsl-2.4/matrix/oper.c /tmp/GSL/gsl-2.4/matrix/prop.c /tmp/GSL/gsl-2.4/matrix/rowcol.c /tmp/GSL/gsl-2.4/matrix/submatrix.c /tmp/GSL/gsl-2.4/matrix/swap.c /tmp/GSL/gsl-2.4/matrix/view.c /tmp/GSL/gsl-2.4/min/bracketing.c /tmp/GSL/gsl-2.4/min/brent.c /tmp/GSL/gsl-2.4/min/convergence.c /tmp/GSL/gsl-2.4/min/fsolver.c /tmp/GSL/gsl-2.4/min/golden.c /tmp/GSL/gsl-2.4/min/quad_golden.c /tmp/GSL/gsl-2.4/monte/miser.c /tmp/GSL/gsl-2.4/monte/plain.c /tmp/GSL/gsl-2.4/multifit/convergence.c /tmp/GSL/gsl-2.4/multifit/covar.c /tmp/GSL/gsl-2.4/multifit/fdfridge.c /tmp/GSL/gsl-2.4/multifit/fdfsolver.c /tmp/GSL/gsl-2.4/multifit/fdjac.c /tmp/GSL/gsl-2.4/multifit/fsolver.c /tmp/GSL/gsl-2.4/multifit/gcv.c /tmp/GSL/gsl-2.4/multifit/gradient.c /tmp/GSL/gsl-2.4/multifit/lmder.c /tmp/GSL/gsl-2.4/multifit/lmniel.c /tmp/GSL/gsl-2.4/multifit/multilinear.c /tmp/GSL/gsl-2.4/multifit/multireg.c /tmp/GSL/gsl-2.4/multifit/multirobust.c /tmp/GSL/gsl-2.4/multifit/multiwlinear.c /tmp/GSL/gsl-2.4/multifit/robust_wfun.c /tmp/GSL/gsl-2.4/multifit/work.c /tmp/GSL/gsl-2.4/multifit_nlinear/cholesky.c /tmp/GSL/gsl-2.4/multifit_nlinear/convergence.c /tmp/GSL/gsl-2.4/multifit_nlinear/covar.c /tmp/GSL/gsl-2.4/multifit_nlinear/dogleg.c /tmp/GSL/gsl-2.4/multifit_nlinear/fdf.c /tmp/GSL/gsl-2.4/multifit_nlinear/fdfvv.c /tmp/GSL/gsl-2.4/multifit_nlinear/fdjac.c /tmp/GSL/gsl-2.4/multifit_nlinear/lm.c /tmp/GSL/gsl-2.4/multifit_nlinear/qr.c /tmp/GSL/gsl-2.4/multifit_nlinear/scaling.c /tmp/GSL/gsl-2.4/multifit_nlinear/subspace2D.c /tmp/GSL/gsl-2.4/multifit_nlinear/svd.c /tmp/GSL/gsl-2.4/multifit_nlinear/trust.c /tmp/GSL/gsl-2.4/multilarge/multilarge.c /tmp/GSL/gsl-2.4/multilarge/normal.c /tmp/GSL/gsl-2.4/multilarge/tsqr.c /tmp/GSL/gsl-2.4/multilarge_nlinear/cgst.c /tmp/GSL/gsl-2.4/multilarge_nlinear/cholesky.c /tmp/GSL/gsl-2.4/multilarge_nlinear/convergence.c /tmp/GSL/gsl-2.4/multilarge_nlinear/dogleg.c /tmp/GSL/gsl-2.4/multilarge_nlinear/dummy.c /tmp/GSL/gsl-2.4/multilarge_nlinear/fdf.c /tmp/GSL/gsl-2.4/multilarge_nlinear/lm.c /tmp/GSL/gsl-2.4/multilarge_nlinear/scaling.c /tmp/GSL/gsl-2.4/multilarge_nlinear/subspace2D.c /tmp/GSL/gsl-2.4/multilarge_nlinear/trust.c /tmp/GSL/gsl-2.4/multimin/conjugate_fr.c /tmp/GSL/gsl-2.4/multimin/conjugate_pr.c /tmp/GSL/gsl-2.4/multimin/convergence.c /tmp/GSL/gsl-2.4/multimin/diff.c /tmp/GSL/gsl-2.4/multimin/fdfminimizer.c /tmp/GSL/gsl-2.4/multimin/fminimizer.c /tmp/GSL/gsl-2.4/multimin/simplex.c /tmp/GSL/gsl-2.4/multimin/simplex2.c /tmp/GSL/gsl-2.4/multimin/steepest_descent.c /tmp/GSL/gsl-2.4/multimin/vector_bfgs.c /tmp/GSL/gsl-2.4/multimin/vector_bfgs2.c /tmp/GSL/gsl-2.4/multiroots/broyden.c /tmp/GSL/gsl-2.4/multiroots/convergence.c /tmp/GSL/gsl-2.4/multiroots/dnewton.c /tmp/GSL/gsl-2.4/multiroots/fdfsolver.c /tmp/GSL/gsl-2.4/multiroots/fdjac.c /tmp/GSL/gsl-2.4/multiroots/fsolver.c /tmp/GSL/gsl-2.4/multiroots/gnewton.c /tmp/GSL/gsl-2.4/multiroots/hybrid.c /tmp/GSL/gsl-2.4/multiroots/hybridj.c /tmp/GSL/gsl-2.4/multiroots/newton.c /tmp/GSL/gsl-2.4/multiset/file.c /tmp/GSL/gsl-2.4/multiset/init.c /tmp/GSL/gsl-2.4/multiset/inline.c /tmp/GSL/gsl-2.4/multiset/multiset.c /tmp/GSL/gsl-2.4/ntuple/ntuple.c /tmp/GSL/gsl-2.4/ode-initval/bsimp.c /tmp/GSL/gsl-2.4/ode-initval/control.c /tmp/GSL/gsl-2.4/ode-initval/cscal.c /tmp/GSL/gsl-2.4/ode-initval/cstd.c /tmp/GSL/gsl-2.4/ode-initval/evolve.c /tmp/GSL/gsl-2.4/ode-initval/gear1.c /tmp/GSL/gsl-2.4/ode-initval/gear2.c /tmp/GSL/gsl-2.4/ode-initval/rk2.c /tmp/GSL/gsl-2.4/ode-initval/rk2imp.c /tmp/GSL/gsl-2.4/ode-initval/rk2simp.c /tmp/GSL/gsl-2.4/ode-initval/rk4.c /tmp/GSL/gsl-2.4/ode-initval/rk4imp.c /tmp/GSL/gsl-2.4/ode-initval/rk8pd.c /tmp/GSL/gsl-2.4/ode-initval/rkck.c /tmp/GSL/gsl-2.4/ode-initval/rkf45.c /tmp/GSL/gsl-2.4/ode-initval/step.c /tmp/GSL/gsl-2.4/ode-initval2/bsimp.c /tmp/GSL/gsl-2.4/ode-initval2/control.c /tmp/GSL/gsl-2.4/ode-initval2/cscal.c /tmp/GSL/gsl-2.4/ode-initval2/cstd.c /tmp/GSL/gsl-2.4/ode-initval2/driver.c /tmp/GSL/gsl-2.4/ode-initval2/evolve.c /tmp/GSL/gsl-2.4/ode-initval2/msadams.c /tmp/GSL/gsl-2.4/ode-initval2/msbdf.c /tmp/GSL/gsl-2.4/ode-initval2/rk1imp.c /tmp/GSL/gsl-2.4/ode-initval2/rk2.c /tmp/GSL/gsl-2.4/ode-initval2/rk2imp.c /tmp/GSL/gsl-2.4/ode-initval2/rk4.c /tmp/GSL/gsl-2.4/ode-initval2/rk4imp.c /tmp/GSL/gsl-2.4/ode-initval2/rk8pd.c /tmp/GSL/gsl-2.4/ode-initval2/rkck.c /tmp/GSL/gsl-2.4/ode-initval2/rkf45.c /tmp/GSL/gsl-2.4/ode-initval2/step.c /tmp/GSL/gsl-2.4/permutation/canonical.c /tmp/GSL/gsl-2.4/permutation/file.c /tmp/GSL/gsl-2.4/permutation/init.c /tmp/GSL/gsl-2.4/permutation/inline.c /tmp/GSL/gsl-2.4/permutation/permutation.c /tmp/GSL/gsl-2.4/permutation/permute.c /tmp/GSL/gsl-2.4/poly/dd.c /tmp/GSL/gsl-2.4/poly/deriv.c /tmp/GSL/gsl-2.4/poly/eval.c /tmp/GSL/gsl-2.4/poly/solve_cubic.c /tmp/GSL/gsl-2.4/poly/solve_quadratic.c /tmp/GSL/gsl-2.4/poly/zsolve.c /tmp/GSL/gsl-2.4/poly/zsolve_cubic.c /tmp/GSL/gsl-2.4/poly/zsolve_init.c /tmp/GSL/gsl-2.4/poly/zsolve_quadratic.c /tmp/GSL/gsl-2.4/qrng/halton.c /tmp/GSL/gsl-2.4/qrng/inline.c /tmp/GSL/gsl-2.4/qrng/niederreiter-2.c /tmp/GSL/gsl-2.4/qrng/qrng.c /tmp/GSL/gsl-2.4/qrng/reversehalton.c /tmp/GSL/gsl-2.4/qrng/sobol.c /tmp/GSL/gsl-2.4/randist/bernoulli.c /tmp/GSL/gsl-2.4/randist/beta.c /tmp/GSL/gsl-2.4/randist/bigauss.c /tmp/GSL/gsl-2.4/randist/binomial.c /tmp/GSL/gsl-2.4/randist/binomial_tpe.c /tmp/GSL/gsl-2.4/randist/cauchy.c /tmp/GSL/gsl-2.4/randist/chisq.c /tmp/GSL/gsl-2.4/randist/dirichlet.c /tmp/GSL/gsl-2.4/randist/discrete.c /tmp/GSL/gsl-2.4/randist/erlang.c /tmp/GSL/gsl-2.4/randist/exponential.c /tmp/GSL/gsl-2.4/randist/exppow.c /tmp/GSL/gsl-2.4/randist/fdist.c /tmp/GSL/gsl-2.4/randist/flat.c /tmp/GSL/gsl-2.4/randist/gamma.c /tmp/GSL/gsl-2.4/randist/gauss.c /tmp/GSL/gsl-2.4/randist/gausstail.c /tmp/GSL/gsl-2.4/randist/gausszig.c /tmp/GSL/gsl-2.4/randist/geometric.c /tmp/GSL/gsl-2.4/randist/gumbel.c /tmp/GSL/gsl-2.4/randist/hyperg.c /tmp/GSL/gsl-2.4/randist/landau.c /tmp/GSL/gsl-2.4/randist/laplace.c /tmp/GSL/gsl-2.4/randist/levy.c /tmp/GSL/gsl-2.4/randist/logarithmic.c /tmp/GSL/gsl-2.4/randist/logistic.c /tmp/GSL/gsl-2.4/randist/lognormal.c /tmp/GSL/gsl-2.4/randist/multinomial.c /tmp/GSL/gsl-2.4/randist/mvgauss.c /tmp/GSL/gsl-2.4/randist/nbinomial.c /tmp/GSL/gsl-2.4/randist/pareto.c /tmp/GSL/gsl-2.4/randist/pascal.c /tmp/GSL/gsl-2.4/randist/poisson.c /tmp/GSL/gsl-2.4/randist/rayleigh.c /tmp/GSL/gsl-2.4/randist/shuffle.c /tmp/GSL/gsl-2.4/randist/sphere.c /tmp/GSL/gsl-2.4/randist/tdist.c /tmp/GSL/gsl-2.4/randist/weibull.c /tmp/GSL/gsl-2.4/rng/borosh13.c /tmp/GSL/gsl-2.4/rng/cmrg.c /tmp/GSL/gsl-2.4/rng/coveyou.c /tmp/GSL/gsl-2.4/rng/default.c /tmp/GSL/gsl-2.4/rng/file.c /tmp/GSL/gsl-2.4/rng/fishman18.c /tmp/GSL/gsl-2.4/rng/fishman20.c /tmp/GSL/gsl-2.4/rng/fishman2x.c /tmp/GSL/gsl-2.4/rng/gfsr4.c /tmp/GSL/gsl-2.4/rng/inline.c /tmp/GSL/gsl-2.4/rng/knuthran.c /tmp/GSL/gsl-2.4/rng/knuthran2.c /tmp/GSL/gsl-2.4/rng/knuthran2002.c /tmp/GSL/gsl-2.4/rng/lecuyer21.c /tmp/GSL/gsl-2.4/rng/minstd.c /tmp/GSL/gsl-2.4/rng/mrg.c /tmp/GSL/gsl-2.4/rng/mt.c /tmp/GSL/gsl-2.4/rng/r250.c /tmp/GSL/gsl-2.4/rng/ran0.c /tmp/GSL/gsl-2.4/rng/ran1.c /tmp/GSL/gsl-2.4/rng/ran2.c /tmp/GSL/gsl-2.4/rng/ran3.c /tmp/GSL/gsl-2.4/rng/rand.c /tmp/GSL/gsl-2.4/rng/rand48.c /tmp/GSL/gsl-2.4/rng/random.c /tmp/GSL/gsl-2.4/rng/randu.c /tmp/GSL/gsl-2.4/rng/ranf.c /tmp/GSL/gsl-2.4/rng/ranlux.c /tmp/GSL/gsl-2.4/rng/ranlxd.c /tmp/GSL/gsl-2.4/rng/ranlxs.c /tmp/GSL/gsl-2.4/rng/ranmar.c /tmp/GSL/gsl-2.4/rng/rng.c /tmp/GSL/gsl-2.4/rng/slatec.c /tmp/GSL/gsl-2.4/rng/taus.c /tmp/GSL/gsl-2.4/rng/taus113.c /tmp/GSL/gsl-2.4/rng/transputer.c /tmp/GSL/gsl-2.4/rng/tt.c /tmp/GSL/gsl-2.4/rng/types.c /tmp/GSL/gsl-2.4/rng/uni.c /tmp/GSL/gsl-2.4/rng/uni32.c /tmp/GSL/gsl-2.4/rng/vax.c /tmp/GSL/gsl-2.4/rng/waterman14.c /tmp/GSL/gsl-2.4/rng/zuf.c /tmp/GSL/gsl-2.4/roots/bisection.c /tmp/GSL/gsl-2.4/roots/brent.c /tmp/GSL/gsl-2.4/roots/convergence.c /tmp/GSL/gsl-2.4/roots/falsepos.c /tmp/GSL/gsl-2.4/roots/fdfsolver.c /tmp/GSL/gsl-2.4/roots/fsolver.c /tmp/GSL/gsl-2.4/roots/newton.c /tmp/GSL/gsl-2.4/roots/secant.c /tmp/GSL/gsl-2.4/roots/steffenson.c /tmp/GSL/gsl-2.4/rstat/rquantile.c /tmp/GSL/gsl-2.4/rstat/rstat.c /tmp/GSL/gsl-2.4/siman/siman.c /tmp/GSL/gsl-2.4/sort/sort.c /tmp/GSL/gsl-2.4/sort/sortind.c /tmp/GSL/gsl-2.4/sort/sortvec.c /tmp/GSL/gsl-2.4/sort/sortvecind.c /tmp/GSL/gsl-2.4/sort/subset.c /tmp/GSL/gsl-2.4/sort/subsetind.c /tmp/GSL/gsl-2.4/spblas/spdgemm.c /tmp/GSL/gsl-2.4/spblas/spdgemv.c /tmp/GSL/gsl-2.4/specfunc/airy.c /tmp/GSL/gsl-2.4/specfunc/airy_der.c /tmp/GSL/gsl-2.4/specfunc/airy_zero.c /tmp/GSL/gsl-2.4/specfunc/atanint.c /tmp/GSL/gsl-2.4/specfunc/bessel.c /tmp/GSL/gsl-2.4/specfunc/bessel_I0.c /tmp/GSL/gsl-2.4/specfunc/bessel_I1.c /tmp/GSL/gsl-2.4/specfunc/bessel_In.c /tmp/GSL/gsl-2.4/specfunc/bessel_Inu.c /tmp/GSL/gsl-2.4/specfunc/bessel_J0.c /tmp/GSL/gsl-2.4/specfunc/bessel_J1.c /tmp/GSL/gsl-2.4/specfunc/bessel_Jn.c /tmp/GSL/gsl-2.4/specfunc/bessel_Jnu.c /tmp/GSL/gsl-2.4/specfunc/bessel_K0.c /tmp/GSL/gsl-2.4/specfunc/bessel_K1.c /tmp/GSL/gsl-2.4/specfunc/bessel_Kn.c /tmp/GSL/gsl-2.4/specfunc/bessel_Knu.c /tmp/GSL/gsl-2.4/specfunc/bessel_Y0.c /tmp/GSL/gsl-2.4/specfunc/bessel_Y1.c /tmp/GSL/gsl-2.4/specfunc/bessel_Yn.c /tmp/GSL/gsl-2.4/specfunc/bessel_Ynu.c /tmp/GSL/gsl-2.4/specfunc/bessel_amp_phase.c /tmp/GSL/gsl-2.4/specfunc/bessel_i.c /tmp/GSL/gsl-2.4/specfunc/bessel_j.c /tmp/GSL/gsl-2.4/specfunc/bessel_k.c /tmp/GSL/gsl-2.4/specfunc/bessel_olver.c /tmp/GSL/gsl-2.4/specfunc/bessel_sequence.c /tmp/GSL/gsl-2.4/specfunc/bessel_temme.c /tmp/GSL/gsl-2.4/specfunc/bessel_y.c /tmp/GSL/gsl-2.4/specfunc/bessel_zero.c /tmp/GSL/gsl-2.4/specfunc/beta.c /tmp/GSL/gsl-2.4/specfunc/beta_inc.c /tmp/GSL/gsl-2.4/specfunc/clausen.c /tmp/GSL/gsl-2.4/specfunc/coulomb.c /tmp/GSL/gsl-2.4/specfunc/coulomb_bound.c /tmp/GSL/gsl-2.4/specfunc/coupling.c /tmp/GSL/gsl-2.4/specfunc/dawson.c /tmp/GSL/gsl-2.4/specfunc/debye.c /tmp/GSL/gsl-2.4/specfunc/dilog.c /tmp/GSL/gsl-2.4/specfunc/elementary.c /tmp/GSL/gsl-2.4/specfunc/ellint.c /tmp/GSL/gsl-2.4/specfunc/elljac.c /tmp/GSL/gsl-2.4/specfunc/erfc.c /tmp/GSL/gsl-2.4/specfunc/exp.c /tmp/GSL/gsl-2.4/specfunc/expint.c /tmp/GSL/gsl-2.4/specfunc/expint3.c /tmp/GSL/gsl-2.4/specfunc/fermi_dirac.c /tmp/GSL/gsl-2.4/specfunc/gamma.c /tmp/GSL/gsl-2.4/specfunc/gamma_inc.c /tmp/GSL/gsl-2.4/specfunc/gegenbauer.c /tmp/GSL/gsl-2.4/specfunc/hermite.c /tmp/GSL/gsl-2.4/specfunc/hyperg.c /tmp/GSL/gsl-2.4/specfunc/hyperg_0F1.c /tmp/GSL/gsl-2.4/specfunc/hyperg_1F1.c /tmp/GSL/gsl-2.4/specfunc/hyperg_2F0.c /tmp/GSL/gsl-2.4/specfunc/hyperg_2F1.c /tmp/GSL/gsl-2.4/specfunc/hyperg_U.c /tmp/GSL/gsl-2.4/specfunc/laguerre.c /tmp/GSL/gsl-2.4/specfunc/lambert.c /tmp/GSL/gsl-2.4/specfunc/legendre_H3d.c /tmp/GSL/gsl-2.4/specfunc/legendre_Qn.c /tmp/GSL/gsl-2.4/specfunc/legendre_con.c /tmp/GSL/gsl-2.4/specfunc/legendre_poly.c /tmp/GSL/gsl-2.4/specfunc/log.c /tmp/GSL/gsl-2.4/specfunc/mathieu_angfunc.c /tmp/GSL/gsl-2.4/specfunc/mathieu_charv.c /tmp/GSL/gsl-2.4/specfunc/mathieu_coeff.c /tmp/GSL/gsl-2.4/specfunc/mathieu_radfunc.c /tmp/GSL/gsl-2.4/specfunc/mathieu_workspace.c /tmp/GSL/gsl-2.4/specfunc/poch.c /tmp/GSL/gsl-2.4/specfunc/pow_int.c /tmp/GSL/gsl-2.4/specfunc/psi.c /tmp/GSL/gsl-2.4/specfunc/result.c /tmp/GSL/gsl-2.4/specfunc/shint.c /tmp/GSL/gsl-2.4/specfunc/sinint.c /tmp/GSL/gsl-2.4/specfunc/synchrotron.c /tmp/GSL/gsl-2.4/specfunc/transport.c /tmp/GSL/gsl-2.4/specfunc/trig.c /tmp/GSL/gsl-2.4/specfunc/zeta.c /tmp/GSL/gsl-2.4/splinalg/gmres.c /tmp/GSL/gsl-2.4/splinalg/itersolve.c /tmp/GSL/gsl-2.4/spmatrix/spcompress.c /tmp/GSL/gsl-2.4/spmatrix/spio.c /tmp/GSL/gsl-2.4/spmatrix/spoper.c /tmp/GSL/gsl-2.4/spmatrix/spprop.c /tmp/GSL/gsl-2.4/statistics/absdev.c /tmp/GSL/gsl-2.4/statistics/covariance.c /tmp/GSL/gsl-2.4/statistics/kurtosis.c /tmp/GSL/gsl-2.4/statistics/lag1.c /tmp/GSL/gsl-2.4/statistics/mean.c /tmp/GSL/gsl-2.4/statistics/median.c /tmp/GSL/gsl-2.4/statistics/p_variance.c /tmp/GSL/gsl-2.4/statistics/quantiles.c /tmp/GSL/gsl-2.4/statistics/skew.c /tmp/GSL/gsl-2.4/statistics/ttest.c /tmp/GSL/gsl-2.4/statistics/variance.c /tmp/GSL/gsl-2.4/statistics/wabsdev.c /tmp/GSL/gsl-2.4/statistics/wkurtosis.c /tmp/GSL/gsl-2.4/statistics/wmean.c /tmp/GSL/gsl-2.4/statistics/wskew.c /tmp/GSL/gsl-2.4/statistics/wvariance.c /tmp/GSL/gsl-2.4/sum/levin_u.c /tmp/GSL/gsl-2.4/sum/levin_utrunc.c /tmp/GSL/gsl-2.4/sum/work_u.c /tmp/GSL/gsl-2.4/sum/work_utrunc.c /tmp/GSL/gsl-2.4/sys/coerce.c /tmp/GSL/gsl-2.4/sys/expm1.c /tmp/GSL/gsl-2.4/sys/fcmp.c /tmp/GSL/gsl-2.4/sys/fdiv.c /tmp/GSL/gsl-2.4/sys/hypot.c /tmp/GSL/gsl-2.4/sys/infnan.c /tmp/GSL/gsl-2.4/sys/invhyp.c /tmp/GSL/gsl-2.4/sys/ldfrexp.c /tmp/GSL/gsl-2.4/sys/log1p.c /tmp/GSL/gsl-2.4/sys/minmax.c /tmp/GSL/gsl-2.4/sys/pow_int.c /tmp/GSL/gsl-2.4/sys/prec.c /tmp/GSL/gsl-2.4/test/results.c /tmp/GSL/gsl-2.4/utils/placeholder.c /tmp/GSL/gsl-2.4/vector/copy.c /tmp/GSL/gsl-2.4/vector/file.c /tmp/GSL/gsl-2.4/vector/init.c /tmp/GSL/gsl-2.4/vector/oper.c /tmp/GSL/gsl-2.4/vector/prop.c /tmp/GSL/gsl-2.4/vector/reim.c /tmp/GSL/gsl-2.4/vector/subvector.c /tmp/GSL/gsl-2.4/vector/swap.c /tmp/GSL/gsl-2.4/vector/vector.c /tmp/GSL/gsl-2.4/vector/view.c /tmp/GSL/gsl-2.4/version.c /tmp/GSL/gsl-2.4/wavelet/bspline.c /tmp/GSL/gsl-2.4/wavelet/daubechies.c /tmp/GSL/gsl-2.4/wavelet/dwt.c /tmp/GSL/gsl-2.4/wavelet/haar.c /tmp/GSL/gsl-2.4/wavelet/wavelet.c
