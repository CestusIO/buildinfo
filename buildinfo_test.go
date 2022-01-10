package buildinfo

import (
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Public Api test", func() {
	BeforeEach(func() {
		version = ""
		buildDate = ""
		name = ""
	})
	It("Shows an empty version when not initialized", func() {
		bi := ProvideBuildInfo()
		Expect(bi).To(Equal(BuildInfo{}))
	})
	Describe("No version yaml", func() {
		It("Fills name and version correctly when passing only name", func() {
			GenerateVersion("testapp")
			bi := ProvideBuildInfo()
			Expect(bi.Name).To(Equal("testapp"))
			c := strings.Contains(bi.Version, "-localbuild")
			Expect(c).To(BeTrue())
		})
		It("Fills name and version correctly when passing only name", func() {
			GenerateVersion("")
			bi := ProvideBuildInfo()
			Expect(bi.Name).To(Equal(""))
			c := strings.Contains(bi.Version, "-localbuild")
			Expect(c).To(BeTrue())
		})
	})
	Describe("With version yaml", func() {
		It("It works with empty yaml", func() {
			GenerateVersionFromVersionYaml(nil, "testapp")
			bi := ProvideBuildInfo()
			Expect(bi.Name).To(Equal("testapp"))
			c := strings.Contains(bi.Version, "-localbuild")
			Expect(c).To(BeTrue())
		})
		It("It works with an rc", func() {
			file, err := os.ReadFile("./testdata/rcversion.yml")
			Expect(err).ToNot(HaveOccurred())
			GenerateVersionFromVersionYaml(file, "testapp")
			bi := ProvideBuildInfo()
			Expect(bi.Name).To(Equal("testapp"))
			c := strings.Contains(bi.Version, "-localbuild")
			Expect(c).To(BeTrue())
			v := strings.Contains(bi.Version, "v0.0.1-rc")
			Expect(v).To(BeTrue())
		})
		It("It works with a release", func() {
			file, err := os.ReadFile("./testdata/version.yml")
			Expect(err).ToNot(HaveOccurred())
			GenerateVersionFromVersionYaml(file, "testapp")
			bi := ProvideBuildInfo()
			Expect(bi.Name).To(Equal("testapp"))
			Expect(bi.Version).To(Equal("v0.0.1"))
			// other fields are dynamic so only test for them not being empty
			_, err = time.Parse(time.RFC3339, bi.BuildDate)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(bi.GoVersion)).ToNot(Equal(0))
			Expect(len(bi.OS)).ToNot(Equal(0))
			Expect(len(bi.Platform)).ToNot(Equal(0))
		})
		It("It does not overwrite values set through LDFLAGs", func() {
			version = "v1.0.0"
			name = "other app"
			buildDate = "2022-01-10T20:34:51Z"
			file, err := os.ReadFile("./testdata/version.yml")
			Expect(err).ToNot(HaveOccurred())
			GenerateVersionFromVersionYaml(file, "testapp")
			bi := ProvideBuildInfo()
			Expect(bi.Name).To(Equal("other app"))
			Expect(bi.Version).To(Equal("v1.0.0"))
			// other fields are dynamic so only test for them not being empty
			_, err = time.Parse(time.RFC3339, bi.BuildDate)
			Expect(err).ToNot(HaveOccurred())
			Expect(bi.BuildDate).To(Equal("2022-01-10T20:34:51Z"))
			Expect(len(bi.GoVersion)).ToNot(Equal(0))
			Expect(len(bi.OS)).ToNot(Equal(0))
			Expect(len(bi.Platform)).ToNot(Equal(0))
		})
	})
})
