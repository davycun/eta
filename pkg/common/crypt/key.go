package crypt

const (
	rsaPublickKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqUk8l/CsslrfWm0hDNg+
DSPHOT+qNqTnJPec5f5noO/ET1J2msst/vRx0C1uxTOpmZmZOkIitcZWXVyoLK4j
upLk2olrgPvJNf1MJEn0AVgLewnLoIbH4hV9k211DRfN/3VXXZ0ig7udHjJqGFhT
mPDPE9/kqviXdO6XoUyYPQNPryQvtsaAvbl2Q0mYe3xPjs9MI005mQmnYnkQFmzK
4N/mJepS9TfQZuddu2707BVSsA2JTMN6l8mVov92b534XUeKW9RkQfeyzeI6R7QH
1Ba5pRYtcF67w8A8fz+71NCFpdctfHpyIU6BGiV125oOn78w0rVZQZroJPJkO+b/
1wIDAQAB
-----END PUBLIC KEY-----`
	rsaPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCpSTyX8KyyWt9a
bSEM2D4NI8c5P6o2pOck95zl/meg78RPUnaayy3+9HHQLW7FM6mZmZk6QiK1xlZd
XKgsriO6kuTaiWuA+8k1/UwkSfQBWAt7CcughsfiFX2TbXUNF83/dVddnSKDu50e
MmoYWFOY8M8T3+Sq+Jd07pehTJg9A0+vJC+2xoC9uXZDSZh7fE+Oz0wjTTmZCadi
eRAWbMrg3+Yl6lL1N9Bm5127bvTsFVKwDYlMw3qXyZWi/3ZvnfhdR4pb1GRB97LN
4jpHtAfUFrmlFi1wXrvDwDx/P7vU0IWl1y18enIhToEaJXXbmg6fvzDStVlBmugk
8mQ75v/XAgMBAAECggEALm2NlsZFNu1BUJWZeOJdslDbtNHHJxF262aVu2ZYmYTo
vDCLosySotf71vJ+7MrMevnrUlUNG/l3ekeNQCPKXMMozN7fgxKLDqmXlmRJ7Yxu
KaJ4HGCatWDbffGJJrwenS9bdKtB6gssfjmpa1/eHZX91R+UoWdocoN1RqGxJxos
1G3HQm4bH54bE+e8QYjPNtY9kSMNNVsWYvr+C4Np/V+D/pEuTFwYn9yL6N88sKdJ
v14t0wZsj1CSKabOBsuLScE6DEGjOiSMPPsqhFAOgOGOAYvKsh7Fx1ZER1H+WymH
VlF4mq8ihPJw1VI3HUfr3rF1Ug3IfBlMVZTBQ+spMQKBgQDhhKj9GF2YYTGxRQBP
Gf9tVSvr12Losa5vFdGtEMJAe0AaNHcnbreH8wkYpTSTYOMRKdpAyxAThoHMfBnT
nYPKZFwJMf7I2+zF44TuX3BbSHTxZMBDg4c6vst9AffCFlLyC9drERW6cXJUhucH
SRQIDAaKL6+3FBsR4X7WIVDo2wKBgQDAKtfHxDDTbbCNDCWtZ3RiYwkXqsq/w4u0
wIN6+KftUPNh49++qPaKii8sFbLwM31PTUmDdGT82TRk/lHbe/fL2XbqdBnPExiJ
QEkUMmd4ATDwLCRnAbNqB6laJL568DMUYxiHXd3SNLQDN60e58Qi6VR+0iZiXycg
bWxzO3sntQKBgQDMa8H62fFNRR0UQSXMjhZY9tF/UJgZsaYaj76mkABlDtPGbTRt
DBFVLFpcerQzu6lYT4XFIcyKxmw1/XAzwB5MgbhjpWv8ik4P+vLuWJiyRfWrMtaM
3FbiSzyNLhF1X2MEgPNd2/jELj2inT4h/n4n5S6waBVxcMow5SRh2YVL7QJ/GkXB
Oo6YHk/g02fVvt9mCg2AMLwo/A7ACvZHA4j7bHakz71X2bk/+7Dowh41WKGxgCYQ
5ugC5o7LmwMqLdfFCRmyKFu3K5hwwlMHqrs5ai6/ghaG445X+ScTAEAFyskMMr1Z
nSg4K4UJF2eFQ7RsHFnKM4yHsutPbh+HtDltnQKBgQCu1ResuGTN3tjVMNSEFdjP
KbZ1ROy4Xyb/TtFM2ahyNiGY02UsSEc5QoNgg4gdb9dFtrQs0fxa6blT9ieurT7T
pHf/eNuxIysqPzOWEY+UaeUmnWFKzmXasZNSRszVHWmUpJd8+VBuGeW1qc6Gqmdi
3etuyDHd0zWi0QVau7fUJA==
-----END PRIVATE KEY-----`
	//sm2PublicKey  = `BMB8vZ4fS6ZbyMKtdHoqflMdabaLk192pihn9gMN4DMWfkL7lbYPY+CW/elQMlse5hwBnWqQWvMdWmPSs8NxvMI=`
	//sm2PrivateKey = `TVMRS3HT7zId/M1QOFRdwqC+sl2nevZEs/acDHuxvnM=`
	aesKey = "f327f6e6c482e21e4c0bdd0764b41d1e"
	sm4Key = "8eadb267ecd6e860"

	sm2PrivateKey = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgB/IIvg4ijmPUOR8p
wKVpVTRXRqXFPKrJj2hg8fqBWTSgCgYIKoEcz1UBgi2hRANCAAQQOlEWQAsN1DDi
zfIKg9xCgWb7BGHJHc+YN+rkhNAEMV/7Qfg8AwIAqEI++tkb8bgVX7knmMPjBbL8
N57dip/u
-----END PRIVATE KEY-----`
	sm2PublicKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEEDpRFkALDdQw4s3yCoPcQoFm+wRh
yR3PmDfq5ITQBDFf+0H4PAMCAKhCPvrZG/G4FV+5J5jD4wWy/Dee3Yqf7g==
-----END PUBLIC KEY-----`
)

var (
	publicKey = map[string]string{
		AlgoAsymRsaPKCS1v15:  rsaPublickKey,
		AlgoAsymSm2Pkcs8C132: sm2PublicKey,
	}
	privateKey = map[string]string{
		AlgoAsymRsaPKCS1v15:  rsaPrivateKey,
		AlgoAsymSm2Pkcs8C132: sm2PrivateKey,
	}
	keyMap = map[string]string{
		AlgoSymAesEcbPkcs7padding: aesKey,
		AlgoSymSm4EcbPkcs7padding: sm4Key,
	}
)

func GetPublicKey(algo string) string {
	return publicKey[algo]
}
func GetPrivateKey(algo string) string {
	return privateKey[algo]
}
