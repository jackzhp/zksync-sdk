package zkscrypto

/*
#cgo LDFLAGS: -L${SRCDIR}/libs -lzks_crypto

#include "zks_crypto.h"
*/
import "C"
import (
	"encoding/hex"
	"errors"
	"unsafe"
)

// MaxSignedMessageLen Maximum byte length of the message that can be signed.
const MaxSignedMessageLen = 92

// PackedSignatureLen Byte length of the signature. Signature contains r and s points.
const PackedSignatureLen = 64

// PrivateKeyLen Byte length of the private key.
const PrivateKeyLen = 32

// PubkeyHashLen Byte length of the public key hash.
const PubkeyHashLen = 20

// PublicKeyLen Byte length of the public key.
const PublicKeyLen = 32

var (
	errSeedLen      = errors.New("Given seed is too short, length must be greater than 32")
	errSignedMsgLen = errors.New("Musig message length must not be larger than 92")
)

func init() {
	C.zks_crypto_init()
}

/*
************************************************************************************************
Private key implementation
************************************************************************************************
*/

// NewPrivateKey generates private key from seed
func NewPrivateKey(seed []byte) (*PrivateKey, error) {
	pointer := C.struct_ZksPrivateKey{}
	rawSeed := C.CBytes(seed)
	defer C.free(rawSeed)
	result := C.zks_crypto_private_key_from_seed((*C.uint8_t)(rawSeed), C.ulong(len(seed)), &pointer)
	if result != 0 {
		switch result {
		case 1:
			return nil, errSeedLen
		}
	}
	data := unsafe.Pointer(&pointer.data)
	return &PrivateKey{data: C.GoBytes(data, PrivateKeyLen)}, nil
}

// Sign message with musig Schnorr signature scheme
func (pk *PrivateKey) Sign(message []byte) (*Signature, error) {
	privateKeyC := C.struct_ZksPrivateKey{}
	rawMessage := C.CBytes(message)
	defer C.free(rawMessage)
	for i := range pk.data {
		privateKeyC.data[i] = C.uint8_t(pk.data[i])
	}
	signatureC := C.struct_ZksSignature{}
	result := C.zks_crypto_sign_musig(&privateKeyC, (*C.uint8_t)(rawMessage), C.ulong(len(message)), &signatureC)
	if result != 0 {
		switch result {
		case 1:
			return nil, errSignedMsgLen
		}
	}
	data := unsafe.Pointer(&signatureC.data)
	return &Signature{data: C.GoBytes(data, PackedSignatureLen)}, nil
}

// PublicKey generates public key from private key
func (pk *PrivateKey) PublicKey() (*PublicKey, error) {
	privateKeyC := C.struct_ZksPrivateKey{}
	for i := range pk.data {
		privateKeyC.data[i] = C.uint8_t(pk.data[i])
	}
	pointer := C.struct_ZksPackedPublicKey{}
	result := C.zks_crypto_private_key_to_public_key(&privateKeyC, &pointer)
	if result != 0 {
		return nil, errors.New("Error on public key generation")
	}
	data := unsafe.Pointer(&pointer.data)
	return &PublicKey{data: C.GoBytes(data, PublicKeyLen)}, nil
}

// HexString creates a hex string representation of a private key
func (pk *PrivateKey) HexString() string {
	if pk.data == nil || len(pk.data) == 0 {
		return "0x"
	}
	return hex.EncodeToString(pk.data)
}

/*
************************************************************************************************
Public key implementation
************************************************************************************************
*/

// Hash generates hash from public key
func (pk *PublicKey) Hash() (*PublicKeyHash, error) {
	publicKeyC := C.struct_ZksPackedPublicKey{}
	for i := range pk.data {
		publicKeyC.data[i] = C.uint8_t(pk.data[i])
	}
	pointer := C.struct_ZksPubkeyHash{}
	result := C.zks_crypto_public_key_to_pubkey_hash(&publicKeyC, &pointer)
	if result != 0 {
		return nil, errors.New("Error on public key hash generation")
	}
	data := unsafe.Pointer(&pointer.data)
	return &PublicKeyHash{data: C.GoBytes(data, PubkeyHashLen)}, nil
}

// HexString creates a hex string representation of a public key
func (pk *PublicKey) HexString() string {
	if pk.data == nil || len(pk.data) == 0 {
		return "0x"
	}
	return hex.EncodeToString(pk.data)
}

/*
************************************************************************************************
Private key Hash implementation
************************************************************************************************
*/

// HexString creates a hex string representation of a public key hash
func (pk *PublicKeyHash) HexString() string {
	if pk.data == nil || len(pk.data) == 0 {
		return "0x"
	}
	return hex.EncodeToString(pk.data)
}

/*
************************************************************************************************
Signature implementation
************************************************************************************************
*/

// HexString creates a hex string representation of a signature
func (pk *Signature) HexString() string {
	if pk.data == nil || len(pk.data) == 0 {
		return "0x"
	}
	return hex.EncodeToString(pk.data)
}
