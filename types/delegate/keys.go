package types

import (
	"fmt"
)

const (
	ModuleName   = "delegate"
	RouterKey    = ModuleName
	StoreKey     = ModuleName
	QuerierRoute = ModuleName

	DelegateBindKey      = "delegate_bind_%v_%v_%v"
	DelegateValidatorKey = "delegate_validator_%v_%v"
	DelegateEpochKey     = "delegate_epoch_%v"

	CysicVeTokenKey = "cysic_ve_token_%v"
)

// prefix bytes for the delegate persistent store
const (
	prefixStorage = iota + 1
	prefixParams
)

// KVStore key prefixes
var (
	KeyPrefixStorage = []byte{prefixStorage}
)

func GetDelegateBindKey(epoch int64, tokenID, worker string) string {
	return fmt.Sprintf(DelegateBindKey, epoch, tokenID, worker)
}

func GetDelegateValidatorKey(epoch int64, validator string) string {
	return fmt.Sprintf(DelegateValidatorKey, epoch, validator)
}

func GetDelegateEpochKey(epoch int64) string {
	return fmt.Sprintf(DelegateEpochKey, epoch)
}

func GetTokenKey(token string) string {
	return fmt.Sprintf(CysicVeTokenKey, token)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// EpochDelegateStoragePrefix returns a prefix to iterate over a given account storage.
func EpochDelegateStoragePrefix(epochID int64, validator string) []byte {
	return append(KeyPrefixStorage, []byte(GetDelegateValidatorKey(epochID, validator))...)
}
