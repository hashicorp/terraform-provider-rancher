package rancher

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancherCertificate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancherCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rancher_certificate.foo", "name", "foo"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckRancherCertificateConfig = `
resource "rancher_certificate" "foo" {
	name = "foo"
	environment_id = "1a5"
	cert = <<EOT
-----BEGIN CERTIFICATE-----
MIIB4TCCAYugAwIBAgIUEJP/YW4QmK/dt85L+gaNsfaJUhQwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMDAxMTUxMjMyMDVaFw0yMDAy
MTQxMjMyMDVaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwXDANBgkqhkiG9w0BAQEF
AANLADBIAkEA6YGnZfYxwQOEk4L2ZcbsghBDTt7MD2+STKshmSv0yUfI0lhmogT+
NzsHjGbP2onZV5Pw8mMy4Snu7D+0zm0q/wIDAQABo1MwUTAdBgNVHQ4EFgQURhK1
Sh9akzaFPxpvsrB27AXKcVgwHwYDVR0jBBgwFoAURhK1Sh9akzaFPxpvsrB27AXK
cVgwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAANBAH6Qf8dD8SHP9rNI
JTsRiUYfDIJfGzCni3on3EdscBanEbb3LAAmWCI0fJ/tMbzAPdGcTyuK5mSrVhBr
+2EQTnw=
-----END CERTIFICATE-----
EOT
    key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEA6YGnZfYxwQOEk4L2
ZcbsghBDTt7MD2+STKshmSv0yUfI0lhmogT+NzsHjGbP2onZV5Pw8mMy4Snu7D+0
zm0q/wIDAQABAkEAlKwPaDTzgr/5pm4o8a5RIZK3OD1U0bMpBBWlo7+/8HLDoJnk
d5ssmuzg8xOLadXhZnjgkUPk55cxyfzzKCNI4QIhAPyzEdn8GjxlKSgzfZ1pZ4gl
JHkLpBKPYB5KwJpnGlTxAiEA7I5o94ypM7yf5rKMTvF/TkGRrXKha0HH/ZHHqYqd
vu8CIHdAypvkrTzzQIkIQ6+VnpZRcPTu2W8o2mNxQ5OaNIMBAiBlrvWJ64nT9nHZ
jchoKsDpV6ASKaMfYsBfzClCRJZ4OwIhALqR1ZVh2ktONmMdBk4h/Hpjuim0Y0EZ
0cwepge3Oz6M
-----END PRIVATE KEY-----
EOT
}
`
