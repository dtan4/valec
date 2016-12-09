package secret

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Secret represents key=value pair
type Secret struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// CompareSecretList compares two secret lists and returns the differences between them
func CompareSecretList(src, dst []*Secret) ([]*Secret, []*Secret) {
	added, deleted := []*Secret{}, []*Secret{}
	srcMap, dstMap := SecretsToMap(src), SecretsToMap(dst)

	for _, c := range src {
		v, ok := dstMap[c.Key]
		if !ok || v != c.Value {
			added = append(added, c)
		}
	}

	for _, c := range dst {
		v, ok := srcMap[c.Key]
		if !ok || v != c.Value {
			deleted = append(deleted, c)
		}
	}

	return added, deleted
}

// SecretsToMap converts secret list to map
func SecretsToMap(secrets []*Secret) map[string]string {
	secretMap := map[string]string{}

	for _, secret := range secrets {
		secretMap[secret.Key] = secret.Value
	}

	return secretMap
}

// LoadSecretYAML loads secrets from the given YAML file
func LoadSecretYAML(filename string) ([]*Secret, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return []*Secret{}, errors.Wrapf(err, "Failed to read secret file. filename=%s", filename)
	}

	var secrets []*Secret

	if err := yaml.Unmarshal(body, &secrets); err != nil {
		return []*Secret{}, errors.Wrapf(err, "Failed to parse secret file as YAML. filename=%s", filename)
	}

	return secrets, nil
}

// MapToSecrets converts map to secret list
func MapToSecrets(secretMap map[string]string) []*Secret {
	secrets := []*Secret{}

	for key, value := range secretMap {
		secrets = append(secrets, &Secret{
			Key:   key,
			Value: value,
		})
	}

	return secrets
}

// SaveAsYAML saves secrets to local secret file
func SaveAsYAML(secrets []*Secret, filename string) error {
	body, err := yaml.Marshal(secrets)
	if err != nil {
		return errors.Wrap(err, "Failed to convert secrets as YAML.")
	}

	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		return errors.Wrapf(err, "Failed to save file. filename=%s", filename)
	}

	return nil
}
