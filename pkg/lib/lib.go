package lib

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// UserHome builds path with user home path
func UserHome(s ...string) string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	dir := user.HomeDir

	if len(s) == 0 {
		return dir
	}

	if []rune(s[0])[0] != '/' {
		return dir + "/" + s[0]
	}

	return dir + s[0]
}

// BookInfoBasic contains basic book information
type BookInfoBasic struct {
	Title  string
	Author string
	Volume int
	Images int
}

// StringSliceContain search the string slice and see if it contains the searched word, match from first character
func StringSliceContain(strSlice []string, strSearch string) bool {
	for _, str := range strSlice {
		if strings.Index(strSearch, str) == 0 {
			return true
		}
	}

	return false
}

// valid characters for the session id
const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateString create random new string
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func GenerateString(n int) string {
	// slightly less deterministic randomness
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

// SHA256Iter iterate password and salt x times, return hex result
func SHA256Iter(password, salt string, iter int) string {
	str := password + ":" + salt
	bstr := []byte(str)

	var bstrx [sha256.Size]byte
	for i := 0; i < iter; i++ {
		bstrx = sha256.Sum256(bstr)
		bstr = bstrx[0:32]
	}

	return fmt.Sprintf("%x", bstr)
}

// NewUUIDV4 generate uuid version 4
func NewUUIDV4() (string, error) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", errors.New("failed to generate uuid")
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}

var chunkifyRegexp = regexp.MustCompile(`(\d+|\D+)`)

func chunkify(s string) []string {
	return chunkifyRegexp.FindAllString(s, -1)
}

// AlphaNumCaseCompare returns true if the first string precedes the second one according to natural order
// https://github.com/facette/natsort/blob/master/natsort.go
func AlphaNumCaseCompare(a, b string) bool {
	chunksA := chunkify(a)
	chunksB := chunkify(b)

	nChunksA := len(chunksA)
	nChunksB := len(chunksB)

	for i := range chunksA {
		if i >= nChunksB {
			return false
		}

		aInt, aErr := strconv.Atoi(chunksA[i])
		bInt, bErr := strconv.Atoi(chunksB[i])

		// If both chunks are numeric, compare them as integers
		if aErr == nil && bErr == nil {
			if aInt == bInt {
				if i == nChunksA-1 {
					// We reached the last chunk of A, thus B is greater than A
					return true
				} else if i == nChunksB-1 {
					// We reached the last chunk of B, thus A is greater than B
					return false
				}

				continue
			}

			return aInt < bInt
		}

		// So far both strings are equal, continue to next chunk
		if chunksA[i] == chunksB[i] {
			if i == nChunksA-1 {
				// We reached the last chunk of A, thus B is greater than A
				return true
			} else if i == nChunksB-1 {
				// We reached the last chunk of B, thus A is greater than B
				return false
			}

			continue
		}

		return chunksA[i] < chunksB[i]
	}

	return false
}

// IsFileExists returns whether the given file or directory exists
func IsFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}