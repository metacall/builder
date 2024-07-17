package staging

import (
	"errors"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

func validateArgs(args []string) (string, error) {
	cmdArgs := []string{}
	for _, arg := range args {
		lang, ok := languageMap[arg]
		if !ok {
			return "", errors.New("Invalid language: " + arg)
		}
		isExists := false
		for _, str := range cmdArgs {
			if lang == str {
				isExists = true
			}
		}
		if !isExists {
			cmdArgs = append(cmdArgs, lang)
		}
	}
	return strings.Join(cmdArgs, " "), nil
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

// func copyForStates(src llb.State, dst llb.State, srcpaths []string, dstpath string) llb.State {
// 	return dst.With(
// 		copyMultiple(src, srcpaths, dstpath),
// 	)
// }

// func copyMultiple(src llb.State, srcPaths []string, destPath string) llb.StateOption {
// 	var stateOptions []llb.StateOption
// 	for _, srcPath := range srcPaths {
// 		stateOptions = append(stateOptions, copyFrom(src, srcPath, destPath))
// 	}

// 	return func(s llb.State) llb.State {
// 		for _, stateOption := range stateOptions {
// 			s = stateOption(s)
// 		}
// 		return s
// 	}
// }
