package netbox

import (
	"errors"
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/virtualization"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetboxClusterType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetboxClusterTypeRead,
		Schema: map[string]*schema.Schema{
			"cluster_type_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceNetboxClusterTypeRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	name := d.Get("name").(string)
	params := virtualization.NewVirtualizationClusterTypesListParams()
	params.Name = &name
	limit := int64(2) // Limit of 2 is enough
	params.Limit = &limit

	res, err := api.Virtualization.VirtualizationClusterTypesList(params, nil)
	if err != nil {
		return err
	}

	if *res.GetPayload().Count > int64(1) {
		return errors.New("More than one result. Specify a more narrow filter")
	}
	if *res.GetPayload().Count == int64(0) {
		return errors.New("No result")
	}
	result := res.GetPayload().Results[0]
	d.Set("cluster_type_id", result.ID)
	d.SetId(strconv.FormatInt(result.ID, 10))
	d.Set("name", result.Name)
	return nil
}
