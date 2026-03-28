package user

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"exchange-go/internal/platform"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	UbCaptchaPrefix         = "ub-captcha_"
	UbCaptchaKeyPath        = "./config/ub-captcha/"
	TestUbCaptchaKeyPath    = "../../config/ub-captcha/"
	UbCaptchaPrivateKeyFile = "private.pem"
	UbCaptchaPublicKeyFile  = "public.pem"
	BeginMessage            = "-----BEGIN MESSAGE-----"
	EndMessage              = "-----END MESSAGE-----"
)

// UbCaptchaManager provides a custom RSA-based CAPTCHA mechanism for anti-bot protection.
// It handles key generation, encryption/decryption, and CAPTCHA validation.
type UbCaptchaManager interface {
	// CheckUbCaptcha validates a UbCaptcha response string by decrypting and verifying
	// the embedded timestamp and interaction time.
	CheckUbCaptcha(ubCaptchaStr string) (bool, error)
	// Encrypt encrypts the given plaintext using the RSA public key stored on disk.
	Encrypt(plainText string) (string, error)
	// Decrypt decrypts the given PEM-encoded ciphertext using the RSA private key stored on disk.
	Decrypt(encryptedText string) (string, error)
	// NewKey generates a new 2048-bit RSA key pair.
	NewKey() (UbCaptchaKey, error)
	// GetKey loads the existing RSA key pair from PEM files on disk.
	GetKey() (UbCaptchaKey, error)
	// SaveKeyToPemFile persists the RSA key pair to private.pem and public.pem files.
	SaveKeyToPemFile(key UbCaptchaKey) error
}

type UbCaptchaKey struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type DecryptedMsg struct {
	Timestamp int64 `json:"timestamp"`
	SpentTime int64 `json:"spent_time"`
}

type ubCaptchaManager struct {
	logger platform.Logger
}

func (ubc *ubCaptchaManager) CheckUbCaptcha(ubCaptchaStr string) (bool, error) {
	//prepare encrypted message
	encryptedMsg, err := ubc.prepareEncryptedMessage(ubCaptchaStr)
	if err != nil {
		ubc.logger.Error2("can not prepare encrypyed message", err,
			zap.String("service", "ubCaptchaManager"),
			zap.String("method", "CheckUbCaptcha"),
			zap.String("ubCaptchaStr", ubCaptchaStr),
		)
		return false, err
	}

	//decrypt message
	decryptedMsg, err := ubc.Decrypt(encryptedMsg)
	if err != nil {
		ubc.logger.Error2("can not decrypt encrypyed message", err,
			zap.String("service", "ubCaptchaManager"),
			zap.String("method", "CheckUbCaptcha"),
			zap.String("ubCaptchaStr", ubCaptchaStr),
		)
		return false, err
	}

	decryptedMsgObj := &DecryptedMsg{}
	err = json.Unmarshal([]byte(decryptedMsg), decryptedMsgObj)
	if err != nil {
		ubc.logger.Error2("can not unmarshal decrypted message", err,
			zap.String("service", "ubCaptchaManager"),
			zap.String("method", "CheckUbCaptcha"),
			zap.String("ubCaptchaStr", ubCaptchaStr),
			zap.String("decryptedMsg", decryptedMsg),
		)
		return false, err
	}

	isUbCaptchaValid := ubc.isValid(decryptedMsgObj)

	return isUbCaptchaValid, nil
}

func (ubc *ubCaptchaManager) isValid(ubCaptchaObj *DecryptedMsg) bool {

	if ubCaptchaObj.SpentTime < 1 {
		return false
	}

	currentTimeStamp := time.Now().Unix()

	if currentTimeStamp-ubCaptchaObj.Timestamp > 20 && flag.Lookup("test.v") == nil {
		return false
	}

	return true
}

func (ubc *ubCaptchaManager) NewKey() (UbCaptchaKey, error) {
	var key UbCaptchaKey

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return key, err
	}

	key.PublicKey = &privateKey.PublicKey
	key.PrivateKey = privateKey

	return key, nil
}

func (ubc *ubCaptchaManager) GetKey() (UbCaptchaKey, error) {
	var key UbCaptchaKey
	keyPath := UbCaptchaKeyPath

	//in test env
	if flag.Lookup("test.v") != nil {
		keyPath = TestUbCaptchaKeyPath
	}

	//read private key from private.pem file
	privateKeyFile, err := os.ReadFile(keyPath + UbCaptchaPrivateKeyFile)
	if err != nil {
		return key, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	privateKeyBlockBytes := privateKeyBlock.Bytes

	/*ok := x509.IsEncryptedPEMBlock(privateKeyBlock)
	if ok {
		privateKeyBlockBytes, err = x509.DecryptPEMBlock(privateKeyBlock, nil)
		if err != nil {
			return key, err
		}
	}*/

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlockBytes)
	if err != nil {
		return key, err
	}

	key.PrivateKey = privateKey

	//read public key from public.pem file
	publicKeyFile, err := os.ReadFile(keyPath + UbCaptchaPublicKeyFile)
	if err != nil {
		return key, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyFile)
	publicKeyBlockBytes := publicKeyBlock.Bytes

	/*ok = x509.IsEncryptedPEMBlock(publicKeyBlock)
	if ok {
		publicKeyBlockBytes, err = x509.DecryptPEMBlock(publicKeyBlock, nil)
		if err != nil {
			return key, err
		}
	}*/

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlockBytes)
	if err != nil {
		return key, err
	}

	key.PublicKey = publicKey

	return key, nil
}

func (ubc *ubCaptchaManager) Encrypt(plainText string) (string, error) {

	key, err := ubc.GetKey()
	if err != nil {
		return "", err
	}

	/*cipher, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, key.PublicKey, []byte(plainText), nil)
	if err != nil {
		return "", err
	}*/

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, key.PublicKey, []byte(plainText))
	if err != nil {
		return "", err
	}

	return ubc.cipherToPemString(cipher), nil

}

func (ubc *ubCaptchaManager) Decrypt(encryptedText string) (string, error) {

	key, err := ubc.GetKey()
	if err != nil {
		return "", err
	}

	/*plainText, err := rsa.DecryptOAEP(
		sha512.New(),
		rand.Reader,
		key.PrivateKey,
		ubc.pemStringToCipher(encryptedText),
		nil,
	)*/

	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, key.PrivateKey, ubc.pemStringToCipher(encryptedText))

	return string(plainText), err

}

func (ubc *ubCaptchaManager) SaveKeyToPemFile(key UbCaptchaKey) error {
	//in test env
	if flag.Lookup("test.v") != nil {
		return nil
	}

	//Save private key to pem file
	err := ubc.savePrivateKeyToPemFile(key)
	if err != nil {
		return err
	}

	//Save public key to pem file
	err = ubc.savePublicKeyToPemFile(key)
	if err != nil {
		return err
	}

	return nil
}

func (ubc *ubCaptchaManager) savePublicKeyToPemFile(key UbCaptchaKey) error {
	fileContent := ubc.publicKeyToPemString(key)
	filePath := UbCaptchaKeyPath + UbCaptchaPublicKeyFile

	return os.WriteFile(filePath, []byte(fileContent), 0500)
}

func (ubc *ubCaptchaManager) savePrivateKeyToPemFile(key UbCaptchaKey) error {
	fileContent := ubc.privateKeyToPemString(key)
	filePath := UbCaptchaKeyPath + UbCaptchaPrivateKeyFile

	return os.WriteFile(filePath, []byte(fileContent), 0500)
}

func (ubc *ubCaptchaManager) publicKeyToPemString(key UbCaptchaKey) string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(key.PublicKey),
		}))
}

func (ubc *ubCaptchaManager) privateKeyToPemString(key UbCaptchaKey) string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key.PrivateKey),
		}))
}

func (ubc *ubCaptchaManager) pemStringToCipher(encryptedMessage string) []byte {
	b, _ := pem.Decode([]byte(encryptedMessage))
	return b.Bytes
}

func (ubc *ubCaptchaManager) cipherToPemString(cipher []byte) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "MESSAGE",
				Bytes: cipher,
			},
		),
	)
}

func (ubc *ubCaptchaManager) prepareEncryptedMessage(encryptedMsg string) (string, error) {
	//remove 'ub-captcha_' from the first
	encryptedMsg = strings.Replace(encryptedMsg, UbCaptchaPrefix, "", 1)

	//validate length
	if len(encryptedMsg) != 344 {
		return "", fmt.Errorf("encrypted message length is not valid")
	}

	sliceSize := 64
	sliceCount := 6
	var slices []string

	for i := 1; i <= sliceCount; i++ {
		if i < sliceCount {
			slice := encryptedMsg[(i-1)*sliceSize : i*sliceSize]
			slices = append(slices, slice)
		} else {
			slice := encryptedMsg[(i-1)*sliceSize:]
			slices = append(slices, slice)
		}
	}

	newEncryptedMsg := ""
	for _, slice := range slices {
		newEncryptedMsg = newEncryptedMsg + slice + "\n"
	}

	newEncryptedMsg = BeginMessage + "\n" + newEncryptedMsg + EndMessage

	return newEncryptedMsg, nil

}

func NewUbCaptchaManager(logger platform.Logger) UbCaptchaManager {
	return &ubCaptchaManager{logger}
}
