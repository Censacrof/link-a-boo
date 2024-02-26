import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";

export class LinkABooApi extends Construct {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    const api = new apigateway.RestApi(this, "LinkABooApi", {
      restApiName: "Link-a-boo api",
    });

    const shortenLambda = new lambda.Function(this, "shortenLambda", {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromDockerBuild("lambda/shorten"),
      handler: "shorten",
    });

    const shortenIntegration = new apigateway.LambdaIntegration(shortenLambda);

    const shortenResource = api.root.addResource("shorten");
    shortenResource.addMethod("POST", shortenIntegration);
  }
}
