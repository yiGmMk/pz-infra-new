package encryptUtil

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	Test_AES_KEY    = "5ca2aff64a3211e6a797567a637f5735"
	Test_AES_KEY_2  = "5ca2aff64a3211e6a797567a637f5736"
	Invalid_AES_KEY = "5ca2aff64a3211e6a797567a637f5735-INVALID"
)

func TestAesEncryption(t *testing.T) {
	Convey("Testing Encryption and Decryption", t, func() {
		Convey("Test with normal plain text", func() {
			plainText := "This is for testing string"
			encryptText, err := AesEncrypt(Test_AES_KEY, plainText)
			So(err, ShouldBeNil)
			So(encryptText, ShouldNotBeEmpty)
			So(encryptText, ShouldNotContainSubstring, plainText)
			fmt.Println("encryptText is: " + encryptText)
			//
			newPlainText, err := AesDecrypt(Test_AES_KEY, encryptText)
			So(err, ShouldBeNil)
			So(newPlainText, ShouldNotBeEmpty)
			So(newPlainText, ShouldEqual, plainText)
		})

		Convey("Test with Invalid AES Key", func() {
			plainText := "This is for testing string"
			encryptText, err := AesEncrypt(Invalid_AES_KEY, plainText)
			So(err, ShouldNotBeNil)
			So(encryptText, ShouldBeEmpty)
			So(err.Error(), ShouldContainSubstring, "invalid key size")
		})

		Convey("Test with Incorrect AES Key", func() {
			plainText := "This is for testing string"
			encryptText, err := AesEncrypt(Test_AES_KEY, plainText)
			So(err, ShouldBeNil)
			So(encryptText, ShouldNotBeEmpty)
			newPlainText, err := AesDecrypt(Test_AES_KEY_2, encryptText)
			So(err, ShouldBeNil)
			So(newPlainText, ShouldNotEqual, plainText)
		})

		Convey("Test with Incorrect AES Encryption Content", func() {
			plainText := "This is for testing string"
			encryptText, err := AesEncrypt(Test_AES_KEY, plainText)
			So(err, ShouldBeNil)
			So(encryptText, ShouldNotBeEmpty)
			// Change Encrypt Text
			newEncryptText := encryptText[:len(encryptText)-2] + "00"
			newPlainText, err := AesDecrypt(Test_AES_KEY, newEncryptText)
			So(err, ShouldBeNil)
			So(newPlainText, ShouldNotEqual, plainText)
			// Change Encrypt Text to short
			newEncryptText2 := encryptText[:4]
			newPlainText2, err := AesDecrypt(Test_AES_KEY, newEncryptText2)
			So(err, ShouldNotBeNil)
			So(newPlainText2, ShouldBeEmpty)
			So(err.Error(), ShouldContainSubstring, "encryptText too short")
		})
	})
}
