// Copyright © 2013 Hraban Luyat <hraban@0brg.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// Create a new temporary directory with a hard-to-guess name. Go
// implementation of mkdtemp (3) from the C stdlib.h.
package mkdtemp

import (
	crand "crypto/rand"
	"math"
	mrand "math/rand"
	"os"
	"strings"
	"time"
)

const randchars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

func randtempidxMath() int {
	mrand.Seed(time.Now().UnixNano())
	return mrand.Intn(len(randchars))
}

func randtempidxCrypt() (int, error) {
	buf := make([]byte, 1)
	_, err := crand.Read(buf)
	if err != nil {
		return 0, err
	}
	b := int(buf[0])
	// Ratio between random integer domain and range
	factor := (1 + math.MaxUint8) / len(randchars)
	// Maximum domain that can be uniformly mapped onto target range
	limit := factor * len(randchars)
	// Discard anything outside that domain
	if b > limit {
		return randtempidxCrypt()
	}
	return b / factor, nil
}

func randtempc() string {
	i, err := randtempidxCrypt()
	if err != nil {
		i = randtempidxMath()
	}
	return string(randchars[i])
}

func filltemplate(template string) string {
	if !strings.Contains(template, "X") {
		return template
	}
	return filltemplate(strings.Replace(template, "X", randtempc(), 1))
}

// Create a new directory in the system's temporary folder. The argument is a
// directory name template, e.g.: "mynameXXXXXX". All `X' characters will be
// replaced by random characters.
// 
// The directory is directly created in the system's temporary folder. Anyone
// with read access to that temporary folder can see your temporary directory
// name. The directory is not removed, cleared, or in any way altered after
// usage: the only thing this function does is create a new directory with mode
// 0700 and return its name.
func Mkdtemp(template string) (string, error) {
	name := os.TempDir() + "/" + filltemplate(template)
	err := os.Mkdir(name, 0700)
	if err == os.ErrExist {
		return Mkdtemp(template)
	}
	return name, err
}
