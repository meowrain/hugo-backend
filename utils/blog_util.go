package utils

import (
	"hugo_backend/models"
	"os"
	"path/filepath"
	"strings"
)

func BuildDirectoryTree(root, dir string) (models.Directory, error) {
	var tree models.Directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return tree, err
	}

	tree.Name = filepath.Base(dir)
	tree.Path = strings.TrimPrefix(dir, root+"/")
	for _, file := range files {
		if file.IsDir() {
			childTree, err := BuildDirectoryTree(root, filepath.Join(dir, file.Name()))
			if err != nil {
				return tree, err
			}
			tree.Children = append(tree.Children, childTree)
		} else if strings.HasSuffix(file.Name(), ".md") {
			post := models.Post{
				Name: file.Name(),
				Path: strings.TrimPrefix(filepath.Join(dir, file.Name()), root+"/"),
			}
			tree.Posts = append(tree.Posts, post)
		}
	}
	return tree, nil
}
