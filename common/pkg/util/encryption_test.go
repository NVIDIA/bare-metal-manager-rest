// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package util

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestCreateHash(t *testing.T) {
	key := "key"
	expectedHash := sha256.Sum256([]byte(key))

	hash := CreateHash(key)

	if !bytes.Equal(hash, expectedHash[:]) {
		t.Errorf("Expected hash %x but got %x", expectedHash, hash)
	}
}

func TestEncryptAndDecryptData(t *testing.T) {
	data := []byte("this is some data to encrypt")
	passphrase := "testpassphrase"

	// Encrypt the data
	encryptedData := EncryptData(data, passphrase)

	if len(encryptedData) == 0 {
		t.Fatal("Expected encrypted data, but got empty byte slice")
	}

	// Decrypt the data
	decryptedData := DecryptData(encryptedData, passphrase)

	// Verify that decrypted data matches the original data
	if !bytes.Equal(decryptedData, data) {
		t.Errorf("Expected decrypted data to be %s but got %s", data, decryptedData)
	}
}

func TestEncryptAndDecryptWithWrongPassphrase(t *testing.T) {
	data := []byte("this is some data to encrypt")
	passphrase := "testpassphrase"
	wrongPassphrase := "wrongpassphrase"

	// Encrypt the data
	encryptedData := EncryptData(data, passphrase)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic but did not get one")
		}
	}()

	// Attempt to decrypt with wrong passphrase (should cause a panic)
	DecryptData(encryptedData, wrongPassphrase)
}
