package utils_test

import (
	"testing"

	"github.com/Thales-s-Orgs/VideoEncoder-GoLang/framework/utils"
	"github.com/stretchr/testify/require"
)

var test_correct_json = `
{
	"id": "525b5fd9-700d-4feb-89c0-415a1e6e148c",
	"file_path": "convite.mp4",
	"status": "pending"
}
`

var test_incorrect_json = `Thales`

func TestIsJson(t *testing.T) {
	err := utils.IsJson(test_correct_json)
	require.Nil(t, err)

	err = utils.IsJson(test_incorrect_json)
	require.Error(t, err)
}
