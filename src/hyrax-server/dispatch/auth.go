package dispatch

import (
    "crypto/sha1"
    "encoding/hex"
    "hyrax/types"
    "bytes"
    "sort"
    "strings"
)

var secretkeys [][]byte

// SetSecretkeys sets the list of secret keys to a different list
func SetSecretKeys(keys [][]byte) {
    secretkeys = keys
}

// GetSecretKeys returns the list of keys currently in use
func GetSecretKeys() [][]byte {
    return secretkeys
}

// Given a command payload (which presumably has a secret set), checks
// whether that secret checks out for one of the secret keys
func CheckAuth(cmdP *types.Payload) bool {
    h := sha1.New()
    for i := range secretkeys {
        h.Write( authMsg(cmdP.Domain,cmdP.Name,secretkeys[i]) )
        sum := h.Sum(nil)
        sumencodedsize := hex.EncodedLen(len(sum))
        sumencoded := make([]byte,sumencodedsize)
        hex.Encode(sumencoded,sum)
        if bytes.Equal(sumencoded,cmdP.Secret) { return true }
        h.Reset()
    }
    return false
}

// Similar to CheckAuth, will check whether a secret works with a secret
// key but this sorts the payload state keys to sign that data
func CheckAuthState(cmdP *types.Payload) bool {
    h := sha1.New()
    for i := range secretkeys {
        h.Write(authState(cmdP.State, secretkeys[i]))
        sum := h.Sum(nil)
        sumencodedsize := hex.EncodedLen(len(sum))
        sumencoded := make([]byte, sumencodedsize)
        hex.Encode(sumencoded, sum)
        if bytes.Equal(sumencoded, cmdP.Secret) {
            return true
        }
        h.Reset()
    }
    return false
}

func authMsg(domain, name, secret []byte) []byte {
    dl := len(domain)
    dnl := dl + len(name)
    buf := make([]byte,dnl+len(secret))
    copy(buf,domain)
    copy(buf[dl:],name)
    copy(buf[dnl:],secret)
    return buf
}

func authState(state map[string] string, secret []byte) []byte {
    keys := append([]string{}, string(secret))
    vals := append([]string{}, string(secret))
    for key, val := range state {
        keys = append(keys, key)
        vals = append(keys, val)
    }
    sort.Strings(keys)
    sort.Strings(vals)
    result := append([]string{}, keys...)
    result = append(result, "|")
    result = append(result, vals...)
    return []byte(strings.Join(result, ""))
}
