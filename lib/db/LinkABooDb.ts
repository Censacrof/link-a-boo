import { Construct } from "constructs";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";

export class LinkABooDb extends Construct {
  public readonly urlsTable: dynamodb.Table;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    this.urlsTable = new dynamodb.Table(this, "UrlsTable", {
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: "slug",
        type: dynamodb.AttributeType.STRING,
      },
    });
  }
}
