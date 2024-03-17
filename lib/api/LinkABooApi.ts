import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import { LinkABooDb } from "../db/LinkABooDb";

export type LinkABooApiProps = {
  db: LinkABooDb;
};

export class LinkABooApi extends Construct {
  constructor(scope: Construct, id: string, props: LinkABooApiProps) {
    super(scope, id);

    const api = new apigateway.RestApi(this, "LinkABooApi", {
      restApiName: "Link-a-boo api",
    });

    // shorten
    const shortenLambda = new lambda.Function(this, "shortenLambda", {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromDockerBuild("lambda/go", {
        buildArgs: {
          CMD_NAME: "shorten",
        },
      }),
      handler: "shorten",
      environment: {
        URLS_TABLE_NAME: props.db.urlsTable.tableName,
      },
    });

    props.db.urlsTable.grantReadWriteData(shortenLambda);

    const shortenIntegration = new apigateway.LambdaIntegration(shortenLambda);
    const shortenResource = api.root.addResource("shorten");
    shortenResource.addMethod("POST", shortenIntegration);

    // validate
    const validateLambda = new lambda.Function(this, "validateLambda", {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromDockerBuild("lambda/go", {
        buildArgs: {
          CMD_NAME: "validate",
        },
      }),
      handler: "validate",
      environment: {
        URLS_TABLE_NAME: props.db.urlsTable.tableName,
      },
    });

    props.db.urlsTable.grantReadWriteData(validateLambda);

    const validateIntegration = new apigateway.LambdaIntegration(
      validateLambda,
    );
    const validateResource = api.root.addResource("validate");
    validateResource.addMethod("GET", validateIntegration);

    // redirect
    const redirectLambda = new lambda.Function(this, "redirectLambda", {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromDockerBuild("lambda/go", {
        buildArgs: {
          CMD_NAME: "redirect",
        },
      }),
      handler: "redirect",
      environment: {
        URLS_TABLE_NAME: props.db.urlsTable.tableName,
      },
    });

    props.db.urlsTable.grantReadData(redirectLambda);

    const redirectIntegration = new apigateway.LambdaIntegration(
      redirectLambda,
    );
    const redirectResource = api.root.addResource("{slug}");
    redirectResource.addMethod("GET", redirectIntegration);
  }
}
