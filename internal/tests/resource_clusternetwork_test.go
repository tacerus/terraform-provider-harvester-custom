package tests

import (
	"context"
	"fmt"
	"testing"

	harvsternetworkv1 "github.com/harvester/harvester-network-controller/pkg/apis/network.harvesterhci.io/v1beta1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tacerus/terraform-provider-harvester-custom/pkg/client"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/constants"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/helper"
)

const (
	testAccClusterNetworkName         = "vlan"
	testAccClusterNetworkResourceName = constants.ResourceTypeClusterNetwork + "." + testAccClusterNetworkName
	testAccClusterNetworkDescription  = "Terraform Harvester ClusterNetwork acceptance test"

	testAccClusterNetworkEnable             = "true"
	testAccClusterNetworkDefaultPhysicalNIC = "eth0"

	testAccClusterNetworkConfigTemplate = `
resource %s "%s" {
	%s = "%s"
	%s = "%s"
	%s = %s
	%s = "%s"
}
`
)

func buildClusterNetworkConfig(name, description, enable, defaultPhysicalNIC string) string {
	return fmt.Sprintf(testAccClusterNetworkConfigTemplate, constants.ResourceTypeClusterNetwork, name,
		constants.FieldCommonName, name,
		constants.FieldCommonDescription, description,
		constants.FieldClusterNetworkEnable, enable,
		constants.FieldClusterNetworkDefaultPhysicalNIC, defaultPhysicalNIC)
}

func TestAccClusterNetwork_basic(t *testing.T) {
	var (
		clusterNetwork *harvsternetworkv1.ClusterNetwork
		ctx            = context.Background()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ResourceName:  testAccClusterNetworkResourceName,
				ImportState:   true,
				ImportStateId: testAccClusterNetworkName,
				Destroy:       false,
				Config:        buildClusterNetworkConfig(testAccClusterNetworkName, testAccClusterNetworkDescription, testAccClusterNetworkEnable, testAccClusterNetworkDefaultPhysicalNIC),
				Check: resource.ComposeTestCheckFunc(
					testAccClusterNetworkExists(ctx, testAccClusterNetworkResourceName, clusterNetwork),
					resource.TestCheckResourceAttr(testAccClusterNetworkResourceName, constants.FieldCommonName, testAccClusterNetworkName),
					resource.TestCheckResourceAttr(testAccClusterNetworkResourceName, constants.FieldCommonDescription, testAccClusterNetworkDescription),
					resource.TestCheckResourceAttr(testAccClusterNetworkResourceName, constants.FieldClusterNetworkEnable, testAccClusterNetworkEnable),
					resource.TestCheckResourceAttr(testAccClusterNetworkResourceName, constants.FieldClusterNetworkDefaultPhysicalNIC, testAccClusterNetworkDefaultPhysicalNIC),
				),
			},
		},
	})
}

func testAccClusterNetworkExists(ctx context.Context, n string, clusterNetwork *harvsternetworkv1.ClusterNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource %s ID not set. ", n)
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*client.Client)

		_, name, err := helper.IDParts(id)
		if err != nil {
			return err
		}
		foundClusterNetwork, err := c.HarvesterNetworkClient.NetworkV1beta1().ClusterNetworks().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		clusterNetwork = foundClusterNetwork
		return nil
	}
}
