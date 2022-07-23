package clusternetwork

import (
	"context"
	"errors"

	harvsternetworkv1 "github.com/harvester/harvester-network-controller/pkg/apis/network.harvesterhci.io/v1beta1"
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

func ResourceClusterNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterNetworkCreate,
		ReadContext:   resourceClusterNetworkRead,
		DeleteContext: resourceClusterNetworkDelete,
		UpdateContext: resourceClusterNetworkUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: Schema(),
	}
}

func resourceClusterNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get(constants.FieldCommonName).(string)
	switch name {
	case "vlan":
		return diag.FromErr(errors.New("can not create the existing vlan clusternetwork, to avoid this error and continue with the plan, use `terraform import harvester_clusternetwork.vlan vlan` to import it first"))
	default:
		return diag.FromErr(errors.New("can not create clusternetwork, to avoid this error and continue with the plan, either move clusternetwork to another module or reduce the scope of the plan using the -target flag"))
	}
}

func resourceClusterNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	_, name, err := helper.IDParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	obj, err := c.HarvesterNetworkClient.NetworkV1beta1().ClusterNetworks().Get(ctx, name, metav1.GetOptions{})
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
	_, err = c.HarvesterNetworkClient.NetworkV1beta1().ClusterNetworks().Update(ctx, toUpdate.(*harvsternetworkv1.ClusterNetwork), metav1.UpdateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceClusterNetworkRead(ctx, d, meta)
}

func resourceClusterNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	_, name, err := helper.IDParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	obj, err := c.HarvesterNetworkClient.NetworkV1beta1().ClusterNetworks().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	return diag.FromErr(resourceClusterNetworkImport(d, obj))
}

func resourceClusterNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.FromErr(errors.New("clusternetwork should not be destroyed, to avoid this error and continue with the plan, either move clusternetwork to another module or reduce the scope of the plan using the -target flag"))
}

func resourceClusterNetworkImport(d *schema.ResourceData, obj *harvsternetworkv1.ClusterNetwork) error {
	stateGetter, err := importer.ResourceClusterNetworkStateGetter(obj)
	if err != nil {
		return err
	}
	return util.ResourceStatesSet(d, stateGetter)
}
