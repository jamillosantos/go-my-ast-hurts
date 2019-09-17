package myasthurts

import (
	"go/ast"
)

func Parse(file *ast.File, definitions *Environment) {

	definitions.Packages[0] = &Package{
		Name:    file.Name.Name,
		Comment: "",
	}

	/*definitions.Packages[0] = &Package{
		Files: []*File{
			{
				FileName: "teste",
				Package:
			},
		},
	}*/

	//definitions.Packages = make([]*Package, len(file.Imports))
	/*if file.Imports != nil {
		for idx, i := range file.Imports {

				definitions.Packages = append(definitions.Packages, &Package{
					Comment: i.Comment.Text(),
				})

			// 2
			definitions.Packages[idx] = &Package{Comment: i.Comment.Text()}
		}
	}*/
}
