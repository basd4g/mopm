// vim:set noexpandtab :
package main

import (
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type Package struct {
	Name         string
	Url          string
	Description  string
	Environments []Environment
}

type PackageFile struct {
	Package *Package
	Path    string
}

func (pkg Package) String() string {
	out := new(bytes.Buffer)
	fmt.Fprintf(out, "name:         %s\n", pkg.Name)
	fmt.Fprintf(out, "url:          %s\n", pkg.Url)
	fmt.Fprintf(out, "description:  %s\n", pkg.Description)
	fmt.Fprintf(out, "environments: ")

	for i, env := range pkg.Environments {
		if i != 0 {
			fmt.Fprint(out, ", ")
		}
		fmt.Fprint(out, env)
	}
	return string(out.Bytes())
}

func (pkgFile PackageFile) String() string {
	return fmt.Sprintf("path:         %s\n%s", pkgFile.Path, pkgFile.Package)
}

func findAllPackageFile(packageName string) ([]PackageFile, error) {
	var pkgFiles []PackageFile
	for _, repo := range repositories() {
		path := repo.dir + "/definitions/" + packageName + ".yaml"
		pkgFile, err := readPackageFile(path)
		if err == nil {
			pkgFiles = append(pkgFiles, pkgFile)
		}
	}
	return pkgFiles, nil
}

func readPackageFile(path string) (PackageFile, error) {
	_, err := os.Stat(path)
	if err != nil {
		return PackageFile{}, fmt.Errorf("The package does not exist: %s\nWrapped: %w", path, err)
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return PackageFile{}, err
	}

	pkg := Package{}
	err = yaml.Unmarshal(buf, &pkg)
	if err != nil {
		return PackageFile{}, fmt.Errorf("Failed to parse yaml file: %s\nWrapped: %w", path, err)
	}
	err = lintPackage(&pkg)
	if err != nil {
		return PackageFile{}, err
	}
	return PackageFile{
		Package: &pkg,
		Path:    path,
	}, nil
}
