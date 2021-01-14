package bsdiffUtil

import (
	"errors"
	"fmt"
	"os"

	"github.com/yiGmMk/pz-infra-new/fileUtil"
	"github.com/yiGmMk/pz-infra-new/log"

	"github.com/hsinhoyeh/binarydist"
)

//生成差分包patchFile
func Diff(oldFile, newFile, patchFile string) error {
	if fileUtil.ExistFile(patchFile) {
		os.Remove(patchFile)
	}

	if !fileUtil.ExistFile(oldFile) {
		err := errors.New(fmt.Sprintf("%s is not exist", oldFile))
		log.Error(err)
		return err
	}

	if !fileUtil.ExistFile(newFile) {
		err := errors.New(fmt.Sprintf("%s is not exist", newFile))
		log.Error(err)
		return err
	}

	old, _ := os.Open(oldFile)
	new, _ := os.Open(newFile)
	patch, _ := os.OpenFile(patchFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	err := binarydist.Diff(old, new, patch)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("bsdiff done")
	return nil
}

//根据老版本和差分包,来生成新版本newFile
func Patch(oldFile, newFile, patchFile string) error {
	if fileUtil.ExistFile(newFile) {
		os.Remove(newFile)
	}

	if !fileUtil.ExistFile(oldFile) {
		err := errors.New(fmt.Sprintf("%s is not exist", oldFile))
		log.Error(err)
		return err
	}

	if !fileUtil.ExistFile(patchFile) {
		err := errors.New(fmt.Sprintf("%s is not exist", patchFile))
		log.Error(err)
		return err
	}

	old, _ := os.Open(oldFile)
	new, _ := os.OpenFile(newFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	patch, _ := os.Open(patchFile)
	err := binarydist.Patch(old, new, patch)
	if err != nil {
		log.Error("failed. reason:", err.Error())
		return err
	}

	log.Debug("patch done")
	return nil
}
