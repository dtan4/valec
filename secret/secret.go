package secret

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	separatorRegExp = regexp.MustCompile(`^#(?:-+|=+)$`)
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
func (ss Secrets) CompareList(old Secrets) (added, updated, deleted Secrets) {
	newMap, oldMap := ss.ListToMap(), old.ListToMap()

	for _, c := range ss {
		v, ok := oldMap[c.Key]
		if !ok {
			added = append(added, c)
		} else if v != c.Value {
			updated = append(updated, c)
		}
	}

	for _, c := range old {
		_, ok := newMap[c.Key]
		if !ok {
			deleted = append(deleted, c)
		}
	}

	return added, updated, deleted
}

// ListToMap converts secret list to map
func (ss Secrets) ListToMap() map[string]string {
	secretMap := map[string]string{}

	for _, secret := range ss {
		secretMap[secret.Key] = secret.Value
	}

	return secretMap
}

// SaveAsDotenv saves secrets as dotenv format
func (ss Secrets) SaveAsDotenv(filename string, preserve bool) error {
	sslice := []string{}

	for _, secret := range ss {
		sslice = append(sslice, fmt.Sprintf("%s=%s", secret.Key, secret.Value))
	}

	body := []byte(strings.Join(sslice, "\n") + "\n")

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			if err2 := saveEntirely(filename, body); err2 != nil {
				return errors.Wrapf(err, "Failed to create and save file. filename=%s", filename)
			}
		} else {
			return errors.Errorf("File has something wrong. filename=%s", filename)
		}
	} else {
		if preserve {
			if err2 := saveWithSection(filename, body); err2 != nil {
				return errors.Wrapf(err, "Failed to overwrite file. filename=%s", filename)
			}
		} else {
			if err2 := saveEntirely(filename, body); err2 != nil {
				return errors.Wrapf(err, "Failed to overwrite file. filename=%s", filename)
			}
		}
	}

	return nil
}

func saveEntirely(filename string, body []byte) error {
	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		return errors.Wrapf(err, "Failed to save file. filename=%s", filename)
	}

	return nil
}

func saveWithSection(filename string, body []byte) error {
	fp, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to open file. filename=%s", filename)
	}

	sc := bufio.NewScanner(fp)
	preserve := false
	preserveLines := []string{""}

	for sc.Scan() {
		line := sc.Text()

		if !preserve && separatorRegExp.Match([]byte(line)) {
			preserve = true
		}

		if preserve {
			preserveLines = append(preserveLines, line)
		}
	}

	body = append(body, []byte(strings.Join(preserveLines, "\n")+"\n")...)

	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		return errors.Wrapf(err, "Failed to save file. filename=%s", filename)
	}

	return nil
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
