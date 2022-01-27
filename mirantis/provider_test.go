package mirantis

import (
	//"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"MIRANTIS": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

// func testAccPreCheck(t *testing.T) {
// 	if err := os.Getenv("MIRANTIS_USERNAME"); err == "" {
// 		t.Fatal("MIRANTIS_USERNAME must be set for acceptance tests")
// 	}
// 	if err := os.Getenv("MIRANTIS_PASSWORD"); err == "" {
// 		t.Fatal("MIRANTIS_PASSWORD must be set for acceptance tests")
// 	}
// }
