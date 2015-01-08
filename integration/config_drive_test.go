package integration_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	. "github.com/cloudfoundry/bosh-agent/integration"
)

var _ = Describe("ConfigDrive", func() {
	var (
		testEnvironment TestEnvironment
	)

	BeforeEach(func() {
		logger := boshlog.NewLogger(boshlog.LevelDebug)
		cmdRunner := boshsys.NewExecCmdRunner(logger)
		testEnvironment = NewTestEnvironment(cmdRunner)
	})

	Context("when infrastructure is openstack", func() {
		BeforeEach(func() {
			err := testEnvironment.SetInfrastructure("openstack")
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when vm is using config drive", func() {
			BeforeEach(func() {
				err := testEnvironment.SetupConfigDrive()
				Expect(err).ToNot(HaveOccurred())

				err = testEnvironment.RemoveAgentSettings()
				Expect(err).ToNot(HaveOccurred())

				registrySettings := RegistrySettings{
					AgentID: "fake-agent-id",
				}

				err = testEnvironment.StartRegistry(registrySettings)
				Expect(err).ToNot(HaveOccurred())

				err = testEnvironment.UpdateAgentConfig("config-drive-agent.json")
				Expect(err).ToNot(HaveOccurred())

				err = testEnvironment.RestartAgent()
				Expect(err).ToNot(HaveOccurred())
			})

			It("using config drive to get registry URL", func() {
				settingsJSON, err := testEnvironment.GetFileContents("/var/vcap/bosh/settings.json")
				Expect(err).ToNot(HaveOccurred())
				Expect(settingsJSON).To(ContainSubstring("fake-agent-id"))
			})

			It("config drive is being unmounted", func() {
				Eventually(func() string {
					result, _ := testEnvironment.RunCommand("sudo mount | grep -c /dev/loop2")
					return strings.TrimSpace(result)
				}, 5*time.Second, 1*time.Second).Should(Equal("0"))
			})
		})
	})
})