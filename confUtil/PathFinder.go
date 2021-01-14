package confUtil

import (
	"os"
	"path/filepath"

	//"yproject/github.com/yiGmMk/pz-infra-new/log"
	gopath "path"
	"runtime"

	"github.com/yiGmMk/pz-infra-new/log"
)

func FileHierarchyFind(path string) string {
	return hierarchyFind(path, true)
}

func DirectoryHierarchyFind(path string) string {
	return hierarchyFind(path, false)
}

func hierarchyFind(path string, isFilePath bool) string {
	if path == "" {
		//log.Error("empty path")
		return path
	}
	rootDir, err := filepath.Abs("/")
	if err != nil {
		//log.Error(err.Error())
		return ""
	}
	currentDir, err := filepath.Abs(".")
	log.Debug("currentDir is ", currentDir)
	for {
		testPath := filepath.Join(currentDir, path)
		if stat, err := os.Stat(testPath); err == nil {
			if stat.IsDir() == !isFilePath {
				log.Debug("find ", path, " at ", testPath)
				return testPath
			} else {
				log.Debug("find ", path, " at ", testPath, ", but not match isFilePath:", isFilePath)
			}
		} else {
			log.Debug("not find ", testPath)
		}
		if currentDir == rootDir {
			break
		}
		currentDir, err = filepath.Abs(filepath.Join(currentDir, "../"))
		if err != nil {
			log.Error(err.Error())
			return ""
		}
	}

	_, filename, _, _ := runtime.Caller(1)
	currentCallerDir := gopath.Dir(filename)
	log.Debug("currentDir Caller gopath is ", currentDir)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	for {
		testPath := filepath.Join(currentCallerDir, path)
		if stat, err := os.Stat(testPath); err == nil {
			if stat.IsDir() == !isFilePath {
				log.Debug("find ", path, " at ", testPath)
				return testPath
			} else {
				log.Debug("find ", path, " at ", testPath, ", but not match isFilePath:", isFilePath)
			}
		} else {
			log.Debug("not find ", testPath)
		}
		if currentCallerDir == rootDir {
			break
		}
		currentCallerDir, err = filepath.Abs(filepath.Join(currentCallerDir, "../"))
		if err != nil {
			log.Error(err.Error())
			return ""
		}
	}

	log.Debug("not found ", path, " in current or parent directories")
	return ""
}
