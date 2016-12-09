package secret

import (
	"io/ioutil"
	"sort"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Secret represents key=value pair
type Secret struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Secrets represents the array of Secret
type Secrets []*Secret

// Len returns the length of the array
func (ss Secrets) Len() int {
	return len(ss)
}

// Less returns Secrets[i] is less than Secrets[j]
func (ss Secrets) Less(i, j int) bool {
	si, sj := ss[i], ss[j]

	if si.Key < sj.Key {
		return true
	}

	if si.Key > sj.Key {
		return false
	}

	if si.Value < sj.Value {
		return true
	}

	if si.Value > sj.Value {
		return false
	}

	return false
}

// Swap swaps Secrets[i] and Secrets[j]
func (ss Secrets) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

// CompareList compares two secret lists and returns the differences between them
func (ss Secrets) CompareList(target Secrets) (Secrets, Secrets) {
	added, deleted := Secrets{}, Secrets{}
	srcMap, dstMap := ss.ListToMap(), target.ListToMap()

	for _, c := range ss {
		v, ok := dstMap[c.Key]
		if !ok || v != c.Value {
			added = append(added, c)
		}
	}

	for _, c := range target {
		v, ok := srcMap[c.Key]
		if !ok || v != c.Value {
			deleted = append(deleted, c)
		}
	}

	return added, deleted
}

// ListToMap converts secret list to map
func (ss Secrets) ListToMap() map[string]string {
	secretMap := map[string]string{}

	for _, secret := range ss {
		secretMap[secret.Key] = secret.Value
	}

	return secretMap
}

// SaveAsYAML saves secrets to local secret file
func (ss Secrets) SaveAsYAML(filename string) error {
	body, err := yaml.Marshal(ss)
	if err != nil {
		return errors.Wrap(err, "Failed to convert secrets as YAML.")
	}

	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		return errors.Wrapf(err, "Failed to save file. filename=%s", filename)
	}

	return nil
}

// LoadFromYAML loads secrets from the given YAML file
func LoadFromYAML(filename string) (Secrets, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return Secrets{}, errors.Wrapf(err, "Failed to read secret file. filename=%s", filename)
	}

	var secrets Secrets

	if err := yaml.Unmarshal(body, &secrets); err != nil {
		return Secrets{}, errors.Wrapf(err, "Failed to parse secret file as YAML. filename=%s", filename)
	}

	return secrets, nil
}

// MapToList converts map to secret list
func MapToList(secretMap map[string]string) Secrets {
	secrets := Secrets{}

	for key, value := range secretMap {
		secrets = append(secrets, &Secret{
			Key:   key,
			Value: value,
		})
	}

	sort.Sort(secrets)

	return secrets
}
