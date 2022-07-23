package clusternetwork

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tacerus/tf-harvester-custom/internal/util"
	"github.com/tacerus/tf-harvester-custom/pkg/constants"
)

func Schema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		constants.FieldClusterNetworkEnable: {
			Type:     schema.TypeBool,
			Required: true,
		},
		constants.FieldClusterNetworkDefaultPhysicalNIC: {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	util.NonNamespacedSchemaWrap(s)
	return s
}

func DataSourceSchema() map[string]*schema.Schema {
	return util.DataSourceSchemaWrap(Schema())
}
