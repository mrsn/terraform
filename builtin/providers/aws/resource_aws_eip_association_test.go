package aws

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAWSEipAssociation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAWSEipAssociationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSEipAssociationConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccAWSEipAssociationConfig = `
resource "aws_vpc" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "aws_instance" "foo" {
    ami = "ami-ef5b69df"
    instance_type = "t1.micro"
}

resource "aws_eip" "bar" {
    vpc = true
}

resource "aws_eip_association" "eip-assoc-foobar" {
    allocation_id = "${aws_eip.bar.id}"
    instance_id = "${aws_instance.foo.id}"
}
`
