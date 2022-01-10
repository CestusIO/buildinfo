package buildinfo

import (
	"errors"
	"time"

	"github.com/Scardiecat/svermaker"
	"github.com/Scardiecat/svermaker/semver"
	"gopkg.in/yaml.v3"
)

// serializer implements the Serializer interface
type serializer struct {
	versionyaml []byte
	// Services
	projectVersionService semver.ProjectVersionService
}

// newSerializer creates a Serializer
func newSerializer(versionyaml []byte) *serializer {
	s := &serializer{
		versionyaml: versionyaml,
	}
	s.projectVersionService.Serializer = s
	return s
}

// Exists checks if the file exists
func (s *serializer) Exists() bool {
	return len(s.versionyaml) > 0
}

type projectVersion struct {
	Current string
	Next    string
}

// Serialize writes the ProjectVersion to a yml
func (s *serializer) Serialize(p svermaker.ProjectVersion) error {
	return errors.New("not implemented")
}

// Deserialize reads a projectcersion from a yml
func (s *serializer) Deserialize() (*svermaker.ProjectVersion, error) {
	v := projectVersion{}
	m := semver.Manipulator{}
	projectVersion := svermaker.ProjectVersion{}
	if s.Exists() {
		if err := yaml.Unmarshal(s.versionyaml, &v); err == nil {

		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("version.yaml does not exist")
	}

	if current, err := m.Create(v.Current); err == nil {
		projectVersion.Current = *current
	} else {
		return nil, err
	}

	if next, err := m.Create(v.Next); err == nil {
		projectVersion.Next = *next
	} else {
		return nil, err
	}

	return &projectVersion, nil
}

func generate(versionyaml []byte, appName string) (version string, name string, date string) {
	name = appName
	var serializer = newSerializer(versionyaml)
	var pvs = semver.ProjectVersionService{Serializer: serializer}
	var meta []string
	v, err := pvs.Get()
	if err != nil {
		return
	}
	version, err = buildVersionString(*v, meta)
	if err != nil {
		return
	}
	date = time.Now().UTC().Format(time.RFC3339)
	return
}

func buildVersionString(p svermaker.ProjectVersion, buildMetadata []string) (string, error) {
	m := semver.Manipulator{}

	isRelease := m.Compare(p.Current, p.Next) == 0
	c := p.Current
	if !isRelease {
		md, err := m.SetMetadata(c, buildMetadata)
		if err != nil {
			return "", err
		}
		c = md
		pre := c.Pre
		pre = append(pre, svermaker.PRVersion{VersionStr: localBuild(), VersionNum: 0, IsNum: false})
		c.Pre = pre

	}
	return "v" + c.String(), nil
}
