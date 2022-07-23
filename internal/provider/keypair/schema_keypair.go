package keypair

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/tacerus/terraform-provider-harvester-custom/internal/util"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/constants"
)

func Schema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		constants.FieldKeyPairPublicKey: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		constants.FieldKeyPairFingerPrint: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	util.NamespacedSchemaWrap(s, false)
	return s
}

func DataSourceSchema() map[string]*schema.Schema {
	return util.DataSourceSchemaWrap(Schema())
}
