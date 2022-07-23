package virtualmachine

import (
	"github.com/harvester/harvester/pkg/builder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/tacerus/terraform-provider-harvester-custom/pkg/constants"
)

func resourceNetworkInterfaceSchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		constants.FiledNetworkInterfaceName: {
			Type:     schema.TypeString,
			Required: true,
		},
		constants.FiledNetworkInterfaceType: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				builder.NetworkInterfaceTypeBridge,
				builder.NetworkInterfaceTypeMasquerade,
				"",
			}, false),
		},
		constants.FiledNetworkInterfaceModel: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "virtio",
			ValidateFunc: validation.StringInSlice([]string{
				"virtio",
				"e1000",
				"e1000e",
				"ne2k_pco",
				"pcnet",
				"rtl8139",
			}, false),
		},
		constants.FiledNetworkInterfaceMACAddress: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		constants.FiledNetworkInterfaceIPAddress: {
			Type:     schema.TypeString,
			Computed: true,
		},
		constants.FiledNetworkInterfaceWaitForLease: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "wait for this network interface to obtain an IP address",
		},
		constants.FiledNetworkInterfaceInterfaceName: {
			Type:     schema.TypeString,
			Computed: true,
		},
		constants.FiledNetworkInterfaceNetworkName: {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}
