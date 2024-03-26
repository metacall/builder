package staging

import (
	"errors"
	"github.com/moby/buildkit/client/llb"
)

func copyForStates(src llb.State, dst llb.State, srcpaths []string, dstpath string) llb.State {
	return dst.With(
		copyMultiple(src, srcpaths, dstpath),
	)
}

func validateArgs(args []string) error {
	for _, arg := range args {
		if _, ok := languageMap[arg]; !ok {
			return errors.New("Invalid language")
		}
	}
	return nil
}

func copyMultiple(src llb.State, srcPaths []string, destPath string) llb.StateOption {
	var stateOptions []llb.StateOption
	for _, srcPath := range srcPaths {
		stateOptions = append(stateOptions, copyFrom(src, srcPath, destPath))
	}

	return func(s llb.State) llb.State {
		for _, stateOption := range stateOptions {
			s = stateOption(s)
		}
		return s
	}
}

func copyFrom(src llb.State, srcPath, destPath string) llb.StateOption {
	return func(s llb.State) llb.State {
		return copy(src, srcPath, s, destPath)
	}
}

func copy(src llb.State, srcPath string, dest llb.State, destPath string) llb.State {
	return dest.File(llb.Copy(src, srcPath, destPath, &llb.CopyInfo{
		AllowWildcard:  true,
		AttemptUnpack:  true,
		CreateDestPath: true,
	}))
}
