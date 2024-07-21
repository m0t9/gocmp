package main

import (
	"flag"
	"fmt"
	"go-compressor/internal"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

const (
	msgArgsMissing          = "(⁎˃ᆺ˂) source and (or) destination files are missing\n"
	msgSrcFileNotOpen       = "(⁎˃ᆺ˂) source file '%s' can not be open: %s\n"
	msgDstFileNotCreated    = "(⁎˃ᆺ˂) output file '%s' can not be created: %s\n"
	msgCompressionFailed    = "(⁎˃ᆺ˂) can not compress: %s\n"
	msgDecompressionFailed  = "(⁎˃ᆺ˂) can not decompress: %s\n"
	msgCompressionSuccess   = "(=^ ◡ ^=) successfully compressed to file '%s'\n"
	msgDecompressionSuccess = "(=^ ◡ ^=) successfully decompressed to file '%s'\n"
	msgCompressionRate      = "( ^..^)ﾉ  compression rate is %.2f\n"
	msgRuntime              = "(^･o･^)ﾉ  gocmp running time is %s\n"
)

var (
	cpuprofile     = flag.String("cpuprofile", "", "write cpu profile to this file")
	decompressMode = flag.Bool("d", false, "enable decompression mode")
)

func main() {
	flag.Parse()

	// debugging features
	if *cpuprofile != "" {
		cpuf, _ := os.Create(*cpuprofile)
		pprof.StartCPUProfile(cpuf)
		defer cpuf.Close()
		defer pprof.StopCPUProfile()
	}

	args := flag.Args()
	if len(args) != 2 {
		fmt.Print(msgArgsMissing)
		os.Exit(-1)
	}

	srcPath := args[0]
	dstPath := args[1]

	srcName := filepath.Base(srcPath)
	dstName := filepath.Base(dstPath)

	inf, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf(msgSrcFileNotOpen, srcName, err)
		os.Exit(-1)
	}

	outf, err := os.Create(dstPath)
	if err != nil {
		fmt.Printf(msgDstFileNotCreated, dstName, err)
		os.Exit(-1)
	}

	enc := internal.NewHuffmanEncoderDecoder()

	startTime := time.Now()
	if *decompressMode {
		if err = enc.Decode(inf, outf); err != nil {
			fmt.Printf(msgDecompressionFailed, err)
			_ = os.Remove(dstPath)
			os.Exit(-1)
		}
		fmt.Printf(msgDecompressionSuccess, dstName)
	} else {
		if err = enc.Encode(inf, outf); err != nil {
			fmt.Printf(msgCompressionFailed, err)
			_ = os.Remove(dstPath)
			os.Exit(-1)
		}
		fmt.Printf(msgCompressionSuccess, dstName)
		infStat, infErr := inf.Stat()
		outfStat, outfErr := outf.Stat()
		if infErr == nil && outfErr == nil {
			fmt.Printf(msgCompressionRate, float64(infStat.Size())/float64(outfStat.Size()))
		}
	}
	finishTime := time.Now()
	fmt.Printf(msgRuntime, finishTime.Sub(startTime))
}
