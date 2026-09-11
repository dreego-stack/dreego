package dreegotest

import (
	"fmt"
	"path/filepath"
)

func testModuleFile(repoRoot string, withSSR bool) string {
	return testModuleFileNamed("t", repoRoot, withSSR)
}

func testModuleFileNamed(module, repoRoot string, withSSR bool) string {
	requires := "\tgithub.com/dreego-stack/dreego/core v0.0.0\n"
	replaces := fmt.Sprintf("replace github.com/dreego-stack/dreego/core => %s\n", filepath.Join(repoRoot, "core"))
	if withSSR {
		requires += "\tgithub.com/dreego-stack/dreego/adapter/ssr v0.0.0\n"
		replaces += fmt.Sprintf("replace github.com/dreego-stack/dreego/adapter/ssr => %s\n", filepath.Join(repoRoot, "adapter", "ssr"))
	}
	return "module " + module + "\n\ngo 1.27\n\nrequire (\n" + requires + "\tgolang.org/x/text v0.22.0 // indirect\n)\n\n" +
		fmt.Sprintf("replace github.com/dreego-stack/dreego => %s\n", repoRoot) + replaces
}
