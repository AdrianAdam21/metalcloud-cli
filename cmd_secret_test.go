package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	mock_metalcloud "github.com/metalsoft-io/metalcloud-cli/helpers"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

func TestSecretsListCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	secret := metalcloud.Secret{
		SecretName: "test",
	}

	list := map[string]metalcloud.Secret{
		"secret": secret,
	}

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	client.EXPECT().
		Secrets("").
		Return(&list, nil).
		AnyTimes()

	//test json
	format := "json"
	emptyStr := ""
	cmd := Command{
		Arguments: map[string]interface{}{
			"usage":  &emptyStr,
			"format": &format,
		},
	}

	ret, err := secretsListCmd(&cmd, client)
	Expect(err).To(BeNil())

	var m []interface{}
	err = json.Unmarshal([]byte(ret), &m)

	Expect(err).To(BeNil())

	r := m[0].(map[string]interface{})
	Expect(int(r["ID"].(float64))).To(Equal(0))
	Expect(r["NAME"].(string)).To(Equal(secret.SecretName))

	//test plaintext
	format = ""
	cmd = Command{
		Arguments: map[string]interface{}{
			"usage":  &emptyStr,
			"format": &format,
		},
	}

	ret, err = secretsListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	//test csv
	format = "csv"

	cmd = Command{
		Arguments: map[string]interface{}{
			"usage":  &emptyStr,
			"format": &format,
		},
	}

	ret, err = secretsListCmd(&cmd, client)
	Expect(err).To(BeNil())
	Expect(ret).NotTo(BeEmpty())

	reader := csv.NewReader(strings.NewReader(ret))

	csv, err := reader.ReadAll()
	Expect(csv[1][0]).To(Equal(fmt.Sprintf("%d", 0)))
	Expect(csv[1][1]).To(Equal("test"))

}

func TestSecretsDeleteCmd(t *testing.T) {
	RegisterTestingT(t)
	ctrl := gomock.NewController(t)

	client := mock_metalcloud.NewMockMetalCloudClient(ctrl)

	secret := metalcloud.Secret{
		SecretID:   10,
		SecretName: "test",
	}

	client.EXPECT().
		SecretGet(10).
		Return(&secret, nil).
		AnyTimes()

	client.EXPECT().
		SecretDelete(10).
		Return(nil).
		AnyTimes()

	list := map[string]metalcloud.Secret{
		"secret": secret,
	}
	client.EXPECT().
		Secrets("").
		Return(&list, nil).
		AnyTimes()

	//test json

	id := "test"
	bTrue := true
	cmd := Command{
		Arguments: map[string]interface{}{
			"secret_id_or_name": &id,
			"autoconfirm":       &bTrue,
		},
	}

	_, err := secretDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())

	//check with int id
	idint := 10
	cmd = Command{
		Arguments: map[string]interface{}{
			"secret_id_or_name": &idint,
			"autoconfirm":       &bTrue,
		},
	}

	_, err = secretDeleteCmd(&cmd, client)
	Expect(err).To(BeNil())

}
