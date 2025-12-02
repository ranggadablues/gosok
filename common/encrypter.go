package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

const key = "iWYiY{{G:w(rU!hHFfMnPrC9Wfam@e}GX/;hhz;!;3W=;3&3K.!3Hg$4E9WGdbJ.uNm&wH-]bh4bGwVgbDfv]djqBN&%y-1xYVn$H.!wUu51fMkLt@.BB&gu/RGJ0q+#1VU!!}K:ND:12)Q-EaYjkfn=#D}Mueqqn9kEim0!0+,9wz0xCMa?;t,/JLJn&[Sfv]3ERV:x}5/DqShWnjj27v1YBLx8yKE{a)jBzGzxJS;}k[!0$mt!:HA$gG/fmzY(mcW5W*;&8163L{8U1,2GBJ*GbmRgVU(EeSYhS!$*jn%=%ht@]Q1=Y!L(*SK90Xn&JBGZ(AJP2eVjPg82Ayg?A(Y(&KNy.VX2R{_gyZmp_b%G2+FX)wW@E_65VffjN6;]42U4ppvAqub2ZEX8Cw,mezHMaqBuv6wPG7eRV+Wq3QB6LBA.C(eeCU)Xw4gdma[GH5BwP3XfCb5G7=&ViT&iUkcZ44D8a06d4BF(,QHFjVD$hkW0VHdJ7(n#1f2:N!Axbq81%uu/+@(ZP&31C(HQE_-c6=kLKxnTWK+TapGH2,fV%73G$]iXXP4ZZDYfny]@{ZJgJ/E*98Za8[w_q/}U)?Yhea&aWG{q(6b}n}MCi$=G#/zr?!:hju_0PV!q.te+R9uinq_U-QZywyz%3=ZA]x!!*8@QwtM&p*h[8qptZ/QZ@uiuFg,3Jzi4*%?4FX&S70UYadbq03Jq%Ey//jU-f@mMt!#Nd[kt%BnPW=?_&wU{k8$!4j+kM)jMG,[3zE#M,9@PdUF3)h6PW-zMtkq2+AvFU}Zd_2:v*Gxi,bN@a=+1q(f2Vww}UxaitRwj+cBA457B90yP=$5nay2fK[=[e$!C6T=QBji$W2B[Q4p{J@0S2.Hg+(&=L8E6c9nh_7gQ/(@]ZZt*K#gDYyUyEy9u+p+yJ_hh-/@DA+VD$W!tYr{Q9N0U!.?vDFG4d6}YfGQrYi_@a,:&kGE}?,X1DBYL9(Y-?uxQJaE+eY};k6FV"

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encryptedText string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(encryptedText)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
