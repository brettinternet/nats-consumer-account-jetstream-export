//go:build ignore

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Basically the example from https://github.com/nats-io/jwt README
func run() error {
	okp, err := nkeys.CreateOperator()
	if err != nil {
		return err
	}
	opk, err := okp.PublicKey()
	if err != nil {
		return err
	}

	oc := jwt.NewOperatorClaims(opk)
	oc.Name = "O"
	oskp, err := nkeys.CreateOperator()
	if err != nil {
		return err
	}
	ospk, err := oskp.PublicKey()
	if err != nil {
		return err
	}
	oc.SigningKeys.Add(ospk)
	operatorJWT, err := oc.Encode(okp)
	if err != nil {
		return err
	}

	// ACCOUNT
	sysRes, sysCreds, err := createSystemAccount(oskp)
	a1Res, u1Creds, err := createUserAccount(oskp, func(ac *jwt.AccountClaims) {
		ac.Name = "A1"
		ac.DefaultPermissions.Pub.Allow.Add(">")
		ac.DefaultPermissions.Sub.Allow.Add(">")
		ac.Limits.JetStreamLimits.DiskStorage = -1
		ac.Limits.JetStreamLimits.MemoryStorage = -1
		ac.Limits.JetStreamLimits.Streams = -1

		ac.Exports.Add(
			&jwt.Export{
				Subject: "$JS.API.>",
				Type:    jwt.Service,
			},
			&jwt.Export{
				Subject: "$JS.ACK.>",
				Type:    jwt.Service,
			},
		)
	})
	if err != nil {
		return err
	}
	a2Res, u2Creds, err := createUserAccount(oskp, func(ac *jwt.AccountClaims) {
		ac.Name = "A2"
		ac.DefaultPermissions.Pub.Allow.Add(">")
		ac.DefaultPermissions.Sub.Allow.Add(">")
		ac.Limits.JetStreamLimits.DiskStorage = -1
		ac.Limits.JetStreamLimits.MemoryStorage = -1
		ac.Limits.JetStreamLimits.Streams = -1

		ac.Imports.Add(
			&jwt.Import{
				Type:    jwt.Service,
				Subject: "$JS.API.>",
				Account: a1Res.PK,
			},
			&jwt.Import{
				Type:    jwt.Service,
				Subject: "$JS.ACK.>",
				Account: a1Res.PK,
			},
		)
	})
	if err != nil {
		return err
	}

	dir := "./nats"
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	var perm fs.FileMode = 0777
	if err := os.Mkdir("./nats", perm); err != nil {
		return err
	}

	// Create account resolver in a nats server config
	resolver := fmt.Sprintf(`operator: %s
system_account: %s
jetstream: enabled
resolver: MEMORY
resolver_preload: {
  %s: %s
	%s: %s
	%s: %s
}
`, operatorJWT, sysRes.PK, sysRes.PK, sysRes.JWT, a1Res.PK, a1Res.JWT, a2Res.PK, a2Res.JWT)
	if err := os.WriteFile(path.Join(dir, "nats.conf"),
		[]byte(resolver), perm); err != nil {
		return err
	}
	fmt.Printf("Created nats.conf at %s\n", path.Join(dir, "nats.conf"))

	// store the user credentials
	sysCredsPath := path.Join(dir, "sys.creds")
	if err := os.WriteFile(sysCredsPath, sysCreds, perm); err != nil {
		return err
	}
	fmt.Printf("Created sys.creds at %s\n", sysCredsPath)
	u1CredsPath := path.Join(dir, "u1.creds")
	if err := os.WriteFile(u1CredsPath, u1Creds, perm); err != nil {
		return err
	}
	fmt.Printf("Created u1.creds at %s\n", u1CredsPath)
	u2CredsPath := path.Join(dir, "u2.creds")
	if err := os.WriteFile(u2CredsPath, u2Creds, perm); err != nil {
		return err
	}
	fmt.Printf("Created u2.creds at %s\n", u2CredsPath)
	return nil
}

type accountResolution struct {
	PK  string
	JWT string
}

func createUserAccount(oskp nkeys.KeyPair, claim func(ac *jwt.AccountClaims)) (*accountResolution, []byte, error) {
	// ACCOUNT
	akp, err := nkeys.CreateAccount()
	if err != nil {
		return nil, nil, err
	}
	apk, err := akp.PublicKey()
	if err != nil {
		return nil, nil, err
	}
	ac := jwt.NewAccountClaims(apk)
	askp, err := nkeys.CreateAccount()
	if err != nil {
		return nil, nil, err
	}
	aspk, err := askp.PublicKey()
	if err != nil {
		return nil, nil, err
	}
	ac.SigningKeys.Add(aspk)
	claim(ac)
	aJWT, err := ac.Encode(oskp)
	if err != nil {
		return nil, nil, err
	}

	// USER
	ukp, err := nkeys.CreateUser()
	if err != nil {
		return nil, nil, err
	}
	upk, err := ukp.PublicKey()
	if err != nil {
		return nil, nil, err
	}
	uc := jwt.NewUserClaims(upk)
	uc.IssuerAccount = apk
	userJwt, err := uc.Encode(askp)
	if err != nil {
		return nil, nil, err
	}
	useed, err := ukp.Seed()
	if err != nil {
		return nil, nil, err
	}
	creds, err := jwt.FormatUserConfig(userJwt, useed)
	if err != nil {
		return nil, nil, err
	}
	return &accountResolution{apk, aJWT}, creds, nil
}

func createSystemAccount(oskp nkeys.KeyPair) (*accountResolution, []byte, error) {
	return createUserAccount(oskp, func(ac *jwt.AccountClaims) {
		ac.Name = "SYS"
		ac.DefaultPermissions.Pub.Allow.Add(">")
		ac.DefaultPermissions.Sub.Allow.Add(">")

		ac.Exports.Add(&jwt.Export{
			Name:                 "account-monitoring-services",
			Subject:              "$SYS.REQ.ACCOUNT.*.>",
			Type:                 jwt.Service,
			ResponseType:         jwt.ResponseTypeStream,
			AccountTokenPosition: 4,
			Info: jwt.Info{
				Description: `Request account specific monitoring services for: SUBSZ, CONNZ, LEAFZ, JSZ and INFO`,
				InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
			},
		}, &jwt.Export{
			Name:                 "account-monitoring-streams",
			Subject:              "$SYS.ACCOUNT.*.>",
			Type:                 jwt.Stream,
			AccountTokenPosition: 3,
			Info: jwt.Info{
				Description: `Account specific monitoring stream`,
				InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
			},
		}, &jwt.Export{
			Name:                 "user-monitoring-services",
			Subject:              "$SYS.REQ.USER.*.*",
			Type:                 jwt.Stream,
			AccountTokenPosition: 5,
			Info: jwt.Info{
				Description: `Request user specific monitoring services for: SUBSZ, CONNZ, LEAFZ, JSZ and INFO`,
				InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
			},
		}, &jwt.Export{
			Name:                 "server-monitoring-services",
			Subject:              "$SYS.REQ.SERVER.*.*",
			Type:                 jwt.Stream,
			AccountTokenPosition: 5,
			Info: jwt.Info{
				Description: `Request server specific monitoring services for: SUBSZ, CONNZ, LEAFZ, JSZ and INFO`,
				InfoURL:     "https://docs.nats.io/nats-server/configuration/sys_accounts",
			},
		})
	})
}
