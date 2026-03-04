resource "terraform_data" "test" {}

resource "aws_instance" "web1" {
  count = 1
  ami = 4321
  lifecycle {
    ignore_changes = [tags]
  }
  depends_on = []
}

resource "aws_instance" "web2" {
  ami = 1234
  availability_zone = "custom-az"
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = 0.002
    }
  }
}
