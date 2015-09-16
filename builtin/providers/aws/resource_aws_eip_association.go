package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAwsEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsEipAssociationCreate,
		Read:   resourceAwsEipAssociationRead,
		Update: resourceAwsEipAssociationUpdate,
		Delete: resourceAwsEipAssociationDelete,

		Schema: map[string]*schema.Schema{
			"allocation_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAwsEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn

	opts := getAwsEipAssociationInput(d)
	resp, err := ec2conn.AssociateAddress(&opts)

	if err != nil {
		return err
	}

	d.SetId(*resp.AssociationId)
	log.Printf("[INFO] AssociationId: %v", *resp.AssociationId)

	return resourceAwsEipAssociationRead(d, meta)
}

func resourceAwsEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn

	EIPAssociationFilter := &ec2.Filter{
		Name:   aws.String("association-id"),
		Values: []*string{aws.String(d.Id())},
	}

	opts := &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{EIPAssociationFilter},
	}

	resp, err := ec2conn.DescribeAddresses(opts)

	if err != nil {
		return err
	}

	if len(resp.Addresses) != 1 {
		return fmt.Errorf("[ERROR] Error finding EIPAssociation: %s", d.Id())
	}

	if resp == nil {
		d.SetId("")
		return nil
	}

	address := resp.Addresses[0]
	d.Set("allocation_id", address.AllocationId)
	d.Set("public_ip", address.PublicIp)
	d.Set("instance_id", address.InstanceId)
	d.Set("network_interface_id", address.NetworkInterfaceId)
	d.Set("private_ip", address.PrivateIpAddress)

	return nil
}

func resourceAwsEipAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn

	log.Printf("[INFO] Disassociate an elastic ip: %s", d.Id())

	DisAssociateAddressInput := &ec2.DisassociateAddressInput{
		AssociationId: aws.String(d.Id()),
	}

	_, err := ec2conn.DisassociateAddress(DisAssociateAddressInput)

	if err != nil {
		return err
	}

	return nil
}

func getAwsEipAssociationInput(d *schema.ResourceData) ec2.AssociateAddressInput {
	EIPAssociationOpts := ec2.AssociateAddressInput{}

	if allocation_id_value, ok := d.GetOk("allocation_id"); ok {
		EIPAssociationOpts.AllocationId = aws.String(allocation_id_value.(string))
	}

	if public_ip_value, ok := d.GetOk("public_ip"); ok {
		EIPAssociationOpts.PublicIp = aws.String(public_ip_value.(string))
	}

	if instance_id_value, ok := d.GetOk("instance_id"); ok {
		EIPAssociationOpts.InstanceId = aws.String(instance_id_value.(string))
	}

	if network_interface_id_value, ok := d.GetOk("network_interface_id"); ok {
		EIPAssociationOpts.NetworkInterfaceId = aws.String(network_interface_id_value.(string))
	}

	if private_ip_value, ok := d.GetOk("private_ip"); ok {
		EIPAssociationOpts.PrivateIpAddress = aws.String(private_ip_value.(string))
	}

	return EIPAssociationOpts
}
