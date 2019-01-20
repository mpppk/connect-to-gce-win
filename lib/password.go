//  Copyright 2018 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mpppk/connect-to-gce-win/spinner"

	daisyCompute "github.com/GoogleCloudPlatform/compute-image-tools/daisy/compute"
	"google.golang.org/api/compute/v1"
)

func getInstanceMetadata(client daisyCompute.Client, i, z, p string) (*compute.Metadata, error) {
	ins, err := client.GetInstance(p, z, i)
	if err != nil {
		return nil, fmt.Errorf("error getting instance: %v", err)
	}

	return ins.Metadata, nil
}

type windowsKeyJSON struct {
	ExpireOn string
	Exponent string
	Modulus  string
	UserName string
}

func generateKey(priv *rsa.PublicKey, u string) (*windowsKeyJSON, error) {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(priv.E))

	return &windowsKeyJSON{
		ExpireOn: time.Now().Add(5 * time.Minute).Format(time.RFC3339),
		// This is different than what the other tools produce,
		// AQAB vs AQABAA==, both are decoded as 65537.
		Exponent: base64.StdEncoding.EncodeToString(bs),
		Modulus:  base64.StdEncoding.EncodeToString(priv.N.Bytes()),
		UserName: u,
	}, nil
}

type credsJSON struct {
	ErrorMessage      string `json:"errorMessage,omitempty"`
	EncryptedPassword string `json:"encryptedPassword,omitempty"`
	Modulus           string `json:"modulus,omitempty"`
}

func getEncryptedPassword(client daisyCompute.Client, i, z, p, mod string) (string, error) {
	out, err := client.GetSerialPortOutput(p, z, i, 4, 0)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(out.Contents, "\n") {
		var creds credsJSON
		if err := json.Unmarshal([]byte(line), &creds); err != nil {
			continue
		}
		if creds.Modulus == mod {
			if creds.ErrorMessage != "" {
				return "", fmt.Errorf("error from agent: %s", creds.ErrorMessage)
			}
			return creds.EncryptedPassword, nil
		}
	}
	return "", errors.New("password not found in serial output")
}

func decryptPassword(priv *rsa.PrivateKey, ep string) (string, error) {
	bp, err := base64.StdEncoding.DecodeString(ep)
	if err != nil {
		return "", fmt.Errorf("error decoding password: %v", err)
	}
	pwd, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, priv, bp, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting password: %v", err)
	}
	return string(pwd), nil
}

func ResetPassword(client daisyCompute.Client, instanceName, zone, project, userName string) (string, error) {
	md, err := getInstanceMetadata(client, instanceName, zone, project)
	if err != nil {
		return "", fmt.Errorf("error getting instance metadata: %v", err)
	}

	spinner.ChangeStatus("Generating public/private key pair")
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	winKey, err := generateKey(&key.PublicKey, userName)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(winKey)
	if err != nil {
		return "", err
	}

	winKeys := string(data)
	var found bool
	for _, mdi := range md.Items {
		if mdi.Key == "windows-keys" {
			val := fmt.Sprintf("%s\n%s", *mdi.Value, winKeys)
			mdi.Value = &val
			found = true
			break
		}
	}
	if !found {
		md.Items = append(md.Items, &compute.MetadataItems{Key: "windows-keys", Value: &winKeys})
	}

	spinner.ChangeStatus("Setting new 'windows-keys' metadata")
	if err := client.SetInstanceMetadata(project, zone, instanceName, md); err != nil {
		return "", err
	}

	spinner.ChangeStatus("Fetching encrypted password")

	var trys int
	var ep string
	for {
		time.Sleep(1 * time.Second)
		ep, err = getEncryptedPassword(client, instanceName, zone, project, winKey.Modulus)
		if err == nil {
			break
		}
		if trys > 10 {
			return "", err
		}
		trys++
	}

	spinner.ChangeStatus("Decrypting password")

	password, err := decryptPassword(key, ep)
	spinner.CompleteStatus()
	return password, err
}
