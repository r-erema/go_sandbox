package example1_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangingFileOffsetInCaseOfMultipleReading(t *testing.T) {
	t.Parallel()

	testDataFilePath := "./test_data/test_text_data"
	tempFile, err := os.OpenFile(testDataFilePath, os.O_RDONLY, 0o755)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = tempFile.Close()
		require.NoError(t, err)
	})

	fileDescriptor := tempFile.Fd()

	file := os.NewFile(fileDescriptor, "")
	buf := make([]byte, 1)
	bytesCount, err := file.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 1, bytesCount)
	assert.Equal(t, []byte("1"), buf)

	file2 := os.NewFile(fileDescriptor, "")
	buf = make([]byte, 2)
	bytesCount, err = file2.Read(buf)
	assert.Equal(t, 2, bytesCount)
	assert.Equal(t, []byte("23"), buf)

	file3 := os.NewFile(fileDescriptor, "")
	buf = make([]byte, 3)
	bytesCount, err = file3.Read(buf)
	assert.Equal(t, 3, bytesCount)
	assert.Equal(t, []byte("456"), buf)
}
