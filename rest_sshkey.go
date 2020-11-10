package ne

import (
	"net/url"

	"github.com/equinix/ne-go/internal/api"
	"github.com/go-resty/resty/v2"
)

//GetSSHPublicKeys retrieves list of available SSH public keys
func (c RestClient) GetSSHPublicKeys() ([]SSHPublicKey, error) {
	path := "/ne/v1/device/public-keys"
	respBody := make([]api.SSHPublicKey, 0)
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, resty.MethodGet, path); err != nil {
		return nil, err
	}
	return mapSSHPublicKeysAPIToDomain(respBody), nil
}

//GetSSHPublicKey retrieves SSH public key with a given identifier
func (c RestClient) GetSSHPublicKey(uuid string) (*SSHPublicKey, error) {
	path := "/ne/v1/device/public-keys/" + url.PathEscape(uuid)
	respBody := api.SSHPublicKey{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, resty.MethodGet, path); err != nil {
		return nil, err
	}
	mapped := mapSSHPublicKeyAPIToDomain(respBody)
	return &mapped, nil
}

//CreateSSHPublicKey creates new SSH public key with a given details
func (c RestClient) CreateSSHPublicKey(key SSHPublicKey) (string, error) {
	path := "/ne/v1/device/public-keys"
	reqBody := mapSSHPublicKeyDomainToAPI(key)
	respBody := api.SSHPublicKeyCreateResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, resty.MethodPost, path); err != nil {
		return "", err
	}
	return respBody.UUID, nil
}

//DeleteSSHPublicKey removes SSH Public key with given identifier
func (c RestClient) DeleteSSHPublicKey(uuid string) error {
	path := "/ne/v1/device/public-keys/" + url.PathEscape(uuid)
	if err := c.Execute(c.R(), resty.MethodDelete, path); err != nil {
		return err
	}
	return nil
}

func mapSSHPublicKeysAPIToDomain(apiKeys []api.SSHPublicKey) []SSHPublicKey {
	transformed := make([]SSHPublicKey, len(apiKeys))
	for i := range apiKeys {
		transformed[i] = mapSSHPublicKeyAPIToDomain(apiKeys[i])
	}
	return transformed
}

func mapSSHPublicKeyAPIToDomain(apiKey api.SSHPublicKey) SSHPublicKey {
	return SSHPublicKey{
		UUID:  apiKey.UUID,
		Name:  apiKey.KeyName,
		Value: apiKey.KeyValue,
	}
}

func mapSSHPublicKeyDomainToAPI(key SSHPublicKey) api.SSHPublicKey {
	return api.SSHPublicKey{
		UUID:     key.UUID,
		KeyName:  key.Name,
		KeyValue: key.Value,
	}
}