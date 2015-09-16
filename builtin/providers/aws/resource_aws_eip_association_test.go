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
				Config: testAccAWSEipAssociationWithVPCConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccAWSEipAssociationWithVPCConfig = `
resource "aws_vpc" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "aws_internet_gateway" "gw" {
	vpc_id = "${aws_vpc.foo.id}"
}

resource "aws_subnet" "test" {
    vpc_id = "${aws_vpc.foo.id}"
    cidr_block = "10.1.0.0/24"
}

resource "aws_instance" "foo" {
    ami = "ami-ef5b69df"
    instance_type = "t1.micro"
    subnet_id = "${aws_subnet.test.id}"
}

resource "aws_eip" "bar" {
    vpc = true
}

resource "aws_eip_association" "eip-assoc-foobar" {
	allocation_id = "${aws_eip.bar.id}"
	instance_id = "${aws_instance.foo.id}"
}
`
