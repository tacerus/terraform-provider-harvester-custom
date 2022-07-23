package keypair

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tacerus/tf-harvester-custom/pkg/client"
	"github.com/tacerus/tf-harvester-custom/pkg/constants"
)

func DataSourceKeypair() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeypairRead,
		Schema:      DataSourceSchema(),
	}
}

func dataSourceKeypairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	namespace := d.Get(constants.FieldCommonNamespace).(string)
	name := d.Get(constants.FieldCommonName).(string)
	keyPair, err := c.HarvesterClient.HarvesterhciV1beta1().KeyPairs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(resourceKeyPairImport(d, keyPair))
}
