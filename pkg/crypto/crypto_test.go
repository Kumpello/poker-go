package crypto

import (
	"errors"
	"testing"
)

func Test_HashAndVerify(t *testing.T) {
	type tc struct {
		name string
		pass string
		hash string // hash is one of hashes that matches the password
		err  error
	}

	tcs := []tc{
		{
			name: "empty password",
			pass: "",
			hash: "$2a$14$Pxc9Eyl3bKxyMAvvetH/iujpX3gzCrSUyr1ux7u6yRZiRQsQmrgxO",
			err:  nil,
		},
		{
			name: "pass_1",
			pass: "pass_1",
			hash: "$2a$14$Pxc9Eyl3bKxyMAvvetH/iujpX3gzCrSUyr1ux7u6yRZiRQsQmrgxO",
			err:  nil,
		},
	}

	for _, test := range tcs {
		t.Run(test.name, func(t *testing.T) {
			_, err := HashPassword(test.pass)

			if err != nil {
				if errors.Is(err, test.err) {
					return // expected err
				}
				t.Fatalf("unexpected err, is: %v, should be: %v", err, test.err)
				return
			}

			err = VerifyPassword(test.hash, test.pass)
			if err != nil {
				if !errors.Is(err, test.err) {
					return
				}
				t.Fatalf("cannot verify the password, is: %s", test.err)
				return
			}
		})
	}
}
