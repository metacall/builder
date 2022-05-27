package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/moby/buildkit/client/llb"
	_ "github.com/moby/buildkit/util/progress"
)

//installs the dependencies for metacall_base_runtime
func builDeps(baseImg llb.State) llb.State {
	deps := baseImg.Run(llb.Shlex("apt-get update")).
		Run(llb.Shlex("apt-get install -y --no-install-recommends build-essential git cmake libgtest-dev wget apt-utils apt-transport-https gnupg dirmngr ca-certificates")).Root()
	return deps
}

//builds metacall base runtime without any languages
func buildmetacallbaseRuntime(depsImg llb.State) llb.State {
	dev := depsImg.Run(llb.Shlex("git clone --depth 1 --single-branch --branch=develop https://github.com/metacall/core.git")).
		Run(llb.Shlex("mkdir core/build")).
		Dir("core/build").
		Run(llb.Shlex("cmake -Wno-dev -DOPTION_BUILD_LOG_PRETTY=Off -DOPTION_BUILD_LOADERS=On -DOPTION_GIT_HOOKS=Off ..")).
		Run(llb.Shlexf("make -j %v ", runtime.NumCPU())).Root()
	return dev
}

//builds metacall_python_runtime to support python scripts at runtime
func buildMetacallPythonRuntime(metacallBaseRuntime llb.State) llb.State {
	metacallPyRuntime := metacallBaseRuntime.Run(llb.Shlex("apt-get install -y --no-install-recommends python3-dev python3-pip")).
		Run(llb.Shlex("cmake -DOPTION_BUILD_LOADERS_PY=On -DOPTION_BUILD_PORTS_PY=On ..")).
		Run(llb.Shlexf("make -j %v ", runtime.NumCPU())).Root()
	return metacallPyRuntime
}

//tests have been excluded for now
//Run(llb.Shlexf("ctest -j %v --output-on-failure --test-output-size-failed 3221000000 -C Release", runtime.NumCPU())).

//installs the runtime images
func buildDevinstall(devImage llb.State) llb.State {
	devInstall := devImage.Run(llb.Shlex("make install")).Root()
	return devInstall
}

//installs python runtime dependencies
func buildPyRuntime(baseupdateImg llb.State) llb.State {
	pyruntime := baseupdateImg.Run(llb.Shlex("apt-get -y install --no-install-recommends libpython3.9")).Root()
	return pyruntime
}

//to get the apt-get update layer segregated from the other
func buildBaseUpdate(baseImg llb.State) llb.State {
	return baseImg.Run(llb.Shlex("apt-get update")).Root()
}

func main() {
	base := llb.Image("docker.io/library/debian:bullseye-slim")

	//devImage

	depsImage := builDeps(base)
	metacall_base_runtime := buildmetacallbaseRuntime(depsImage)
	metacall_base_runtime_install := buildDevinstall(metacall_base_runtime)

	metacall_python_runtime := buildMetacallPythonRuntime(metacall_base_runtime)
	metacall_python_install := buildDevinstall(metacall_python_runtime)

	baseUpdate := buildBaseUpdate(base)

	//contains addtional files that come with apt-get update
	pyRuntimeDepsbuilt := buildPyRuntime(baseUpdate)

	//contains only python runtime dependencies nothing else
	pyRuntimeDeps := llb.Diff(baseUpdate, pyRuntimeDepsbuilt)

	/*TODO

	Use this to copy only .so files and runtimeDeps to an Image leaving the core repo


	scratch := llb.Image("scratch")

	c := b.Copy(a.WithState(llb.Scratch().Dir("/ced")), "./foo", "./baz") // /abc/baz
	*/

	//only contains *.so files for python and related python files
	metacall_py_runtime_final := llb.Diff(metacall_base_runtime_install, metacall_python_install)

	finImage := llb.Merge([]llb.State{metacall_py_runtime_final, pyRuntimeDeps})

	resImg, err := finImage.Marshal(context.TODO(), llb.LinuxAmd64)

	if err != nil {
		log.Fatal(err)
	}
	llb.WriteTo(resImg, os.Stdout)
}
