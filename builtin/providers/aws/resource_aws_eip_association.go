package aws

import (
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
			// Allocation ID for the VPC Elastic IP address to be associated
			"allocation_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAwsEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	ec2conn := meta.(*AWSClient).ec2conn

	log.Printf(
		"[INFO] Eip association: %s => %s",
		d.Get("allocation_id").(string),
		d.Get("instance_id").(string))

	AssociateAddressInput := &ec2.AssociateAddressInput{
		AllocationId: aws.String(d.Get("allocation_id").(string)),
		InstanceId:   aws.String(d.Get("instance_id").(string)),
	}

	resp, err := ec2conn.AssociateAddress(AssociateAddressInput)

	if err != nil {
		return err
	}

	d.SetId(*resp.AssociationId)
	log.Printf("[INFO] Association ID: %s", d.Id())

	return nil
}

func resourceAwsEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
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
