/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chartutil

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Expand uncompresses and extracts a chart into the specified directory.
func Expand(dir string, r io.Reader) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fmt.Println(header.Name)
		path := filepath.Clean(filepath.Join(dir, header.Name))
		fmt.Println(path)
		info := header.FileInfo()
		fmt.Printf("%#v", &info)
		fmt.Printf("%v\n", info.Name())
		fmt.Printf("%v\n", info.IsDir())
		//TODO: put in some logic. if file is nested in a dir that doesn't exist, then create it
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			fmt.Println("something with op.open")
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tr)
		if err != nil {
			fmt.Println("copying")
			return err
		}
	}
	return nil
}
