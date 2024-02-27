import { Construct } from "constructs";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";

export class LinkABooDb extends Construct {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new dynamodb.Table(this, "UrlsTable", {
      tableName: "urls",
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: "shortenedUrl",
        type: dynamodb.AttributeType.STRING,
      },
    });
  }
}
