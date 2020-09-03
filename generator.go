package goinsta

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"strconv"
	"time"
	mrand "math/rand"
)

const (
	volatileSeed = "12345"
)

func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateHMAC(text, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateDeviceID(seed string) string {
	hash := generateMD5Hash(seed + volatileSeed)
	return "android-" + hash[:16]
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func generateUUID() string {
	uuid, err := newUUID()
	if err != nil {
		return "cb479ee7-a50d-49e7-8b7b-60cc1a105e22" // default value when error occurred
	}
	return uuid
}

func generateSignature(data string) map[string]string {
	m := make(map[string]string)
	m["ig_sig_key_version"] = goInstaSigKeyVersion
	m["signed_body"] = fmt.Sprintf(
		"%s.%s", generateHMAC(data, goInstaIGSigKey), data,
	)
	return m
}

func generateUserAgent() string {
	mrand.Seed(time.Now().Unix())
	return fmt.Sprintf(goInstaUserAgent,goInstaAppVersion,
		devices[mrand.Intn(len(devices))],goInstaLanguage,goInstaAppVersionCode)
}

func GenerateEncPassword(password, publicKeyId, publicKey string) (enc string, err error) {
	sessionKey := make([]byte, 32)
	_, err = rand.Read(sessionKey)
	if len(sessionKey) != 32 || err != nil {
		return
	}
	fmt.Printf("randKey: %x\n", sessionKey)

	iv := make([]byte, 12)
	_, err = rand.Read(iv)
	if len(iv) != 12 || err != nil {
		return
	}

	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return
	}

	pubPem, _ := pem.Decode(decodedPublicKey)
	if pubPem == nil {
		err = fmt.Errorf("pem.Decode has empty result")
		return
	}
	if pubPem.Type != "PUBLIC KEY" {
		err = fmt.Errorf("key is not public key but: %s",pubPem.Type)
		return
	}
	parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		return
	}
	pubKey, ok := parsedKey.(*rsa.PublicKey);
	if !ok {
		err = fmt.Errorf("failed to parse public key")
		return
	}

	encSessionKey, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, sessionKey)
	if err != nil {
		return
	}

	fmt.Printf("rsaEncrypted: %x\n",encSessionKey)

	timestring := fmt.Sprintf("%d",time.Now().Unix())
	//timestring := "1596435414"
	fmt.Printf("time: %s\n",timestring)


	cipherAes, err := aes.NewCipher(sessionKey)
	if err != nil {
		return
	}

	cipherAesGcm, err := cipher.NewGCM(cipherAes)
	if err != nil {
		return
	}
	cipherWithTag := cipherAesGcm.Seal(iv,iv,[]byte(password),[]byte(timestring))
	cipherText := cipherWithTag[:len(cipherWithTag)+-cipherAesGcm.Overhead()]
	cipherText = cipherText[12:]
	cipherTag := cipherWithTag[len(cipherWithTag)-cipherAesGcm.Overhead():]

	fmt.Printf("aesEncrypted: %x\n",cipherText)
	fmt.Printf("aesEncrypted.byteLength: %d\n",len(cipherText))

	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, uint16(len(encSessionKey)))

	ii, _ := strconv.Atoi(publicKeyId)
	kId := make([]byte, 2)
	binary.BigEndian.PutUint16(kId, uint16(ii))
	kId[0] = 0x01

	fmt.Printf("rsaEncrypted.byteLength: %d\n",len(encSessionKey))
	fmt.Printf("sizeBuffer: %x\n",bs)
	fmt.Printf("authTag: %x\n",cipherTag)

	fmt.Printf("kId: %x\n",kId)

	var data []byte
	data = append(data,kId...)
	data = append(data,iv...)
	data = append(data,bs...)
	data = append(data,encSessionKey...)
	data = append(data,cipherTag...)
	data = append(data,cipherText...)

	fmt.Printf("data: %x\n",data)
	payload := base64.StdEncoding.EncodeToString(data)
	fmt.Printf("datalen: %d\n",len(data))

	// return f"#PWD_INSTAGRAM:4:{time}:{payload.decode()}"
	enc = fmt.Sprintf("#PWD_INSTAGRAM:4:%s:%s",timestring,payload)

	return
}

func GenerateJazoest(input string) string {
	buf := []byte(input)
	sum := 0
	for _, x := range buf {
		sum += int(x)
	}
	return fmt.Sprintf("2%d",sum)
}