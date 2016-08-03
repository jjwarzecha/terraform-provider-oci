package main

import (
	"github.com/MustWin/baremtlclient"
	"github.com/hashicorp/terraform/helper/schema"
)

// ResourceIdentityUser exposes a IdentityUser Resource
func ResourceIdentityUser() *schema.Resource {
	return &schema.Resource{
		Create: createUser,
		Read:   readUser,
		Update: updateUser,
		Delete: destroyUser,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"compartment_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_modified": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createUser(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(BareMetalClient)

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	var res *baremtlsdk.Resource
	if res, e = client.CreateUser(name, description); e != nil {
		return
	}

	// Set the id and set any fields that were returned by the API.
	d.SetId(res.ID)
	setResourceData(d, res)

	if res.State != baremtlsdk.ResourceCreated {
		res, e = waitForStateRefresh(d, client, client.GetUser)
	}

	return
}

func readUser(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(BareMetalClient)

	var res *baremtlsdk.Resource
	if res, e = client.GetUser(d.Id()); e != nil {
		return
	}

	setResourceData(d, res)

	return
}

func updateUser(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(BareMetalClient)

	desc := d.Get("description").(string)

	d.Partial(true)
	var res *baremtlsdk.Resource
	if res, e = client.UpdateUser(d.Id(), desc); e != nil {
		return
	}
	d.Partial(false)

	setResourceData(d, res)

	return
}

func destroyUser(d *schema.ResourceData, m interface{}) (e error) {
	client := m.(BareMetalClient)

	if e = client.DeleteUser(d.Id()); e != nil {
		return
	}

	return
}