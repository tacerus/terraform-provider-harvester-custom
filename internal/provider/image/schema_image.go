package image

import (
	harvsterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/tacerus/terraform-provider-harvester-custom/internal/util"
	"github.com/tacerus/terraform-provider-harvester-custom/pkg/constants"
)

const (
	URLDescription = "supports the `raw` and `qcow2` image formats which are supported by [qemu](https://www.qemu.org/docs/master/system/images.html#disk-image-file-formats). Bootable ISO images can also be used and are treated like `raw` images."
)

func Schema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		constants.FieldImageDisplayName: {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		constants.FieldImageURL: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			Description:  URLDescription,
		},
		constants.FieldImagePVCNamespace: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: util.IsValidName,
		},
		constants.FieldImagePVCName: {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: util.IsValidName,
		},
		constants.FieldImageSourceType: {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				harvsterv1.VirtualMachineImageSourceTypeDownload,
				harvsterv1.VirtualMachineImageSourceTypeUpload,
				harvsterv1.VirtualMachineImageSourceTypeExportVolume,
			}, false),
		},
		constants.FieldImageProgress: {
			Type:     schema.TypeInt,
			Computed: true,
		},
		constants.FieldImageSize: {
			Type:     schema.TypeInt,
			Computed: true,
		},
		constants.FieldImageStorageClassName: {
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
