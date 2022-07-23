package keypair

import (
	"context"

	harvsterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tacerus/terraform-provider-harvester-custom/internal/util"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/client"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/constants"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/helper"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/importer"
)

func ResourceKeypair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeypairCreate,
		ReadContext:   resourceKeypairRead,
		DeleteContext: resourceKeypairDelete,
		UpdateContext: resourceKeypairUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: Schema(),
	}
}

func resourceKeypairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	namespace := d.Get(constants.FieldCommonNamespace).(string)
	name := d.Get(constants.FieldCommonName).(string)
	toCreate, err := util.ResourceConstruct(d, Creator(namespace, name))
	if err != nil {
		return diag.FromErr(err)
	}
	obj, err := c.HarvesterClient.HarvesterhciV1beta1().KeyPairs(namespace).Create(ctx, toCreate.(*harvsterv1.KeyPair), metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(resourceKeyPairImport(d, obj))
}

func resourceKeypairUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	namespace, name, err := helper.IDParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	obj, err := c.HarvesterClient.HarvesterhciV1beta1().KeyPairs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	toUpdate, err := util.ResourceConstruct(d, Updater(obj))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.HarvesterClient.HarvesterhciV1beta1().KeyPairs(namespace).Update(ctx, toUpdate.(*harvsterv1.KeyPair), metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeypairRead(ctx, d, meta)
}

func resourceKeypairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	namespace, name, err := helper.IDParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	obj, err := c.HarvesterClient.HarvesterhciV1beta1().KeyPairs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	return diag.FromErr(resourceKeyPairImport(d, obj))
}

func resourceKeypairDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	namespace, name, err := helper.IDParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = c.HarvesterClient.HarvesterhciV1beta1().KeyPairs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceKeyPairImport(d *schema.ResourceData, obj *harvsterv1.KeyPair) error {
	stateGetter, err := importer.ResourceKeyPairStateGetter(obj)
	if err != nil {
		return err
	}
	return util.ResourceStatesSet(d, stateGetter)
}
