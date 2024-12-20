package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfigFile(t *testing.T) {
	testFlagNetAddr := "localhost:8082"
	testFlagRestore := true
	testFlagStoreInterval := 1
	testFlagFileStoragePath := "test/file/path"
	testFlagDatabaseDsn := "test dsn"
	testFlagCryptoKey := "test crypto key"

	createFile := func(name string) {
		data := fmt.Sprintf("{\"address\": \"%s\",\"restore\": %t,\"store_interval\": \"%ds\",\"store_file\": \"%s\",\"database_dsn\": \"%s\",\"crypto_key\": \"%s\"}",
			testFlagNetAddr, testFlagRestore, testFlagStoreInterval, testFlagFileStoragePath, testFlagDatabaseDsn, testFlagCryptoKey)
		f, err := os.Create(name)
		require.NoError(t, err)
		_, err = f.Write([]byte(data))
		require.NoError(t, err)
	}
	nameFile := "./test_config.json"
	createFile(nameFile)

	configs, err := ParseConfigFile(nameFile)
	require.NoError(t, err)

	assert.Equal(t, testFlagNetAddr, configs.Address)
	assert.Equal(t, testFlagRestore, configs.Restore)
	assert.Equal(t, testFlagStoreInterval, int(configs.StoreInterval.Seconds()))
	assert.Equal(t, testFlagFileStoragePath, configs.StoreFile)
	assert.Equal(t, testFlagDatabaseDsn, configs.DatabaseDSN)
	assert.Equal(t, testFlagCryptoKey, configs.CryptoKey)

	err = os.Remove(nameFile)
	require.NoError(t, err)
}
