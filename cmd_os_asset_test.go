package main

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	. "github.com/onsi/gomega"
)

func TestAssetsListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	assetList := map[string]metalcloud.OSAsset{
		"test": {
			OSAssetID:    10,
			OSAssetUsage: "test",
		},
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		OSAssets().
		Return(&assetList, nil).
		AnyTimes()

	//test json

	expectedFirstRow := map[string]interface{}{
		"ID":    10,
		"USAGE": "test",
	}

	testListCommand(assetsListCmd, nil, client, expectedFirstRow, t)

}

func TestCreateAssetCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"1": tmpl,
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID: 100,
	}

	client.EXPECT().
		OSAssetCreate(gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSAssetMakePublic(gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateAddOSAsset(tmpl.VolumeTemplateID, asset.OSAssetID, gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf1",
				"usage":                  "testf1",
				"read_content_from_pipe": true,
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good2, associate a template (id)",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf2",
				"usage":                  "testf2",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
				"path":                   "test2",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good3, associate a template",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf3",
				"usage":                  "testf3",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test3",
				"variables_json":         "['1': 'test']",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good4, associate a template (name)",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf4",
				"usage":                  "testf4",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test4",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good5, delete asset if it exists",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf6",
				"usage":                  "testf6",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
				"path":                   "test5",
				"delete_if_exists":       false,
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, non-existant template",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf5",
				"usage":                  "testf5",
				"read_content_from_pipe": true,
				"template_id_or_name":    "tmpl1",
				"path":                   "test5",
			}),
			good: false,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, missing path",
			cmd: MakeCommand(map[string]interface{}{
				"filename":               "testf6",
				"usage":                  "testf6",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
			}),
			good: false,
			id:   asset.OSAssetID,
		},
	}

	testCreateCommand(assetCreateCmd, cases, client, t)
}

func TestDeleteAssetCmd(t *testing.T) {
	RegisterTestingT(t)
	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	asset := metalcloud.OSAsset{
		OSAssetID: 100,
	}

	client.EXPECT().
		OSAssetGet(asset.OSAssetID).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSAssetDelete(asset.OSAssetID).
		Return(nil).
		MinTimes(1)

	cmd := MakeCommand(map[string]interface{}{"asset_id_or_name": asset.OSAssetID})
	testCommandWithConfirmation(assetDeleteCmd, cmd, client, t)
}

func TestEditAssetCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"1": tmpl,
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID:       100,
		OSAssetFileName: "test",
	}

	assetf := metalcloud.OSAsset{
		OSAssetID: 101,
	}

	assetl := map[string]metalcloud.OSAsset{
		"1": asset,
	}

	client.EXPECT().
		OSAssets().
		Return(&assetl, nil).
		AnyTimes()

	client.EXPECT().
		OSAssetGet(asset.OSAssetID).
		Return(&asset, nil).
		AnyTimes()

	client.EXPECT().
		OSAssetMakePublic(gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSAssetGet(asset.OSAssetFileName).
		Return(&asset, nil).
		AnyTimes()

	client.EXPECT().
		OSAssetGet(assetf.OSAssetID).
		Return(nil, fmt.Errorf("test")).
		Times(1)

	client.EXPECT().
		OSAssetUpdate(gomock.Any(), gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateAddOSAsset(tmpl.VolumeTemplateID, asset.OSAssetID, gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf1",
				"usage":                  "testf1",
				"read_content_from_pipe": true,
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good2, associate a template (id)",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf2",
				"usage":                  "testf2",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
				"path":                   "test2",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good3, associate a template",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf3",
				"usage":                  "testf3",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test3",
				"variables_json":         "['1': 'test']",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "good4, associate a template (name)",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf4",
				"usage":                  "testf4",
				"read_content_from_pipe": true,
				"template_id_or_name":    "test",
				"path":                   "test4",
			}),
			good: true,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, non-existant template",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf5",
				"usage":                  "testf5",
				"read_content_from_pipe": true,
				"template_id_or_name":    "tmpl1",
				"path":                   "test5",
			}),
			good: false,
			id:   asset.OSAssetID,
		},
		{
			name: "associate a template, missing path",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       100,
				"filename":               "testf6",
				"usage":                  "testf6",
				"read_content_from_pipe": true,
				"template_id_or_name":    10,
			}),
			good: false,
			id:   asset.OSAssetID,
		},
		{
			name: "asset not found",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       101,
				"filename":               "testf1",
				"usage":                  "testf1",
				"read_content_from_pipe": true,
			}),
			good: false,
			id:   0,
		},
		{
			name: "asset not found",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":       "file-test",
				"filename":               "testf1",
				"usage":                  "testf1",
				"read_content_from_pipe": true,
			}),
			good: false,
			id:   0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := assetEditCmd(&c.cmd, client)
			if c.good && err != nil {
				//t.Error(err)
			}
		})
	}
}

func TestAssociateAssetCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"1": tmpl,
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID: 100,
	}

	client.EXPECT().
		OSAssetGet(gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateAddOSAsset(tmpl.VolumeTemplateID, asset.OSAssetID, gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":    100,
				"template_id_or_name": 10,
				"path":                "test",
				"variables_json":      "['1': 'test', '2': 'test1']",
			}),
			good: true,
			id:   0,
		},
		{
			name: "good2, associate a template (id)",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":    100,
				"template_id_or_name": 10,
				"path":                "test",
				"variables_json":      "['1': 'test', '2': 'test1']",
			}),
			good: true,
			id:   0,
		},
		{
			name: "good3, associate a template without variables_json",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":    100,
				"template_id_or_name": 10,
				"path":                "test",
			}),
			good: true,
			id:   0,
		},
		{
			name: "associate a template, non-existant template",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":    100,
				"template_id_or_name": "tmpl1",
				"path":                "test",
				"variables_json":      "['1': 'test', '2': 'test1']",
			}),
			good: false,
			id:   0,
		},
		{
			name: "associate a template, missing path",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":    100,
				"template_id_or_name": "tmpl1",
				"variables_json":      "['1': 'test', '2': 'test1']",
			}),
			good: false,
			id:   0,
		},
	}

	testCreateCommand(associateAssetCmd, cases, client, t)
}

func TestAssociateAssetMissingVariablesCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	tmpl := metalcloud.OSTemplate{
		VolumeTemplateID:    10,
		VolumeTemplateLabel: "test",
	}

	tmpls := map[string]metalcloud.OSTemplate{
		"1": tmpl,
	}

	client.EXPECT().
		OSTemplateGet(gomock.Any(), false).
		Return(&tmpl, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplates().
		Return(&tmpls, nil).
		AnyTimes()

	asset := metalcloud.OSAsset{
		OSAssetID: 100,
	}

	client.EXPECT().
		OSAssetGet(gomock.Any()).
		Return(&asset, nil).
		MinTimes(1)

	client.EXPECT().
		OSTemplateAddOSAsset(tmpl.VolumeTemplateID, asset.OSAssetID, gomock.Any(), "[]").
		Return(nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name":    100,
				"template_id_or_name": 10,
				"path":                "test",
			}),
			good: true,
			id:   0,
		},
	}

	testCreateCommand(associateAssetCmd, cases, client, t)
}

func TestOSAssetMakePrivateCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	asset := metalcloud.OSAsset{
		OSAssetID:       100,
		OSAssetFileName: "test",
	}

	assets := map[string]metalcloud.OSAsset{
		"test": asset,
	}

	user := metalcloud.User{
		UserID: 1,
	}

	user1 := metalcloud.User{
		UserEmail: "test",
	}

	client.EXPECT().
		OSAssetGet(gomock.Any()).
		Return(&asset, nil).
		AnyTimes()

	client.EXPECT().
		OSAssets().
		Return(&assets, nil).
		MinTimes(1)

	client.EXPECT().
		UserGet(gomock.Any()).
		Return(&user, nil).
		AnyTimes()

	client.EXPECT().
		UserGetByEmail(gomock.Any()).
		Return(&user1, nil).
		MinTimes(1)

	client.EXPECT().
		OSAssetMakePrivate(gomock.Any(), gomock.Any()).
		Return(&asset, nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": 100,
				"user_id":          1,
			}),
			good: true,
			id:   0,
		},
		{
			name: "good2",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": "test",
				"user_id":          1,
			}),
			good: true,
			id:   0,
		},
		{
			name: "good3",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": 100,
				"user_id":          "test",
			}),
			good: true,
			id:   0,
		},
		{
			name: "asset not found",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": "test1",
				"user_id":          1,
			}),
			good: false,
			id:   0,
		},
		{
			name: "missing asset id or name",
			cmd: MakeCommand(map[string]interface{}{
				"user_id": 1,
			}),
			good: false,
			id:   0,
		},
		{
			name: "missing user id or email",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": "test",
			}),
			good: false,
			id:   0,
		},
	}

	testCreateCommand(assetMakePrivateCmd, cases, client, t)
}

func TestOSAssetMakePublicCmd(t *testing.T) {

	client := mock_metalcloud.NewMockMetalCloudClient(gomock.NewController(t))

	asset := metalcloud.OSAsset{
		OSAssetID:       100,
		OSAssetFileName: "test",
	}

	assets := map[string]metalcloud.OSAsset{
		"test": asset,
	}

	client.EXPECT().
		OSAssetGet(gomock.Any()).
		Return(&asset, nil).
		AnyTimes()

	client.EXPECT().
		OSAssets().
		Return(&assets, nil).
		MinTimes(1)

	client.EXPECT().
		OSAssetMakePublic(gomock.Any()).
		Return(&asset, nil).
		AnyTimes()

	cases := []CommandTestCase{
		{
			name: "good1",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": 100,
			}),
			good: true,
			id:   0,
		},
		{
			name: "good2",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": "test",
			}),
			good: true,
			id:   0,
		},
		{
			name: "asset not found",
			cmd: MakeCommand(map[string]interface{}{
				"asset_id_or_name": "test1",
			}),
			good: false,
			id:   0,
		},
		{
			name: "missing asset id or name",
			cmd:  MakeCommand(map[string]interface{}{}),
			good: false,
			id:   0,
		},
	}

	testCreateCommand(assetMakePublicCmd, cases, client, t)
}
