import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { LinkABooApi } from "./api/LinkABooApi";
// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class LinkABooStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new LinkABooApi(this, id);
  }
}
