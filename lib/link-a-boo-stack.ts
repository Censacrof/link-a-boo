import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class LinkABooStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const handler = new lambda.Function(this, "root-handler-lambda", {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromDockerBuild("lambda/shorten"),
      handler: "shorten",
    });
  }
}
