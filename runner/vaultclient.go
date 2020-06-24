package runner

import (
	"log"

	"github.com/hashicorp/vault/api"
)

func vaultClient(wrapper_config Config, token string) *api.Client {
	config := &api.Config{
		Address: wrapper_config.VaultAddr,
	}

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	client.SetToken(token)
	return client
}

func getSecret(client *api.Client, path string, key string) string {
	c := client.Logical()
	secret, err := c.Read(path)
	if err != nil {
		log.Fatal(err)
	}
	if secret == nil {
		log.Fatal("[VaultWrapper][ERROR] Vault path: " + path + " does not exist")
	}
	return secret.Data[key].(string)
}
