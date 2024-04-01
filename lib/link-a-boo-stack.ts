import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { LinkABooApi } from "./api/LinkABooApi";
import { LinkABooDb } from "./db/LinkABooDb";
import { LinkABooRedoc } from "./redoc/Redoc";
// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class LinkABooStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const db = new LinkABooDb(this, "LinkABooDb");
    new LinkABooApi(this, "LinkABooApi", {
      db,
    });

    new LinkABooRedoc(this, "LinkABooRedoc");
  }
}
